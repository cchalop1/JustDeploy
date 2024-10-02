package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/api/dto"
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
			Servers:  []domain.Server{},
			Deploys:  []domain.Deploy{},
			Services: []domain.Service{},
			Projects: []domain.Project{}, // Add this line
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

	for idx, existingServer := range databaseModels.Servers {
		if existingServer.Id == newServer.Id {
			databaseModels.Servers[idx] = newServer
			return d.writeDeployConfigInDataBaseFile(databaseModels)
		}
	}

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

func (d *DatabaseAdapter) GetServicesByDeployId(deployId string) []domain.Service {
	// databaseModels := d.readDeployConfigInDataBaseFile()
	// serviceList := []domain.Service{}
	// for _, s := range databaseModels.Services {
	// 	if *s.DeployId == deployId {
	// 		serviceList = append(serviceList, s)
	// 	}
	// }
	// return serviceList
	return []domain.Service{}
}

func (d *DatabaseAdapter) GetLocalService() []domain.Service {
	// databaseModels := d.readDeployConfigInDataBaseFile()
	// serviceList := []domain.Service{}
	// for _, s := range databaseModels.Services {
	// 	if s.DeployId == nil {
	// 		serviceList = append(serviceList, s)
	// 	}
	// }
	// return serviceList
	return []domain.Service{}
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
	var newServices []domain.Service
	for _, s := range databaseModels.Services {
		if s.Id != id {
			newServices = append(newServices, s)
		}
	}

	// TODO: move to other methode
	for i, p := range databaseModels.Projects {
		for j, s := range p.Services {
			if s.Id == id {
				databaseModels.Projects[i].Services = append(databaseModels.Projects[i].Services[:j], databaseModels.Projects[i].Services[j+1:]...)
				break
			}
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

func (d *DatabaseAdapter) SaveLogs(deployId string, logs []dto.Logs) error {
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

func (d *DatabaseAdapter) GetLogs(deployId string) ([]dto.Logs, error) {
	file, err := os.ReadFile(internal.JUSTDEPLOY_FOLDER + "/" + deployId + ".log")
	if err != nil {
		return []dto.Logs{}, err
	}

	var logs []dto.Logs
	err = json.Unmarshal(file, &logs)
	if err != nil {
		return []dto.Logs{}, err
	}
	return logs, nil
}

func (d *DatabaseAdapter) DeleteLogFile(deployId string) error {
	return os.Remove(internal.JUSTDEPLOY_FOLDER + "/" + deployId + ".log")
}

// refactor with new database
func (d *DatabaseAdapter) GetServerFromService(serviceId string) (domain.Server, error) {
	// databaseModels := d.readDeployConfigInDataBaseFile()
	// for _, s := range databaseModels.Services {
	// 	if s.Id == serviceId {
	// 		if s.DeployId == nil {
	// 			return domain.Server{}, errors.New("deploy not found")
	// 		}
	// 		if s.DeployId == nil {
	// 			// return Localhost server
	// 			return domain.Server{
	// 				Id:     "localhost",
	// 				Name:   "localhost",
	// 				Ip:     "localhost",
	// 				Domain: "localhost",
	// 			}, nil
	// 		} else {
	// 			deploy, err := d.GetDeployById(*s.DeployId)
	// 			if err != nil {
	// 				return domain.Server{}, errors.New("deploy not found")
	// 			}
	// 			server, err := d.GetServerById(deploy.ServerId)

	// 			if err != nil {
	// 				return domain.Server{}, errors.New("server not found")
	// 			}
	// 			return server, nil

	// 		}

	// 	}
	// }
	// return domain.Server{}, errors.New("service not found")
	return domain.Server{}, nil
}

// Project

func (d *DatabaseAdapter) SaveProject(project domain.Project) error {
	databaseModels := d.readDeployConfigInDataBaseFile()

	// Check if the project already exists
	for i, existingProject := range databaseModels.Projects {
		if existingProject.Id == project.Id {
			// Update existing project
			log.Printf("Updating existing project with ID: %s", project.Id)
			databaseModels.Projects[i] = project
			return d.writeDeployConfigInDataBaseFile(databaseModels)
		}
	}

	// Add new project
	log.Printf("Adding new project with ID: %s", project.Id)
	databaseModels.Projects = append(databaseModels.Projects, project)
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

// Project
func (d *DatabaseAdapter) GetProjectById(id string) (*domain.Project, error) {
	databaseModels := d.readDeployConfigInDataBaseFile()
	for _, project := range databaseModels.Projects {
		if project.Id == id {
			return &project, nil
		}
	}
	return nil, errors.New("project not found")
}

func (d *DatabaseAdapter) GetProjectByPath(path string) *domain.Project {
	databaseModels := d.readDeployConfigInDataBaseFile()

	for _, project := range databaseModels.Projects {
		if project.Path == path {
			return &project
		}
	}
	return nil
}

func (d *DatabaseAdapter) GetServicesByProjectId(projectId string) []domain.Service {
	// databaseModels := d.readDeployConfigInDataBaseFile()
	// serviceList := []domain.Service{}
	// for _, s := range databaseModels.Services {
	// 	if s.ProjectId != nil && *s.ProjectId == projectId {
	// 		serviceList = append(serviceList, s)
	// 	}
	// }
	// return serviceList
	return []domain.Service{}
}

func (d *DatabaseAdapter) SaveServiceByProjectId(projectId string, service domain.Service) error {
	databaseModels := d.readDeployConfigInDataBaseFile()
	for i, p := range databaseModels.Projects {
		if p.Id == projectId {
			databaseModels.Projects[i].Services = append(databaseModels.Projects[i].Services, service)
			break
		}
	}
	return d.writeDeployConfigInDataBaseFile(databaseModels)
}

func (d *DatabaseAdapter) GetProjects() []domain.Project {
	databaseModels := d.readDeployConfigInDataBaseFile()
	return databaseModels.Projects
}
