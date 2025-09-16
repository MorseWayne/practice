
#include <sys/epoll.h>
#include <sys/timerfd.h>
#include <chrono>
#include <cstdint>
#include <functional>
#include <memory>
#include <thread>

// fd 删除器
struct FdCloser {
    void operator()(int* fd) const
    {
        if (fd && *fd != -1) {
            close(*fd);
            delete fd;
        }
    }
};

// 创建智能指针管理的 fd, 使用RAII编程风格防止fd泄露
using UniqueFd = std::unique_ptr<int, FdCloser>;

class TimerWithEpoll {
public:
    /// @brief Constructor
    /// @param initialDelay: 首次启动定时器的相对时间, 单位s
    /// @param interval: 定时器触发的周期间隔，单位s
    TimerWithEpoll(std::chrono::seconds initialDelay, std::chrono::seconds interval,
                   std::function<void(uint64_t)> callback);

    /// @brief Destructor
    ~TimerWithEpoll() = default;

    /// @brief 启动定时器
    /// @return 创建成功返回true
    bool Start();

    /// @brief 停止定时器
    void Stop();

private:
    UniqueFd CreateTimer();
    UniqueFd CreateEpoll();
    bool Listen();

private:
    UniqueFd timerFd_;
    UniqueFd epollFd_;
    std::chrono::seconds initialDelay_;
    std::chrono::seconds interval_;
    std::function<void(uint64_t)> callback_;
    std::thread listenThread_;
    bool running_ { false };
};