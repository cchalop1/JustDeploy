package domain

import (
	"strings"

	"cchalop1.com/deploy/internal/api/dto"
)

type ServiceExposeSettings struct {
	IsExposed bool   `json:"isExposed"`
	SubDomain string `json:"subDomain"`
}

type Service struct {
	Id          string    `json:"id"`
	Type        string    `json:"type"`
	Url         string    `json:"url"`
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

	ExposeSettings ServiceExposeSettings `json:"exposeSettings"`
}

func (s *Service) GetDockerName() string {
	return strings.ToLower(s.Name + "-" + s.Id)
}
