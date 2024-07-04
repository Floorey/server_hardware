package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

const (
	defaultInterval    = 5 * time.Second
	defaultLogFileName = "hardware_log.txt"
)

func logData(file *os.File, interval time.Duration, wg *sync.WaitGroup, stop <-chan struct{}) {
	defer wg.Done()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			// CPU Usage
			cpuPercents, err := cpu.Percent(0, false)
			if err != nil {
				logLine := fmt.Sprintf("Error getting CPU percent: %v\n", err)
				writer.WriteString(logLine)
				fmt.Print(logLine)
			} else {
				logLine := fmt.Sprintf("CPU Usage: %v%%\n", cpuPercents[0])
				writer.WriteString(logLine)
				fmt.Print(logLine)
			}

			// Memory Usage
			vmStat, err := mem.VirtualMemory()
			if err != nil {
				logLine := fmt.Sprintf("Error getting memory info: %v\n", err)
				writer.WriteString(logLine)
				fmt.Print(logLine)
			} else {
				logLine := fmt.Sprintf("Memory Usage: %v%%\n", vmStat.UsedPercent)
				writer.WriteString(logLine)
				fmt.Print(logLine)
			}

			// Disk Usage
			diskStat, err := disk.Usage("/")
			if err != nil {
				logLine := fmt.Sprintf("Error getting disk usage: %v\n", err)
				writer.WriteString(logLine)
				fmt.Print(logLine)
			} else {
				logLine := fmt.Sprintf("Disk Usage: %v%%\n", diskStat.UsedPercent)
				writer.WriteString(logLine)
				fmt.Print(logLine)
			}

			writer.Flush()
		}
	}
}

func main() {
	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Get log file name from environment variable or use default
	logFileName := os.Getenv("LOG_FILE_NAME")
	if logFileName == "" {
		logFileName = defaultLogFileName
	}

	// Open log file
	logFile, err := os.Create(logFileName)
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// Get log interval from environment variable or use default
	intervalStr := os.Getenv("LOG_INTERVAL")
	interval := defaultInterval
	if intervalStr != "" {
		if i, err := strconv.Atoi(intervalStr); err == nil {
			interval = time.Duration(i) * time.Second
		} else {
			fmt.Printf("Invalid LOG_INTERVAL, using default: %v\n", defaultInterval)
		}
	}

	// Goroutine for monitoring and logging system stats
	wg.Add(1)
	go logData(logFile, interval, &wg, stop)

	// Goroutine for capturing user input
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(os.Stdin)
		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				continue
			}
			if char == 'q' {
				close(stop)
				return
			}
		}
	}()

	// Goroutine for handling system interrupts
	wg.Add(1)
	go func() {
		defer wg.Done()
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		close(stop)
	}()

	wg.Wait()
	fmt.Println("Program terminated")
}
