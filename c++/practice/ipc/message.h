#pragma once
#include <cstdint>
#include <string>

struct Message {
    int32_t id;
    char data[128] {};
};