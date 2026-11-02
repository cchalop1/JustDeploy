package application

import "cchalop1.com/deploy/internal/api/service"

func RemoveApplicationById(ds *service.DeployService, id string) error { return DeleteService(ds, id) }
