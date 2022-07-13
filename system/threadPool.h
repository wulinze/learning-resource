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

namespace ThreadPool{
    class ThreadPool{
    public:
        ThreadPool(size_t n);
        template<typename F, typename... Args>
        auto add(F&& f, Args&&... args) 
            -> std::future<typename std::result_of<F(Args...)>::type>;
        ~ThreadPool();
    private:
        bool isStopped = false;
        std::mutex mu_;
        std::queue<std::function<void>> task_queue;
        std::vector<std::thread> workers;
        std::condition_variable cv_;
    };


    inline ThreadPool::ThreadPool(size_t n){
        workers.reserve(n);

        for(int i=0; i<n; i++){
            workers.emplace_back(
                [this]{
                    for(;;){
                        std::function<void()> task;

                        {
                            std::unique_lock<std::mutex> lock(this->mu_);
                            cv_.wait(lock, 
                            [this]{return this->isStopped || !this->task_queue.empty();});
                            if(this->isStopped && this->task_queue.empty()){
                                return;
                            }

                            task = std::move(this->task_queue.front());
                            this->task_queue.pop();
                        }

                        task();
                    }
                }
            );
        }
    }

    template<typename F, typename... Args>
    auto ThreadPool::add(F&& f, Args&&... args) 
        -> std::future<typename std::result_of<F(Args...)>::type>
    {
        using return_type = std::result_of<F(Args...)>::type;

        auto task = std::make_shared<std::packaged_task<return_type()>(
            std::bind(std::forward<F>(f), std::forward<Args>(args)...)
        );

        std::future<return_type> res = task->get_future();
        {
            std::unique_lock<std::mutex> lock(mu_);

            if(isStopped){
                throw std::runtime_error("ThreadPool already stopped");
            }

            task_queue.emplace([task](){(*task)();})
        }
        cv_.notify_one();

        return res;
    }

    inline ThreadPool::~ThreadPool(){
        {
            std::unique_lock<std::mutex> lock(mu_);
            isStopped = true;
        }

        cv_.notify_all();

        for(auto&& worker : workers){
            worker.join();
        }
    }
}

