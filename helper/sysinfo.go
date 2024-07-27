package helper

type SysInfo struct {
	Hostname      string
	CurrentUser   string
	OSName        string
	OSVersion     string
	OSCodename    string
	KernelVersion string
	Shell         string
	ShellVersion  string
	Architecture  string
	Uptime        string
	CPU           CPUInfo
	GPU           GPUInfo
	Motherboard   MotherboardInfo
	Memory        MemoryInfo
	Storage       []StorageInfo
	Network       []NetworkInfo
	Battery       BatteryInfo
	Peripherals   PeripheralInfo
	Software      SoftwareInfo
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
	OSDetails          string // Operating System details (Distro name or Windows edition)
	DesktopEnvironment string // Desktop Environment name and version
	WindowManager      string // Window Manager name and version
	WMTheme            string // Window Manager theme
	GTKTheme           string // GTK theme
	IconsTheme         string // Icons theme
	Font               string // Font used in the system
}
