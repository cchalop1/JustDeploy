package application

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/domain"
)

// StoreBuildLog stores build log entries for a service
func StoreBuildLog(serviceId string, logCollector domain.LogCollector) error {
	// Create logs directory if it doesn't exist
	logsDir := filepath.Join(internal.JUSTDEPLOY_FOLDER, "build-logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("error creating logs directory: %w", err)
	}

	// Read existing logs
	existingLogs := []domain.Logs{}
	logFile := filepath.Join(logsDir, serviceId+".json")

	if _, err := os.Stat(logFile); err == nil {
		// File exists, read it
		data, err := os.ReadFile(logFile)
		if err != nil {
			return fmt.Errorf("error reading log file: %w", err)
		}
		if err := json.Unmarshal(data, &existingLogs); err != nil {
			return fmt.Errorf("error unmarshaling logs: %w", err)
		}
	}

	// Append new log entries from the collector
	for _, log := range logCollector.GetLogs() {
		existingLogs = append(existingLogs, domain.Logs{
			Date:    log.Date,
			Message: log.Message,
		})
	}

	// Write back to file
	data, err := json.Marshal(existingLogs)
	if err != nil {
		return fmt.Errorf("error marshaling logs: %w", err)
	}

	if err := os.WriteFile(logFile, data, 0644); err != nil {
		return fmt.Errorf("error writing log file: %w", err)
	}

	return nil
}
