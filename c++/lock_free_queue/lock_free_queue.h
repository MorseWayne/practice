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
        std::atomic<std::shared_ptr<Node>> next;
    };

public:
    LockFreeQueue(/* args */) { dummy_ = std::make_shared<Node>(); };
    ~LockFreeQueue();

    void push(T val);
    void pop();
    T front();
    T end();
    int size();
    bool empty() { return tail_ == dummy_; };

private:
    std::atomic<std::shared_ptr<Node>> head_;
    std::atomic<std::shared_ptr<Node>> tail_;
    std::shared_ptr<Node> dummy_;
    std::atomic<int> cnts_;
};