package main

import (
	"fmt"
	"sync"
	"time"
)

// Task 表示一个任务（函数）
type Task func() time.Duration

// Scheduler 任务调度器
type Scheduler struct {
	wg sync.WaitGroup
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// AddTask 添加任务
func (s *Scheduler) AddTask(task Task, name string) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		start := time.Now()
		duration := task() // 执行任务
		elapsed := time.Since(start)
		fmt.Printf("任务 [%s] 执行完成，耗时: %v (任务返回: %v)\n", name, elapsed, duration)
	}()
}

// Wait 等待所有任务完成
func (s *Scheduler) Wait() {
	s.wg.Wait()
	fmt.Println("所有任务执行完毕！")
}

// 示例任务
func task1() time.Duration {
	time.Sleep(1 * time.Second)
	return 100 * time.Millisecond // 模拟返回值
}

func task2() time.Duration {
	time.Sleep(500 * time.Millisecond)
	return 200 * time.Millisecond
}

func task3() time.Duration {
	for i := 0; i < 1e8; i++ { // 模拟 CPU 密集
	}
	return 300 * time.Millisecond
}

func main() {
	scheduler := NewScheduler()

	scheduler.AddTask(task1, "下载文件")
	scheduler.AddTask(task2, "处理数据")
	scheduler.AddTask(task3, "加密计算")

	scheduler.Wait()
}
