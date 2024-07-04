package main

import (
	"testing"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func TestCPUUsage(t *testing.T) {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		t.Errorf("Error getting CPU percent: %v", err)
	}
	if len(cpuPercents) == 0 {
		t.Error("No CPU usage data returned")
	}

}
func TestMemoryUsage(t *testing.T) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		t.Errorf("Error getting memory info: %v", err)
	}
	if vmStat.Total == 0 {
		t.Error("Nomemory usage data returned")
	}
}
func TestDiskUs(t *testing.T) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		t.Errorf("Error getting disk usage: %v", err)
	}
	if diskStat.Total == 0 {
		t.Errorf("No disk usage data returned")
	}
}
