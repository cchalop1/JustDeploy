package dto

type CreateProjectDto struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type CreateAppDto struct {
	Path      string `json:"path"`
	ProjectId string `json:"projectId"`
}

type PathDto struct {
	Name     string `json:"name"`
	FullPath string `json:"fullPath"`
	// Folders  []PathDto `json:"folders"`
}

type ProjectSettingsDto struct {
	CurrentPath       string    `json:"currentPath"`
	CurrentFolderName string    `json:"currentFolderName"`
	Folders           []PathDto `json:"folders"`
}

type DeployProjectDto struct {
	ProjectId   string  `json:"projectId"`
	ServerId    string  `json:"serverId"`
	Domain      *string `json:"domain"`
	IsTLSDomain bool    `json:"isTLSDomain"`
}
