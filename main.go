package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
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

type SystemStats struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `jason:"memory_usage"`
	DiskUsage   float64 `jason:"disk_usage"`
}

var currentStats SystemStats
var statsMutex sync.RWMutex

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
				statsMutex.Lock()
				currentStats.CPUUsage = cpuPercents[0]
				statsMutex.Unlock()
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
				statsMutex.Lock()
				currentStats.MemoryUsage = vmStat.UsedPercent
				statsMutex.Unlock()
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
				statsMutex.Lock()
				currentStats.DiskUsage = diskStat.UsedPercent
				statsMutex.Unlock()
			}

			writer.Flush()
		}
	}
}
func statsHandler(w http.ResponseWriter, r *http.Request) {
	statsMutex.RLock()
	defer statsMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentStats)
}

func main() {
	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Open log file
	logFileName := os.Getenv("LOG_FILE_NAME")
	if logFileName == "" {
		logFileName = "hardware_log.txt"
	}
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

	// Start HTTP server
	http.HandleFunc("/stats", statsHandler)
	server := &http.Server{Addr: ":8080"}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	wg.Wait()
	fmt.Println("Program terminated")
}
