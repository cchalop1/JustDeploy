package service

import (
	"cchalop1.com/deploy/internal/adapter"
	"sync"
)

type DeployService struct {
	Operation         sync.Mutex
	DockerAdapter     *adapter.DockerAdapter
	DatabaseAdapter   *adapter.DatabaseAdapter
	FilesystemAdapter *adapter.FilesystemAdapter
	NetworkAdapter    *adapter.NetworkAdapter
	GithubAdapter     *adapter.GithubAdapter
	GitAdapter        *adapter.GitAdapter
}
