package service

import (
	"cchalop1.com/deploy/internal/adapter"
)

type DeployService struct {
	DockerAdapter     *adapter.DockerAdapter
	DatabaseAdapter   *adapter.DatabaseAdapter
	FilesystemAdapter *adapter.FilesystemAdapter
	EventAdapter      *adapter.AdapterEvent
	NetworkAdapter    *adapter.NetworkAdapter
	GithubAdapter     *adapter.GithubAdapter
}
