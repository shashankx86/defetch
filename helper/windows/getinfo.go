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

func GetWindowsInfo() helper.SysInfo {
	// Hostname
	hostname, _ := os.Hostname()

	// Current User
	currentUser, _ := user.Current()

	// Operating System Info
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
