#pragma once

#include <iostream>
#include <memory>
#include <queue>
#include <condition_variable>
#include <functional>
#include <thread>
#include <stdexcept>
#include <future>
#include <mutex>
#include <vector>

namespace threadPool{
    template<typename F>
    struct task{
        static void run(){

        }
    };

    class threadPool{
    public:
        threadPool(size_t n);
        template<typename F, typename... Args> 
        void add(task<>)
        ~threadPool();
    private:
        bool isStopped = false;
        std::mutex mu_;
        std::queue<task<>> task_queue;
        std::vector<std::thread> workers;
        std::condition_variable cv_;
    };
}
