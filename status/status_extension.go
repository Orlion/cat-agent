package status

import (
	"fmt"
	"strconv"

	"github.com/Orlion/cat-agent/pkg/stringx"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type CpuStatusExtension struct {
	lastTime    *cpu.TimesStat
	lastCPUTime float64
}

func newCpuStatusExtension() *CpuStatusExtension {
	return &CpuStatusExtension{
		lastTime:    new(cpu.TimesStat),
		lastCPUTime: 0,
	}
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
}

func newMemStatusExtension() *MemStatusExtension {
	return &MemStatusExtension{}
}

func (ext *MemStatusExtension) GetId() string {
	return "mem"
}

func (ext *MemStatusExtension) GetDesc() string {
	return "mem"
}

func (ext *MemStatusExtension) GetProperties() map[string]string {
	m := make(map[string]string)
	if stat, err := mem.VirtualMemory(); err == nil {
		m = map[string]string{
			"total":     strconv.FormatUint(stat.Total, 10),
			"available": strconv.FormatUint(stat.Available, 10),
			"used":      strconv.FormatUint(stat.Used, 10),
			"free":      strconv.FormatUint(stat.Free, 10),
			"shared":    strconv.FormatUint(stat.Shared, 10),
			"buffers":   strconv.FormatUint(stat.Buffers, 10),
			"cached":    strconv.FormatUint(stat.Cached, 10),
		}
	}

	return m
}

type NetStatusExtension struct {
}

func newNetStatusExtension() *NetStatusExtension {
	return &NetStatusExtension{}
}

func (ext *NetStatusExtension) GetId() string {
	return "net"
}

func (ext *NetStatusExtension) GetDesc() string {
	return "net"
}

func (ext *NetStatusExtension) GetProperties() map[string]string {
	m := make(map[string]string)

	if stats, err := net.IOCounters(false); err == nil && len(stats) > 0 {
		m[fmt.Sprintf("net.%s.sent_bytes", stats[0].Name)] = strconv.FormatUint(stats[0].BytesSent, 10)
		m[fmt.Sprintf("net.%s.recv_bytes", stats[0].Name)] = strconv.FormatUint(stats[0].BytesRecv, 10)
		m[fmt.Sprintf("net.%s.sent_packets", stats[0].Name)] = strconv.FormatUint(stats[0].PacketsSent, 10)
		m[fmt.Sprintf("net.%s.recv_packets", stats[0].Name)] = strconv.FormatUint(stats[0].PacketsRecv, 10)
		m[fmt.Sprintf("net.%s.errin", stats[0].Name)] = strconv.FormatUint(stats[0].Errin, 10)
		m[fmt.Sprintf("net.%s.errout", stats[0].Name)] = strconv.FormatUint(stats[0].Errout, 10)
		m[fmt.Sprintf("net.%s.dropin", stats[0].Name)] = strconv.FormatUint(stats[0].Dropin, 10)
		m[fmt.Sprintf("net.%s.dropout", stats[0].Name)] = strconv.FormatUint(stats[0].Dropout, 10)
		m[fmt.Sprintf("net.%s.fifoin", stats[0].Name)] = strconv.FormatUint(stats[0].Fifoin, 10)
		m[fmt.Sprintf("net.%s.fifoout", stats[0].Name)] = strconv.FormatUint(stats[0].Fifoout, 10)
	}

	return m
}
