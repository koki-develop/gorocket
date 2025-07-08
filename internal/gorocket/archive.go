package gorocket

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// createArchive creates an archive from build result
func createArchive(result *BuildResult) (string, error) {
	// Extract module name from binary path
	binaryName := filepath.Base(result.Binary)
	moduleName := strings.TrimSuffix(binaryName, filepath.Ext(binaryName))

	// Determine archive name
	var archiveName string
	if result.OS == "windows" {
		archiveName = fmt.Sprintf("%s_%s_%s_%s.zip", moduleName, result.Version, result.OS, result.Arch)
		return createZip(result.Binary, archiveName, moduleName, result.Version, result.OS, result.Arch)
	} else {
		archiveName = fmt.Sprintf("%s_%s_%s_%s.tar.gz", moduleName, result.Version, result.OS, result.Arch)
		return createTarGz(result.Binary, archiveName, moduleName, result.Version, result.OS, result.Arch)
	}
}

// createTarGz creates a tar.gz archive
func createTarGz(src, archiveName, moduleName, version, osName, arch string) (string, error) {
	archivePath := filepath.Join("dist", archiveName)

	// Create archive file
	file, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to create archive file: %w", err)
	}
	defer file.Close()

	// gzip writer
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	// tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Get file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Directory structure in archive
	dirName := fmt.Sprintf("%s_%s_%s_%s", moduleName, version, osName, arch)
	binaryNameInArchive := filepath.Join(dirName, moduleName)

	// Create tar header
	header := &tar.Header{
		Name: binaryNameInArchive,
		Mode: 0755,
		Size: srcInfo.Size(),
	}

	// Write header
	if err := tarWriter.WriteHeader(header); err != nil {
		return "", fmt.Errorf("failed to write tar header: %w", err)
	}

	// Copy file content
	if _, err := io.Copy(tarWriter, srcFile); err != nil {
		return "", fmt.Errorf("failed to write file to tar: %w", err)
	}

	return archivePath, nil
}

// createZip creates a zip archive
func createZip(src, archiveName, moduleName, version, osName, arch string) (string, error) {
	archivePath := filepath.Join("dist", archiveName)

	// Create archive file
	file, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to create archive file: %w", err)
	}
	defer file.Close()

	// zip writer
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Directory structure in archive
	dirName := fmt.Sprintf("%s_%s_%s_%s", moduleName, version, osName, arch)
	binaryNameInArchive := filepath.Join(dirName, moduleName+".exe")

	// Create zip entry
	writer, err := zipWriter.Create(binaryNameInArchive)
	if err != nil {
		return "", fmt.Errorf("failed to create zip entry: %w", err)
	}

	// Copy file content
	if _, err := io.Copy(writer, srcFile); err != nil {
		return "", fmt.Errorf("failed to write file to zip: %w", err)
	}

	return archivePath, nil
}
