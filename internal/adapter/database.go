package adapter

import (
	"encoding/json"
	"errors"
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
	fileContent, err := json.MarshalIndent(databaseModels, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	err = os.WriteFile(internal.DATABASE_FILE_PATH, fileContent, 0644)
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
			Server:   domain.Server{},
			Services: []domain.Service{},
			Settings: domain.Settings{},
		}
		d.writeDeployConfigInDataBaseFile(databaseData)
	}
}

func (d *DatabaseAdapter) DeleteServer(server domain.Server) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	databaseModels.Server = domain.Server{}
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) SaveServer(newServer domain.Server) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	databaseModels.Server = newServer

	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) GetServer() domain.Server {
	databaseModels := d.readDeployConfigInDataBaseFile()
	return databaseModels.Server
}

// Services
func (d *DatabaseAdapter) SaveService(service domain.Service) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	serviceExists := false

	for i, s := range databaseModels.Services {
		if s.Id == service.Id {
			databaseModels.Services[i] = service
			serviceExists = true
			break
		}
	}

	if !serviceExists {
		databaseModels.Services = append(databaseModels.Services, service)
	}

	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) GetServices() []domain.Service {
	databaseModels := d.readDeployConfigInDataBaseFile()
	return databaseModels.Services
}

func (d *DatabaseAdapter) GetServiceById(id string) (*domain.Service, error) {
	databaseModels := d.readDeployConfigInDataBaseFile()
	for _, s := range databaseModels.Services {
		if s.Id == id {
			return &s, nil
		}
	}
	return &domain.Service{}, errors.New("service not found")
}

func (d *DatabaseAdapter) DeleteServiceById(id string) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	var newServices = []domain.Service{}
	for _, s := range databaseModels.Services {
		if s.Id != id {
			newServices = append(newServices, s)
		}
	}

	databaseModels.Services = newServices
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

// Logs
func (d *DatabaseAdapter) CreateDeployLogsFileIfNot(deployId string) error {
	filePath := internal.JUSTDEPLOY_FOLDER + "/" + deployId + ".log"
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}
	_, err := os.Create(filePath)
	return err
}

func (d *DatabaseAdapter) SaveLogs(deployId string, logs []domain.Logs) error {
	file, err := os.OpenFile(internal.JUSTDEPLOY_FOLDER+"/"+deployId+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fileContent, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	file.Truncate(0)
	_, err = file.Write(fileContent)
	return err
}

func (d *DatabaseAdapter) GetLogs(deployId string) ([]domain.Logs, error) {
	file, err := os.ReadFile(internal.JUSTDEPLOY_FOLDER + "/" + deployId + ".log")
	if err != nil {
		return []domain.Logs{}, err
	}

	var logs []domain.Logs
	err = json.Unmarshal(file, &logs)
	if err != nil {
		return []domain.Logs{}, err
	}
	return logs, nil
}

func (d *DatabaseAdapter) DeleteLogFile(deployId string) error {
	return os.Remove(internal.JUSTDEPLOY_FOLDER + "/" + deployId + ".log")
}

// Setting
func (d *DatabaseAdapter) GetSettings() domain.Settings {
	databaseModels := d.readDeployConfigInDataBaseFile()
	return databaseModels.Settings
}

func (d *DatabaseAdapter) SaveSettings(settings domain.Settings) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	databaseModels.Settings = settings
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) SaveInstallationToken(installationId string, token string) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	databaseModels.Settings.GithubToken = token
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}
