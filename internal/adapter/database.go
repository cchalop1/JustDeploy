package adapter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/api/dto"
)

type DatabaseAdapter struct {
}

func NewDatabaseAdapter() *DatabaseAdapter {
	return &DatabaseAdapter{}
}

func (d *DatabaseAdapter) databaseFileIsCreated() bool {
	_, err := os.Stat(internal.DATABASE_FILE_PATH)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	log.Printf("Error checking if file exists: %v", err)
	return false
}

func (d *DatabaseAdapter) createFoldeJustDeployFolderIfDontExist() error {
	err := os.MkdirAll(internal.JUSTDEPLOY_FOLDER, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	fmt.Println("Directory created successfully")
	return nil
}

func (d *DatabaseAdapter) writeDeployConfigInDataBaseFile(deployConfig dto.DeployConfigDto) error {
	file, err := json.MarshalIndent(deployConfig, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	err = os.WriteFile(internal.DATABASE_FILE_PATH, file, 0644)
	fmt.Println(err)
	return err
}

func (d *DatabaseAdapter) readDeployConfigInDataBaseFile() dto.DeployConfigDto {
	deployConfig := dto.DeployConfigDto{}

	file, err := os.ReadFile(internal.DATABASE_FILE_PATH)

	if err != nil {
		log.Fatalf("Error for read the database file: %v", err)
	}

	err = json.Unmarshal(file, &deployConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	return deployConfig
}

func (d *DatabaseAdapter) GetState() dto.DeployConfigDto {
	if !d.databaseFileIsCreated() {
		d.createFoldeJustDeployFolderIfDontExist()
		d.writeDeployConfigInDataBaseFile(dto.DeployConfigDto{
			PathToProject:   "",
			DockerFileValid: false,
			DeployStatus:    "serverconfig",
			ServerConfig:    dto.ConnectServerDto{},
			AppConfig:       dto.AppConfigDto{},
			AppStatus:       "",
		})
	}
	deployConfig := d.readDeployConfigInDataBaseFile()
	return deployConfig
}

func (d *DatabaseAdapter) SaveState(deployConfig dto.DeployConfigDto) error {
	d.writeDeployConfigInDataBaseFile(deployConfig)
	return nil
}
