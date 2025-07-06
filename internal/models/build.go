package models

type BuildTarget struct {
	OS   string
	Arch string
}

type BuildResult struct {
	Target     BuildTarget
	BinaryPath string
	Error      error
}

type ArchiveResult struct {
	Target      BuildTarget
	ArchivePath string
	Error       error
}

type BuildInfo struct {
	ModuleName string
	Version    string
}
