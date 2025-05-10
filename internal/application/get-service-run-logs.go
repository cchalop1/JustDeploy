package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

// GetServiceRunLogs retrieves the run logs for a specific service
func GetServiceRunLogs(deployService *service.DeployService, serviceId string) ([]domain.Logs, error) {
	logsCollector := domain.LogCollector{}

	// Get the service from the database
	service, err := deployService.DatabaseAdapter.GetServiceById(serviceId)
	if err != nil {
		fmt.Printf("Error getting service with ID %s: %v\n", serviceId, err)
		return logsCollector.GetLogs(), err
	}

	// Get actual container logs
	dockerLogs, err := deployService.DockerAdapter.GetLogsOfContainer(service.GetDockerName())

	if err != nil {
		fmt.Printf("Error getting logs for container %s: %v\n", service.GetDockerName(), err)

		return logsCollector.GetLogs(), err
	}

	if logsCollector.Count() == 0 {
		logsCollector.AddLog(fmt.Sprintf("No logs available for %s", service.Name))
	}

	logsCollector.AddListLogs(dockerLogs)

	// If no logs were found, add a placeholder message

	return logsCollector.GetLogs(), nil
}
