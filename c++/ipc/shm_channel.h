/**
 * 使用共享内存实现的数据读写
 */
#pragma once
#include <semaphore.h>
#include <sys/mman.h>
#include <sys/shm.h>
#include <sys/stat.h>
#include <memory>
#include "channel.h"
#include "unique_fd.h"

constexpr const char* const SHM_NAME = "my_shared_memory";  // 共享内存名称，系统唯一
constexpr const char* const SEM_READ = "my_sem_read";  // 信号量名称，系统唯一
constexpr const char* const SEM_WRITE = "my_sem_write";  // 信号量名称，系统唯一
constexpr uint32_t SHM_SIZE = sizeof(Message);  // 共享内存大小，单位字节

struct MmapDeleter {
    void operator()(void* p) const
    {
        if (p) {
            munmap(p, SHM_SIZE);
        }
    }
};

struct SemDeleter {
    void operator()(sem_t* p) const
    {
        if (p) {
            sem_close(p);
        }
    }
};

using MmapPtr = std::unique_ptr<void, MmapDeleter>;
using SemPtr = std::unique_ptr<sem_t, SemDeleter>;

class ShmChannel : public Channel {
public:
    int Init(bool isWriter)
    {
        isWriter_ = isWriter;
        // 0666是一个比较宽泛的权限，它允许其他任意知道本共享内存名称的进程对这块儿数据进行读写
        int fd = shm_open(SHM_NAME, O_CREAT | O_RDWR, 0666);
        if (fd == -1) {
            perror("open shm failed");
            return -1;
        }

        shmFd_ = UniqueFd(new int(fd));
        if (ftruncate(fd, SHM_SIZE) == -1) {
            perror("set shm size failed");
            return -1;
        }

        void* mapped_addr = mmap(NULL, SHM_SIZE, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
        if (mapped_addr == MAP_FAILED) {
            perror("mmap failed");
            return -1;
        }
        shm_ = MmapPtr(mapped_addr);

        sem_t* readSem {};
        sem_t* writeSem {};
        if (isWriter_) {
            readSem = sem_open(SEM_READ, O_CREAT, 0666, 0);  // 信号量初始值设置为1，代表可写
            writeSem = sem_open(SEM_WRITE, O_CREAT, 0666, 1);  // 信号量初始值设置为0，代表不可读
        } else {
            readSem = sem_open(SEM_READ, 0);
            writeSem = sem_open(SEM_WRITE, 0);
        }
        if (readSem == SEM_FAILED || writeSem == SEM_FAILED) {
            perror("open sem failed");
            return -1;
        }
        semWrite_ = SemPtr(writeSem);
        semRead_ = SemPtr(readSem);

        return 0;
    };

    bool Write(Message& msg) override
    {
        if (!semWrite_ || !semRead_ || !shm_) {
            return false;
        }

        sem_wait(semWrite_.get());
        auto ptr = reinterpret_cast<Message*>(shm_.get());
        *ptr = msg;
        sem_post(semRead_.get());

        return true;
    }

    std::optional<Message> Read() override
    {
        if (!semWrite_ || !semRead_ || !shm_) {
            return std::nullopt;
        }

        sem_wait(semRead_.get());
        auto ptr = reinterpret_cast<Message*>(shm_.get());
        sem_post(semWrite_.get());
        return *ptr;
    }

private:
    UniqueFd shmFd_;
    MmapPtr shm_;
    SemPtr semRead_;
    SemPtr semWrite_;
    bool isWriter_ { false };
};
