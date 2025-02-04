package linux

import (
	"defetch/helper"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
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

	// Battery Information
	// batteryInfo := getBatteryInfo()

	// Peripherals Information
	peripheralsInfo := getPeripheralsInfo()

	// Software Information
	softwareInfo := getSoftwareInfo()

	// Performance Information
	performanceInfo := getPerformanceInfo()

	// Package Management Information
	packageManagementInfo := getPackageManagementInfo()

	// Other Information
	otherInfo := getOtherInfo()

	return helper.SysInfo{
		Hostname:          hostname,
		CurrentUser:       currentUser.Username,
		OSName:            osName,
		OSVersion:         osVersion,
		OSCodename:        osCodename,
		KernelVersion:     kernelVersion,
		Shell:             shellBinary,
		ShellVersion:      shellVersion,
		Architecture:      architecture,
		Uptime:            uptimeStr,
		CPU:               cpuInfo,
		GPU:               gpuInfo,
		Motherboard:       motherboardInfo,
		Memory:            memoryInfo,
		Storage:           storageInfo,
		Network:           networkInfo,
		Peripherals:       peripheralsInfo,
		Software:          softwareInfo,
		Performance:       performanceInfo,
		PackageManagement: packageManagementInfo,
		OtherInfo:         otherInfo,
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

// NOTE FOR ME: /sys/class/power_supply/BAT0/ NOT WORKING, TODO LATER
// Helper function to get Battery information
// func getBatteryInfo() helper.BatteryInfo {
// 	// Placeholder values, actual implementation will depend on available battery data source
// 	status := "Unknown"
// 	capacity := "Unknown"
// 	percentage := "Unknown"
// 	manufacturer := "Unknown"
// 	model := "Unknown"

// 	// Retrieve battery information from system files (example: /sys/class/power_supply/BAT0/)
// 	batteryPath := "/sys/class/power_supply/BAT0/"
// 	statusBytes, err := os.ReadFile(batteryPath + "status")
// 	if err == nil {
// 		status = strings.TrimSpace(string(statusBytes))
// 	}
// 	capacityBytes, err := os.ReadFile(batteryPath + "energy_full")
// 	if err == nil {
// 		capacity = strings.TrimSpace(string(capacityBytes)) + " µWh"
// 	}
// 	percentageBytes, err := os.ReadFile(batteryPath + "capacity")
// 	if err == nil {
// 		percentage = strings.TrimSpace(string(percentageBytes)) + "%"
// 	}
// 	manufacturerBytes, err := os.ReadFile(batteryPath + "manufacturer")
// 	if err == nil {
// 		manufacturer = strings.TrimSpace(string(manufacturerBytes))
// 	}
// 	modelBytes, err := os.ReadFile(batteryPath + "model_name")
// 	if err == nil {
// 		model = strings.TrimSpace(string(modelBytes))
// 	}

// 	return helper.BatteryInfo{
// 		Status:       status,
// 		Capacity:     capacity,
// 		Percentage:   percentage,
// 		Manufacturer: manufacturer,
// 		Model:        model,
// 	}
// }

// Helper function to get Peripherals information
func getPeripheralsInfo() helper.PeripheralInfo {
	var peripherals helper.PeripheralInfo

	// Connected devices (using xinput for example)
	connectedDevicesOutput, err := exec.Command("xinput", "--list", "--name-only").Output()
	if err == nil {
		peripherals.ConnectedDevices = strings.Split(strings.TrimSpace(string(connectedDevicesOutput)), "\n")
	}

	// USB devices
	usbDevicesOutput, err := exec.Command("lsusb").Output()
	if err == nil {
		usbLines := strings.Split(strings.TrimSpace(string(usbDevicesOutput)), "\n")
		for _, line := range usbLines {
			parts := strings.Fields(line)
			if len(parts) > 5 {
				vendorID := parts[5][:4]
				productID := parts[5][5:]
				vendor := parts[6]
				name := strings.Join(parts[7:], " ")
				peripherals.USBDevices = append(peripherals.USBDevices, helper.USBDeviceInfo{
					Name:      name,
					Vendor:    vendor,
					ProductID: productID,
					VendorID:  vendorID,
				})
			}
		}
	}

	// Audio devices (using aplay -l for example)
	audioDevicesOutput, err := exec.Command("aplay", "-l").Output()
	if err == nil {
		audioLines := strings.Split(strings.TrimSpace(string(audioDevicesOutput)), "\n")
		for _, line := range audioLines {
			if strings.Contains(line, "card") {
				peripherals.AudioDevices = append(peripherals.AudioDevices, line)
			}
		}
	}

	// Printer details (using lpstat -p for example)
	printerOutput, err := exec.Command("lpstat", "-p").Output()
	if err == nil {
		printerLines := strings.Split(strings.TrimSpace(string(printerOutput)), "\n")
		peripherals.PrinterDetails = printerLines
	}

	return peripherals
}

// Helper function to get Software information
func getSoftwareInfo() helper.SoftwareInfo {
	osDetails := getOSDetails()
	desktopEnvironment := getDesktopEnvironment()
	windowManager := getWindowManager()
	wmTheme := getWMTheme()
	gtkTheme := getGTKTheme()
	iconsTheme := getIconsTheme()
	font := getFont()
	browsers := getInstalledBrowsers()
	processes := getRunningProcesses()
	startupPrograms := getStartupPrograms()

	return helper.SoftwareInfo{
		OSDetails:          osDetails,
		DesktopEnvironment: desktopEnvironment,
		WindowManager:      windowManager,
		WMTheme:            wmTheme,
		GTKTheme:           gtkTheme,
		IconsTheme:         iconsTheme,
		Font:               font,
		Browser:            browsers,
		RunningProcesses:   processes,
		StartupPrograms:    startupPrograms,
	}
}

// Helper function to get installed browsers and their versions
func getInstalledBrowsers() []helper.BrowserInfo {
	// List of known browser executables (add more as needed)
	browsers := []string{"firefox", "google-chrome", "chromium", "brave", "opera"}

	var browserInfos []helper.BrowserInfo
	for _, browser := range browsers {
		output, err := exec.Command(browser, "--version").Output()
		if err == nil {
			versionInfo := strings.TrimSpace(string(output))
			parts := strings.Fields(versionInfo)
			if len(parts) > 1 {
				browserInfos = append(browserInfos, helper.BrowserInfo{Name: parts[0], Version: parts[1]})
			}
		}
	}

	return browserInfos
}

// Helper function to get running processes information
func getRunningProcesses() []helper.ProcessInfo {
	// Use ps command to get processes info
	output, err := exec.Command("ps", "axo", "pid,comm,pcpu,pmem", "--sort=-pcpu").Output()
	if err != nil {
		panic("Cannot execute ps command")
	}

	lines := strings.Split(string(output), "\n")
	var processes []helper.ProcessInfo
	for _, line := range lines[1:] { // Skip the header
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		pid, _ := strconv.Atoi(fields[0])
		name := fields[1]
		cpuUsage, _ := strconv.ParseFloat(fields[2], 64)
		memUsage, _ := strconv.ParseFloat(fields[3], 64)

		processes = append(processes, helper.ProcessInfo{
			PID:         pid,
			Name:        name,
			CPUUsage:    cpuUsage,
			MemoryUsage: memUsage,
		})
	}

	return processes
}

// Helper function to get startup programs
func getStartupPrograms() []helper.StartupProgram {
	// This implementation will vary based on the desktop environment and OS

	// Example for .desktop files in autostart directory
	autostartDir := filepath.Join(os.Getenv("HOME"), ".config", "autostart")
	files, err := os.ReadDir(autostartDir)
	if err != nil {
		return nil
	}

	var startupPrograms []helper.StartupProgram
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".desktop" {
			// Reading .desktop file content
			content, err := os.ReadFile(filepath.Join(autostartDir, file.Name()))
			if err == nil {
				var name, command string
				for _, line := range strings.Split(string(content), "\n") {
					if strings.HasPrefix(line, "Name=") {
						name = strings.TrimSpace(strings.TrimPrefix(line, "Name="))
					}
					if strings.HasPrefix(line, "Exec=") {
						command = strings.TrimSpace(strings.TrimPrefix(line, "Exec="))
					}
				}
				startupPrograms = append(startupPrograms, helper.StartupProgram{Name: name, Command: command})
			}
		}
	}

	return startupPrograms
}

// Helper function to get Window Manager theme
func getWMTheme() string {
	// Assuming `gsettings` or `xfconf-query` can be used, this is dependent on the environment
	output, err := exec.Command("gsettings", "get", "org.gnome.desktop.wm.preferences", "theme").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(output))
}

