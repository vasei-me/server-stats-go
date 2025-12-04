// main.go - Advanced Server Stats Tool in Go (100% Working with gopsutil v4)
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

// ANSI Colors
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Bold    = "\033[1m"
)

var jsonOutput = flag.Bool("json", false, "Output as JSON")

func main() {
	flag.Parse()

	if *jsonOutput {
		printJSON()
		return
	}

	printBanner()
	printSystemOverview()
	printCPUInfo()
	printMemoryInfo()
	printDiskInfo()
	printNetworkInfo()
	printTopProcesses()
	printLoadAverage()
	fmt.Printf("%sGenerated at %s%s\n\n", Cyan, time.Now().Format("2006-01-02 15:04:05"), Reset)
}

func printBanner() {
	fmt.Printf("%s%s╔══════════════════════════════════════════════════════════╗%s\n", Bold, Blue, Reset)
	fmt.Printf("%s║            SERVER PERFORMANCE STATISTICS              ║%s\n", Bold, Reset)
	fmt.Printf("%s║          %s • %s                                 ║%s\n", Bold, time.Now().Format("2006-01-02 15:04:05"), getHostname(), Reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════╝%s\n\n", Blue, Reset)
}

func printSystemOverview() {
	info, err := host.Info()
	if err != nil {
		fmt.Printf("Error getting host info: %v\n", err)
		return
	}
	uptime := time.Duration(info.Uptime) * time.Second
	bootTime := time.Unix(int64(info.BootTime), 0)

	users, _ := host.Users()

	fmt.Printf("%sSystem Information%s\n", Bold, Reset)
	fmt.Printf("   Hostname   : %s%s%s\n", Cyan, getHostname(), Reset)
	fmt.Printf("   OS         : %s%s %s%s\n", Green, info.Platform, info.PlatformVersion, Reset)
	fmt.Printf("   Kernel     : %s\n", info.KernelVersion)
	fmt.Printf("   Arch       : %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("   Uptime     : %s%s%s (%.1f days)\n", Yellow, uptime.Round(time.Minute), Reset, uptime.Hours()/24)
	fmt.Printf("   Boot Time  : %s\n", bootTime.Format("2006-01-02 15:04"))
	fmt.Printf("   Users      : %d logged in\n\n", len(users))
}

func printCPUInfo() {
	count, _ := cpu.Counts(true)
	percent, _ := cpu.Percent(time.Second, false)

	fmt.Printf("%sCPU Usage%s\n", Bold, Reset)
	fmt.Printf("   Cores      : %d\n", count)
	fmt.Printf("   Usage      : %s%.1f%%%s\n\n", getColor(percent[0]), percent[0], Reset)
}

func printMemoryInfo() {
	v, _ := mem.VirtualMemory()
	s, _ := mem.SwapMemory()

	fmt.Printf("%sMemory & Swap%s\n", Bold, Reset)
	fmt.Printf("   RAM Total  : %s%s%s\n", Green, formatBytes(v.Total), Reset)
	fmt.Printf("   RAM Used   : %s%s%s (%.1f%%)\n", getColor(v.UsedPercent), formatBytes(v.Used), Reset, v.UsedPercent)
	fmt.Printf("   RAM Free   : %s%s%s\n", Cyan, formatBytes(v.Free), Reset)
	if s.Total > 0 {
		fmt.Printf("   Swap Used  : %s%s%s / %s (%.1f%%)\n",
			getColor(s.UsedPercent), formatBytes(s.Used), Reset, formatBytes(s.Total), s.UsedPercent)
	}
	fmt.Println()
}

func printDiskInfo() {
	parts, _ := disk.Partitions(false)
	fmt.Printf("%sDisk Usage%s\n", Bold, Reset)
	fmt.Printf("   %-12s %8s %8s %8s %6s\n", "Mount", "Total", "Used", "Free", "Use%")
	fmt.Printf("   %s\n", strings.Repeat("─", 56))
	for _, p := range parts {
		if contains([]string{"tmpfs", "devtmpfs", "overlay", "squashfs", "/snap"}, p.Fstype) || strings.HasPrefix(p.Mountpoint, "/snap") {
			continue
		}
		u, err := disk.Usage(p.Mountpoint)
		if err != nil || u.Total == 0 {
			continue
		}
		color := getColor(u.UsedPercent)
		fmt.Printf("   %-12s %8s %8s %8s %s%5.1f%%%s\n",
			p.Mountpoint, formatBytes(u.Total), formatBytes(u.Used), formatBytes(u.Free),
			color, u.UsedPercent, Reset)
	}
	fmt.Println()
}

func printNetworkInfo() {
	ip := getLocalIP()
	publicIP := getPublicIP()

	fmt.Printf("%sNetwork%s\n", Bold, Reset)
	fmt.Printf("   Local IP   : %s%s%s\n", Cyan, ip, Reset)
	if publicIP != "N/A" && publicIP != "" {
		fmt.Printf("   Public IP  : %s%s%s\n", Magenta, publicIP, Reset)
	}
	fmt.Println()
}

func printTopProcesses() {
	ps, _ := process.Processes()
	type proc struct {
		Pid  int32
		Name string
		CPU  float64
		Mem  float32
	}
	var list []proc
	for _, p := range ps {
		name, _ := p.Name()
		cpuP, _ := p.CPUPercent()
		memP, _ := p.MemoryPercent()
		if cpuP > 0.1 {
			list = append(list, proc{Pid: p.Pid, Name: name, CPU: cpuP, Mem: memP})
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].CPU > list[j].CPU })

	fmt.Printf("%sTop 5 Processes by CPU%s\n", Bold, Reset)
	fmt.Printf("   %-6s %-8s %-6s %s\n", "PID", "%CPU", "%MEM", "Command")
	fmt.Printf("   %s\n", strings.Repeat("─", 50))
	for i := 0; i < 5 && i < len(list); i++ {
		p := list[i]
		fmt.Printf("   %-6d %s%6.1f%s %5.1f  %s\n", p.Pid, getColor(p.CPU), p.CPU, Reset, p.Mem, truncate(p.Name, 30))
	}
	fmt.Println()
}

func printLoadAverage() {
	l, _ := load.Avg()
	fmt.Printf("%sLoad Average (1/5/15 min)%s\n", Bold, Reset)
	fmt.Printf("   %s%.2f  %.2f  %.2f%s\n\n", getLoadColor(l.Load1), l.Load1, l.Load5, l.Load15, Reset)
}

// Helpers
func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "KMGT"[exp])
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if strings.Contains(val, item) {
			return true
		}
	}
	return false
}

