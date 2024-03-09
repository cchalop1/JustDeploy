package adapter

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenBrowser(urlString string) error {
	var err error

	if urlString == "" {
		return fmt.Errorf("empty url string")
	}

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", urlString}
	case "darwin":
		cmd = "open"
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if cmd != "" {
		args = append([]string{urlString}, args...)
		err = exec.Command(cmd, args...).Start()
	}

	return err
}
