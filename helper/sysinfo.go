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
