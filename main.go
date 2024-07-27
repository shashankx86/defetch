package main

import (
	"defetch/helper/linux"
	"fmt"
)

func main() {
	sysInfo := linux.GetLinuxInfo()

	fmt.Printf("Hostname: %s\n", sysInfo.Hostname)
	fmt.Printf("Current User: %s\n", sysInfo.CurrentUser)
	fmt.Printf("Operating System: %s %s (%s)\n", sysInfo.OSName, sysInfo.OSVersion, sysInfo.OSCodename)
	fmt.Printf("Kernel Version: %s\n", sysInfo.KernelVersion)
	fmt.Printf("Shell: %s\n", sysInfo.Shell)
	fmt.Printf("Architecture: %s\n", sysInfo.Architecture)
	fmt.Printf("Uptime: %s\n", sysInfo.Uptime)
}
