package main

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// var stats *Stats

type Stats struct {
	CPU         float64
	DiskUsed    uint64
	DiskTotal   uint64
	DiskPercent float64
	MemoryUsed  uint64
	MemoryTotal uint64
	GoMemory    uint64
	GoRoutines  int

	lastCPUStat *cpu.TimesStat
}

func (s *Stats) DisplayStats() {
	fmt.Printf("%+v\n", s)
}

func (s *Stats) LoadStats() {

	// count cpu cycles between last count
	if stats, err := cpu.Times(false); err == nil {
		stat := stats[0]
		total := totalCPUTime(stat)
		last := s.lastCPUStat
		if last != nil {
			lastTotal := totalCPUTime(*last)
			if lastTotal != 0 {
				totalDelta := total - lastTotal
				if totalDelta > 0 {
					idleDelta := (stat.Iowait + stat.Idle) - (last.Iowait + last.Idle)
					usedDelta := (totalDelta - idleDelta)
					s.CPU = 100 * usedDelta / totalDelta
				}
			}
		}
		s.lastCPUStat = &stat
	}

	// count disk usage
	if stat, err := disk.Usage("."); err == nil {
		s.DiskUsed = stat.Used
		s.DiskTotal = stat.Total
		s.DiskPercent = stat.UsedPercent
	}

	// count memory usage
	if stat, err := mem.VirtualMemory(); err == nil {
		s.MemoryUsed = stat.Used
		s.MemoryTotal = stat.Total
	}

	// count total bytes allocated by the go runtime
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)
	s.GoMemory = memStats.Alloc

	// count current number of goroutines
	s.GoRoutines = runtime.NumGoroutine()

}

func totalCPUTime(t cpu.TimesStat) float64 {
	total := t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Guest + t.GuestNice + t.Idle
	return total
}
