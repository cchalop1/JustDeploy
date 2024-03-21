package dto

type EditDeployementDto struct {
	DeployOnCommit bool `json:"deployOnCommit"`
}

type Env struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type NewDeployDto struct {
	Name           string `json:"name"`
	ServerId       string `json:"serverId"`
	EnableTls      bool   `json:"enableTls"`
	Email          string `json:"email"`
	PathToSource   string `json:"pathToSource"`
	Envs           []Env  `json:"envs"`
	DeployOnCommit bool   `json:"deployOnCommit"`
}

type DeployDto struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Server       ServerDto `json:"server"`
	EnableTls    bool      `json:"enableTls"`
	Email        string    `json:"email"`
	PathToSource string    `json:"pathToSource"`
	Envs         []Env     `json:"envs"`
	Status       string    `json:"status"`
	Url          string    `json:"url"`
}

type DeployConfigDto struct {
	SourceType       string `json:"sourceType"`
	PathToSource     string `json:"pathToSource"`
	DockerFileFound  bool   `json:"dockerFileFound"`
	ComposeFileFound bool   `json:"composeFileFound"`
	EnvFileFound     bool   `json:"envFileFound"`
	Envs             []Env  `json:"envs"`
}
