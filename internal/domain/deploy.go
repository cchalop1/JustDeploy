package domain

import "cchalop1.com/deploy/internal/api/dto"

// TODO: replace the status by a enum
type Deploy struct {
	Id             string
	Name           string
	ServerId       string
	EnableTls      bool
	Email          string
	PathToSource   string
	Envs           []dto.Env
	DeployOnCommit bool
	Status         string
	Url            string
}
