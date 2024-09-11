package service

import (
	"cchalop1.com/deploy/internal/adapter"
)

// TODO: move to application

type DeployService struct {
	DockerAdapter     *adapter.DockerAdapter
	DatabaseAdapter   *adapter.DatabaseAdapter
	FilesystemAdapter *adapter.FilesystemAdapter
	EventAdapter      *adapter.AdapterEvent
}
