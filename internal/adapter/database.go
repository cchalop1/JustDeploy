package adapter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/models"
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

func (d *DatabaseAdapter) writeDeployConfigInDataBaseFile(deployConfig models.DeployConfigDto) error {
	file, err := json.MarshalIndent(deployConfig, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	err = os.WriteFile(internal.DATABASE_FILE_PATH, file, 0644)
	fmt.Println(err)
	return err
}

func (d *DatabaseAdapter) readDeployConfigInDataBaseFile() models.DeployConfigDto {
	deployConfig := models.DeployConfigDto{}

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

func (d *DatabaseAdapter) GetState() models.DeployConfigDto {
	if !d.databaseFileIsCreated() {
		d.createFoldeJustDeployFolderIfDontExist()
		d.writeDeployConfigInDataBaseFile(models.DeployConfigDto{
			PathToProject:   "",
			DockerFileValid: false,
			DeployStatus:    "serverconfig",
			ServerConfig:    models.ConnectServerDto{},
			AppConfig:       models.AppConfigDto{},
			AppStatus:       "",
		})
	}
	deployConfig := d.readDeployConfigInDataBaseFile()
	return deployConfig
}

func (d *DatabaseAdapter) SaveState(deployConfig models.DeployConfigDto) error {
	d.writeDeployConfigInDataBaseFile(deployConfig)
	return nil
}
