/**
 * 自定义vector实现
 */

#pragma once
#include <cstdint>
#include <memory>
#include <optional>
#include <stdexcept>

template <typename T>
class Vector {
public:
    Vector(size_t cnt = 0)
    {
        capacity_ = cnt == 0 ? 2 : cnt;  // 容量至少为2
        data_ = std::make_unique<T[]>(capacity_);
        cnt_ = cnt;  // 实际元素数量等于传入的cnt
    }

    Vector(const Vector<T>& other)
    {
        data_ = std::make_unique<T[]>(other.capacity_);
        cnt_ = other.cnt_;
        capacity_ = other.capacity_;
        for (size_t i = 0; i < cnt_; ++i) {
            data_[i] = other.data_[i];
        }
    }

    Vector(Vector<T>&& other)
    {
        data_ = std::move(other.data_);
        cnt_ = other.cnt_;
        capacity_ = other.capacity_;

        other.cnt_ = 0;
        other.capacity_ = 0;
    }

    Vector<T>& operator=(const Vector<T>& other)
    {
        if (this == &other) {
            return *this;
        }

        data_ = std::make_unique<T[]>(other.capacity_);
        cnt_ = other.cnt_;
        capacity_ = other.capacity_;
        for (size_t i = 0; i < cnt_; ++i) {
            data_[i] = other.data_[i];
        }
        return *this;
    }

    Vector<T>& operator=(Vector<T>&& other)
    {
        if (this == &other) {
            return *this;
        }
        
        data_ = std::move(other.data_);
        cnt_ = other.cnt_;
        capacity_ = other.capacity_;

        other.cnt_ = 0;
        other.capacity_ = 0;
        return *this;
    }
    ~Vector() = default;

    void push_back(const T& item)
    {
        if (cnt_ >= capacity_) {
            realloc();
        }
        data_[cnt_] = item;
        ++cnt_;
    }

    T& at(size_t index)
    {
        if (index >= cnt_) {
            throw std::out_of_range("Vector::at: index " + std::to_string(index) + " >= size " + std::to_string(cnt_));
        }
        return data_[index];
    }

    // 下标访问运算符 - 非const版本
    T& operator[](size_t index)
    {
        return data_[index];  // 不做边界检查，与STL行为一致
    }

    // 下标访问运算符 - const版本
    const T& operator[](size_t index) const
    {
        return data_[index];  // 不做边界检查，与STL行为一致
    }

    T& front()
    {
        if (empty()) {
            throw std::runtime_error("Vector::front: vector is empty");
        }
        return data_[0];
    }

    T& back()
    {
        if (empty()) {
            throw std::runtime_error("Vector::back: vector is empty");
        }
        return data_[cnt_ - 1];
    }

    void pop_back()
    {
        if (empty()) {
            throw std::runtime_error("Vector::pop_back: vector is empty");
        }
        data_[cnt_ - 1].~T();
        --cnt_;
    }

    void resize(size_t n)
    {
        if (n > capacity_) {
            realloc(n);
        }

        if (n == 0) {
            clear();
            return;
        }

        // 如果缩小，析构多余的元素
        if (n < cnt_) {
            for (size_t i = n; i < cnt_; ++i) {
                data_[i].~T();
            }
        }

        cnt_ = n;
    }

    T* data() { return data_.get(); }

    bool empty() { return cnt_ == 0; }

    size_t size() const { return cnt_; }

    void reserve(size_t n) {
        if (n <= capacity_) {
            return;
        }

        realloc(n);
    }

    size_t capacity() const { return capacity_; }

    void clear()
    {
        // 调用所有元素的析构函数
        for (size_t i = 0; i < cnt_; ++i) {
            data_[i].~T();
        }
        cnt_ = 0;
        // 保持capacity不变，不释放内存
    }

private:
    void realloc(std::optional<size_t> newCap = std::nullopt)
    {
        auto newCapacity = newCap.has_value() ? newCap.value() : capacity_ * 2;
        auto newData = std::make_unique<T[]>(newCapacity);
        if constexpr (std::is_move_constructible_v<T>) {
            for (size_t i = 0; i < cnt_; ++i) {
                newData[i] = std::move(data_[i]);
            }
        } else {
            std::copy(data_.get(), data_.get() + cnt_, newData.get());
        }

        data_.swap(newData);
        capacity_ = newCapacity;
    }

private:
    std::unique_ptr<T[]> data_;
    size_t cnt_ { 0 };
    size_t capacity_;
};
