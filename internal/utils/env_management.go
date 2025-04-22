package utils

import (
	"fmt"
	"strings"

	"cchalop1.com/deploy/internal/api/dto"
)

// GenerateEnvValues generates values for environment variables based on their type.
// For secret variables (IsSecret=true), it generates a random string.
// For non-secret variables, it uses the lowercase name of the variable.
func GenerateEnvValues(envs []dto.Env) []dto.Env {
	generatedEnvs := make([]dto.Env, len(envs))

	for i, env := range envs {
		newEnv := dto.Env{
			Name:     env.Name,
			IsSecret: env.IsSecret,
		}

		if env.IsSecret {
			// Generate a random string for secret variables
			newEnv.Value = GenerateRandomPassword(12)
		} else {
			// Use lowercase name for non-secret variables
			newEnv.Value = strings.ToLower(env.Name)
		}

		generatedEnvs[i] = newEnv
	}

	return generatedEnvs
}

// EnvToSlice converts a slice of dto.Env to a slice of strings in the format "NAME=VALUE"
// Only includes environment variables that have both a name and a value
func EnvToSlice(envVars []dto.Env) []string {
	envSlice := make([]string, 0, len(envVars))
	for _, value := range envVars {
		if value.Name != "" && value.Value != "" {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", value.Name, value.Value))
		}
	}
	return envSlice
}
