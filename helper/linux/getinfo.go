package linux

import (
	"defetch/helper"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func GetLinuxInfo() helper.SysInfo {
	// Hostname
	hostname, _ := os.Hostname()

	// Operating System
	osNameOutput, _ := exec.Command("lsb_release", "-ds").Output()
	osName := strings.TrimSpace(string(osNameOutput))
	osVersion, _ := exec.Command("lsb_release", "-rs").Output()
	osCodename, _ := exec.Command("lsb_release", "-cs").Output()

	// Kernel Version
	kernelVersion, _ := exec.Command("uname", "-r").Output()

	// Architecture
	architecture := runtime.GOARCH

	// Uptime
	uptimeOutput, _ := exec.Command("uptime", "-p").Output()

	return helper.SysInfo{
		Hostname:      hostname,
		OSName:        osName,
		OSVersion:     strings.TrimSpace(string(osVersion)),
		OSCodename:    strings.TrimSpace(string(osCodename)),
		KernelVersion: strings.TrimSpace(string(kernelVersion)),
		Architecture:  architecture,
		Uptime:        strings.TrimSpace(string(uptimeOutput)),
	}
}
