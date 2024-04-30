package domain

import "cchalop1.com/deploy/internal/api/dto"

// TODO: change to save serviceConfig
type Service struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Envs        []dto.Env `json:"envs"`
	VolumsNames []string  `json:"volumsNames"`
	Status      string    `json:"status"`
	Host        string    `json:"host"`
	DeployId    string    `json:"deployId"`
	ImageName   string    `json:"imageName"`
}