func getHostname() string {
	h, _ := os.Hostname()
	return h
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "N/A"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

func getPublicIP() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "N/A"
	}
	defer resp.Body.Close()
	ip, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(ip))
}

func getColor(percent float64) string {
	if percent < 30 {
		return Green
	} else if percent < 70 {
		return Yellow
	}
	return Red
}

func getLoadColor(load float64) string {
	cpus, _ := cpu.Counts(true)
	if load < float64(cpus)*0.7 {
		return Green
	} else if load < float64(cpus)*1.5 {
		return Yellow
	}
	return Red
}

func printJSON() {
	data := map[string]interface{}{
		"timestamp":   time.Now().Unix(),
		"hostname":    getHostname(),
		"cpu_percent": getCPUPercent(),
		"memory":      getMemoryInfo(),
		"disks":       getDiskUsage(),
		"load_avg":    getLoadAvg(),
	}
	json.NewEncoder(os.Stdout).Encode(data)
}

func getCPUPercent() float64 {
	p, _ := cpu.Percent(0, false)
	if len(p) > 0 {
		return p[0]
	}
	return 0
}

func getMemoryInfo() map[string]interface{} {
	v, _ := mem.VirtualMemory()
	return map[string]interface{}{
		"total":       v.Total,
		"used":        v.Used,
		"free":         v.Free,
		"used_percent": v.UsedPercent,
	}
}

func getLoadAvg() map[string]float64 {
	l, _ := load.Avg()
	return map[string]float64{
		"load1":  l.Load1,
		"load5":  l.Load5,
		"load15": l.Load15,
	}
}

func getDiskUsage() []map[string]interface{} {
	var result []map[string]interface{}
	parts, _ := disk.Partitions(false)
	for _, p := range parts {
		if contains([]string{"tmpfs", "devtmpfs", "overlay"}, p.Fstype) {
			continue
		}
		u, _ := disk.Usage(p.Mountpoint)
		if u != nil && u.Total > 0 {
			result = append(result, map[string]interface{}{
				"mountpoint": p.Mountpoint,
				"total":      u.Total,
				"used":       u.Used,
				"free":       u.Free,
				"percent":    u.UsedPercent,
			})
		}
	}
	return result
}