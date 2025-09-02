/**
 * 无锁队列实现
 */
#pragma once

#include <atomic>
#include <cstdint>
#include <memory>
#include <queue>

template <typename T>
class LockFreeQueue {
private:
    struct Node {
        T val;
        std::atomic<Node*> next;
    };

public:
    LockFreeQueue()
    {
        auto dummy = new Node { -1, nullptr };
        head_.store(dummy);
        tail_.store(dummy);
        cnts_ = 0;
    };

    ~LockFreeQueue() = default;

    void push(T val)
    {
        auto newNode = new Node { val, nullptr };
        while (true) {
            // peek 当前时刻的队列尾部
            auto tailPeek = tail_.load();
            auto next = tailPeek->next.load();
            if (next == nullptr) {
                // 执行完上一句后，线程可能会发生切出又切回，导致队列尾已经变化，需要利用CAS先查看本线程看到的是不是真的队列尾部
                // 如果是，更新队列尾部，不是，重新寻找队尾，这一系列行为通过 compare_exchange_weak 完成
                if (tailPeek->next.compare_exchange_weak(next, newNode)) {
                    // 更新tail_: 执行完上一步cas后，还是可能会发生线程切换，更新tail_指针仍然要用cas操作
                    tail_.compare_exchange_weak(tailPeek, newNode);
                    break;
                }
            } else {
                // 利用cas操作帮助其他线程更新队尾指针
                tail_.compare_exchange_weak(tailPeek, next);
            }
        }
    }

    bool pop()
    {
        while (true) {
            auto headPeek = head_.load();
            auto toDelete = headPeek->next.load();
            // 需要判断队列是否为空，防止其他线程已经将队列清空了之后这里拿到一个空指针继续向下执行
            if (toDelete == nullptr) {
                return false;
            }
            // 利用CAS将要删除的节点作为新的dummy节点
            if (head_.compare_exchange_weak(headPeek, toDelete)) {
                cnts_--;
                // 原先的head指针需要被释放, 但是不能直接delete, 因为headPeek可能是其他线程正在使用的节点
                safe_delete(headPeek);
                break;
            }
        }

        return true;
    }

    T front() {
        auto head = head_.load();
        auto next = head->next.load();
        if (next == nullptr) {
            throw std::runtime_error("queue is empty");
        }
        return next->val;
    }

    T end() {
        auto tail = tail_.load();
        if (tail == nullptr) {
            throw std::runtime_error("queue is empty");
        }
        return tail->val;
    };
    int size() { return cnts_; };
    bool empty() { return head_ == tail_; };
    void safe_delete(Node* ptr)
    {
        // TO DO
    }

private:
    std::atomic<Node*> head_;
    std::atomic<Node*> tail_;
    std::atomic<int> cnts_;
};