package domain

import (
	"strings"

	"cchalop1.com/deploy/internal/api/dto"
)

type ServiceExposeSettings struct {
	IsExposed  bool   `json:"isExposed"`
	SubDomain  string `json:"subDomain"`
	Tls        bool   `json:"tls"`
	ExposePort string `json:"exposePort"`
}

type Service struct {
	Id       string    `json:"id"`
	Type     string    `json:"type"`
	Url      string    `json:"url"`
	Name     string    `json:"name"`
	FullName string    `json:"full_name"`
	Envs     []dto.Env `json:"envs"`
	// VolumsNames []string  `json:"volumsNames"`
	Status string `json:"status"`
	// Host      string `json:"host"`
	ImageName string `json:"imageName"`
	ImageUrl  string `json:"imageUrl"`

	CurrentPath  string `json:"currentPath"`
	DockerHubUrl string `json:"dockerHubUrl"`

	ExposeSettings ServiceExposeSettings `json:"exposeSettings"`
}

func (s *Service) GetDockerName() string {
	return strings.ToLower(s.Name + "-" + s.Id)
}

func (s *Service) SetUrl(serverDomain string) {
	if !s.ExposeSettings.IsExposed {
		return
	}

	protocol := "http"
	if s.ExposeSettings.Tls {
		protocol = "https"
	}

	subDomain := s.ExposeSettings.SubDomain
	if subDomain != "" {
		subDomain += "."
	}

	s.Url = protocol + "://" + subDomain + serverDomain
}
