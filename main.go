package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	var wg sync.WaitGroup
	stop := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				cpuPercents, err := cpu.Percent(0, false)
				if err != nil {
					fmt.Printf("Error getting CPU percent: %v\n", err)
				} else {
					fmt.Printf("CPU Usage: %v%%\n", cpuPercents[0])
				}

				vmStat, err := mem.VirtualMemory()
				if err != nil {
					fmt.Printf("Error getting memory info: %v\n", err)
				} else {
					fmt.Printf("Memory usage: %v%%\n", vmStat.UsedPercent)
				}

				time.Sleep(5 * time.Second)

			}
		}
	}()

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
	wg.Wait()
	fmt.Println("Exit program")
}
