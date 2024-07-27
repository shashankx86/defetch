package linux

import (
	"defetch/helper"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func GetLinuxInfo() helper.SysInfo {
	// Hostname
	hostname, _ := os.Hostname()

	// Current User
	currentUser, _ := user.Current()

	// Operating System
	osNameOutput, _ := exec.Command("lsb_release", "-ds").Output()
	osName := strings.TrimSpace(string(osNameOutput))
	osVersion, _ := exec.Command("lsb_release", "-rs").Output()
	osCodename, _ := exec.Command("lsb_release", "-cs").Output()

	// Kernel Version
	kernelVersion, _ := exec.Command("uname", "-r").Output()

	// Shell Name and Version
	shell := os.Getenv("SHELL")
	shellVersionOutput, _ := exec.Command(shell, "--version").Output()
	shellVersion := strings.SplitN(strings.TrimSpace(string(shellVersionOutput)), "\n", 2)[0]

	// Architecture
	architecture := runtime.GOARCH

	// Uptime
	uptimeOutput, _ := exec.Command("uptime", "-p").Output()

	return helper.SysInfo{
		Hostname:      hostname,
		CurrentUser:   currentUser.Username,
		OSName:        osName,
		OSVersion:     strings.TrimSpace(string(osVersion)),
		OSCodename:    strings.TrimSpace(string(osCodename)),
		KernelVersion: strings.TrimSpace(string(kernelVersion)),
		Shell:         shell,
		ShellVersion:  shellVersion,
		Architecture:  architecture,
		Uptime:        strings.TrimSpace(string(uptimeOutput)),
	}
}
