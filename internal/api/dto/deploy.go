package dto

type EditDeployDto struct {
	Id             string `json:"id"`
	DeployOnCommit bool   `json:"deployOnCommit"`
	Envs           []Env  `json:"envs"`
	SubDomain      string `json:"subDomain"`
}

type Env struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
}

func (e Env) Equals(env Env) bool {
	return e.Name == env.Name && e.Value == env.Value && e.IsSecret == env.IsSecret
}

type NewDeployDto struct {
	ServerId       string `json:"serverId" validate:"required"`
	EnableTls      bool   `json:"enableTls"`
	Email          string `json:"email" validate:"omitempty,email"`
	PathToSource   string `json:"pathToSource" validate:"required"`
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

type ParamsDeployConfigDto struct {
	DeployId string `json:"deployId"`
	Path     string `json:"path"`
}

type DeployConfigDto struct {
	DeployName       string `json:"deployName"`
	SourceType       string `json:"sourceType"`
	PathToSource     string `json:"pathToSource"`
	DockerFileFound  bool   `json:"dockerFileFound"`
	ComposeFileFound bool   `json:"composeFileFound"`
	EnvFileFound     bool   `json:"envFileFound"`
	Envs             []Env  `json:"envs"`
}

type CreateDeployResponse struct {
	Id string `json:"id"`
}

type Logs struct {
	Message string `json:"message"`
	Date    string `json:"date"`
}

type NewDomain struct {
	Domain string `json:"domain"`
}
