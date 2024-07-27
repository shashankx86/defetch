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

	// Storage Information
	storageInfo := getStorageInfo()

	// Network Information
	networkInfo := getNetworkInfo()

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
		Storage:       storageInfo,
		Network:       networkInfo,
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
	manufacturer := readSysFileOrFallback("/sys/class/dmi/id/board_vendor", "dmidecode -s baseboard-manufacturer")
	model := readSysFileOrFallback("/sys/class/dmi/id/board_name", "dmidecode -s baseboard-product-name")
	biosVersion := readSysFileOrFallback("/sys/class/dmi/id/bios_version", "dmidecode -s bios-version")
	serialNumber := readSysFileOrFallback("/sys/class/dmi/id/board_serial", "dmidecode -s baseboard-serial-number")

	return helper.MotherboardInfo{
		Manufacturer: manufacturer,
		Model:        model,
		BIOSVersion:  biosVersion,
		SerialNumber: serialNumber,
	}
}

// parseSize converts a memory size string (e.g., "4096 kB") into an integer value in kilobytes.
func parseSize(sizeStr string) int64 {
	sizeParts := strings.Fields(sizeStr)
	size, err := strconv.ParseInt(sizeParts[0], 10, 64)
	if err != nil {
		return 0
	}
	return size
}

// Helper function to get Memory information
func getMemoryInfo() helper.MemoryInfo {
	var totalSize, usedSize, freeSize string
	var slots []helper.MemorySlotInfo

	// Use /proc/meminfo to get memory usage
	meminfoOutput, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		panic("Cannot read /proc/meminfo")
	}

	lines := strings.Split(string(meminfoOutput), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			totalSize = strings.TrimSpace(strings.Fields(line)[1]) + " kB"
		}
		if strings.HasPrefix(line, "MemFree:") {
			freeSize = strings.TrimSpace(strings.Fields(line)[1]) + " kB"
		}
	}

	// Calculate used size
	usedSize = fmt.Sprintf("%d kB", parseSize(totalSize)-parseSize(freeSize))

	// Use dmidecode as a fallback for memory slot information if sysfs is unavailable
	dmidecodeOutput, err := exec.Command("dmidecode", "-t", "memory").Output()
	if err == nil {
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
	}

	return helper.MemoryInfo{
		TotalSize: totalSize,
		UsedSize:  usedSize,
		FreeSize:  freeSize,
		Slots:     slots,
	}
}

// Helper function to read from sysfs or fallback to dmidecode command
func readSysFileOrFallback(sysFilePath string, dmidecodeCmd string) string {
	content, err := os.ReadFile(sysFilePath)
	if err == nil {
		return strings.TrimSpace(string(content))
	}

	// Fallback to using dmidecode
	output, err := exec.Command("sh", "-c", dmidecodeCmd).Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(output))
}

// Helper function to get Storage information
func getStorageInfo() []helper.StorageInfo {
	var storages []helper.StorageInfo

	// Use lsblk command to get block device information
	output, err := exec.Command("lsblk", "-d", "-o", "NAME,MODEL,SIZE,ROTA,RM").Output()
	if err != nil {
		panic("Cannot execute lsblk command")
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		device := parts[0]
		model := parts[1]
		capacity := parts[2]

		// Use df command to get used and available space, and file system type
		dfOutput, err := exec.Command("df", "-hT", "/dev/"+device).Output()
		if err != nil {
			continue
		}

		dfLines := strings.Split(string(dfOutput), "\n")
		if len(dfLines) < 2 {
			continue
		}

		dfParts := strings.Fields(dfLines[1])
		if len(dfParts) < 7 {
			continue
		}

		filesystem := dfParts[1]
		used := dfParts[2]
		available := dfParts[3]
		mountPoint := dfParts[6]

		storages = append(storages, helper.StorageInfo{
			Device:     device,
			Model:      model,
			Capacity:   capacity,
			Used:       used,
			Available:  available,
			FileSystem: filesystem,
			MountPoint: mountPoint,
			ReadSpeed:  "Unknown", // Placeholder, as getting read speed dynamically is complex
			WriteSpeed: "Unknown", // Placeholder, as getting write speed dynamically is complex
		})
	}

	return storages
}

// Helper function to get Network information
func getNetworkInfo() []helper.NetworkInfo {
	var networks []helper.NetworkInfo

	// Use ip command to get network interfaces
	output, err := exec.Command("ip", "-o", "addr", "show").Output()
	if err != nil {
		panic("Cannot execute ip command")
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		interfaceName := fields[1]
		ipAddress := fields[3]

		// Get MAC address
		macOutput, err := exec.Command("cat", fmt.Sprintf("/sys/class/net/%s/address", interfaceName)).Output()
		macAddress := "Unknown"
		if err == nil {
			macAddress = strings.TrimSpace(string(macOutput))
		}

		// Check if interface is active
		active := strings.Contains(fields[len(fields)-1], "state UP")

		// Get network speed (placeholder, as getting actual speed is complex)
		speed := "Unknown"

		// Get default gateway
		gatewayOutput, err := exec.Command("ip", "route", "show", "default").Output()
		defaultGateway := "Unknown"
		if err == nil {
			gatewayLines := strings.Split(string(gatewayOutput), "\n")
			if len(gatewayLines) > 0 {
				gatewayFields := strings.Fields(gatewayLines[0])
				if len(gatewayFields) > 2 {
					defaultGateway = gatewayFields[2]
				}
			}
		}

		networks = append(networks, helper.NetworkInfo{
			InterfaceName:  interfaceName,
			IPAddress:      ipAddress,
			MACAddress:     macAddress,
			Speed:          speed,
			Active:         active,
			DefaultGateway: defaultGateway,
		})
	}

	return networks
}