// Helper function to get GTK theme
func getGTKTheme() string {
	output, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(output))
}

// Helper function to get Icons theme
func getIconsTheme() string {
	output, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "icon-theme").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(output))
}

// Helper function to get the font used in the system
func getFont() string {
	output, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "font-name").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(output))
}

// Helper function to get OS details
func getOSDetails() string {
	platform, _, version, err := host.PlatformInformation()
	if err != nil {
		return "Unknown"
	}
	return fmt.Sprintf("%s %s", platform, version)
}

// Helper function to get Desktop Environment information
func getDesktopEnvironment() string {
	de := os.Getenv("XDG_CURRENT_DESKTOP")
	version, err := exec.Command("sh", "-c", "echo $XDG_SESSION_DESKTOP").Output()
	if err != nil {
		return de
	}
	return fmt.Sprintf("%s %s", de, strings.TrimSpace(string(version)))
}

// Helper function to get Window Manager information
func getWindowManager() string {
	wm := os.Getenv("XDG_SESSION_DESKTOP")
	version, err := exec.Command("sh", "-c", "wmctrl -m | grep 'Name|Version'").Output()
	if err != nil {
		return wm
	}
	return fmt.Sprintf("%s %s", wm, strings.TrimSpace(string(version)))
}

// Helper function to get system performance information
func getPerformanceInfo() helper.PerformanceInfo {
	cpuUsage, perCoreUsage := getCPUUsage()
	memoryUsage, perAppMemoryUsage := getMemoryUsage()

	return helper.PerformanceInfo{
		CPUUsage:          cpuUsage,
		PerCoreUsage:      perCoreUsage,
		MemoryUsage:       memoryUsage,
		PerAppMemoryUsage: perAppMemoryUsage,
	}
}

