// Package docs JustDeploy API.
//
// @title JustDeploy API
// @version 1.0
// @description JustDeploy API - Simplify Development & Deployment
//
// @contact.name JustDeploy Support
// @contact.url https://github.com/cchalop1/just-deploy
// @contact.email support@justdeploy.com
//
// @license.name MIT
// @license.url https://github.com/cchalop1/just-deploy/blob/main/LICENSE
//
// @host localhost:8080
// @BasePath /
// @schemes http https
//
// @tag.name Server
// @tag.description Server management endpoints
//
// @tag.name Deploy
// @tag.description Deployment related endpoints
//
// @tag.name Service
// @tag.description Service management endpoints
//
// @tag.name GitHub
// @tag.description GitHub integration endpoints
//
// @tag.name Database
// @tag.description Database management endpoints
//
// @tag.name Version
// @tag.description Version information endpoints
package docs

// Server represents a server instance
type Server struct {
	ID   string `json:"id" example:"server-123"`
	Name string `json:"name" example:"production-server"`
}

// Service represents a service instance
type Service struct {
	ID       string `json:"id" example:"service-123"`
	Name     string `json:"name" example:"web-app"`
	Status   string `json:"status" example:"running"`
	ServerID string `json:"serverId" example:"server-123"`
}

// Deploy represents a deployment
type Deploy struct {
	ID        string `json:"id" example:"deploy-123"`
	ServiceID string `json:"serviceId" example:"service-123"`
	Status    string `json:"status" example:"completed"`
	CreatedAt string `json:"createdAt" example:"2023-01-01T12:00:00Z"`
}

// ResponseApi is a generic response
type ResponseApi struct {
	Message string `json:"message" example:"Operation successful"`
}

// VersionDto contains version information
type VersionDto struct {
	TagName   string `json:"tagName" example:"v1.0.0"`
	GithubUrl string `json:"githubUrl" example:"https://github.com/cchalop1/just-deploy/releases/tag/v1.0.0"`
}

// GithubIsConnected indicates if GitHub is connected
type GithubIsConnected struct {
	IsConnected bool `json:"isConnected" example:"true"`
}

// SaveGithubToken contains GitHub token information
type SaveGithubToken struct {
	Token string `json:"token" example:"github_pat_123456789"`
}

// CreateDatabaseRequest represents a request to create a database
type CreateDatabaseRequest struct {
	Name     string `json:"name" example:"my-database"`
	Type     string `json:"type" example:"postgres"`
	ServerID string `json:"serverId" example:"server-123"`
}

// CreateRepoRequest represents a request to create a repository
type CreateRepoRequest struct {
	Name     string `json:"name" example:"my-repo"`
	URL      string `json:"url" example:"https://github.com/user/repo"`
	ServerID string `json:"serverId" example:"server-123"`
}

// DeployRequest represents a request to deploy an application
type DeployRequest struct {
	ServiceID string `json:"serviceId" example:"service-123"`
	RepoID    string `json:"repoId" example:"repo-123"`
}

// AddDomainRequest represents a request to add a domain to a server
type AddDomainRequest struct {
	Domain string `json:"domain" example:"example.com"`
}

// @Summary Get version information
// @Description Get the current version of the application
// @Tags Version
// @Accept json
// @Produce json
// @Success 200 {object} VersionDto
// @Router /api/version [get]
func GetVersionHandler() {}

// @Summary Check if GitHub is connected
// @Description Check if the application is connected to GitHub
// @Tags GitHub
// @Accept json
// @Produce json
// @Success 200 {object} GithubIsConnected
// @Router /api/github/is-connected [get]
func GetGithubIsConnectedHandler() {}

// @Summary Connect GitHub app
// @Description Connect the application to GitHub using a code
// @Tags GitHub
// @Accept json
// @Produce json
// @Param code path string true "GitHub code"
// @Success 200 {object} ResponseApi
// @Router /api/github/connect/{code} [post]
func PostConnectGithubAppHandler() {}

// @Summary Get GitHub repositories
// @Description Get a list of GitHub repositories
// @Tags GitHub
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /api/github/repos [get]
func GetGithubRepos() {}

// @Summary Save GitHub access token
// @Description Save a GitHub access token for a specific installation
// @Tags GitHub
// @Accept json
// @Produce json
// @Param installationId path string true "GitHub installation ID"
// @Param token body SaveGithubToken true "GitHub token"
// @Success 200 {object} ResponseApi
// @Router /api/github/save-access-token/{installationId} [post]
func PostSaveAccessTokenHandler() {}

// @Summary Create a database
// @Description Create a new database
// @Tags Database
// @Accept json
// @Produce json
// @Param database body CreateDatabaseRequest true "Database information"
// @Success 200 {object} ResponseApi
// @Router /api/database/create [post]
func PostCreateDatabaseHandler() {}

// @Summary Create a repository
// @Description Create a new repository
// @Tags GitHub
// @Accept json
// @Produce json
// @Param repo body CreateRepoRequest true "Repository information"
// @Success 200 {object} ResponseApi
// @Router /api/repo/create [post]
func PostCreateRepoHandler() {}

