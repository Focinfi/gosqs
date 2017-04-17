package psutil

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var virtualMemoryStat *mem.VirtualMemoryStat

func init() {
	stat, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}

	virtualMemoryStat = stat
}

// MemoryUsedPercent returns the memory used percent
func MemoryUsedPercent() int {
	return int(virtualMemoryStat.UsedPercent)
}

// CPUUsedPercent returns the cpu used percent
func CPUUsedPercent() (int, error) {
	res, err := cpu.Percent(time.Millisecond*100, false)
	if err != nil || len(res) == 0 {
		return 100, err
	}

	return int(res[0]), nil
}