// Helper function to get CPU usage information
func getCPUUsage() (float64, []float64) {
	// Use top command to get overall CPU usage and per-core usage
	output, err := exec.Command("top", "-bn1").Output()
	if err != nil {
		panic("Cannot execute top command")
	}

	lines := strings.Split(string(output), "\n")
	var cpuUsage float64
	var perCoreUsage []float64

	for _, line := range lines {
		if strings.HasPrefix(line, "%Cpu(s):") {
			// Overall CPU usage
			parts := strings.Fields(line)
			usage, _ := strconv.ParseFloat(parts[1], 64)
			cpuUsage = usage
		} else if strings.HasPrefix(line, "%Cpu") {
			// Per-core usage
			parts := strings.Fields(line)
			usage, _ := strconv.ParseFloat(parts[1], 64)
			perCoreUsage = append(perCoreUsage, usage)
		}
	}

	return cpuUsage, perCoreUsage
}

// Helper function to get memory usage information
func getMemoryUsage() (helper.MemoryUsageInfo, []helper.AppMemoryUsage) {
	// Use free command to get overall memory usage
	freeOutput, err := exec.Command("free", "-h").Output()
	if err != nil {
		panic("Cannot execute free command")
	}

	freeLines := strings.Split(string(freeOutput), "\n")
	var memoryUsageInfo helper.MemoryUsageInfo
	for _, line := range freeLines {
		if strings.HasPrefix(line, "Mem:") {
			parts := strings.Fields(line)
			memoryUsageInfo.Total = parts[1]
			memoryUsageInfo.TotalUsed = parts[2]
			memoryUsageInfo.Free = parts[3]
		}
	}

	// Use ps command to get per-application memory usage
	psOutput, err := exec.Command("ps", "axo", "pid,comm,%mem", "--sort=-%mem").Output()
	if err != nil {
		panic("Cannot execute ps command")
	}

	psLines := strings.Split(string(psOutput), "\n")
	var perAppMemoryUsage []helper.AppMemoryUsage
	for _, line := range psLines[1:] { // Skip the header
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		pid, _ := strconv.Atoi(fields[0])
		name := fields[1]
		memUsage, _ := strconv.ParseFloat(fields[2], 64)

		perAppMemoryUsage = append(perAppMemoryUsage, helper.AppMemoryUsage{
			PID:         pid,
			Name:        name,
			MemoryUsage: memUsage,
		})
	}

	return memoryUsageInfo, perAppMemoryUsage
}

