#include "lock_free_queue.h"
#include <atomic>
#include <catch2/catch.hpp>
#include <chrono>
#include <thread>
#include <vector>

// 测试单线程基本功能
TEST_CASE("LockFreeQueue single thread basic operations", "[LFQ]")
{
    LockFreeQueue<int> queue;

    SECTION("Test empty queue")
    {
        REQUIRE(queue.empty() == true);
        REQUIRE(queue.size() == 0);
    }

    SECTION("Test push and front")
    {
        queue.push(42);
        REQUIRE(queue.empty() == false);
        REQUIRE(queue.front() == 42);
    }

    SECTION("Test push and pop")
    {
        queue.push(42);
        REQUIRE(queue.pop() == true);
        REQUIRE(queue.empty() == true);
    }

    SECTION("Test multiple push and pop")
    {
        for (int i = 0; i < 10; i++) {
            queue.push(i);
        }

        for (int i = 0; i < 10; i++) {
            REQUIRE(queue.front() == i);
            REQUIRE(queue.pop() == true);
        }

        REQUIRE(queue.empty() == true);
    }
}

// 测试多生产者单消费者场景
TEST_CASE("LockFreeQueue multiple producers single consumer", "[LFQ]")
{
    LockFreeQueue<int> queue;
    const int producer_count = 4;
    const int items_per_producer = 1000;

    std::vector<std::thread> producers;

    // 创建多个生产者线程
    for (int p = 0; p < producer_count; p++) {
        producers.emplace_back([&queue, p, items_per_producer]() {
            for (int i = 0; i < items_per_producer; i++) {
                queue.push(p * items_per_producer + i);
            }
        });
    }

    // 创建一个消费者线程
    std::atomic<int> consumed_count(0);
    std::thread consumer([&queue, &consumed_count, total = producer_count * items_per_producer]() {
        while (consumed_count < total) {
            if (!queue.empty()) {
                queue.front();  // 读取前端元素但不移除
                queue.pop();
                consumed_count++;
            } else {
                // 避免忙等待
                std::this_thread::yield();
            }
        }
    });

    // 等待所有线程完成
    for (auto& producer : producers) {
        producer.join();
    }
    consumer.join();

    REQUIRE(queue.empty() == true);
    REQUIRE(consumed_count == producer_count * items_per_producer);
}

// 测试单生产者多消费者场景
TEST_CASE("LockFreeQueue single producer multiple consumers", "[LFQ]")
{
    LockFreeQueue<int> queue;
    const int consumer_count = 4;
    const int total_items = 10000;

    // 创建多个消费者线程
    std::vector<std::thread> consumers;
    std::atomic<int> produced_count(0);
    std::atomic<int> consumed_count(0);

    for (int c = 0; c < consumer_count; c++) {
        consumers.emplace_back([&queue, &consumed_count, total_items]() {
            while (consumed_count < total_items) {
                if (!queue.empty()) {
                    queue.front();  // 读取前端元素但不移除
                    if (queue.pop()) {
                        consumed_count++;
                    }
                } else {
                    // 避免忙等待
                    std::this_thread::yield();
                }
            }
        });
    }

    // 创建一个生产者线程
    std::thread producer([&queue, &produced_count, total_items]() {
        for (int i = 0; i < total_items; i++) {
            queue.push(i);
            produced_count++;
        }
    });

    // 等待所有线程完成
    producer.join();
    for (auto& consumer : consumers) {
        consumer.join();
    }

    REQUIRE(queue.empty() == true);
    REQUIRE(produced_count == total_items);
    REQUIRE(consumed_count == total_items);
}

