package status

import (
	"runtime"
	"strconv"

	"github.com/Orlion/cat-agent/pkg/stringx"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

type CpuStatusExtension struct {
	lastTime    *cpu.TimesStat
	lastCPUTime float64
}

func (ext *CpuStatusExtension) GetId() string {
	return "cpu"
}

func (ext *CpuStatusExtension) GetDesc() string {
	return "cpu"
}

func (ext *CpuStatusExtension) GetProperties() map[string]string {
	m := make(map[string]string)

	if avg, err := load.Avg(); err == nil {
		m["load.1min"] = stringx.F642str(avg.Load1)
		m["load.5min"] = stringx.F642str(avg.Load5)
		m["load.15min"] = stringx.F642str(avg.Load15)
		m["system.load.average"] = m["load.1min"]
	}

	if times, err := cpu.Times(false); err == nil {
		if len(times) > 0 {
			currentTime := times[0]

			currentCpuTime := 0.0 +
				currentTime.User +
				currentTime.System +
				currentTime.Idle +
				currentTime.Nice +
				currentTime.Iowait +
				currentTime.Irq +
				currentTime.Softirq +
				currentTime.Steal +
				currentTime.Guest +
				currentTime.GuestNice

			if ext.lastCPUTime > 0 {
				cpuTime := currentCpuTime - ext.lastCPUTime

				if cpuTime > 0.0 {
					user := currentTime.User - ext.lastTime.User
					system := currentTime.System - ext.lastTime.System
					nice := currentTime.Nice - ext.lastTime.Nice
					idle := currentTime.Idle - ext.lastTime.Idle
					iowait := currentTime.Iowait - ext.lastTime.Iowait
					softirq := currentTime.Softirq - ext.lastTime.Softirq
					irq := currentTime.Irq - ext.lastTime.Irq
					steal := currentTime.Steal - ext.lastTime.Steal

					m["cpu.user"] = stringx.F642str(user)
					m["cpu.sys"] = stringx.F642str(system)
					m["cpu.nice"] = stringx.F642str(nice)
					m["cpu.idle"] = stringx.F642str(idle)
					m["cpu.iowait"] = stringx.F642str(iowait)
					m["cpu.softirq"] = stringx.F642str(softirq)
					m["cpu.irq"] = stringx.F642str(irq)
					m["cpu.steal"] = stringx.F642str(steal)

					m["cpu.user.percent"] = stringx.F642str(user / cpuTime * 100)
					m["cpu.sys.percent"] = stringx.F642str(system / cpuTime * 100)
					m["cpu.nice.percent"] = stringx.F642str(nice / cpuTime * 100)
					m["cpu.idle.percent"] = stringx.F642str(idle / cpuTime * 100)
					m["cpu.iowait.percent"] = stringx.F642str(iowait / cpuTime * 100)
					m["cpu.softirq.percent"] = stringx.F642str(softirq / cpuTime * 100)
					m["cpu.irq.percent"] = stringx.F642str(irq / cpuTime * 100)
					m["cpu.steal.percent"] = stringx.F642str(steal / cpuTime * 100)
				}
			}
			ext.lastCPUTime = currentCpuTime
			ext.lastTime = &currentTime
		}
	}

	return m
}

type MemStatusExtension struct {
	m runtime.MemStats

	alloc,
	mallocs,
	lookups,
	frees uint64
}

func (ext *MemStatusExtension) GetId() string {
	return "mem.runtime"
}

func (ext *MemStatusExtension) GetDesc() string {
	return "mem.runtime"
}

func (ext *MemStatusExtension) GetProperties() map[string]string {
	runtime.ReadMemStats(&ext.m)

	m := map[string]string{
		"mem.sys": stringx.B2kbstr(ext.m.Sys),

		// heap
		"mem.heap.alloc":    stringx.B2kbstr(ext.m.HeapAlloc),
		"mem.heap.sys":      stringx.B2kbstr(ext.m.HeapSys),
		"mem.heap.idle":     stringx.B2kbstr(ext.m.HeapIdle),
		"mem.heap.inuse":    stringx.B2kbstr(ext.m.HeapInuse),
		"mem.heap.released": stringx.B2kbstr(ext.m.HeapReleased),
		"mem.heap.objects":  strconv.Itoa(int(ext.m.HeapObjects)),

		// stack
		"mem.stack.inuse": stringx.B2kbstr(ext.m.StackInuse),
		"mem.stack.sys":   stringx.B2kbstr(ext.m.StackSys),
	}

	if ext.alloc > 0 {
		m["mem.alloc"] = stringx.B2kbstr(ext.m.TotalAlloc - ext.alloc)
		m["mem.mallocs"] = strconv.Itoa(int(ext.m.Mallocs - ext.mallocs))
		m["mem.lookups"] = strconv.Itoa(int(ext.m.Lookups - ext.lookups))
		m["mem.frees"] = strconv.Itoa(int(ext.m.Frees - ext.frees))
	}
	ext.alloc = ext.m.TotalAlloc
	ext.mallocs = ext.m.Mallocs
	ext.lookups = ext.m.Lookups
	ext.frees = ext.m.Frees

	return m
}