// Helper function to get package management information
func getPackageManagementInfo() helper.PackageManagementInfo {
	packageCount := getPackageCount()
	availableUpdates := getAvailableUpdates()
	packageManagers := getPackageManagers()
	recentlyInstalledPackages := getRecentlyInstalledPackages()

	return helper.PackageManagementInfo{
		PackageCount:              packageCount,
		AvailableUpdates:          availableUpdates,
		PackageManagers:           packageManagers,
		RecentlyInstalledPackages: recentlyInstalledPackages,
	}
}

// Function to get the number of installed packages
func getPackageCount() int {
	output, err := exec.Command("dpkg-query", "-f", ".", "-W").Output()
	if err != nil {
		panic("Cannot execute dpkg-query command")
	}
	return len(output)
}

// Function to get the number of available updates
func getAvailableUpdates() int {
	output, err := exec.Command("apt", "list", "--upgradable").Output()
	if err != nil {
		panic("Cannot execute apt command to check for updates")
	}
	lines := strings.Split(string(output), "\n")
	return len(lines) - 1 // Exclude header line
}

// Function to list the used package managers
func getPackageManagers() []string {
	// List common package managers on Linux systems
	packageManagers := []string{}
	if _, err := exec.LookPath("apt"); err == nil {
		packageManagers = append(packageManagers, "apt")
	}
	if _, err := exec.LookPath("dnf"); err == nil {
		packageManagers = append(packageManagers, "dnf")
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		packageManagers = append(packageManagers, "pacman")
	}
	if _, err := exec.LookPath("yum"); err == nil {
		packageManagers = append(packageManagers, "yum")
	}
	if _, err := exec.LookPath("zypper"); err == nil {
		packageManagers = append(packageManagers, "zypper")
	}
	return packageManagers
}

// Helper function to get recently installed packages
func getRecentlyInstalledPackages() []helper.PackageInfo {
	logFiles := []string{"/var/log/dpkg.log", "/var/log/dpkg.log.1"}
	var packages []helper.PackageInfo

	for _, logFile := range logFiles {
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			continue // Skip if the log file does not exist
		}

		output, err := exec.Command("zgrep", "install ", logFile).Output()
		if err != nil {
			fmt.Printf("Warning: Cannot read dpkg logs from %s: %v\n", logFile, err)
			continue
		}

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				packages = append(packages, helper.PackageInfo{
					Name:          parts[3],
					Version:       parts[4],
					InstalledDate: fmt.Sprintf("%s %s", parts[0], parts[1]),
				})
			}
		}
	}

	return packages
}

