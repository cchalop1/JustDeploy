package application

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

// GetServiceBuildLogs retrieves the build logs for a specific service
func GetServiceBuildLogs(deployService *service.DeployService, serviceId string) ([]domain.Logs, error) {
	logs := []domain.Logs{}

	// Read logs from file
	logFile := filepath.Join(internal.JUSTDEPLOY_FOLDER, "build-logs", serviceId+".json")

	if _, err := os.Stat(logFile); err != nil {
		// If file doesn't exist, return empty logs
		return logs, nil
	}

	// Read and parse the log file
	data, err := os.ReadFile(logFile)
	if err != nil {
		fmt.Printf("Error reading log file for service %s: %v\n", serviceId, err)
		return logs, err
	}

	if err := json.Unmarshal(data, &logs); err != nil {
		fmt.Printf("Error parsing log file for service %s: %v\n", serviceId, err)
		return logs, err
	}

	return logs, nil
}
