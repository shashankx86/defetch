package linux

import (
	"defetch/helper"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/host"
	"golang.org/x/sys/unix"
)

func GetLinuxInfo() helper.SysInfo {
	// Hostname
	hostname, _ := os.Hostname()

	// Current User
	currentUser, _ := user.Current()

	// Operating System
	platform, family, version, _ := host.PlatformInformation()
	osName := platform
	osVersion := version
	osCodename := family

	// Kernel Version
	var uname unix.Utsname
	_ = unix.Uname(&uname)
	kernelVersion := charsToString(uname.Release[:])

	// Shell and Version
	shell := os.Getenv("SHELL")
	shellBinary := getShellBinary(shell)
	shellVersion := getShellVersion(shellBinary)

	// Architecture
	architecture := runtime.GOARCH

	// Uptime
	uptime, _ := host.Uptime()
	uptimeStr := formatUptime(uptime)

	// CPU Information
	cpuInfo := getCPUInfo()

	// GPU Information
	gpuInfo := getGPUInfo()

	// Motherboard Information
	motherboardInfo := getMotherboardInfo()

	// Memory Information
	memoryInfo := getMemoryInfo()

	return helper.SysInfo{
		Hostname:      hostname,
		CurrentUser:   currentUser.Username,
		OSName:        osName,
		OSVersion:     osVersion,
		OSCodename:    osCodename,
		KernelVersion: kernelVersion,
		Shell:         shellBinary,
		ShellVersion:  shellVersion,
		Architecture:  architecture,
		Uptime:        uptimeStr,
		CPU:           cpuInfo,
		GPU:           gpuInfo,
		Motherboard:   motherboardInfo,
		Memory:        memoryInfo,
	}
}

// Helper function to convert Unix Utsname to a string
func charsToString(ca []byte) string {
	s := make([]byte, len(ca))
	var lens int
	for lens = 0; lens < len(ca); lens++ {
		if ca[lens] == 0 {
			break
		}
		s[lens] = ca[lens]
	}
	return string(s[:lens])
}

// Helper function to format uptime from seconds
func formatUptime(uptime uint64) string {
	days := uptime / (24 * 3600)
	hours := (uptime % (24 * 3600)) / 3600
	minutes := (uptime % 3600) / 60
	return strings.TrimSpace(strings.Join([]string{
		formatPlural(days, "day"),
		formatPlural(hours, "hour"),
		formatPlural(minutes, "minute"),
	}, " "))
}

func formatPlural(value uint64, unit string) string {
	if value == 0 {
		return ""
	}
	if value == 1 {
		return fmt.Sprintf("%d %s", value, unit)
	}
	return fmt.Sprintf("%d %ss", value, unit)
}

// Helper function to get the shell binary name
func getShellBinary(shellPath string) string {
	parts := strings.Split(shellPath, "/")
	return parts[len(parts)-1]
}

// Helper function to get the shell version
func getShellVersion(shellBinary string) string {
	var output []byte
	var err error

	switch shellBinary {
	case "bash":
		output, err = exec.Command(shellBinary, "--version").Output()
	case "zsh":
		output, err = exec.Command(shellBinary, "--version").Output()
	case "fish":
		output, err = exec.Command(shellBinary, "--version").Output()
	case "sh":
		output, err = exec.Command(shellBinary, "--version").Output()
	default:
		return "Unknown shell or version not available"
	}

	if err != nil {
		return "Version not available"
	}

	firstLine := strings.SplitN(strings.TrimSpace(string(output)), "\n", 2)[0]
	return strings.TrimSpace(firstLine)
}

