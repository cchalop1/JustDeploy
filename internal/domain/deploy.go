package domain

import (
	"strings"
)

// TODO: replace the status by a enum
type Deploy struct {
	Id             string    `json:"id"`
	ServerId       string    `json:"serverId"`
	ProjectId      string    `json:"projectId"`
	EnableTls      bool      `json:"enableTls"`
	Email          string    `json:"email"`
	ServicesDeploy []Service `json:"servicesDeploy"`
}

func (d *Deploy) GetDockerName() string {
	// TODO: put the server id in the docker name
	return strings.ToLower(d.Id) + "-" + strings.ToLower(d.ServerId)
}
