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
	Type        string `json:"type"`
}

var servicesConfigMap = map[string]ServicesConfig{
	"Postgres": {
		DefaultPort: 5432,
		Name:        "Postgres",
		Icon:        "https://upload.wikimedia.org/wikipedia/commons/thumb/2/29/Postgresql_elephant.svg/640px-Postgresql_elephant.svg.png",
		Type:        "database",
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
		Type:        "database",
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
		Type:        "database",
		Icon:        "https://grafikart.fr/uploads/icons/redis.svg",
		Env: []dto.Env{
			{Name: "REDIS_PASSWORD", Value: "", IsSecret: true},
		},
		Config: container.Config{
			Image: "redis:latest",
			Cmd:   []string{"--requirepass", "$REDIS_PASSWORD"},
		},
	},
	// "Keycloak": {
	// 	DefaultPort: 8080,
	// 	Name:        "Keycloak",
	// 	Type:        "database",
	// 	Icon:        "https://www.keycloak.org/resources/images/logo.svg",
	// 	Env: []dto.Env{
	// 		{Name: "KEYCLOAK_USER", Value: "", IsSecret: false},
	// 		{Name: "KEYCLOAK_PASSWORD", Value: "", IsSecret: true},
	// 	},
	// 	Config: container.Config{
	// 		Image: "quay.io/keycloak/keycloak",
	// 		Cmd:   []string{"start-dev"},
	// 		Env: []string{
	// 			"KEYCLOAK_USER=$KEYCLOAK_USER",
	// 			"KEYCLOAK_PASSWORD=$KEYCLOAK_PASSWORD",
	// 		},
	// 	},
	// },
	"Minio": {
		DefaultPort: 9000,
		Name:        "Minio",
		Type:        "database",
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
	// Ollama LLM Models
	"Gemma3": {
		DefaultPort: 11434,
		Name:        "Gemma3",
		Type:        "llm",
		Icon:        "https://ollama.com/public/ollama.png",
		Env:         []dto.Env{},
		Config: container.Config{
			Image: "ollama/ollama",
			Cmd:   []string{"ollama", "run", "gemma3"},
		},
	},
	"Llama3.2": {
		DefaultPort: 11434,
		Name:        "Llama3.2",
		Type:        "llm",
		Icon:        "https://ollama.com/public/ollama.png",
		Env:         []dto.Env{},
		Config: container.Config{
			Image: "ollama/ollama",
			Cmd:   []string{"ollama", "run", "llama3.2"},
		},
	},
	"Mistral": {
		DefaultPort: 11434,
		Name:        "Mistral",
		Type:        "llm",
		Icon:        "https://ollama.com/public/ollama.png",
		Env:         []dto.Env{},
		Config: container.Config{
			Image: "ollama/ollama",
			Cmd:   []string{"ollama", "run", "mistral"},
		},
	},
	"Qwen2.5": {
		DefaultPort: 11434,
		Name:        "Qwen2.5",
		Type:        "llm",
		Icon:        "https://ollama.com/public/ollama.png",
		Env:         []dto.Env{},
		Config: container.Config{
			Image: "ollama/ollama",
			Cmd:   []string{"ollama", "run", "qwen2.5"},
		},
	},
	// "Umami": {
	// 	DefaultPort: 3000,
	// 	Name:        "Umami",
	// 	Type:        "database",
	// 	Icon:        "https://umami.is/logo-icon-dark.svg",
	// 	Env: []dto.Env{
	// 		{Name: "DATABASE_URL", Value: "postgresql://umami:umami@localhost:5432/umami", IsSecret: false},
	// 		{Name: "HASH_SALT", Value: "", IsSecret: true},
	// 	},
	// 	Config: container.Config{
	// 		Image: "ghcr.io/umami-software/umami:postgresql-latest",
	// 		Env: []string{
	// 			"DATABASE_URL=$DATABASE_URL",
	// 			"HASH_SALT=$HASH_SALT",
	// 		},
	// 	},
	// },
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
