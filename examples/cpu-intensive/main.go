package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <app-name>")
	}

	appName := os.Args[1]
	pid := os.Getpid()

	fmt.Printf("🔥 Starting CPU-intensive app: %s (PID: %d)\n", appName, pid)
	fmt.Printf("📊 Using %d CPU cores\n", runtime.NumCPU())

	// 启动 CPU 密集型任务
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU() // 使用所有可用 CPU 核心

	fmt.Printf("🚀 Starting %d CPU workers...\n", numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			cpuIntensiveWork(workerID, appName)
		}(i)
	}

	// 启动状态报告协程
	go statusReporter(appName, pid)

	// 等待所有工作协程完成 (永远不会完成，除非被信号终止)
	wg.Wait()
}

// CPU 密集型工作函数
func cpuIntensiveWork(workerID int, appName string) {
	iteration := 0
	for {
		// 执行一些 CPU 密集型计算
		result := 0
		for i := 0; i < 1000000; i++ {
			// 随机数学运算来消耗 CPU
			x := rand.Float64() * 100
			y := rand.Float64() * 100

			// 复杂的数学计算
			result += int(x*y) % 1000
			result = (result * 17) % 999983 // 大质数取模
		}

		iteration++

		// 每 100 次迭代休息一小段时间，避免完全占用 CPU
		if iteration%100 == 0 {
			time.Sleep(10 * time.Millisecond)
		}

		// 每 1000 次迭代输出一次状态
		if iteration%1000 == 0 {
			fmt.Printf("🔥 Worker %d: %s completed %d iterations (result: %d)\n",
				workerID, appName, iteration, result)
		}
	}
}

// 状态报告函数
func statusReporter(appName string, pid int) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-ticker.C:
			uptime := time.Since(startTime)
			fmt.Printf("📊 %s Status Report:\n", appName)
			fmt.Printf("   - PID: %d\n", pid)
			fmt.Printf("   - Uptime: %v\n", uptime.Round(time.Second))
			fmt.Printf("   - Goroutines: %d\n", runtime.NumGoroutine())
			fmt.Printf("   - CPU Workers: %d (high CPU usage expected)\n", runtime.NumCPU())
			fmt.Printf("   - Memory Usage: CPU-intensive workload running\n")
			fmt.Println("🔥 High CPU usage is intentional for testing PMGO monitoring")
		}
	}
}
