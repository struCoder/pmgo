package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <app-name>")
		os.Exit(1)
	}

	appName := os.Args[1]
	pid := os.Getpid()

	fmt.Printf("🔥 Starting simple CPU test: %s (PID: %d)\n", appName, pid)

	iteration := 0
	for {
		// 简单但有效的 CPU 密集型计算
		sum := 0.0
		for i := 0; i < 500000; i++ {
			sum += math.Sin(float64(i)) * math.Cos(float64(i))
		}

		iteration++

		if iteration%100 == 0 {
			fmt.Printf("📊 %s: iteration %d, sum=%.2f\n", appName, iteration, sum)
		}

		// 短暂休息，让系统有机会记录 CPU 使用
		time.Sleep(50 * time.Millisecond)
	}
}
