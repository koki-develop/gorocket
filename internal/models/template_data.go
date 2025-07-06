package models

// TemplateData contains data available for template processing in configuration files
type TemplateData struct {
	// Version is the current git tag version (e.g., "v1.0.0")
	Version string
	// Module is the Go module name from go.mod
	Module string
}
