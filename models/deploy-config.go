package models

type DeployConfigDto struct {
	PathToProject   string           `json:"pathToProject"`
	DockerFileValid bool             `json:"dockerFileValid"`
	ServerConfig    ConnectServerDto `json:"serverConfig"`
	AppConfig       AppConfigDto     `json:"appConfig"`
	DeployStatus    string           `json:"deployStatus"`
	Url             string           `json:"url"`
	AppStatus       string           `json:"appStatus"`
}

type ConnectServerDto struct {
	Domain   string  `json:"domain"`
	SshKey   *string `json:"sshKey"`
	Password *string `json:"password"`
	User     string  `json:"user"`
}

type Env struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type AppConfigDto struct {
	Name         string `json:"name"`
	EnableTls    bool   `json:"enableTls"`
	Email        string `json:"email"`
	PathToSource string `json:"pathToSource"`
	Envs         []Env  `json:"envs"`
}
