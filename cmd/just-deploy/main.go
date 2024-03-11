package main

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/application"
)

func main() {
	// filesystemAdapter := adapter.NewFilesystemAdapter()
	databaseAdapter := adapter.NewDatabaseAdapter()

	deployConfig := application.GetDeployConfig(databaseAdapter)
	httpAdapter := application.NewHttpAdapter(deployConfig)

	httpAdapter.StartServer(true)
}
