package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
func TestStatsHandler(t *testing.T) {
	// set initial data
	statsMutex.Lock()
	currentStats = SystemStats{
		CPUUsage:    10.5,
		MemoryUsage: 65.3,
		DiskUsage:   75.4,
	}
	statsMutex.Unlock()

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatalf("Cloud not create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(statsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returened wrong ststus code: got %v want %v",
			status, http.StatusOK)
	}
	expected := SystemStats{
		CPUUsage:    10.5,
		MemoryUsage: 65.3,
		DiskUsage:   75.4,
	}
	var actual SystemStats
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("Cloud not decode responce: %v", err)
	}
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}
}
