package main

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/application"
)

func main() {
	filesystemAdapter := adapter.NewFilesystemAdapter()

	deployConfig := application.GetFormDetails(filesystemAdapter)
	httpAdapter := application.NewHttpAdapter(deployConfig)

	httpAdapter.StartServer(true)
}
