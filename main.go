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

	// GPU Information
	fmt.Printf("GPU Model: %s\n", sysInfo.GPU.ModelName)
	fmt.Printf("GPU Driver Version: %s\n", sysInfo.GPU.DriverVersion)
	fmt.Printf("GPU Memory Size: %s\n", sysInfo.GPU.MemorySize)

	// Motherboard Information
	fmt.Printf("Motherboard Manufacturer: %s\n", sysInfo.Motherboard.Manufacturer)
	fmt.Printf("Motherboard Model: %s\n", sysInfo.Motherboard.Model)
	fmt.Printf("BIOS/UEFI Version: %s\n", sysInfo.Motherboard.BIOSVersion)
	fmt.Printf("Motherboard Serial Number: %s\n", sysInfo.Motherboard.SerialNumber)

	// Memory Information
	fmt.Printf("Total Memory: %s\n", sysInfo.Memory.TotalSize)
	fmt.Printf("Used Memory: %s\n", sysInfo.Memory.UsedSize)
	fmt.Printf("Free Memory: %s\n", sysInfo.Memory.FreeSize)
	fmt.Printf("Memory Slots: %v\n", sysInfo.Memory.Slots)

	// Storage Information
	for _, storage := range sysInfo.Storage {
		fmt.Printf("Device: %s\n", storage.Device)
		fmt.Printf("Model: %s\n", storage.Model)
		fmt.Printf("Capacity: %s\n", storage.Capacity)
		fmt.Printf("Used: %s\n", storage.Used)
		fmt.Printf("Available: %s\n", storage.Available)
		fmt.Printf("File System: %s\n", storage.FileSystem)
		fmt.Printf("Mount Point: %s\n", storage.MountPoint)
		fmt.Printf("Read Speed: %s\n", storage.ReadSpeed)
		fmt.Printf("Write Speed: %s\n", storage.WriteSpeed)
	}
}
