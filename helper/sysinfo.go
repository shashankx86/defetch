package helper

type SysInfo struct {
	Hostname          string
	CurrentUser       string
	OSName            string
	OSVersion         string
	OSCodename        string
	KernelVersion     string
	Shell             string
	ShellVersion      string
	Architecture      string
	Uptime            string
	CPU               CPUInfo
	GPU               GPUInfo
	Motherboard       MotherboardInfo
	Memory            MemoryInfo
	Storage           []StorageInfo
	Network           []NetworkInfo
	Battery           BatteryInfo
	Peripherals       PeripheralInfo
	Software          SoftwareInfo
	Performance       PerformanceInfo
	PackageManagement PackageManagementInfo
	OtherInfo         OtherInfo
}

type CPUInfo struct {
	ModelName    string
	Cores        int
	Threads      int
	Architecture string
	Frequency    float64
	CacheSize    int32
	Flags        string
}

type GPUInfo struct {
	ModelName     string
	DriverVersion string
	MemorySize    string
}

type MotherboardInfo struct {
	Manufacturer string
	Model        string
	BIOSVersion  string
	SerialNumber string
}

type MemoryInfo struct {
	TotalSize string
	UsedSize  string
	FreeSize  string
	Slots     []MemorySlotInfo
}

type MemorySlotInfo struct {
	Size       string
	FormFactor string
	Type       string
	Speed      string
}

type StorageInfo struct {
	Device     string
	Model      string
	Capacity   string
	Used       string
	Available  string
	FileSystem string
	MountPoint string
	ReadSpeed  string
	WriteSpeed string
}

type NetworkInfo struct {
	InterfaceName  string
	IPAddress      string
	MACAddress     string
	Speed          string
	Active         bool
	DefaultGateway string
}

type BatteryInfo struct {
	Status       string // Charging or discharging status
	Capacity     string // Battery capacity
	Percentage   string // Battery percentage remaining
	Manufacturer string // Battery manufacturer
	Model        string // Battery model
}

type PeripheralInfo struct {
	ConnectedDevices []string        // List of connected devices (e.g., mouse, keyboard, monitors)
	USBDevices       []USBDeviceInfo // List of USB devices
	AudioDevices     []string        // List of audio devices
	PrinterDetails   []string        // List of printer details
}

type USBDeviceInfo struct {
	Name      string // Name of the USB device
	Vendor    string // Vendor name
	ProductID string // Product ID
	VendorID  string // Vendor ID
}

type SoftwareInfo struct {
	OSDetails          string           // Operating System details (Distro name or Windows edition)
	DesktopEnvironment string           // Desktop Environment name and version
	WindowManager      string           // Window Manager name and version
	WMTheme            string           // Window Manager theme
	GTKTheme           string           // GTK theme
	IconsTheme         string           // Icons theme
	Font               string           // Font used in the system
	Browser            []BrowserInfo    // List of installed browsers and their versions
	RunningProcesses   []ProcessInfo    // Information on running processes
	StartupPrograms    []StartupProgram // List of programs that run on startup
}

type BrowserInfo struct {
	Name    string // Browser name
	Version string // Browser version
}

type ProcessInfo struct {
	PID         int     // Process ID
	Name        string  // Process name
	CPUUsage    float64 // CPU usage percentage
	MemoryUsage float64 // Memory usage percentage
}

type StartupProgram struct {
	Name    string // Program name
	Command string // Command or path to the executable
}

type PerformanceInfo struct {
	CPUUsage          float64          // Overall CPU usage percentage
	PerCoreUsage      []float64        // CPU usage percentage per core
	MemoryUsage       MemoryUsageInfo  // Memory usage information
	PerAppMemoryUsage []AppMemoryUsage // Memory usage per application
}

type MemoryUsageInfo struct {
	TotalUsed string // Total used memory
	Free      string // Free memory
	Total     string // Total memory
}

type AppMemoryUsage struct {
	Name        string  // Application name
	PID         int     // Process ID
	MemoryUsage float64 // Memory usage percentage
}

type PackageManagementInfo struct {
	PackageCount              int           // Number of installed packages (Linux) or programs (Windows)
	AvailableUpdates          int           // Number of available updates (Linux)
	PackageManagers           []string      // List of used package managers (Linux)
	RecentlyInstalledPackages []PackageInfo // List of recently installed packages (Linux)
}

type PackageInfo struct {
	Name          string // Package name
	Version       string // Package version
	InstalledDate string // Installation date
}

type OtherInfo struct {
	PublicIP         string          // Public IP address
	Timezone         string          // Timezone
	Locale           string          // Locale
	Temperature      TemperatureInfo // Temperature sensors information
	SystemLanguage   string          // System language
	ScreenResolution []ScreenInfo    // Screen resolution and monitor details
	DiskPartitions   []PartitionInfo // Disk partitions
}

type TemperatureInfo struct {
	CPU         float64 // CPU temperature
	GPU         float64 // GPU temperature
	Motherboard float64 // Motherboard temperature
}

type ScreenInfo struct {
	Model       string // Monitor model
	Resolution  string // Resolution
	RefreshRate int    // Refresh rate in Hz
}

type PartitionInfo struct {
	Device     string // Partition device name
	MountPoint string // Mount point
	Filesystem string // Filesystem type
	Size       string // Partition size
	Used       string // Used space
	Available  string // Available space
}
