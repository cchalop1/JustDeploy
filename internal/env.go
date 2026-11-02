package internal

import (
	"os"
	"path/filepath"
)

var JUSTDEPLOY_FOLDER = dataDir()
var CERT_DOCKER_FOLDER = filepath.Join(JUSTDEPLOY_FOLDER, "cert-docker")
var DATABASE_SQLITE_PATH = filepath.Join(JUSTDEPLOY_FOLDER, "database.db")

func dataDir() string {
	if dir := os.Getenv("JUSTDEPLOY_DATA_DIR"); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic("cannot determine home directory: set JUSTDEPLOY_DATA_DIR")
	}
	return filepath.Join(home, ".config", "justdeploy")
}
