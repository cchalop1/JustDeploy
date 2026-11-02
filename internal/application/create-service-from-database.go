package application

import (
	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
	"strconv"
	"strings"
)

func CreateServiceFromDatabase(ds *service.DeployService, name string) (domain.Service, error) {
	if !ds.Operation.TryLock() {
		return domain.Service{}, domain.ErrConflict
	}
	defer ds.Operation.Unlock()
	config, err := database.GetServiceByName(name)
	if err != nil {
		return domain.Service{}, err
	}
	id := utils.GenerateRandomPassword(16)
	envs := utils.GenerateEnvValues(config.Env)
	s := domain.Service{Id: id, Name: strings.ToLower(name) + "-" + id, Type: config.Type, ImageName: config.Config.Image, ImageUrl: config.Icon, Status: "ready_to_deploy", Envs: envs, Cmd: utils.ReplaceEnvVariablesInCmd(append([]string(nil), config.Config.Cmd...), envs), DockerHubUrl: utils.GetDockerHubUrl(config.Config.Image), ExposeSettings: domain.ServiceExposeSettings{ExposePort: strconv.Itoa(config.DefaultPort)}, ServiceRuntime: domain.ServiceRuntime{ProjectID: id}}
	if err := ValidateService(s); err != nil {
		return s, err
	}
	if err := ds.DockerAdapter.EnsureImage(s.ImageName); err != nil {
		return s, err
	}
	// Credentials are returned to the administrator, never copied to unrelated services.
	return s, ds.DatabaseAdapter.SaveService(s)
}
