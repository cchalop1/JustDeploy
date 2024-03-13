package service

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
)

type DeployService struct {
	DockerAdapter   *adapter.DockerAdapter
	DatabaseAdapter *adapter.DatabaseAdapter
	DeployConfig    *dto.DeployConfigDto
}
