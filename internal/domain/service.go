package domain

import (
	"strings"

	"cchalop1.com/deploy/internal/api/dto"
)

// TODO: change to save serviceConfig
type Service struct {
	Id          string    `json:"id"`
	HostName    string    `json:"hostName"`
	Name        string    `json:"name"`
	Envs        []dto.Env `json:"envs"`
	VolumsNames []string  `json:"volumsNames"`
	Status      string    `json:"status"`
	Host        string    `json:"host"`
	ImageName   string    `json:"imageName"`
	ImageUrl    string    `json:"imageUrl"`
	// TODO: rethink this
	IsDevContainer bool   `json:"isDevContainer"`
	CurrentPath    string `json:"currentPath"`
	ExposePort     string `json:"exposePort"`
}

func (s *Service) GetDockerName() string {
	return strings.ToLower(s.HostName + "-" + s.Id)
}
