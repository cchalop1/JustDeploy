package domain

import (
	"strings"

	"cchalop1.com/deploy/internal/api/dto"
)

// TODO: replace the status by a enum
type Deploy struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	ServerId       string    `json:"serverId"`
	EnableTls      bool      `json:"enableTls"`
	Email          string    `json:"email"`
	PathToSource   string    `json:"pathToSource"`
	Envs           []dto.Env `json:"envs"`
	DeployOnCommit bool      `json:"deployOnCommit"`
	Status         string    `json:"status"`
	Url            string    `json:"url"`
	SubDomain      string    `json:"subDomain"`
}

func (d *Deploy) GetDockerName() string {
	return strings.ToLower(d.Id) + "-" + strings.ToLower(d.ServerId)
}
