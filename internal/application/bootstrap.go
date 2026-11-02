package application

import (
	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/utils"
	"crypto/subtle"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EnsureBootstrapToken(db *adapter.DatabaseAdapter) error {
	settings, err := db.ReadSettings()
	if err != nil {
		return err
	}
	if settings.AdminEmail != "" {
		return nil
	}
	path := filepath.Join(internal.JUSTDEPLOY_FOLDER, "setup-token")
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f, utils.GenerateRandomPassword(48))
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}
func CheckBootstrapToken(value string) bool {
	data, err := os.ReadFile(filepath.Join(internal.JUSTDEPLOY_FOLDER, "setup-token"))
	if err != nil {
		return false
	}
	expected := strings.TrimSpace(string(data))
	return len(expected) >= 32 && subtle.ConstantTimeCompare([]byte(value), []byte(expected)) == 1
}