// @Summary Deploy an application
// @Description Deploy an application
// @Tags Deploy
// @Accept json
// @Produce json
// @Param deploy body DeployRequest true "Deploy information"
// @Success 200 {object} ResponseApi
// @Router /api/deploy [post]
func PostDeployHandler() {}

// @Summary Get server information
// @Description Get information about a specific server
// @Tags Server
// @Accept json
// @Produce json
// @Param id path string true "Server ID"
// @Success 200 {object} Server
// @Router /api/server/{id} [get]
func GetServerInfoHandler() {}

// @Summary Add domain to server
// @Description Add a domain to a specific server
// @Tags Server
// @Accept json
// @Produce json
// @Param id path string true "Server ID"
// @Param domain body AddDomainRequest true "Domain information"
// @Success 200 {object} ResponseApi
// @Router /api/server/domain [post]
func PostAddDomainToServerById() {}

// @Summary Delete a service
// @Description Delete a specific service
// @Tags Service
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Success 200 {object} ResponseApi
// @Router /api/service/{serviceId} [delete]
func DeleteServiceHandler() {}

// @Summary Get server loading state
// @Description Subscribe to server loading state updates
// @Tags Server
// @Accept json
// @Produce json
// @Param id path string true "Server ID"
// @Success 200 {object} ResponseApi
// @Router /api/server/{id}/loading [get]
func SubscriptionCreateServerLoadingState() {}

// @Summary Get deploy loading state
// @Description Subscribe to deploy loading state updates
// @Tags Deploy
// @Accept json
// @Produce json
// @Param id path string true "Deploy ID"
// @Success 200 {object} ResponseApi
// @Router /api/deploy/{id}/loading [get]
func SubscriptionCreateDeployLoadingState() {}

// @Summary Get services
// @Description Get a list of services
// @Tags Service
// @Accept json
// @Produce json
// @Success 200 {array} Service
// @Router /api/services [get]
func GetServicesHandler() {}

// @Summary Get services list
// @Description Get a list of services for a specific server
// @Tags Service
// @Accept json
// @Produce json
// @Param serverId path string true "Server ID"
// @Success 200 {array} Service
// @Router /api/services/{serverId} [get]
func GetServicesListHandler() {}

// @Summary Update a service
// @Description Update a specific service
// @Tags Service
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Param service body Service true "Service information"
// @Success 200 {object} ResponseApi
// @Router /api/service/{serviceId} [put]
func UpdateServiceHandler() {}

// @Summary Re-deploy an application
// @Description Re-deploy an existing application
// @Tags Deploy
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Success 200 {object} ResponseApi
// @Router /api/service/{serviceId}/redeploy [post]
func ReDeployAppHandler() {}

// @Summary Remove an application
// @Description Remove an existing application
// @Tags Service
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Success 200 {object} ResponseApi
// @Router /api/service/{serviceId}/remove [post]
func RemoveAppHandler() {}

// @Summary Start an application
// @Description Start an existing application
// @Tags Service
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Success 200 {object} ResponseApi
// @Router /api/service/{serviceId}/start [post]
func StartAppHandler() {}

// @Summary Stop an application
// @Description Stop an existing application
// @Tags Service
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Success 200 {object} ResponseApi
// @Router /api/service/{serviceId}/stop [post]
func StopAppHandler() {}

// @Summary Get server proxy logs
// @Description Get proxy logs for a specific server
// @Tags Server
// @Accept json
// @Produce json
// @Param serverId path string true "Server ID"
// @Success 200 {object} ResponseApi
// @Router /api/server/{serverId}/proxy-logs [get]
func GetServerProxyLogsHandler() {}

// @Summary Get services from Docker Compose
// @Description Get services defined in a Docker Compose file
// @Tags Service
// @Accept json
// @Produce json
// @Param serverId path string true "Server ID"
// @Success 200 {array} Service
// @Router /api/server/{serverId}/docker-compose-services [get]
func GetServicesFromDockerComposeHandler() {}

// @Summary Get deploys by server ID
// @Description Get a list of deployments for a specific server
// @Tags Deploy
// @Accept json
// @Produce json
// @Param serverId path string true "Server ID"
// @Success 200 {array} Deploy
// @Router /api/deploys/{serverId} [get]
func GetDeployByServerIdHandler() {}

// @Summary Get deploy configuration
// @Description Get configuration for a specific deployment
// @Tags Deploy
// @Accept json
// @Produce json
// @Param deployId path string true "Deploy ID"
// @Success 200 {object} Deploy
// @Router /api/deploy/{deployId}/config [get]
func GetDeployConfigHandler() {}

// @Summary Get deploy list
// @Description Get a list of all deployments
// @Tags Deploy
// @Accept json
// @Produce json
// @Success 200 {array} Deploy
// @Router /api/deploys [get]
func GetDeployListHandler() {}

// @Summary Edit a deployment
// @Description Edit an existing deployment
// @Tags Deploy
// @Accept json
// @Produce json
// @Param deployId path string true "Deploy ID"
// @Param deploy body Deploy true "Deploy information"
// @Success 200 {object} ResponseApi
// @Router /api/deploy/{deployId} [put]
func EditDeploymentHandler() {}
