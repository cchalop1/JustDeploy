package application

import "cchalop1.com/deploy/internal/adapter"

func GetApplicationLogs(dockerAdapter *adapter.DockerAdapter, containerName string) []string {
	return dockerAdapter.GetLogsOfContainer(containerName)
}
