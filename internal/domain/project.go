package domain

type Project struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Services []Service `json:"services"`
	Deploy   []Deploy  `json:"deploy"`
}
