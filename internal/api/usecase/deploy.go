package usecase

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/models"
)

type DeployUseCase struct {
	DockerAdapter *adapter.DockerAdapter
	DeployConfig  *models.DeployConfigDto
}
