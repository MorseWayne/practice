#define CATCH_CONFIG_MAIN
#include "timer_with_epoll.h"
#include <atomic>
#include <catch2/catch.hpp>
#include <chrono>
#include <thread>

TEST_CASE("TimerWithEpoll basic periodic callback", "[timer]")
{
    std::atomic<int> call_count { 0 };
    TimerWithEpoll timer(std::chrono::seconds(1), std::chrono::seconds(1),
                         [&](uint64_t expirations) { call_count += expirations; });

    REQUIRE(timer.Start());
    // 等待一段时间，确保定时器触发多次
    std::this_thread::sleep_for(std::chrono::milliseconds(3500));
    timer.Stop();
    // 由于定时器首次1s后触发，后续每1s触发一次，3.5s内应至少触发3次
    REQUIRE(call_count >= 3);
}

TEST_CASE("TimerWithEpoll stop", "[timer]")
{
    std::atomic<int> call_count { 0 };
    TimerWithEpoll timer(std::chrono::seconds(1), std::chrono::seconds(1),
                         [&](uint64_t expirations) { call_count += expirations; });

    REQUIRE(timer.Start());
    std::this_thread::sleep_for(std::chrono::milliseconds(1500));
    timer.Stop();
    int count_after_stop = call_count.load();
    std::this_thread::sleep_for(std::chrono::milliseconds(1500));
    // Stop后不应再有回调
    REQUIRE(call_count == count_after_stop);
}