package graph

import "cchalop1.com/deploy/internal/api/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	deployService *service.DeployService
}

func NewResolver(deployService *service.DeployService) *Resolver {
	return &Resolver{
		deployService: deployService,
	}
}
