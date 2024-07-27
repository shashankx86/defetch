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

	// Network Information
	for _, network := range sysInfo.Network {
		fmt.Printf("Interface: %s\n", network.InterfaceName)
		fmt.Printf("IP Address: %s\n", network.IPAddress)
		fmt.Printf("MAC Address: %s\n", network.MACAddress)
		fmt.Printf("Network Speed: %s\n", network.Speed)
		fmt.Printf("Active: %t\n", network.Active)
		fmt.Printf("Default Gateway: %s\n", network.DefaultGateway)
	}

	// Battery Information
	fmt.Printf("Battery Status: %s\n", sysInfo.Battery.Status)
	fmt.Printf("Battery Capacity: %s\n", sysInfo.Battery.Capacity)
	fmt.Printf("Battery Percentage: %s\n", sysInfo.Battery.Percentage)
	fmt.Printf("Battery Manufacturer: %s\n", sysInfo.Battery.Manufacturer)
	fmt.Printf("Battery Model: %s\n", sysInfo.Battery.Model)

	// Peripherals Information
	fmt.Println("Connected Devices:")
	for _, device := range sysInfo.Peripherals.ConnectedDevices {
		fmt.Printf("- %s\n", device)
	}

	fmt.Println("USB Devices:")
	for _, usbDevice := range sysInfo.Peripherals.USBDevices {
		fmt.Printf("- Name: %s, Vendor: %s, Product ID: %s, Vendor ID: %s\n", usbDevice.Name, usbDevice.Vendor, usbDevice.ProductID, usbDevice.VendorID)
	}

	fmt.Println("Audio Devices:")
	for _, audioDevice := range sysInfo.Peripherals.AudioDevices {
		fmt.Printf("- %s\n", audioDevice)
	}

	fmt.Println("Printer Details:")
	for _, printer := range sysInfo.Peripherals.PrinterDetails {
		fmt.Printf("- %s\n", printer)
	}

	// Software Information
	fmt.Printf("Operating System: %s\n", sysInfo.Software.OSDetails)
	fmt.Printf("Desktop Environment: %s\n", sysInfo.Software.DesktopEnvironment)
	fmt.Printf("Window Manager: %s\n", sysInfo.Software.WindowManager)
	fmt.Printf("WM Theme: %s\n", sysInfo.Software.WMTheme)
	fmt.Printf("GTK Theme: %s\n", sysInfo.Software.GTKTheme)
	fmt.Printf("Icons Theme: %s\n", sysInfo.Software.IconsTheme)
	fmt.Printf("Font: %s\n", sysInfo.Software.Font)
}
