#include <catch2/catch.hpp>
#include <cstring>
#include <semaphore.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

// 在测试中暴露私有成员/方法，便于直接调用 Init
#define private public
#include "shm_channel.h"
#undef private

static void CleanupIpcArtifacts()
{
    // 取消链接共享内存与命名信号量，避免测试间互相影响
    shm_unlink(SHM_NAME);
    sem_unlink(SEM_READ);
    sem_unlink(SEM_WRITE);
}

struct ShmFixture {
    ShmFixture() { CleanupIpcArtifacts(); }
    ~ShmFixture() { CleanupIpcArtifacts(); }

    void InitWriter()
    {
        REQUIRE(writer.Init(true) == 0);
    }

    void InitReader()
    {
        REQUIRE(reader.Init(false) == 0);
    }

    ShmChannel writer;
    ShmChannel reader;
};

TEST_CASE_METHOD(ShmFixture, "ShmChannel init (writer)", "[shm][init]")
{
    InitWriter();
}

TEST_CASE_METHOD(ShmFixture, "ShmChannel init (reader)", "[shm][init]")
{
    InitWriter();
    InitReader();
}

TEST_CASE_METHOD(ShmFixture, "ShmChannel write/read one message", "[shm][rw]")
{
    InitWriter();
    InitReader();

    Message msg {};
    msg.id = 42;
    std::memset(msg.data, 0, sizeof(msg.data));
    std::strncpy(msg.data, "hello shm", sizeof(msg.data) - 1);

    SECTION("写端写入，读端读取")
    {
        REQUIRE(writer.Write(msg) == true);

        auto out = reader.Read();
        REQUIRE(out.has_value());
        REQUIRE(out->id == msg.id);
        REQUIRE(std::strncmp(out->data, msg.data, sizeof(msg.data)) == 0);
    }
}

TEST_CASE_METHOD(ShmFixture, "ShmChannel read/write precondition checks", "[shm][guard]")
{
    ShmChannel ch;
    // 未初始化前，不应读写成功（根据实现预期）
    // 这里的断言可能随实现而调整
    Message msg {1, {0}};
    REQUIRE(ch.Write(msg) == false);
    auto out = ch.Read();
    REQUIRE_FALSE(out.has_value());
}

TEST_CASE_METHOD(ShmFixture, "Repeated initialization and resource cleanup", "[shm][lifecycle]")
{
    ShmChannel writer1;
    REQUIRE(writer1.Init(true) == 0);

    // 再创建一个写端应当也能成功映射同一段共享内存（取决于实现容忍度）
    ShmChannel writer2;
    REQUIRE(writer2.Init(true) == 0);

    // 创建读端
    ShmChannel rd;
    REQUIRE(rd.Init(false) == 0);

    // 退出作用域后析构应当自动释放资源（由系统命名对象引用计数管理），
    // 这里不做显式断言，仅做流程覆盖。
}

TEST_CASE_METHOD(ShmFixture, "ShmChannel cross-process write/read", "[shm][rw][fork]") {
    // 仅父进程清理，避免竞争
    CleanupIpcArtifacts();

    Message msg{};
    msg.id = 100;
    std::strncpy(msg.data, "hello cross proc", sizeof(msg.data) - 1);

    pid_t pid = fork();
    REQUIRE(pid >= 0);

    if (pid == 0) {
        // 子进程：读端
        alarm(3); // 简易超时，防止阻塞
        ShmChannel readerProc;
        if (readerProc.Init(false) != 0) _exit(2);

        auto out = readerProc.Read();
        if (!out.has_value()) _exit(3);
        if (out->id != msg.id) _exit(4);
        if (std::strncmp(out->data, msg.data, sizeof(msg.data)) != 0) _exit(5);

        _exit(0);
    } else {
        // 父进程：写端
        InitWriter(); // 写端初始化
        REQUIRE(writer.Write(msg) == true);

        int status = 0;
        REQUIRE(waitpid(pid, &status, 0) == pid);
        REQUIRE(WIFEXITED(status));
        REQUIRE(WEXITSTATUS(status) == 0);

        // 最后由父进程清理
        CleanupIpcArtifacts();
    }
}