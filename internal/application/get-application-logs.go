package application

import "cchalop1.com/deploy/internal/adapter"

func GetApplicationLogs(containerName string, dockerAdapter *adapter.DockerAdapter) []string {
	return dockerAdapter.GetLogsOfContainer(containerName)
}
