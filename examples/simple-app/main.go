package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	appName := "simple-test-app"
	if len(os.Args) > 1 {
		appName = os.Args[1]
	}

	log.Printf("🚀 Starting %s (PID: %d)", appName, os.Getpid())

	// Print startup information
	fmt.Printf("=== %s ===\n", appName)
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Printf("Start Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("Arguments: %v\n", os.Args)

	// Main loop - simulates a long-running process
	counter := 0
	for {
		counter++
		log.Printf("📊 %s is running... (iteration %d)", appName, counter)

		// Simulate some work
		time.Sleep(5 * time.Second)

		// Print some status every 10 iterations
		if counter%10 == 0 {
			fmt.Printf("Status: %s has been running for %d iterations\n", appName, counter)
		}
	}
}
