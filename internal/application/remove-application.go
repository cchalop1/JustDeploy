package application

import "cchalop1.com/deploy/internal/adapter"

func RemoveApplication(applicationName string, dockerAdapter *adapter.DockerAdapter) error {
	dockerAdapter.Delete(applicationName, true)
	return nil
}
