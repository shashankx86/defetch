package linux

import (
	"defetch/helper"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
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
