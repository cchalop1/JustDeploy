package adapter

import (
	"cchalop1.com/deploy/internal/api/dto"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SshAdapter struct{ client *ssh.Client }

func NewSshAdapter() *SshAdapter { return &SshAdapter{} }
func (s *SshAdapter) Connect(config dto.ConnectNewServerDto) error {
	auth := []ssh.AuthMethod{}
	if config.Password != nil {
		auth = append(auth, ssh.Password(*config.Password))
	} else {
		if config.SshKey == nil {
			return errors.New("SSH password or private key required")
		}
		signer, err := ssh.ParsePrivateKey([]byte(*config.SshKey))
		if err != nil {
			return fmt.Errorf("invalid SSH key: %w", err)
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	path := os.Getenv("JUSTDEPLOY_KNOWN_HOSTS")
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, ".ssh", "known_hosts")
	}
	verify, err := knownhosts.New(path)
	if err != nil {
		return fmt.Errorf("load trusted SSH hosts: %w", err)
	}
	s.client, err = ssh.Dial("tcp", net.JoinHostPort(config.Ip, "22"), &ssh.ClientConfig{User: config.User, Auth: auth, HostKeyCallback: verify, Timeout: 15 * time.Second})
	return err
}
func (s *SshAdapter) CloseConnection() {
	if s.client != nil {
		_ = s.client.Close()
	}
}
func (s *SshAdapter) RunCommand(command string) (string, error) {
	if s.client == nil {
		return "", errors.New("SSH not connected")
	}
	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	timer := time.AfterFunc(2*time.Minute, func() { _ = session.Close() })
	defer timer.Stop()
	output, err := session.CombinedOutput(command)
	return string(output), err
}
func (s *SshAdapter) SaveRemoteFileContentToLocalFile(remote, local string) error {
	if strings.ContainsRune(remote, 0) {
		return errors.New("invalid remote path")
	}
	quoted := "'" + strings.ReplaceAll(remote, "'", "'\"'\"'") + "'"
	content, err := s.RunCommand("cat -- " + quoted)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(local, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if err := file.Chmod(0600); err != nil {
		file.Close()
		return err
	}
	if _, err := file.WriteString(content); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}
