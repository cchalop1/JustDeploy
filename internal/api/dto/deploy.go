package dto

// type DeployConfigDto struct {
// 	ServerConfig    ConnectNewServerDto `json:"serverConfig"`
// 	AppConfig       NewDeployDto        `json:"appConfig"`
// 	DockerFileValid bool                `json:"dockerFileValid"`
// 	// TODO: add enum here
// 	DeployStatus string `json:"deployStatus"`
// 	Url          string `json:"url"`
// 	// TODO: add enum here
// 	AppStatus string `json:"appStatus"`
// }

type EditDeployementDto struct {
	DeployOnCommit bool `json:"deployOnCommit"`
}
