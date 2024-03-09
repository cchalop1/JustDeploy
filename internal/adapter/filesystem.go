package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// FilesystemAdapter provides functions for filesystem operations.
type FilesystemAdapter struct {
	// Any configuration or dependencies needed for filesystem operations
}

// NewFilesystemAdapter creates a new instance of FilesystemAdapter.
func NewFilesystemAdapter() *FilesystemAdapter {
	// Initialize and return a new instance if needed
	return &FilesystemAdapter{}
}

// GetFolderName returns the name of the folder given a path.
func (fs *FilesystemAdapter) GetFolderName(path string) string {
	// Clean the path to handle potential issues with separators
	cleanedPath := filepath.Clean(path)

	// get the absolute path to the project
	projectName := filepath.Base(cleanedPath)

	return projectName
}

func (fs *FilesystemAdapter) GetFullPathToProject(path string) string {
	// Clean the path to handle potential issues with separators
	cleanedPath := filepath.Clean(path)

	// get the absolute path to the project
	fullPathToProject, err := filepath.Abs(cleanedPath)

	if err != nil {
		panic("Error to find the path of the project")

	}
	return fullPathToProject
}

func (fs *FilesystemAdapter) CleanPath(path string) string {
	if path[len(path)-1] != '/' {
		return path + "/"
	}
	return path
}

func (fs *FilesystemAdapter) GetCurrentPath() string {
	return filepath.Base(".")
}

func (fs *FilesystemAdapter) IsWhereIsADockerFileInTheFolder(pathToFolder string) bool {
	entries, err := os.ReadDir(pathToFolder)
	if err != nil {
		panic("Error to read the directory you have")
	}

	for _, e := range entries {
		if e.Name() == "Dockerfile" {
			return true
		}
	}
	return false
}

func (fs *FilesystemAdapter) CopyFileToRemoteServer(sourcePath string, serverIp string) error {
	cmd := exec.Command("scp", sourcePath, "root@"+serverIp+":/root/")
	fmt.Println(cmd.String())
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(stdout))
	return nil
}
