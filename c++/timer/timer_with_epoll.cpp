#include "timer_with_epoll.h"
#include <iostream>

TimerWithEpoll::TimerWithEpoll(std::chrono::seconds initialDelay, std::chrono::seconds interval,
                               std::function<void(uint64_t)> callback)
    : initialDelay_(initialDelay),
      interval_(interval),
      callback_(std::move(callback))
{
}

UniqueFd TimerWithEpoll::CreateTimer()
{
    auto timerFd = timerfd_create(CLOCK_MONOTONIC, 0);
    if (timerFd == -1) {
        std::cerr << "Failed to create timer" << std::endl;
        return nullptr;
    }

    auto timerFdPtr = UniqueFd(new int(timerFd));

    // POSIX定时器配置
    itimerspec timerSpec {};

    // 启动定时器后的首次触发时间，比如这里我希望我的定时器3s后才首次触发
    timerSpec.it_value.tv_sec = initialDelay_.count();
    timerSpec.it_value.tv_nsec = 0;

    // 定时器的周期触发时间设置，比如下面的设置时每隔1s就会触发定时器
    timerSpec.it_interval.tv_sec = interval_.count();
    timerSpec.it_interval.tv_nsec = 0;

    /**
     * 设置或重置 timerfd 定时器的超时时间和周期。
     * 参数1：fd, 由 timerfd_create 返回的定时器文件描述符
     *
     * 参数2：flags，常用为 0，或 TFD_TIMER_ABSTIME（表示
     * it_value是绝对时间[比如这个值时一个系统时钟的时间]，否则为相对时间
     *
     * 参数3：new_value, 新的定时器超时时间和周期
     *
     * 参数4：old_value, 如果不为 nullptr，调用前定时器的设置会被写入这里（可用于获取上一次的定时器设置），
     * 否则可传nullptr
     */
    if (timerfd_settime(*timerFdPtr, 0, &timerSpec, nullptr) == -1) {
        std::cerr << "Failed to set timer" << std::endl;
        return nullptr;
    }

    return timerFdPtr;
}

UniqueFd TimerWithEpoll::CreateEpoll()
{
    auto timerFd = CreateTimer();
    if (timerFd == nullptr) {
        return nullptr;
    }
    timerFd_.swap(timerFd);

    int epollFd = epoll_create1(0);
    if (epollFd == -1) {
        std::cerr << "Failed to create epoll instance" << std::endl;
        return nullptr;
    }

    auto epollFdPtr = UniqueFd(new int(epollFd));
    epoll_event event {};
    event.events = EPOLLIN;
    // 自定义的带回参数，它会在 epoll 事件返回时带回给你，epoll 只负责原样带回
    // 这里你可以设置为任意值
    event.data.fd = *timerFd_;  // Use member variable timerFd_

    auto err = epoll_ctl(epollFd, EPOLL_CTL_ADD, *timerFd_, &event);
    if (err == -1) {
        std::cerr << "Failed to add timer to epoll" << std::endl;
        return nullptr;
    }

    return epollFdPtr;
}

bool TimerWithEpoll::Start()
{
    if (running_) {
        return false;
    }
    auto epollFd = CreateEpoll();
    if (epollFd == nullptr) {
        return false;
    }
    epollFd_.swap(epollFd);
    running_ = true;

    // 创建异步任务执行定时器函数
    listenThread_ = std::thread([this] { Listen(); });
    return true;
}

bool TimerWithEpoll::Listen()
{
    epoll_event events[10] {};
    while (running_) {
        int nfds = epoll_wait(*epollFd_, events, sizeof(events), -1);
        if (nfds == -1) {
            std::cerr << "Failed to wait for epoll events, errno: " << errno << std::endl;
            return false;
        }

        auto now = std::chrono::system_clock::now();
        auto seconds = std::chrono::duration_cast<std::chrono::seconds>(now.time_since_epoch()).count();
        std::cerr << "Epoll event occurres, unix time: " << seconds << std::endl;

        // 处理到来的事件
        for (int i = 0; i < nfds; ++i) {
            if (events[i].data.fd == *timerFd_) {
                /**
                 * 读取定时器的超时事件，获取定时器截止到目前未被读取(也就是下面的read操作)的总次数
                 * 必须执行这个read操作，否则 epoll 会一直报告它可读，定时事件不会被“消耗”掉
                 * 读取的数据类型是 uint64_t，表示"自上次 read以来定时器到期了多少次"。
                 * 如果你的处理慢，可能会积累多次到期。
                 * 举例：
                 *  如果定时器每秒触发一次，你 3 秒后才调用 read，expirations 可能就是 3
                 *  如果你每次都及时 read，expirations 就是 1
                 */
                uint64_t expirations = 0;
                read(*timerFd_, &expirations, sizeof(expirations));
                // 执行业务callback
                callback_(expirations);
            }
        }
    }
    return true;
}

void TimerWithEpoll::Stop()
{
    if (!running_) {
        return;
    }
    running_ = false;

    // 等待监听线程结束
    if (listenThread_.joinable()) {
        listenThread_.join();
    }
}