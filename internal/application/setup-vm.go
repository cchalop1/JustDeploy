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
	log.Println("Starting ConnectAndSetupServer function")
	sshAdapter := adapter.NewSshAdapter()

	// Connect to the server by ssh
	log.Println("Connecting to the server via SSH")
	sshAdapter.Connect(dto.ConnectNewServerDto{
		Ip:       server.Ip,
		SshKey:   server.SshKey,
		Password: server.Password,
		User:     "root",
	})

	// Install packages and docker

	log.Println("Checking if Docker is installed")
	dockerIsInstalled := checkIfDockerIsIntalled(sshAdapter)
	fmt.Println("dockerIsInstalled", dockerIsInstalled)

	var err error

	if !dockerIsInstalled {
		log.Println("Docker not installed. Installing Docker...")
		err = installDocker(sshAdapter)
		if err != nil {
			log.Printf("Error installing Docker: %v", err)
			return nil
		}
		log.Println("Docker installation completed")
	} else {
		log.Println("Docker is already installed")
	}

	log.Println("Setting up Docker certificates")
	err = setupDockerCertificates(sshAdapter, server)
	if err != nil {
		log.Printf("Error setting up Docker certificates: %v", err)
		return nil
	}
	log.Println("Docker certificates setup completed")

	log.Println("Copying certificates")
	copyCertificates(sshAdapter, server.Id)

	log.Println("Opening Docker port")
	err = openPortDockerConfig(sshAdapter)
	if err != nil {
		log.Printf("Error opening Docker port: %v", err)
		return nil
	}
	log.Println("Docker port configuration completed")

	log.Println("Closing SSH connection")
	sshAdapter.CloseConnection()

	log.Println("Creating new Docker adapter")
	adapterDocker := adapter.NewDockerAdapter()

	if err != nil {
		log.Printf("Error connecting to Docker client: %v", err)
		return nil
	}
	log.Println("Successfully connected to Docker client")

	log.Println("Pull the traefik router")

	err = adapterDocker.PullTreafikImage()

	if err != nil {
		log.Printf("Error pulling the traefik router: %v", err)
		return nil
	}

	log.Println("Updating server status")
	server.Status = "Runing"
	deployService.DatabaseAdapter.SaveServer(server)

	return nil
}

func checkIfDockerIsIntalled(sshAdapter *adapter.SshAdapter) bool {
	output, err := sshAdapter.RunCommand("docker --version")
	if err != nil {
		return false
	}

	if len(output) > 0 {
		return true
	}

	return false
}

func checkIsCertificateIsCreate(sshAdapter *adapter.SshAdapter, serverId string) bool {
	certificatePath := "/root/cert-docker/" + serverId
	statCommand := fmt.Sprintf("stat -t -- \"%s\" &>/dev/null", certificatePath)
	output, err := sshAdapter.RunCommand(statCommand)
	if err != nil {
		return false
	}

	// If the output is empty, the certificate file exists
	if len(output) == 0 {
		return true
	}

	// If the output is not empty, the certificate file does not exist
	return false
}

func checkIfDockerPortIsOpen(sshAdapter *adapter.SshAdapter) bool {
	port := "2376"
	ncCommand := fmt.Sprintf("netstat -tuln | grep :%s", port)
	output, err := sshAdapter.RunCommand(ncCommand)
	if err != nil {
		return false
	}

	// If the output is empty, the port is closed
	if len(output) == 0 {
		return false
	}

	// If the output is not empty, the port is open
	return true
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

	// TODO: change to set the ip
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

	cmd = fmt.Sprintf("openssl req -subj \"/CN=%[1]s\" -sha256 -new -key %[2]s/server-key.pem -out %[2]s/server.csr", server.Ip, pathToCert)
	_, err = sshAdapter.RunCommand(cmd)

	if err != nil {
		return err
	}

	// TODO: check if is a domain or a ip
	cmd = fmt.Sprintf("echo subjectAltName = IP:%[1]s,IP:10.10.10.20,IP:127.0.0.1 >> %[2]s/extfile.cnf", server.Ip, pathToCert)
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
