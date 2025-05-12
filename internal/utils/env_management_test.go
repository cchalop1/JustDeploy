package utils

import (
	"reflect"
	"strings"
	"testing"

	"cchalop1.com/deploy/internal/api/dto"
)

func TestReplaceEnvVariablesInCmd(t *testing.T) {
	tests := []struct {
		name     string
		cmd      []string
		envs     []dto.Env
		expected []string
	}{
		{
			name: "replace single variable",
			cmd:  []string{"echo", "$VAR1"},
			envs: []dto.Env{
				{
					Name:  "VAR1",
					Value: "value1",
				},
			},
			expected: []string{"echo", "value1"},
		},
		{
			name: "replace multiple variables",
			cmd:  []string{"echo", "$VAR1", "$VAR2"},
			envs: []dto.Env{
				{
					Name:  "VAR1",
					Value: "value1",
				},
				{
					Name:  "VAR2",
					Value: "value2",
				},
			},
			expected: []string{"echo", "value1", "value2"},
		},
		{
			name: "variable in the middle of a string",
			cmd:  []string{"echo", "prefix-$VAR1-suffix"},
			envs: []dto.Env{
				{
					Name:  "VAR1",
					Value: "value1",
				},
			},
			expected: []string{"echo", "prefix-value1-suffix"},
		},
		{
			name: "no variables to replace",
			cmd:  []string{"echo", "hello", "world"},
			envs: []dto.Env{
				{
					Name:  "VAR1",
					Value: "value1",
				},
			},
			expected: []string{"echo", "hello", "world"},
		},
		{
			name: "variable not found in envs",
			cmd:  []string{"echo", "$UNKNOWN"},
			envs: []dto.Env{
				{
					Name:  "VAR1",
					Value: "value1",
				},
			},
			expected: []string{"echo", "$UNKNOWN"},
		},
		{
			name: "empty cmd",
			cmd:  []string{},
			envs: []dto.Env{
				{
					Name:  "VAR1",
					Value: "value1",
				},
			},
			expected: []string{},
		},
		{
			name:     "empty envs",
			cmd:      []string{"echo", "$VAR1"},
			envs:     []dto.Env{},
			expected: []string{"echo", "$VAR1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceEnvVariablesInCmd(tt.cmd, tt.envs)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReplaceEnvVariablesInCmd() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateEnvValues(t *testing.T) {
	tests := []struct {
		name    string
		envs    []dto.Env
		checkFn func(t *testing.T, input []dto.Env, result []dto.Env)
	}{
		{
			name: "generate values for mixed env variables",
			envs: []dto.Env{
				{
					Name:     "DATABASE_URL",
					IsSecret: true,
				},
				{
					Name:     "PORT",
					IsSecret: false,
				},
			},
			checkFn: func(t *testing.T, input []dto.Env, result []dto.Env) {
				if len(result) != len(input) {
					t.Errorf("Expected %d env variables, got %d", len(input), len(result))
				}

				// Check all names and IsSecret flags are preserved
				for i, env := range result {
					if env.Name != input[i].Name {
						t.Errorf("Expected name %s, got %s", input[i].Name, env.Name)
					}
					if env.IsSecret != input[i].IsSecret {
						t.Errorf("Expected IsSecret %v, got %v", input[i].IsSecret, env.IsSecret)
					}
				}

				// Check non-secret var uses lowercase name
				for _, env := range result {
					if !env.IsSecret && env.Value != strings.ToLower(env.Name) {
						t.Errorf("Non-secret env var %s should have value %s, got %s",
							env.Name, strings.ToLower(env.Name), env.Value)
					}
				}

				// Check secret vars have random 12-char passwords
				for _, env := range result {
					if env.IsSecret {
						if len(env.Value) != 12 {
							t.Errorf("Secret env var %s should have 12-char value, got %d chars: %s",
								env.Name, len(env.Value), env.Value)
						}
						if env.Value == env.Name {
							t.Errorf("Secret env var %s should not have its name as value", env.Name)
						}
					}
				}
			},
		},
		{
			name: "generate values for empty input",
			envs: []dto.Env{},
			checkFn: func(t *testing.T, input []dto.Env, result []dto.Env) {
				if len(result) != 0 {
					t.Errorf("Expected empty result for empty input, got %v", result)
				}
			},
		},
		{
			name: "generate values for only secret vars",
			envs: []dto.Env{
				{
					Name:     "API_KEY",
					IsSecret: true,
				},
				{
					Name:     "JWT_SECRET",
					IsSecret: true,
				},
			},
			checkFn: func(t *testing.T, input []dto.Env, result []dto.Env) {
				if len(result) != len(input) {
					t.Errorf("Expected %d env variables, got %d", len(input), len(result))
				}

				// Check all values are 12 chars and different from each other
				values := make(map[string]bool)
				for _, env := range result {
					if len(env.Value) != 12 {
						t.Errorf("Secret env var %s should have 12-char value, got %d chars",
							env.Name, len(env.Value))
					}
					values[env.Value] = true
				}

				// Check uniqueness of values (should be random)
				if len(values) != len(result) {
					t.Errorf("Expected %d unique values, got %d", len(result), len(values))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateEnvValues(tt.envs)
			tt.checkFn(t, tt.envs, result)
		})
	}
}

func TestEnvToSlice(t *testing.T) {
	tests := []struct {
		name     string
		envVars  []dto.Env
		expected []string
	}{
		{
			name: "convert valid env vars to slice",
			envVars: []dto.Env{
				{
					Name:  "DATABASE_URL",
					Value: "postgres://localhost:5432",
				},
				{
					Name:  "PORT",
					Value: "8080",
				},
			},
			expected: []string{
				"DATABASE_URL=postgres://localhost:5432",
				"PORT=8080",
			},
		},
		{
			name: "filter out env vars with empty name",
			envVars: []dto.Env{
				{
					Name:  "",
					Value: "value1",
				},
				{
					Name:  "PORT",
					Value: "8080",
				},
			},
			expected: []string{
				"PORT=8080",
			},
		},
		{
			name: "filter out env vars with empty value",
			envVars: []dto.Env{
				{
					Name:  "DATABASE_URL",
					Value: "",
				},
				{
					Name:  "PORT",
					Value: "8080",
				},
			},
			expected: []string{
				"PORT=8080",
			},
		},
		{
			name:     "empty input produces empty output",
			envVars:  []dto.Env{},
			expected: []string{},
		},
		{
			name: "all invalid env vars produce empty output",
			envVars: []dto.Env{
				{
					Name:  "",
					Value: "",
				},
				{
					Name:  "PORT",
					Value: "",
				},
				{
					Name:  "",
					Value: "8080",
				},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnvToSlice(tt.envVars)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("EnvToSlice() = %v, want %v", result, tt.expected)
			}
		})
	}
}
