package dto

type ServiceDto struct {
	Name        string   `json:"name"`
	Icon        string   `json:"icon"`
	Image       string   `json:"image"`
	Port        string   `json:"port"`
	Envs        []string `json:"envs"`
	Secrets     []string `json:"secrets"`
	VolumsNames []string `json:"volumsNames"`
}

type CreateServiceDto struct {
	ServiceName       string `json:"serviceName"`
	FromDockerCompose bool   `json:"fromDockerCompose"`
}
