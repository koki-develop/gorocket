package gorocket

import "github.com/koki-develop/gorocket/internal/git"

// version is set by -ldflags during build
var version = "dev"

// GetVersion returns version information
func GetVersion() (string, error) {
	// Return embedded version if available
	if version != "" && version != "dev" {
		return version, nil
	}

	// Try to get version from git tag in development
	gitClient := git.New()
	tag, err := gitClient.GetHeadTag()
	if err != nil {
		return "dev", nil
	}
	return tag, nil
}
