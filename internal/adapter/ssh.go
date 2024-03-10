package adapter

import (
	"fmt"
	"log"
	"os"

	"cchalop1.com/deploy/internal/domain"
	"golang.org/x/crypto/ssh"
)

type SshAdapter struct {
	client *ssh.Client
}

func NewSshAdapter() *SshAdapter {
	// Initialize and return a new instance if needed
	return &SshAdapter{}
}

func (s *SshAdapter) Connect(connectConfig domain.ConnectServerDto) {
	signer, err := ssh.ParsePrivateKey([]byte(*connectConfig.SshKey))
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err)
	}
	config := &ssh.ClientConfig{
		User: connectConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", connectConfig.Domain+":22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	s.client = client
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
