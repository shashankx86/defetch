//go:build windows
// +build windows

package windows

import (
	"defetch/helper"
	"fmt"
	"os"
	"os/user"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows"
)

// GetWindowsInfo retrieves all system information on a Windows machine
func GetWindowsInfo() helper.SysInfo {
	// Hostname
	hostname, _ := os.Hostname()

	// Current User
	currentUser, _ := user.Current()

	// OS Name, Version, Build Number
	osName, osVersion, osBuildNumber := getOSInfo()

	// Manufacturer and Model
	manufacturer, model := getSystemManufacturerAndModel()

	// Kernel Version
	kernelVersion := getKernelVersion()

	// System Uptime
	systemUptime := getSystemUptime()

	// Installed Packages
	installedPackages := getInstalledPackagesCount()

	// Shell Info
	shellName := getShellName()
	shellVersion := getShellVersion()

	// Display Info
	primaryDisplayResolution := getPrimaryDisplayResolution()

	// Desktop Environment and Window Manager
	desktopEnvironment := "Windows Explorer"
	windowManager := "DWM"
	currentTheme := getCurrentTheme()

	// CPU Information
	cpuModel, cpuCores, cpuThreads, cpuSpeed := getCPUInfo()

	// GPU Information
	gpuModel, gpuDriverVersion := getGPUInfo()

	// Memory Information
	totalRAM, usedRAM, freeRAM := getMemoryInfo()

	// Disk Information
	totalDiskSpace, usedDiskSpace, freeDiskSpace := getDiskInfo()

	// Battery Information
	batteryPercentage, batteryStatus := getBatteryInfo()

	return helper.SysInfo{
		Hostname:    hostname,
		CurrentUser: currentUser.Username,
		OSName:      osName,
		OSVersion:   osVersion,
		OtherInfo: helper.OtherInfo{
			KernelVersion:  kernelVersion,
			Uptime:         systemUptime,
			PublicIP:       "",
			Timezone:       "",
			Locale:         "",
			SystemLanguage: "",
			ScreenResolution: []helper.ScreenInfo{
				{
					Resolution: primaryDisplayResolution,
				},
			},
			CurrentTheme: currentTheme,
		},
		Software: helper.SoftwareInfo{
			OSDetails:          osName + " " + osVersion + " (Build " + osBuildNumber + ")",
			DesktopEnvironment: desktopEnvironment,
			WindowManager:      windowManager,
			RunningProcesses:   nil,
			StartupPrograms:    nil,
		},
		PackageManagement: helper.PackageManagementInfo{
			PackageCount: installedPackages,
		},
		Hardware: helper.HardwareInfo{
			CPU: helper.CPUInfo{
				Model:      cpuModel,
				Cores:      cpuCores,
				Threads:    cpuThreads,
				BaseSpeed:  cpuSpeed.Base,
				BoostSpeed: cpuSpeed.Boost,
			},
			GPU: helper.GPUInfo{
				Model:         gpuModel,
				DriverVersion: gpuDriverVersion,
			},
			Memory: helper.MemoryInfo{
				Total: totalRAM,
				Used:  usedRAM,
				Free:  freeRAM,
			},
			Disk: helper.DiskInfo{
				Total: totalDiskSpace,
				Used:  usedDiskSpace,
				Free:  freeDiskSpace,
			},
			Battery: helper.BatteryInfo{
				Percentage: batteryPercentage,
				Status:     batteryStatus,
			},
		},
	}
}

// Get OS Name, Version, and Build Number
func getOSInfo() (string, string, string) {
	var osInfo struct {
		Caption     string
		Version     string
		BuildNumber string
	}

	err := wmi.Query("SELECT Caption, Version, BuildNumber FROM Win32_OperatingSystem", &osInfo)
	if err != nil {
		return "Unknown", "Unknown", "Unknown"
	}

	return osInfo.Caption, osInfo.Version, osInfo.BuildNumber
}

// Get System Manufacturer and Model
func getSystemManufacturerAndModel() (string, string) {
	var systemInfo struct {
		Manufacturer string
		Model        string
	}

	err := wmi.Query("SELECT Manufacturer, Model FROM Win32_ComputerSystem", &systemInfo)
	if err != nil {
		return "Unknown", "Unknown"
	}

	return systemInfo.Manufacturer, systemInfo.Model
}

// Get Kernel Version
func getKernelVersion() string {
	version, err := windows.RtlGetVersion()
	if err != nil {
		return "Unknown"
	}
	return fmt.Sprintf("%d.%d.%d", version.MajorVersion, version.MinorVersion, version.BuildNumber)
}

