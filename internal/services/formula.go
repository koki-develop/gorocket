package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers"
)

type FormulaService interface {
	GenerateFormula(buildInfo models.BuildInfo, archiveResults []models.ArchiveResult, brewConfig models.BrewConfig) error
}

type formulaService struct {
	fsProvider providers.FileSystemProvider
}

func NewFormulaService(fsProvider providers.FileSystemProvider) FormulaService {
	return &formulaService{
		fsProvider: fsProvider,
	}
}

func (s *formulaService) GenerateFormula(buildInfo models.BuildInfo, archiveResults []models.ArchiveResult, brewConfig models.BrewConfig) error {
	formulaInfo, err := s.buildFormulaInfo(buildInfo, archiveResults, brewConfig)
	if err != nil {
		return fmt.Errorf("failed to build formula info: %w", err)
	}

	templateData := s.buildTemplateData(formulaInfo)
	formulaContent, err := providers.GenerateFormula(templateData)
	if err != nil {
		return fmt.Errorf("failed to generate formula: %w", err)
	}

	formulaFileName := fmt.Sprintf("%s.rb", buildInfo.ModuleName)
	formulaPath := filepath.Join("dist", formulaFileName)

	if err := s.fsProvider.WriteFile(formulaPath, []byte(formulaContent), 0644); err != nil {
		return fmt.Errorf("failed to write formula file: %w", err)
	}

	return nil
}

func (s *formulaService) buildFormulaInfo(buildInfo models.BuildInfo, archiveResults []models.ArchiveResult, brewConfig models.BrewConfig) (models.FormulaInfo, error) {
	platformURLs := make(map[string]map[string]models.FormulaURL)

	for _, result := range archiveResults {

		file, err := s.fsProvider.Open(result.ArchivePath)
		if err != nil {
			return models.FormulaInfo{}, fmt.Errorf("failed to open archive %s: %w", result.ArchivePath, err)
		}
		defer func() { _ = file.Close() }()

		sha256, err := s.fsProvider.CalculateSHA256(file)
		if err != nil {
			return models.FormulaInfo{}, fmt.Errorf("failed to calculate SHA256 for %s: %w", result.ArchivePath, err)
		}

		archiveFileName := filepath.Base(result.ArchivePath)
		downloadURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
			brewConfig.Repository.Owner, brewConfig.Repository.Name, buildInfo.Version, archiveFileName)

		os := result.Target.OS
		arch := result.Target.Arch

		if platformURLs[os] == nil {
			platformURLs[os] = make(map[string]models.FormulaURL)
		}

		platformURLs[os][arch] = models.FormulaURL{
			URL:    downloadURL,
			SHA256: sha256,
		}
	}

	version := strings.TrimPrefix(buildInfo.Version, "v")

	return models.FormulaInfo{
		ModuleName:   buildInfo.ModuleName,
		Version:      version,
		Repository:   brewConfig.Repository,
		PlatformURLs: platformURLs,
	}, nil
}

func (s *formulaService) buildTemplateData(formulaInfo models.FormulaInfo) providers.FormulaTemplateData {
	return providers.FormulaTemplateData{
		ClassName:        providers.ToClassName(formulaInfo.ModuleName),
		Version:          formulaInfo.Version,
		ModuleName:       formulaInfo.ModuleName,
		MacOSARM64URL:    s.getURL(formulaInfo.PlatformURLs, "darwin", "arm64"),
		MacOSARM64SHA256: s.getSHA256(formulaInfo.PlatformURLs, "darwin", "arm64"),
		MacOSAMD64URL:    s.getURL(formulaInfo.PlatformURLs, "darwin", "amd64"),
		MacOSAMD64SHA256: s.getSHA256(formulaInfo.PlatformURLs, "darwin", "amd64"),
		LinuxARM64URL:    s.getURL(formulaInfo.PlatformURLs, "linux", "arm64"),
		LinuxARM64SHA256: s.getSHA256(formulaInfo.PlatformURLs, "linux", "arm64"),
		LinuxAMD64URL:    s.getURL(formulaInfo.PlatformURLs, "linux", "amd64"),
		LinuxAMD64SHA256: s.getSHA256(formulaInfo.PlatformURLs, "linux", "amd64"),
	}
}

func (s *formulaService) getURL(platformURLs map[string]map[string]models.FormulaURL, os, arch string) string {
	if osURLs, exists := platformURLs[os]; exists {
		if url, exists := osURLs[arch]; exists {
			return url.URL
		}
	}
	return ""
}

func (s *formulaService) getSHA256(platformURLs map[string]map[string]models.FormulaURL, os, arch string) string {
	if osURLs, exists := platformURLs[os]; exists {
		if url, exists := osURLs[arch]; exists {
			return url.SHA256
		}
	}
	return ""
}
