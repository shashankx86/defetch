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
