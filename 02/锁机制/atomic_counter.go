package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter int64 = 0 // 共享计数器（原子变量）
	var wg sync.WaitGroup

	const numGoroutines = 10
	const incrementsPerGoroutine = 1000

	wg.Add(numGoroutines)

	// 启动 10 个协程
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				atomic.AddInt64(&counter, 1) // 原子递增
			}
			fmt.Printf("协程 %d 完成\n", id)
		}(i)
	}

	// 读取最终值（也用原子读取）
	wg.Wait()
	final := atomic.LoadInt64(&counter)
	fmt.Printf("最终计数器值: %d (预期: %d)\n", final, numGoroutines*incrementsPerGoroutine)
}