// 测试多生产者多消费者场景
TEST_CASE("LockFreeQueue multiple producers multiple consumers", "[LFQ]")
{
    LockFreeQueue<int> queue;
    const int producer_count = 4;
    const int consumer_count = 4;
    const int items_per_producer = 5000;
    const int total_items = producer_count * items_per_producer;

    std::vector<std::thread> producers;
    std::vector<std::thread> consumers;
    std::atomic<int> consumed_count(0);

    // 创建多个生产者线程
    for (int p = 0; p < producer_count; p++) {
        producers.emplace_back([&queue, p, items_per_producer]() {
            for (int i = 0; i < items_per_producer; i++) {
                queue.push(p * items_per_producer + i);
                // 小延迟增加竞争条件出现的概率
                if (i % 100 == 0) {
                    std::this_thread::sleep_for(std::chrono::microseconds(1));
                }
            }
        });
    }

    // 创建多个消费者线程
    for (int c = 0; c < consumer_count; c++) {
        consumers.emplace_back([&queue, &consumed_count, total_items]() {
            while (consumed_count < total_items) {
                if (!queue.empty()) {
                    queue.front();  // 读取前端元素但不移除
                    if (queue.pop()) {
                        consumed_count++;
                    }
                } else {
                    // 避免忙等待
                    std::this_thread::yield();
                }
            }
        });
    }

    // 等待所有线程完成
    for (auto& producer : producers) {
        producer.join();
    }
    for (auto& consumer : consumers) {
        consumer.join();
    }

    REQUIRE(queue.empty() == true);
    REQUIRE(consumed_count == total_items);
}

// 测试队列性能和压力
TEST_CASE("LockFreeQueue performance and stress test", "[LFQ]")
{
    LockFreeQueue<int> queue;
    const int producer_count = 8;
    const int consumer_count = 8;
    const int items_per_producer = 10000;
    const int total_items = producer_count * items_per_producer;

    std::vector<std::thread> producers;
    std::vector<std::thread> consumers;
    std::atomic<int> consumed_count(0);

    auto start_time = std::chrono::high_resolution_clock::now();

    // 创建多个生产者线程
    for (int p = 0; p < producer_count; p++) {
        producers.emplace_back([&queue, items_per_producer]() {
            for (int i = 0; i < items_per_producer; i++) {
                queue.push(i);
            }
        });
    }

    // 创建多个消费者线程
    for (int c = 0; c < consumer_count; c++) {
        consumers.emplace_back([&queue, &consumed_count, total_items]() {
            while (consumed_count < total_items) {
                if (!queue.empty()) {
                    if (queue.pop()) {
                        consumed_count++;
                    }
                } else {
                    // 避免忙等待
                    std::this_thread::yield();
                }
            }
        });
    }

    // 等待所有线程完成
    for (auto& producer : producers) {
        producer.join();
    }
    for (auto& consumer : consumers) {
        consumer.join();
    }

    auto end_time = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end_time - start_time).count();

    CAPTURE(duration);
    REQUIRE(queue.empty() == true);
    REQUIRE(consumed_count == total_items);

    INFO("Processed " << total_items << " items in " << duration << " ms");
}

// 测试边界情况
TEST_CASE("LockFreeQueue edge cases", "[LFQ]")
{
    LockFreeQueue<int> queue;

    SECTION("Test popping from empty queue") { REQUIRE(queue.pop() == false); }

    SECTION("Test front on empty queue") { REQUIRE_THROWS(queue.front()); }

    SECTION("Test end on empty queue")
    {
        // 注意：当前实现中end()方法可能有问题，因为tail_可能指向dummy节点
        REQUIRE_NOTHROW(queue.end());
    }

    SECTION("Test high frequency push and pop")
    {
        std::atomic<bool> running(true);

        std::thread pusher([&queue, &running]() {
            while (running) {
                queue.push(42);
            }
        });

        std::thread popper([&queue, &running]() {
            int count = 0;
            while (count < 100000) {
                if (!queue.empty()) {
                    if (queue.pop()) {
                        count++;
                    }
                } else {
                    std::this_thread::yield();
                }
            }
            running = false;
        });

        pusher.join();
        popper.join();
    }
}
