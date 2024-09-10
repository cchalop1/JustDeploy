package domain

type Project struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Services []Service `json:"services"`
	Apps     []App     `json:"apps"`
}

type App struct {
	Id              string `json:"id"`
	Path            string `json:"path"`
	Name            string `json:"name"`
	IsDockerFile    bool   `json:"isDockerFile"`
	IsDockerCompose bool   `json:"isDockerCompose"`
}
