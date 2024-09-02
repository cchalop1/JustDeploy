package dto

// deprecated to delete:
type ServiceDto struct {
	Name        string   `json:"name"`
	Icon        string   `json:"icon"`
	Image       string   `json:"image"`
	Ports       []string `json:"ports"`
	Envs        []string `json:"envs"`
	Secrets     []string `json:"secrets"`
	VolumsNames []string `json:"volumsNames"`
}

type CreateServiceDto struct {
	ServiceName       string  `json:"serviceName"`
	FromDockerCompose bool    `json:"fromDockerCompose"`
	DeployId          *string `json:"deployId"`
}
