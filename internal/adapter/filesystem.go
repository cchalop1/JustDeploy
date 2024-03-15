package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type FilesystemAdapter struct {
}

func NewFilesystemAdapter() *FilesystemAdapter {
	return &FilesystemAdapter{}
}

func (fs *FilesystemAdapter) GetFolderName(path string) string {
	cleanedPath := filepath.Clean(path)

	projectName := filepath.Base(cleanedPath)

	return projectName
}

func (fs *FilesystemAdapter) GetFullPathToProject(path string) string {
	cleanedPath := filepath.Clean(path)

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

func (fs *FilesystemAdapter) GetCurrentPath() (string, error) {
	return os.Getwd()
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

func (fs *FilesystemAdapter) CreateGitPostCommitHooks(path string) error {
	hooksFilePath := path + ".git/hooks/post-commit"
	fileContent := []byte("#!/bin/sh\njustdeploy redeploy\n")

	err := os.WriteFile(hooksFilePath, fileContent, 0755)

	if err != nil {
		return err
	}

	fmt.Println("Create file ", path+".git/hooks/post-commit")
	return nil
}

func (fs *FilesystemAdapter) DeleteGitPostCommitHooks(path string) error {
	hooksFilePath := path + ".git/hooks/post-commit"

	err := os.Remove(hooksFilePath)

	if err != nil {
		return err
	}

	fmt.Println("Remove file", path+".git/hooks/post-commit")

	return nil
}
