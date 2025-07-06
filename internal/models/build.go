package models

type BuildTarget struct {
	OS   string
	Arch string
}

type BuildResult struct {
	Target     BuildTarget
	BinaryPath string
}

type ArchiveResult struct {
	Target      BuildTarget
	ArchivePath string
}

type BuildInfo struct {
	ModuleName string
	Version    string
}

type FormulaInfo struct {
	ModuleName   string
	Version      string
	Repository   Repository
	PlatformURLs map[string]map[string]FormulaURL
}

type FormulaURL struct {
	URL    string
	SHA256 string
}
