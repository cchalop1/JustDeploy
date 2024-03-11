package adapter

import "cchalop1.com/deploy/internal/domain"

type DatabaseAdapter struct {
}

func NewDatabaseAdapter() *DatabaseAdapter {
	return &DatabaseAdapter{}
}

func (d *DatabaseAdapter) GetState() domain.DeployConfigDto {

	return domain.DeployConfigDto{}
}

func (d *DatabaseAdapter) SaveState(deployConfig domain.DeployConfigDto) error {
	return nil
}
