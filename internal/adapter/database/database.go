package database

import (
	"cchalop1.com/deploy/internal/api/http/dto"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type ServicesConfig struct {
	Name   string    `json:"name"`
	Icon   string    `json:"icon"`
	Env    []dto.Env `json:"env"`
	Config container.Config
}

var servicesConfigMap = map[string]ServicesConfig{
	"Postgres": {
		Name: "Postgres",
		Icon: "https://upload.wikimedia.org/wikipedia/commons/thumb/2/29/Postgresql_elephant.svg/640px-Postgresql_elephant.svg.png",
		Env: []dto.Env{
			{Name: "POSTGRES_USER", Value: "", IsSecret: false},
			{Name: "POSTGRES_PASSWORD", Value: "", IsSecret: true},
			{Name: "POSTGRES_DB", Value: "", IsSecret: false},
		},
		Config: container.Config{
			Image: "postgres:latest",
			Cmd:   []string{"postgres"},
			ExposedPorts: map[nat.Port]struct{}{
				"5432/tcp": {},
			},
			Env: []string{
				"POSTGRES_USER=$POSTGRES_USER",
				"POSTGRES_PASSWORD=$POSTGRES_PASSWORD",
				"POSTGRES_DB=$POSTGRES_DB",
			},
		},
	},
	"Mongo": {
		Name: "Mongo",
		Icon: "https://upload.wikimedia.org/wikipedia/fr/thumb/4/45/MongoDB-Logo.svg/1280px-MongoDB-Logo.svg.png",
		Env: []dto.Env{
			{Name: "MONGO_INITDB_ROOT_USERNAME", Value: "", IsSecret: false},
			{Name: "MONGO_INITDB_ROOT_PASSWORD", Value: "", IsSecret: true},
		},
		Config: container.Config{
			Image: "mongo",
			ExposedPorts: map[nat.Port]struct{}{
				"27017/tcp": {},
			},
			Env: []string{
				"MONGO_INITDB_ROOT_USERNAME=$MONGO_INITDB_ROOT_USERNAME",
				"MONGO_INITDB_ROOT_PASSWORD=$MONGO_INITDB_ROOT_PASSWORD",
			},
		},
	},
	"Redis": {
		Name: "Redis",
		Icon: "https://grafikart.fr/uploads/icons/redis.svg",
		Env: []dto.Env{
			{Name: "REDIS_PASSWORD", Value: "", IsSecret: true},
		},
		Config: container.Config{
			Image: "redis:latest",
			Cmd:   []string{"redis-server", "--requirepass", "$REDIS_PASSWORD"},
			ExposedPorts: map[nat.Port]struct{}{
				"6379/tcp": {},
			},
		},
	},
	"Keycloak": {
		Name: "Keycloak",
		Icon: "https://www.keycloak.org/resources/images/logo.svg",
		Env: []dto.Env{
			{Name: "KEYCLOAK_USER", Value: "", IsSecret: false},
			{Name: "KEYCLOAK_PASSWORD", Value: "", IsSecret: true},
		},
		Config: container.Config{
			Image: "quay.io/keycloak/keycloak",
			Cmd:   []string{"start-dev"},
			ExposedPorts: map[nat.Port]struct{}{
				"8080/tcp": {},
			},
			Env: []string{
				"KEYCLOAK_USER=$KEYCLOAK_USER",
				"KEYCLOAK_PASSWORD=$KEYCLOAK_PASSWORD",
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
