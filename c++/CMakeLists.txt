cmake_minimum_required(VERSION 3.22)
project(pratice)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)
FetchContent_Declare(
    catch2
    GIT_REPOSITORY https://github.com/catchorg/Catch2.git
    GIT_TAG v2.13.7
)
FetchContent_MakeAvailable(catch2)

enable_testing()

set(OUTPUT_NAME cpp_tests)

add_executable(${OUTPUT_NAME})
target_sources(${OUTPUT_NAME} PRIVATE
    timer/timer_with_epoll.cpp
    timer/timer_with_epoll_test.cpp
    lock_free_queue/lock_free_queue_test.cpp
    ipc/channel_test.cpp

    stl/vector_test.cpp
)

target_link_libraries(${OUTPUT_NAME} PRIVATE Catch2::Catch2)
target_include_directories(${OUTPUT_NAME} PRIVATE ./ ./util)

add_test(NAME Catch2Tests COMMAND tests)

