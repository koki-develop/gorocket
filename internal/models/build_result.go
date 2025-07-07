package models

// BuildCommandResult contains all data produced by a build command execution
type BuildCommandResult struct {
	// BuildInfo contains version and module information
	BuildInfo *BuildInfo
	// Config contains the processed configuration
	Config *Config
	// TemplateData contains the template variables used
	TemplateData *TemplateData
	// ArchiveResults contains information about created archives
	ArchiveResults []ArchiveResult
}
