package adapter

import (
	"fmt"
	"log"
	"os"

	"cchalop1.com/deploy/internal/api/dto"
	"golang.org/x/crypto/ssh"
)

type SshAdapter struct {
	client *ssh.Client
}

func NewSshAdapter() *SshAdapter {
	// Initialize and return a new instance if needed
	return &SshAdapter{}
}

func (s *SshAdapter) getAuthMethode(connectConfig dto.ConnectNewServerDto) []ssh.AuthMethod {
	if connectConfig.Password != nil && connectConfig.SshKey == nil {
		return []ssh.AuthMethod{
			ssh.Password(*connectConfig.Password),
		}
	}
	signer, err := ssh.ParsePrivateKey([]byte(*connectConfig.SshKey))
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err)
	}
	return []ssh.AuthMethod{
		ssh.PublicKeys(signer),
	}
}

func (s *SshAdapter) Connect(connectConfig dto.ConnectNewServerDto) error {
	config := &ssh.ClientConfig{
		User:            connectConfig.User,
		Auth:            s.getAuthMethode(connectConfig),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", connectConfig.Domain+":22", config)
	if err != nil {
		return err
	}
	s.client = client
	return nil
}

func (s *SshAdapter) CloseConnection() {
	defer s.client.Close()
}

func (s *SshAdapter) RunCommand(command string) (string, error) {
	fmt.Println("$ " + command)
	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	output, err := session.CombinedOutput(command)
	return string(output), err
}

func (s *SshAdapter) SaveRemoteFileContentToLocalFile(remoteFilePath string, localFilePath string) error {
	catCommand := fmt.Sprintf("cat %s", remoteFilePath)
	content, err := s.RunCommand(catCommand)
	if err != nil {
		return fmt.Errorf("failed to execute 'cat' command: %v", err)
	}

	err = os.WriteFile(localFilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write content to local file: %v", err)
	}

	fmt.Println("File content saved successfully.")
	return nil
}