// Helper function to get CPU information
func getCPUInfo() helper.CPUInfo {
	// Read from /proc/cpuinfo
	output, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		panic("Cannot read /proc/cpuinfo")
	}

	lines := strings.Split(string(output), "\n")
	var modelName, flags string
	var cores, cacheSize int
	var frequency float64

	for _, line := range lines {
		if strings.Contains(line, "model name") && modelName == "" {
			modelName = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "cpu cores") {
			cores, _ = strconv.Atoi(strings.TrimSpace(strings.Split(line, ":")[1]))
		}
		if strings.Contains(line, "cache size") {
			cacheSize, _ = strconv.Atoi(strings.TrimSpace(strings.Split(line, ":")[1]))
		}
		if strings.Contains(line, "cpu MHz") {
			frequency, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64)
		}
		if strings.Contains(line, "flags") && flags == "" {
			flags = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	// Check if model name is empty, fallback to /proc/device-tree/model for ARM systems
	if modelName == "" {
		deviceModel, err := os.ReadFile("/proc/device-tree/model")
		if err == nil {
			modelName = strings.TrimSpace(string(deviceModel))
		}
	}

	// Get number of logical CPUs
	threads := runtime.NumCPU()

	return helper.CPUInfo{
		ModelName:    modelName,
		Cores:        cores,
		Threads:      threads,
		Architecture: runtime.GOARCH,
		Frequency:    frequency,
		CacheSize:    int32(cacheSize),
		Flags:        flags,
	}
}

// Helper function to get GPU information
func getGPUInfo() helper.GPUInfo {
	var modelName, driverVersion, memorySize string

	// Use lshw command to get GPU info
	output, err := exec.Command("lshw", "-C", "display").Output()
	if err != nil {
		panic("Cannot execute lshw command")
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "product:") {
			modelName = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "version:") {
			driverVersion = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "size:") && strings.Contains(line, "memory") {
			memorySize = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	return helper.GPUInfo{
		ModelName:     modelName,
		DriverVersion: driverVersion,
		MemorySize:    memorySize,
	}
}

// Helper function to get Motherboard information
func getMotherboardInfo() helper.MotherboardInfo {
	var manufacturer, model, biosVersion, serialNumber string

	// Use dmidecode command to get motherboard info
	output, err := exec.Command("dmidecode", "-t", "baseboard").Output()
	if err != nil {
		panic("Cannot execute dmidecode command")
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Manufacturer:") {
			manufacturer = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Product Name:") {
			model = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Version:") {
			biosVersion = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Serial Number:") {
			serialNumber = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	return helper.MotherboardInfo{
		Manufacturer: manufacturer,
		Model:        model,
		BIOSVersion:  biosVersion,
		SerialNumber: serialNumber,
	}
}

// Helper function to get Memory information
func getMemoryInfo() helper.MemoryInfo {
	var totalSize, usedSize, freeSize string
	var slots []helper.MemorySlotInfo

	// Use free command to get memory usage
	freeOutput, err := exec.Command("free", "-h").Output()
	if err != nil {
		panic("Cannot execute free command")
	}

	freeLines := strings.Split(string(freeOutput), "\n")
	for _, line := range freeLines {
		if strings.HasPrefix(line, "Mem:") {
			parts := strings.Fields(line)
			totalSize = parts[1]
			usedSize = parts[2]
			freeSize = parts[3]
		}
	}

	// Use dmidecode command to get memory slot information
	dmidecodeOutput, err := exec.Command("dmidecode", "-t", "memory").Output()
	if err != nil {
		panic("Cannot execute dmidecode command")
	}

	dmidecodeLines := strings.Split(string(dmidecodeOutput), "\n")
	var currentSlot helper.MemorySlotInfo
	for _, line := range dmidecodeLines {
		if strings.HasPrefix(line, "Size:") {
			if currentSlot.Size != "" {
				slots = append(slots, currentSlot)
			}
			currentSlot = helper.MemorySlotInfo{Size: strings.TrimSpace(strings.Split(line, ":")[1])}
		}
		if strings.HasPrefix(line, "Form Factor:") {
			currentSlot.FormFactor = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.HasPrefix(line, "Type:") {
			currentSlot.Type = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.HasPrefix(line, "Speed:") {
			currentSlot.Speed = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	if currentSlot.Size != "" {
		slots = append(slots, currentSlot)
	}

	return helper.MemoryInfo{
		TotalSize: totalSize,
		UsedSize:  usedSize,
		FreeSize:  freeSize,
		Slots:     slots,
	}
}
