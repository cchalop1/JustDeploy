package adapter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/domain"
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

func (d *DatabaseAdapter) writeDeployConfigInDataBaseFile(databaseModels domain.DatabaseModelsType) error {
	file, err := json.MarshalIndent(databaseModels, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	err = os.WriteFile(internal.DATABASE_FILE_PATH, file, 0644)
	fmt.Println(err)
	return err
}

func (d *DatabaseAdapter) readDeployConfigInDataBaseFile() domain.DatabaseModelsType {
	databaseModels := domain.DatabaseModelsType{}

	file, err := os.ReadFile(internal.DATABASE_FILE_PATH)

	if err != nil {
		log.Fatalf("Error for read the database file: %v", err)
	}

	err = json.Unmarshal(file, &databaseModels)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	return databaseModels
}

func (d *DatabaseAdapter) Init() {
	if !d.databaseFileIsCreated() {
		d.createFoldeJustDeployFolderIfDontExist()
		databaseData := domain.DatabaseModelsType{
			Servers: []domain.Server{},
			Deploys: []domain.Deploy{},
		}
		d.writeDeployConfigInDataBaseFile(databaseData)
	}
}

// func (d *DatabaseAdapter) GetState() dto.DeployConfigDto {
// 	deployConfig := d.readDeployConfigInDataBaseFile()
// 	return deployConfig
// }

// func (d *DatabaseAdapter) SaveState(deployConfig dto.DeployConfigDto) error {
// 	d.writeDeployConfigInDataBaseFile(deployConfig)
// 	return nil
// }

func (d *DatabaseAdapter) SaveServer(newServer domain.Server) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	databaseModels.Servers = append(databaseModels.Servers, newServer)
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) GetServers() []domain.Server {
	databaseModels := d.readDeployConfigInDataBaseFile()
	return databaseModels.Servers
}

func (d *DatabaseAdapter) SaveDeploy(deploy domain.Deploy) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	databaseModels.Deploys = append(databaseModels.Deploys, deploy)
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) GetDeploys() []domain.Deploy {
	databaseModels := d.readDeployConfigInDataBaseFile()
	return databaseModels.Deploys
}