// Get System Uptime
func getSystemUptime() string {
	var info windows.RTL_OSVERSIONINFOW
	info.OSVersionInfoSize = uint32(unsafe.Sizeof(info))
	if err := windows.RtlGetVersion(&info); err != nil {
		return "Unknown"
	}

	// Uptime calculation using system tick count
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	getTickCount64 := kernel32.NewProc("GetTickCount64")

	var uptime uint64
	_, _, _ = getTickCount64.Call(uintptr(unsafe.Pointer(&uptime)))

	// Convert uptime from milliseconds to readable format
	uptimeSecs := uptime / 1000
	days := uptimeSecs / (24 * 3600)
	hours := (uptimeSecs % (24 * 3600)) / 3600
	minutes := (uptimeSecs % 3600) / 60
	return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
}

// Get Installed Packages Count
func getInstalledPackagesCount() int {
	var softwareList []struct{}
	query := "SELECT * FROM Win32_Product"
	err := wmi.Query(query, &softwareList)
	if err != nil {
		return 0
	}
	return len(softwareList)
}

// Get Shell Name and Version
func getShellName() string {
	return "PowerShell" // Assuming PowerShell as default; could be extended to detect other shells
}

func getShellVersion() string {
	cmd := windows.Cmd{
		// Adjust as needed for actual shell detection
		Name: "powershell",
		Args: []string{"$PSVersionTable.PSVersion.ToString()"},
	}
	output, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}
	return string(output)
}

// Get Primary Display Resolution
func getPrimaryDisplayResolution() string {
	var devMode windows.DEVMODE
	devMode.Size = uint16(unsafe.Sizeof(devMode))
	if !windows.EnumDisplaySettings(nil, windows.ENUM_CURRENT_SETTINGS, &devMode) {
		return "Unknown"
	}
	return fmt.Sprintf("%dx%d", devMode.PelsWidth, devMode.PelsHeight)
}

// Get Current Theme (Dark/Light Mode)
func getCurrentTheme() string {
	var theme int32
	const keyPath = `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	if err := windows.RegQueryValueEx(syscall.HKEY_CURRENT_USER, keyPath, "AppsUseLightTheme", &theme); err != nil {
		return "Unknown"
	}
	if theme == 0 {
		return "Dark Mode"
	}
	return "Light Mode"
}

// Get CPU Information
func getCPUInfo() (model string, cores, threads int, speed struct{ Base, Boost string }) {
	var cpuInfo struct {
		Name          string
		NumberOfCores uint32
		ThreadCount   uint32
		MaxClockSpeed uint32
	}
	err := wmi.Query("SELECT Name, NumberOfCores, ThreadCount, MaxClockSpeed FROM Win32_Processor", &cpuInfo)
	if err != nil {
		return "Unknown", 0, 0, struct{ Base, Boost string }{"Unknown", "Unknown"}
	}
	model = cpuInfo.Name
	cores = int(cpuInfo.NumberOfCores)
	threads = int(cpuInfo.ThreadCount)
	speed = struct{ Base, Boost string }{fmt.Sprintf("%d MHz", cpuInfo.MaxClockSpeed), "Unknown"} // Boost speed not available
	return
}

// Get GPU Information
func getGPUInfo() (model, driverVersion string) {
	var gpuInfo struct {
		Name          string
		DriverVersion string
	}
	err := wmi.Query("SELECT Name, DriverVersion FROM Win32_VideoController", &gpuInfo)
	if err != nil {
		return "Unknown", "Unknown"
	}
	model = gpuInfo.Name
	driverVersion = gpuInfo.DriverVersion
	return
}

// Get Memory Information
func getMemoryInfo() (total, used, free uint64) {
	var memStatus windows.MEMORYSTATUSEX
	memStatus.Length = uint32(unsafe.Sizeof(memStatus))
	err := windows.GlobalMemoryStatusEx(&memStatus)
	if err != nil {
		return 0, 0, 0
	}
	total = memStatus.TotalPhys
	used = total - memStatus.AvailPhys
	free = memStatus.AvailPhys
	return
}

// Get Disk Information
func getDiskInfo() (total, used, free uint64) {
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	err := windows.GetDiskFreeSpaceEx(nil, &freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return 0, 0, 0
	}
	total = totalNumberOfBytes
	used = totalNumberOfBytes - totalNumberOfFreeBytes
	free = totalNumberOfFreeBytes
	return
}

// Get Battery Information
func getBatteryInfo() (percentage int, status string) {
	var batteryInfo struct {
		DesignCapacity     int
		FullChargeCapacity int
		Status             int
	}
	err := wmi.Query("SELECT DesignCapacity, FullChargeCapacity, BatteryStatus FROM Win32_Battery", &batteryInfo)
	if err != nil {
		return 0, "Unknown"
	}
	if batteryInfo.FullChargeCapacity > 0 {
		percentage = (batteryInfo.DesignCapacity * 100) / batteryInfo.FullChargeCapacity
	} else {
		percentage = 0
	}
	switch batteryInfo.Status {
	case 1:
		status = "Discharging"
	case 2:
		status = "Charging"
	case 3:
		status = "Fully Charged"
	default:
		status = "Unknown"
	}
	return
}
