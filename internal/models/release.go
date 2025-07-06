package models

type GitHubRepository struct {
	Owner string
	Name  string
}

type ReleaseAsset struct {
	Name string
	Path string
}
