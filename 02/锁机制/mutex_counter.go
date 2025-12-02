package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int64 = 0 // 共享计数器
	var mu sync.Mutex     // 互斥锁
	var wg sync.WaitGroup // 等待所有协程完成

	const numGoroutines = 10
	const incrementsPerGoroutine = 1000

	wg.Add(numGoroutines)

	// 启动 10 个协程
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				mu.Lock()   // 加锁
				counter++   // 临界区
				mu.Unlock() // 解锁
			}
			fmt.Printf("协程 %d 完成\n", id)
		}(i)
	}

	wg.Wait() // 等待所有协程结束
	fmt.Printf("最终计数器值: %d (预期: %d)\n", counter, numGoroutines*incrementsPerGoroutine)
}
