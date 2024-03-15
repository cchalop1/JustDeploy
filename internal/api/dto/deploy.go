package dto

type DeployConfigDto struct {
	ServerConfig    ConnectServerDto `json:"serverConfig"`
	AppConfig       AppConfigDto     `json:"appConfig"`
	DockerFileValid bool             `json:"dockerFileValid"`
	// TODO: add enum here
	DeployStatus string `json:"deployStatus"`
	Url          string `json:"url"`
	// TODO: add enum here
	AppStatus string `json:"appStatus"`
}
