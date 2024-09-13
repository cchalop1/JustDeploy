package adapter

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"time"
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

func isPortIsUse(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func FindOpenLocalPort(port int) string {
	for {
		if !isPortIsUse("localhost", port) {
			break
		}
		port++
	}
	return strconv.Itoa(port)
}
