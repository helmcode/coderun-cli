package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateBuildContext creates a tar.gz file containing the build context
func CreateBuildContext(contextDir, outputPath string) error {
	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Walk through the context directory
	return filepath.Walk(contextDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories (except .dockerignore)
		if strings.HasPrefix(filepath.Base(path), ".") && filepath.Base(path) != ".dockerignore" {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip the output file itself if it's in the context
		if path == outputPath {
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("failed to create tar header: %w", err)
		}

		// Set the name to the relative path from context directory
		relativePath, err := filepath.Rel(contextDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		header.Name = relativePath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header: %w", err)
		}

		// Write file content if it's a regular file
		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return fmt.Errorf("failed to copy file content: %w", err)
			}
		}

		return nil
	})
}

// ValidateDockerfile checks if the Dockerfile exists in the context
func ValidateDockerfile(contextDir, dockerfilePath string) error {
	fullPath := filepath.Join(contextDir, dockerfilePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("Dockerfile not found at %s", dockerfilePath)
	}
	return nil
}

// GenerateBuildContextPath generates a temporary path for the build context
func GenerateBuildContextPath(appName string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("coderun-build-%s.tar.gz", appName))
}
