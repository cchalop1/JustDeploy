package domain

type Project struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	ServerId *string   `json:"serverId"`
	Services []Service `json:"services"`
}
