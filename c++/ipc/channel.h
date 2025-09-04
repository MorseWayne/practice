/**
 * 进程间数据通道抽象
 * 需要子类重新实现读写方法
 */

#pragma once

#include <optional>
#include "message.h"

class Channel {
public:
    virtual bool Write(Message& msg) = 0;
    virtual std::optional<Message> Read() = 0;
};
