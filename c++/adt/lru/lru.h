/**
 * LRU实现
 */

#pragma once

#include <cstddef>
#include <list>
#include <unordered_map>

struct Node {
    int key {};
    int value {};
};

using Iterator = std::list<Node>::iterator;
class LRUCache {
public:
    LRUCache(int capacity) : capacity_(capacity) {}

    int get(int key)
    {
        if (!keys_.contains(key)) {
            return -1;
        }
        flush(key);
        return keys_.at(key)->value;
    }

    void put(int key, int value)
    {
        if (keys_.contains(key)) {
            keys_[key]->value = value;
            flush(key);
            return;
        }

        if (values_.size() == capacity_) {
            auto item = values_.back();
            keys_.erase(item.key);
            values_.pop_back();
        }

        // 统一插入逻辑：总是插入到最前面
        values_.insert(values_.begin(), Node { key, value });
        keys_[key] = values_.begin();
    }

private:
    void flush(int key)
    {
        auto itr = keys_[key];
        // 使用splice优化：直接移动节点而不是复制
        values_.splice(values_.begin(), values_, itr);
        keys_[key] = values_.begin();
    }

private:
    std::list<Node> values_;
    std::unordered_map<int, Iterator> keys_;
    size_t capacity_ { 0 };
};