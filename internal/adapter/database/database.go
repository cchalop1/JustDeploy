package database

import (
	"cchalop1.com/deploy/internal/api/dto"
	"github.com/docker/docker/api/types/container"
)

type ServicesConfig struct {
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	Env         []dto.Env `json:"env"`
	DefaultPort int       `json:"defaultPort"`
	Config      container.Config
}

var servicesConfigMap = map[string]ServicesConfig{
	"Postgres": {
		DefaultPort: 5432,
		Name:        "Postgres",
		Icon:        "https://upload.wikimedia.org/wikipedia/commons/thumb/2/29/Postgresql_elephant.svg/640px-Postgresql_elephant.svg.png",
		Env: []dto.Env{
			{Name: "POSTGRES_USER", Value: "", IsSecret: false},
			{Name: "POSTGRES_PASSWORD", Value: "", IsSecret: true},
			{Name: "POSTGRES_DB", Value: "", IsSecret: false},
		},
		Config: container.Config{
			Image: "postgres:latest",
			Cmd:   []string{"postgres"},
			Env: []string{
				"POSTGRES_USER=$POSTGRES_USER",
				"POSTGRES_PASSWORD=$POSTGRES_PASSWORD",
				"POSTGRES_DB=$POSTGRES_DB",
			},
		},
	},
	"Mongo": {
		DefaultPort: 27017,
		Name:        "Mongo",
		Icon:        "https://upload.wikimedia.org/wikipedia/fr/thumb/4/45/MongoDB-Logo.svg/1280px-MongoDB-Logo.svg.png",
		Env: []dto.Env{
			{Name: "MONGO_INITDB_ROOT_USERNAME", Value: "", IsSecret: false},
			{Name: "MONGO_INITDB_ROOT_PASSWORD", Value: "", IsSecret: true},
		},
		Config: container.Config{
			Image: "mongo",
			Env: []string{
				"MONGO_INITDB_ROOT_USERNAME=$MONGO_INITDB_ROOT_USERNAME",
				"MONGO_INITDB_ROOT_PASSWORD=$MONGO_INITDB_ROOT_PASSWORD",
			},
		},
	},
	"Redis": {
		DefaultPort: 6379,
		Name:        "Redis",
		Icon:        "https://grafikart.fr/uploads/icons/redis.svg",
		Env: []dto.Env{
			{Name: "REDIS_PASSWORD", Value: "", IsSecret: true},
		},
		Config: container.Config{
			Image: "redis:latest",
			Cmd:   []string{"redis-server", "--requirepass", "$REDIS_PASSWORD"},
		},
	},
	"Keycloak": {
		DefaultPort: 8080,
		Name:        "Keycloak",
		Icon:        "https://www.keycloak.org/resources/images/logo.svg",
		Env: []dto.Env{
			{Name: "KEYCLOAK_USER", Value: "", IsSecret: false},
			{Name: "KEYCLOAK_PASSWORD", Value: "", IsSecret: true},
		},
		Config: container.Config{
			Image: "quay.io/keycloak/keycloak",
			Cmd:   []string{"start-dev"},
			Env: []string{
				"KEYCLOAK_USER=$KEYCLOAK_USER",
				"KEYCLOAK_PASSWORD=$KEYCLOAK_PASSWORD",
			},
		},
	},
	"Minio": {
		DefaultPort: 9000,
		Name:        "Minio",
		Icon:        "https://min.io/resources/img/logo.svg",
		Env: []dto.Env{
			{Name: "MINIO_ACCESS_KEY", Value: "", IsSecret: false},
			{Name: "MINIO_SECRET_KEY", Value: "", IsSecret: true},
		},
		Config: container.Config{
			Image: "minio/minio:latest",
			Cmd:   []string{"server", "/data"},
			Env: []string{
				"MINIO_ACCESS_KEY=$MINIO_ACCESS_KEY",
				"MINIO_SECRET_KEY=$MINIO_SECRET_KEY",
			},
		},
	},
}

func GetListOfDatabasesServices() []ServicesConfig {
	var servicesConfigArray []ServicesConfig
	for _, serviceConfig := range servicesConfigMap {
		servicesConfigArray = append(servicesConfigArray, serviceConfig)
	}
	return servicesConfigArray
}

// TODO: function to get service by name
func GetServiceByName(name string) (ServicesConfig, error) {
	service, ok := servicesConfigMap[name]
	if !ok {
		return ServicesConfig{}, nil
	}
	return service, nil
}
