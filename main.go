package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func logData(file *os.File, wg *sync.WaitGroup, stop <-chan struct{}) {
	defer wg.Done()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for {
		select {
		case <-stop:
			return
		default:
			// CPU Usage
			cpuPercents, err := cpu.Percent(0, false)
			if err != nil {
				fmt.Printf("Error getting CPU percent: %v\n", err)
			} else {
				logLine := fmt.Sprintf("CPU Usage: %v%%\n", cpuPercents[0])
				writer.WriteString(logLine)
				fmt.Print(logLine)
			}

			// Memory Usage
			vmStat, err := mem.VirtualMemory()
			if err != nil {
				fmt.Printf("Error getting memory info: %v\n", err)
			} else {
				logLine := fmt.Sprintf("Memory Usage: %v%%\n", vmStat.UsedPercent)
				writer.WriteString(logLine)
				fmt.Print(logLine)
			}

			// Disk Usage
			diskStat, err := disk.Usage("/")
			if err != nil {
				fmt.Printf("Error getting disk usage: %v\n", err)
			} else {
				logLine := fmt.Sprintf("Disk Usage: %v%%\n", diskStat.UsedPercent)
				writer.WriteString(logLine)
				fmt.Print(logLine)
			}

			writer.Flush()
			time.Sleep(5 * time.Second)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Open log file
	logFile, err := os.Create("hardware_log.txt")
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// Goroutine for monitoring and logging system stats
	wg.Add(1)
	go logData(logFile, &wg, stop)

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
