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
	fmt.Printf("Shell Version: %s\n", sysInfo.ShellVersion)
	fmt.Printf("Architecture: %s\n", sysInfo.Architecture)
	fmt.Printf("Uptime: %s\n", sysInfo.Uptime)

	// CPU Information
	fmt.Printf("CPU Model: %s\n", sysInfo.CPU.ModelName)
	fmt.Printf("CPU Cores: %d\n", sysInfo.CPU.Cores)
	fmt.Printf("CPU Threads: %d\n", sysInfo.CPU.Threads)
	fmt.Printf("CPU Architecture: %s\n", sysInfo.CPU.Architecture)
	fmt.Printf("CPU Frequency: %.2f MHz\n", sysInfo.CPU.Frequency)
	fmt.Printf("CPU Cache Size: %d KB\n", sysInfo.CPU.CacheSize)
	fmt.Printf("CPU Flags: %s\n", sysInfo.CPU.Flags)
}
