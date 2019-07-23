// +build windows

package os

import (
	"os/exec"
	"strconv"
)

const (
	MONITOR_CPU int = 1
	MONITOR_MEM     = 2
)

func CurCpuPercent() float64 {
	return CpuPercent(PID())
}

func CurMemoryUsage() float64 {
	return MemoryUsage(PID())
}

func CpuPercent(pid int) float64 {
	ret, err := doMonitor(pid, MONITOR_CPU)
	if err != nil {
		return -1.0
	}
	result, err := strconv.ParseFloat(string(ret), 64)
	if err != nil {
		return -2.0
	}
	if result < 0.0 {
		result = 0.0
	}
	return result
}

func MemoryUsage(pid int) float64 {
	ret, err := doMonitor(pid, MONITOR_MEM)
	if err != nil {
		return -1.0
	}
	result, err := strconv.ParseFloat(string(ret), 64)
	if err != nil {
		return -2.0
	}
	if result < 0.0 {
		result = 0.0
	}
	return result
}

func doMonitor(pid int, mod int) ([]byte, error) {
	cmd := exec.Command("monitor", strconv.Itoa(pid), strconv.Itoa(mod))
	buf, err := cmd.Output()
	return buf, err
}
