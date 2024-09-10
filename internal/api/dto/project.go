package dto

type CreateProjectDto struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type CreateAppDto struct {
	Path      string `json:"path"`
	ProjectId string `json:"projectId"`
}
