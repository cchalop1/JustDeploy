package dto

type DeployConfigDto struct {
	ServerConfig    ConnectServerDto `json:"serverConfig"`
	AppConfig       AppConfigDto     `json:"appConfig"`
	PathToProject   string           `json:"pathToProject"`
	DockerFileValid bool             `json:"dockerFileValid"`
	// TODO: add enum here
	DeployStatus string `json:"deployStatus"`
	Url          string `json:"url"`
	// TODO: add enum here
	AppStatus string `json:"appStatus"`
}
