package services

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers"
)

type ArchiverService interface {
	CreateArchives(buildInfo *models.BuildInfo, buildResults []models.BuildResult) ([]models.ArchiveResult, error)
	CreateArchive(buildInfo *models.BuildInfo, buildResult models.BuildResult) models.ArchiveResult
}

type archiverService struct {
	fileSystemProvider providers.FileSystemProvider
}

func NewArchiverService(fileSystemProvider providers.FileSystemProvider) ArchiverService {
	return &archiverService{
		fileSystemProvider: fileSystemProvider,
	}
}

func (a *archiverService) CreateArchives(buildInfo *models.BuildInfo, buildResults []models.BuildResult) ([]models.ArchiveResult, error) {
	var results []models.ArchiveResult

	for _, buildResult := range buildResults {
		if buildResult.Error != nil {
			results = append(results, models.ArchiveResult{
				Target: buildResult.Target,
				Error:  buildResult.Error,
			})
			continue
		}

		archiveResult := a.CreateArchive(buildInfo, buildResult)
		results = append(results, archiveResult)
	}

	return results, nil
}

func (a *archiverService) CreateArchive(buildInfo *models.BuildInfo, buildResult models.BuildResult) models.ArchiveResult {
	var archiveName string
	var archivePath string
	var err error

	if buildResult.Target.OS == "windows" {
		archiveName = fmt.Sprintf("%s_%s_%s_%s.zip", buildInfo.ModuleName, buildInfo.Version, buildResult.Target.OS, buildResult.Target.Arch)
		archivePath, err = a.createZipArchive(buildResult.BinaryPath, archiveName, buildInfo, buildResult.Target)
	} else {
		archiveName = fmt.Sprintf("%s_%s_%s_%s.tar.gz", buildInfo.ModuleName, buildInfo.Version, buildResult.Target.OS, buildResult.Target.Arch)
		archivePath, err = a.createTarGzArchive(buildResult.BinaryPath, archiveName, buildInfo, buildResult.Target)
	}

	return models.ArchiveResult{
		Target:      buildResult.Target,
		ArchivePath: archivePath,
		Error:       err,
	}
}

func (a *archiverService) createTarGzArchive(binaryPath, archiveName string, buildInfo *models.BuildInfo, target models.BuildTarget) (string, error) {
	archivePath := filepath.Join("dist", archiveName)

	file, err := a.fileSystemProvider.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to create archive file: %w", err)
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	binaryFile, err := a.fileSystemProvider.Open(binaryPath)
	if err != nil {
		return "", fmt.Errorf("failed to open binary file: %w", err)
	}
	defer binaryFile.Close()

	binaryInfo, err := a.fileSystemProvider.Stat(binaryPath)
	if err != nil {
		return "", fmt.Errorf("failed to get binary file info: %w", err)
	}

	dirName := fmt.Sprintf("%s_%s_%s_%s", buildInfo.ModuleName, buildInfo.Version, target.OS, target.Arch)
	binaryNameInArchive := filepath.Join(dirName, buildInfo.ModuleName)

	header := &tar.Header{
		Name: binaryNameInArchive,
		Mode: 0755,
		Size: binaryInfo.Size(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return "", fmt.Errorf("failed to write tar header: %w", err)
	}

	if _, err := io.Copy(tarWriter, binaryFile); err != nil {
		return "", fmt.Errorf("failed to write binary to tar: %w", err)
	}

	return archivePath, nil
}

func (a *archiverService) createZipArchive(binaryPath, archiveName string, buildInfo *models.BuildInfo, target models.BuildTarget) (string, error) {
	archivePath := filepath.Join("dist", archiveName)

	file, err := a.fileSystemProvider.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to create archive file: %w", err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	binaryFile, err := a.fileSystemProvider.Open(binaryPath)
	if err != nil {
		return "", fmt.Errorf("failed to open binary file: %w", err)
	}
	defer binaryFile.Close()

	dirName := fmt.Sprintf("%s_%s_%s_%s", buildInfo.ModuleName, buildInfo.Version, target.OS, target.Arch)
	binaryNameInArchive := filepath.Join(dirName, strings.TrimSuffix(filepath.Base(binaryPath), filepath.Ext(binaryPath)))

	writer, err := zipWriter.Create(binaryNameInArchive + ".exe")
	if err != nil {
		return "", fmt.Errorf("failed to create zip entry: %w", err)
	}

	if _, err := io.Copy(writer, binaryFile); err != nil {
		return "", fmt.Errorf("failed to write binary to zip: %w", err)
	}

	return archivePath, nil
}
