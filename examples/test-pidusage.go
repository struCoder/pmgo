package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/struCoder/pidusage"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run test-pidusage.go <pid>")
	}

	pidStr := os.Args[1]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Fatal("Invalid PID:", err)
	}

	fmt.Printf("Testing pidusage library with PID: %d\n", pid)
	fmt.Println("Monitoring for 30 seconds with detailed analysis...")
	fmt.Printf("%-8s %-12s %-15s %-10s\n", "Sample", "CPU (%)", "Memory (MB)", "Status")
	fmt.Println("--------------------------------------------------")

	var cpuValues []float64
	var memoryValues []float64

	for i := 0; i < 10; i++ {
		sysInfo, err := pidusage.GetStat(pid)
		if err != nil {
			fmt.Printf("%-8d %-12s %-15s %-10s\n",
				i+1, "ERROR", "ERROR", fmt.Sprintf("Error: %v", err))
			continue
		}

		// Store values for analysis
		cpuValues = append(cpuValues, sysInfo.CPU)
		memoryValues = append(memoryValues, sysInfo.Memory)

		// Format memory in MB
		memoryMB := sysInfo.Memory / (1024 * 1024)

		fmt.Printf("%-8d %-12.2f %-15.1f %-10s\n",
			i+1, sysInfo.CPU, memoryMB, "OK")

		time.Sleep(3 * time.Second)
	}

	// Calculate and display statistics
	if len(cpuValues) > 0 {
		fmt.Println("\n📊 Statistics Summary:")
		fmt.Printf("CPU - Min: %.2f%%, Max: %.2f%%, Avg: %.2f%%\n",
			min(cpuValues), max(cpuValues), avg(cpuValues))
		fmt.Printf("Memory - Min: %.1fMB, Max: %.1fMB, Avg: %.1fMB\n",
			min(memoryValues)/(1024*1024),
			max(memoryValues)/(1024*1024),
			avg(memoryValues)/(1024*1024))
	}
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func avg(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
