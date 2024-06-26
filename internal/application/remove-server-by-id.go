package application

import (
	"errors"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func RemoveServerById(deployService *service.DeployService, serverId string) error {
	server, err := deployService.DatabaseAdapter.GetServerById(serverId)
	if err != nil {
		return err
	}

	deployList := deployService.DatabaseAdapter.GetDeployByServerId(serverId)

	if len(deployList) > 0 {
		return errors.New("you can't remove server with application on it")
	}

	sshAdapter := adapter.NewSshAdapter()

	sshAdapter.Connect(dto.ConnectNewServerDto{
		Ip:       server.Ip,
		SshKey:   server.SshKey,
		Password: server.Password,
		User:     "root",
	})

	removeCertOnServer(sshAdapter)

	deployService.FilesystemAdapter.RemoveDockerCertificateByServerId(server.Id)

	err = deployService.DatabaseAdapter.DeleteServer(server)

	if err != nil {
		return err
	}

	return nil
}

func removeCertOnServer(sshAdapter *adapter.SshAdapter) error {
	_, err := sshAdapter.RunCommand("rm -rf /root/docker-cert")

	if err != nil {
		return err
	}

	_, err = sshAdapter.RunCommand("mkdir -p /etc/systemd/system/docker.service.d")

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
