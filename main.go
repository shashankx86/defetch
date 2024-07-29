package main

import (
	"defetch/helper/linux"
	"defetch/helper/windows"
	"fmt"
	"runtime"
)

func main() {
	switch runtime.GOOS {
	case "linux":
		sysInfo := linux.GetLinuxInfo()
		displaySystemInfo(sysInfo)
	case "windows":
		sysInfo := windows.GetWindowsInfo()
		displaySystemInfo(sysInfo)
	default:
		fmt.Println("Unsupported operating system.")
	}
}

func displaySystemInfo(sysInfo interface{}) {
	if info, ok := sysInfo.(linux.SysInfo); ok {

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

		// // Battery Information (TODO LATER)
		// fmt.Printf("Battery Status: %s\n", sysInfo.Battery.Status)
		// fmt.Printf("Battery Capacity: %s\n", sysInfo.Battery.Capacity)
		// fmt.Printf("Battery Percentage: %s\n", sysInfo.Battery.Percentage)
		// fmt.Printf("Battery Manufacturer: %s\n", sysInfo.Battery.Manufacturer)
		// fmt.Printf("Battery Model: %s\n", sysInfo.Battery.Model)

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

		// Browser Information
		fmt.Println("Browsers:")
		for _, browser := range sysInfo.Software.Browser {
			fmt.Printf("  %s: %s\n", browser.Name, browser.Version)
		}

		// Running Processes Information
		fmt.Printf("\nNumber of Running Processes: %d\n", len(sysInfo.Software.RunningProcesses))
		fmt.Println("Top Processes by CPU Usage:")
		for _, process := range sysInfo.Software.RunningProcesses {
			fmt.Printf("  PID: %d, Name: %s, CPU Usage: %.2f%%, Memory Usage: %.2f%%\n",
				process.PID, process.Name, process.CPUUsage, process.MemoryUsage)
		}

		// Startup Programs Information
		fmt.Println("\nStartup Programs:")
		for _, program := range sysInfo.Software.StartupPrograms {
			fmt.Printf("  Name: %s, Command: %s\n", program.Name, program.Command)
		}

		// System Performance Information
		fmt.Printf("\nOverall CPU Usage: %.2f%%\n", sysInfo.Performance.CPUUsage)
		fmt.Println("Per-Core CPU Usage:")
		for i, usage := range sysInfo.Performance.PerCoreUsage {
			fmt.Printf("  Core %d: %.2f%%\n", i, usage)
		}

		fmt.Printf("\nTotal Memory Used: %s\n", sysInfo.Performance.MemoryUsage.TotalUsed)
		fmt.Printf("Free Memory: %s\n", sysInfo.Performance.MemoryUsage.Free)
		fmt.Printf("Total Memory: %s\n", sysInfo.Performance.MemoryUsage.Total)

		fmt.Println("\nTop Applications by Memory Usage:")
		for _, app := range sysInfo.Performance.PerAppMemoryUsage {
			fmt.Printf("  PID: %d, Name: %s, Memory Usage: %.2f%%\n",
				app.PID, app.Name, app.MemoryUsage)
		}

		// Package Management Information
		fmt.Printf("\nNumber of Installed Packages: %d\n", sysInfo.PackageManagement.PackageCount)
		fmt.Printf("Number of Available Updates: %d\n", sysInfo.PackageManagement.AvailableUpdates)
		fmt.Printf("Used Package Managers: %v\n", sysInfo.PackageManagement.PackageManagers)
		fmt.Println("Recently Installed Packages:")
		for _, pkg := range sysInfo.PackageManagement.RecentlyInstalledPackages {
			fmt.Printf("  Name: %s, Version: %s, Installed Date: %s\n",
				pkg.Name, pkg.Version, pkg.InstalledDate)
		}

		// Other Information
		fmt.Printf("\nPublic IP: %s\n", sysInfo.OtherInfo.PublicIP)
		fmt.Printf("Timezone: %s\n", sysInfo.OtherInfo.Timezone)
		fmt.Printf("Locale: %s\n", sysInfo.OtherInfo.Locale)
		fmt.Printf("System Language: %s\n", sysInfo.OtherInfo.SystemLanguage)
		fmt.Printf("CPU Temperature: %.2f°C\n", sysInfo.OtherInfo.Temperature.CPU)
		fmt.Printf("GPU Temperature: %.2f°C\n", sysInfo.OtherInfo.Temperature.GPU)
		fmt.Printf("Motherboard Temperature: %.2f°C\n", sysInfo.OtherInfo.Temperature.Motherboard)
		fmt.Println("Screen Resolution:")
		for _, screen := range sysInfo.OtherInfo.ScreenResolution {
			fmt.Printf("  Model: %s, Resolution: %s, Refresh Rate: %d Hz\n",
				screen.Model, screen.Resolution, screen.RefreshRate)
		}
		fmt.Println("Disk Partitions:")
		for _, partition := range sysInfo.OtherInfo.DiskPartitions {
			fmt.Printf("  Device: %s, Filesystem: %s, Mount Point: %s, Size: %s, Used: %s, Available: %s\n",
				partition.Device, partition.Filesystem, partition.MountPoint, partition.Size, partition.Used, partition.Available)
		}
	} else if info, ok := sysInfo.(windows.SysInfo); ok {
		// Retrieve and print Windows system information
		sysInfo := windows.GetWindowsInfo()

		fmt.Printf("Hostname: %s\n", sysInfo.Hostname)
		fmt.Printf("Current User: %s\n", sysInfo.CurrentUser)
		fmt.Printf("OS Name: %s\n", sysInfo.OSName)
		fmt.Printf("OS Version: %s\n", sysInfo.OSVersion)
		fmt.Printf("Manufacturer: %s\n", sysInfo.Manufacturer)
		fmt.Printf("Model: %s\n", sysInfo.Model)
		fmt.Printf("Kernel Version: %s\n", sysInfo.OtherInfo.KernelVersion)
		fmt.Printf("System Uptime: %s\n", sysInfo.OtherInfo.Uptime)
		fmt.Printf("Installed Packages: %d\n", sysInfo.PackageManagement.PackageCount)
		fmt.Printf("Shell Name: %s\n", sysInfo.Software.ShellName)
		fmt.Printf("Shell Version: %s\n", sysInfo.Software.ShellVersion)
		fmt.Printf("Primary Display Resolution: %s\n", sysInfo.OtherInfo.ScreenResolution[0].Resolution)
		fmt.Printf("Desktop Environment: %s\n", sysInfo.Software.DesktopEnvironment)
		fmt.Printf("Window Manager: %s\n", sysInfo.Software.WindowManager)
		fmt.Printf("Current Theme: %s\n", sysInfo.OtherInfo.CurrentTheme)
	}
}