// Helper function to get other information
func getOtherInfo() helper.OtherInfo {
	publicIP := getPublicIP()
	timezone := getTimezone()
	locale := getLocale()
	temperature := getTemperature()
	systemLanguage := getSystemLanguage()
	screenResolution := getScreenResolution()
	diskPartitions := getDiskPartitions()

	return helper.OtherInfo{
		PublicIP:         publicIP,
		Timezone:         timezone,
		Locale:           locale,
		Temperature:      temperature,
		SystemLanguage:   systemLanguage,
		ScreenResolution: screenResolution,
		DiskPartitions:   diskPartitions,
	}
}

// Function to get public IP address
func getPublicIP() string {
	output, err := exec.Command("curl", "-s", "ifconfig.me").Output()
	if err != nil {
		return "Unavailable"
	}
	return strings.TrimSpace(string(output))
}

// Function to get timezone
func getTimezone() string {
	output, err := exec.Command("timedatectl", "show", "-p", "Timezone").Output()
	if err != nil {
		return "Unavailable"
	}
	return strings.TrimSpace(strings.Split(string(output), "=")[1])
}

// Function to get locale
func getLocale() string {
	output, err := exec.Command("locale", "|", "grep", "LANG=").Output()
	if err != nil {
		return "Unavailable"
	}
	return strings.TrimSpace(strings.Split(string(output), "=")[1])
}

// Function to get temperature readings
func getTemperature() helper.TemperatureInfo {
	// Example implementation using lm-sensors
	output, err := exec.Command("sensors").Output()
	if err != nil {
		return helper.TemperatureInfo{}
	}

	var cpuTemp, gpuTemp, moboTemp float64
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Core 0:") {
			fmt.Sscanf(line, "Core 0: +%f°C", &cpuTemp)
		}
		if strings.Contains(line, "temp1:") && strings.Contains(line, "GPU") {
			fmt.Sscanf(line, "temp1: +%f°C", &gpuTemp)
		}
		if strings.Contains(line, "temp1:") && strings.Contains(line, "MB") {
			fmt.Sscanf(line, "temp1: +%f°C", &moboTemp)
		}
	}

	return helper.TemperatureInfo{
		CPU:         cpuTemp,
		GPU:         gpuTemp,
		Motherboard: moboTemp,
	}
}

// Function to get system language
func getSystemLanguage() string {
	output, err := exec.Command("locale", "|", "grep", "LANG=").Output()
	if err != nil {
		return "Unavailable"
	}
	return strings.TrimSpace(strings.Split(string(output), "=")[1])
}

// Function to get screen resolution and monitor details
func getScreenResolution() []helper.ScreenInfo {
	var screens []helper.ScreenInfo
	output, err := exec.Command("xrandr", "--query").Output()
	if err != nil {
		return screens
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, " connected") {
			parts := strings.Fields(line)
			model := parts[0]
			resolution := parts[2]
			var refreshRate int
			fmt.Sscanf(parts[4], "%d", &refreshRate)
			screens = append(screens, helper.ScreenInfo{
				Model:       model,
				Resolution:  resolution,
				RefreshRate: refreshRate,
			})
		}
	}
	return screens
}

// Function to get disk partitions information
func getDiskPartitions() []helper.PartitionInfo {
	var partitions []helper.PartitionInfo
	output, err := exec.Command("lsblk", "-o", "NAME,FSTYPE,MOUNTPOINT,SIZE,USED,AVAIL").Output()
	if err != nil {
		return partitions
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		parts := strings.Fields(line)
		if len(parts) >= 6 {
			partitions = append(partitions, helper.PartitionInfo{
				Device:     parts[0],
				Filesystem: parts[1],
				MountPoint: parts[2],
				Size:       parts[3],
				Used:       parts[4],
				Available:  parts[5],
			})
		}
	}
	return partitions
}
