package dto

type CreateServiceDto struct {
	ServiceName       string  `json:"serviceName"`
	FromDockerCompose bool    `json:"fromDockerCompose"`
	DeployId          *string `json:"deployId"`
	ProjectId         *string `json:"projectId"`
	Path              *string `json:"path"`
}
