#pragma once

#include <fcntl.h>
#include <memory>
#include <unistd.h>

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