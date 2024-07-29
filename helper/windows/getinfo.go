package windows

import (
	"defetch/helper"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetWindowsInfo() helper.SysInfo {
	// Hostname
	hostname, _ := os.Hostname()

	// Current User
	currentUser, _ := user.Current()

	// Operating System
	osName := "Windows"
	osVersion := getOSVersion()

	// Architecture
	architecture := runtime.GOARCH

	return helper.SysInfo{
		Hostname:     hostname,
		CurrentUser:  currentUser.Username,
		OSName:       osName,
		OSVersion:    osVersion,
		Architecture: architecture,
		// Additional fields can be filled similarly
	}
}

// Helper function to get the OS version on Windows
func getOSVersion() string {
	var info windows.OSVERSIONINFOEX
	info.OSVersionInfoSize = uint32(unsafe.Sizeof(info))
	err := windows.RtlGetVersion(&info)
	if err != nil {
		return "Unknown"
	}
	return fmt.Sprintf("%d.%d.%d", info.MajorVersion, info.MinorVersion, info.BuildNumber)
}
