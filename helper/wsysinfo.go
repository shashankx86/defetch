package helper

// ScreenResolution holds information about display resolution
type ScreenResolution struct {
	Resolution string
}

// OtherInfo holds miscellaneous system information
type OtherInfo struct {
	KernelVersion    string
	Uptime           string
	ScreenResolution []ScreenResolution
	CurrentTheme     string
}

// PackageManagement holds information about installed packages
type PackageManagement struct {
	PackageCount int
}

// Software holds software-related information
type Software struct {
	ShellName          string
	ShellVersion       string
	DesktopEnvironment string
	WindowManager      string
}

// SysInfo holds the system information
type SysInfo struct {
	Hostname          string
	CurrentUser       string
	OSName            string
	OSVersion         string
	Manufacturer      string
	Model             string
	OtherInfo         OtherInfo
	PackageManagement PackageManagement
	Software          Software
}
