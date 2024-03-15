package dto

type Env struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type AppConfigDto struct {
	Name           string `json:"name"`
	EnableTls      bool   `json:"enableTls"`
	Email          string `json:"email"`
	PathToSource   string `json:"pathToSource"`
	Envs           []Env  `json:"envs"`
	DeployOnCommit bool   `json:"deployOnCommit"`
}
