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

// Server
func (d *DatabaseAdapter) GetServerById(id string) (domain.Server, error) {
	databaseModels := d.readDeployConfigInDataBaseFile()
	for _, s := range databaseModels.Servers {
		if s.Id == id {
			return s, nil
		}
	}
	return domain.Server{}, errors.New("server not found")
}

func (d *DatabaseAdapter) UpdateServer(server domain.Server) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	for i, s := range databaseModels.Servers {
		if s.Name == server.Name {
			databaseModels.Servers[i] = server
			break
		}
	}
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) DeleteServer(server domain.Server) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	var newServer []domain.Server
	for _, s := range databaseModels.Servers {
		if s.Id != server.Id {
			newServer = append(newServer, s)
		}
	}
	databaseModels.Servers = newServer
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) CountServer() int {
	// TODO: change to get it from a index count
	databaseModels := d.readDeployConfigInDataBaseFile()
	return len(databaseModels.Servers)
}

func (d *DatabaseAdapter) SaveServer(newServer domain.Server) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	databaseModels.Servers = append(databaseModels.Servers, newServer)
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) GetServers() []domain.Server {
	databaseModels := d.readDeployConfigInDataBaseFile()
	return databaseModels.Servers
}

// Deploy

func (d *DatabaseAdapter) SaveDeploy(deploy domain.Deploy) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	databaseModels.Deploys = append(databaseModels.Deploys, deploy)
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) GetDeploys() []domain.Deploy {
	databaseModels := d.readDeployConfigInDataBaseFile()
	return databaseModels.Deploys
}

func (d *DatabaseAdapter) GetDeployByServerId(serverId string) []domain.Deploy {
	databaseModels := d.readDeployConfigInDataBaseFile()
	deployList := []domain.Deploy{}
	for _, d := range databaseModels.Deploys {
		if d.ServerId == serverId {
			deployList = append(deployList, d)
		}
	}
	return deployList
}

func (d *DatabaseAdapter) GetDeployById(id string) (domain.Deploy, error) {
	databaseModels := d.readDeployConfigInDataBaseFile()
	for _, d := range databaseModels.Deploys {
		if d.Id == id {
			return d, nil
		}
	}
	return domain.Deploy{}, errors.New("deploy not found")
}

func (d *DatabaseAdapter) UpdateDeploy(deploy domain.Deploy) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	for i, s := range databaseModels.Deploys {
		if s.Id == deploy.Id {
			databaseModels.Deploys[i] = deploy
			break
		}
	}
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) DeleteDeploy(deploy domain.Deploy) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	var newDeploys []domain.Deploy
	for _, s := range databaseModels.Deploys {
		if s.Id != deploy.Id {
			newDeploys = append(newDeploys, s)
		}
	}
	databaseModels.Deploys = newDeploys
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

// Services

func (d *DatabaseAdapter) SaveService(service domain.Service) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	databaseModels.Services = append(databaseModels.Services, service)
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) GetServiceByDeployId(deployId string) []domain.Service {
	databaseModels := d.readDeployConfigInDataBaseFile()
	serviceList := []domain.Service{}
	for _, s := range databaseModels.Services {
		if s.DeployId == deployId {
			serviceList = append(serviceList, s)
		}
	}
	return serviceList
}

func (d *DatabaseAdapter) GetServiceById(id string) (domain.Service, error) {
	databaseModels := d.readDeployConfigInDataBaseFile()
	for _, s := range databaseModels.Services {
		if s.Id == id {
			return s, nil
		}
	}
	return domain.Service{}, errors.New("service not found")
}

func (d *DatabaseAdapter) DeleteServiceById(id string) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	var newServices []domain.Service
	for _, s := range databaseModels.Services {
		if s.Id != id {
			newServices = append(newServices, s)
		}
	}
	databaseModels.Services = newServices
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}
