package main

import (
	"fmt"
	"log"
	"math"
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

	fmt.Printf("🔥 Starting AGGRESSIVE CPU burner: %s (PID: %d)\n", appName, pid)
	fmt.Printf("💥 This will consume maximum CPU on %d cores\n", runtime.NumCPU())

	// 使用所有 CPU 核心
	numWorkers := runtime.NumCPU()

	var wg sync.WaitGroup

	// 启动状态报告
	go statusReporter(appName, pid)

	// 启动CPU燃烧器
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			cpuBurner(workerID, appName)
		}(i)
	}

	fmt.Printf("🚀 Started %d aggressive CPU burners\n", numWorkers)
	fmt.Println("⚠️  This will cause HIGH CPU usage - this is intentional!")

	wg.Wait()
}

// 激进的 CPU 消耗函数 - 无休息时间
func cpuBurner(workerID int, appName string) {
	counter := 0
	result := 0.0

	for {
		// 连续执行大量计算，没有任何休息时间
		for i := 0; i < 100000; i++ {
			// 浮点运算密集型计算
			x := float64(i)
			result += math.Sin(x) * math.Cos(x)
			result = math.Sqrt(math.Abs(result))
			result = math.Pow(result, 1.5)

			// 整数运算
			counter = (counter*31 + i) % 1000000
		}

		// 每 1000 轮输出一次（但不休息）
		if counter%1000000 == 0 {
			fmt.Printf("🔥 Worker %d (%s): burning CPU... result=%.2f\n",
				workerID, appName, result)
		}
	}
}

// 状态报告
func statusReporter(appName string, pid int) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-ticker.C:
			uptime := time.Since(startTime)
			fmt.Printf("\n💥 %s AGGRESSIVE CPU BURNER REPORT:\n", appName)
			fmt.Printf("   🔥 PID: %d\n", pid)
			fmt.Printf("   ⏱️  Uptime: %v\n", uptime.Round(time.Second))
			fmt.Printf("   🚀 Active CPU burners: %d\n", runtime.NumCPU())
			fmt.Printf("   ⚠️  Expected: VERY HIGH CPU usage\n")
			fmt.Printf("   📊 This is designed to test PMGO's CPU monitoring\n\n")
		}
	}
}
