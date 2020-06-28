// +build linux

package os

func CurCpuPercent() float64 {
	return 0
}

func CurMemoryUsage() float64 {
	return 0
}

func CpuPercent(pid int) float64 {
	return 0
}

func MemoryUsage(pid int) float64 {
	return 0
}
