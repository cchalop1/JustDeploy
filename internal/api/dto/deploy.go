package dto

type DeployConfigDto struct {
	ServerConfig    ConnectServerDto `json:"serverConfig"`
	AppConfig       AppConfigDto     `json:"appConfig"`
	PathToProject   string           `json:"pathToProject"`
	DockerFileValid bool             `json:"dockerFileValid"`
	DeployStatus    string           `json:"deployStatus"`
	Url             string           `json:"url"`
	AppStatus       string           `json:"appStatus"`
}
