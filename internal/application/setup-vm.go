package application

import (
	"fmt"
	"log"
	"os"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func ConnectAndSetupServer(deployService *service.DeployService, server domain.Server) *adapter.DockerAdapter {
	sshAdapter := adapter.NewSshAdapter()

	sshAdapter.Connect(dto.ConnectNewServerDto{
		Domain:   server.Domain,
		SshKey:   server.SshKey,
		Password: server.Password,
		User:     "root",
	})

	dockerIsInstalled, err := checkIfDockerIsIntalled(sshAdapter)
	fmt.Println(err)

	if !dockerIsInstalled {
		err = installDocker(sshAdapter)
		fmt.Println(err)
	}

	certificateIsCreated, err := checkIsCertificateIsCreate(sshAdapter, server.Id)
	fmt.Println(err)

	if !certificateIsCreated {
		err = setupDockerCertificates(sshAdapter, server)
		copyCertificates(sshAdapter, server.Id)
		fmt.Println(err)
	}

	portIsOpen, err := checkIfDockerPortIsOpen(sshAdapter)
	fmt.Println("port is open ", portIsOpen)

	if !portIsOpen {
		openPortDockerConfig(sshAdapter)
		fmt.Println(err)
	}

	sshAdapter.CloseConnection()
	adapterDocker := adapter.NewDockerAdapter()
	adapterDocker.ConnectClient(server)

	server.Status = "Runing"

	deployService.DatabaseAdapter.UpdateServer(server)

	return adapterDocker
}

func checkIfDockerIsIntalled(sshAdapter *adapter.SshAdapter) (bool, error) {
	output, err := sshAdapter.RunCommand("docker --version")
	if err != nil {
		return false, err
	}

	if len(output) > 0 {
		return true, nil
	}

	return false, nil
}

func checkIsCertificateIsCreate(sshAdapter *adapter.SshAdapter, serverId string) (bool, error) {
	certificatePath := "/root/cert-docker/" + serverId
	statCommand := fmt.Sprintf("stat -t -- \"%s\" &>/dev/null", certificatePath)
	output, err := sshAdapter.RunCommand(statCommand)
	if err != nil {
		return false, err
	}

	// If the output is empty, the certificate file exists
	if len(output) == 0 {
		return true, nil
	}

	// If the output is not empty, the certificate file does not exist
	return false, nil
}

func checkIfDockerPortIsOpen(sshAdapter *adapter.SshAdapter) (bool, error) {
	port := "2376"
	ncCommand := fmt.Sprintf("netstat -tuln | grep :%s", port)
	output, err := sshAdapter.RunCommand(ncCommand)
	if err != nil {
		return false, err
	}

	// If the output is empty, the port is closed
	if len(output) == 0 {
		return false, nil
	}

	// If the output is not empty, the port is open
	return true, nil
}

func installDocker(sshAdapter *adapter.SshAdapter) error {
	_, err := sshAdapter.RunCommand("sudo apt update")

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("sudo apt-get -y install ca-certificates curl gnupg")

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("sudo install -m 0755 -d /etc/apt/keyrings")

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg")

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("sudo chmod a+r /etc/apt/keyrings/docker.gpg")

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("echo \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable\" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null")

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("sudo apt-get update")

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("sudo apt-get -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin")

	if err != nil {
		return err
	}

	return nil
}

func setupDockerCertificates(sshAdapter *adapter.SshAdapter, server domain.Server) error {
	password := utils.GenerateRandomPassword(15)
	pathToCert := "/root/cert-docker"

	_, err := sshAdapter.RunCommand("mkdir -p " + pathToCert)

	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("openssl req -new -x509 -days 365 -keyout %[3]s/ca-key.pem -passout pass:%[2]s -sha256 -out %[3]s/ca.pem -subj \"/C=/ST=/L=/O=/CN=%[1]s\"", server.Domain, password, pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("openssl genrsa -out %s/server-key.pem 4096", pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("openssl req -subj \"/CN=%[1]s\" -sha256 -new -key %[2]s/server-key.pem -out %[2]s/server.csr", server.Domain, pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("echo subjectAltName = DNS:%[1]s,IP:10.10.10.20,IP:127.0.0.1 >> %[2]s/extfile.cnf", server.Domain, pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("echo extendedKeyUsage = serverAuth >> %s/extfile.cnf", pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("openssl x509 -req -days 365 -sha256 -in %[2]s/server.csr -CA %[2]s/ca.pem -CAkey %[2]s/ca-key.pem -CAcreateserial -out %[2]s/server-cert.pem -extfile %[2]s/extfile.cnf -passin pass:%[1]s", password, pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("openssl genrsa -out %s/key.pem 4096", pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("openssl req -subj '/CN=client' -new -key %[1]s/key.pem -out %[1]s/client.csr", pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("echo extendedKeyUsage = clientAuth > %s/extfile-client.cnf", pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("openssl x509 -req -days 365 -sha256 -in %[2]s/client.csr -CA %[2]s/ca.pem -CAkey %[2]s/ca-key.pem -CAcreateserial -out %[2]s/cert.pem -extfile %[2]s/extfile-client.cnf -passin pass:%[1]s", password, pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("rm -v %[1]s/client.csr %[1]s/server.csr %[1]s/extfile.cnf %[1]s/extfile-client.cnf", pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("chmod -v 0400 %[1]s/ca-key.pem %[1]s/key.pem %[1]s/server-key.pem", pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("chmod -v 0444 %[1]s/ca.pem %[1]s/server-cert.pem %[1]s/cert.pem", pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	return nil
}

func ensureCertDockerDirectory(localDir string) {
	// Check if the directory exists
	if _, err := os.Stat(localDir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(localDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create cert-docker directory: %v", err)
		}
		fmt.Println("Created cert-docker directory.")
	} else if err != nil {
		log.Fatalf("Error checking cert-docker directory: %v", err)
	} else {
		fmt.Println("Cert-docker directory already exists.")
	}
}

func copyCertificates(sshAdapter *adapter.SshAdapter, serverId string) error {
	remoteFiles := []string{"ca.pem", "key.pem", "cert.pem"}
	ensureCertDockerDirectory(internal.CERT_DOCKER_FOLDER)
	pathLocalCertDir := internal.CERT_DOCKER_FOLDER + "/" + serverId + "/"
	os.Mkdir(pathLocalCertDir, 0755)

	for _, remoteFileName := range remoteFiles {
		err := sshAdapter.SaveRemoteFileContentToLocalFile("/root/cert-docker/"+remoteFileName, pathLocalCertDir+remoteFileName)
		if err != nil {
			log.Fatalf("Error saving file content: %v", err)
		}

	}

	return nil
}

func openPortDockerConfig(sshAdapter *adapter.SshAdapter) error {
	pathToCert := "/root/cert-docker"

	_, err := sshAdapter.RunCommand("mkdir -p /etc/systemd/system/docker.service.d")

	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("echo \"[Service]\nExecStart=\nExecStart=/usr/bin/dockerd --tlsverify --tlscacert=%[1]s/ca.pem --tlscert=%[1]s/server-cert.pem --tlskey=%[1]s/server-key.pem -H fd:// -H=0.0.0.0:2376\" > /etc/systemd/system/docker.service.d/override.conf", pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("sudo systemctl daemon-reload")

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("sudo systemctl restart docker.service")

	if err != nil {
		return err
	}

	return nil
}
