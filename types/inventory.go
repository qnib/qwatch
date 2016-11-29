package qtypes

// DiscoveryFS holds the structure for the file
type DiscoveryFS struct {
	Discovery map[string]DiscoveryFilesystem `json:"discovery"`
	Timestamp int64                          `json:"timestamp"`
}

// DiscoveryFilesystem holds one filesystem
type DiscoveryFilesystem struct {
	Device       string            `json:"device"`
	FsOpts       map[string]string `json:"fsopts"`
	FsType       string            `json:"fstype"`
	Mountpoint   string            `json:"mountpoint"`
	Mountoptions []string          `json:"options"`
}

// DiscoveryProcesses holds a list of Processes
type DiscoveryProcesses struct {
	Discovery map[string]DiscoveryProcess `json:"discovery"`
	Timestamp int64                       `json:"timestamp"`
}

// DiscoveryProcess holds a process entry plus it's potential children
type DiscoveryProcess struct {
	Command   string             `json:"cmdline"`
	Exec      string             `json:"exe"`
	Groupname string             `json:"groupname"`
	Username  string             `json:"username"`
	Name      string             `json:"name"`
	Children  []DiscoveryProcess `json:"children"`
}
