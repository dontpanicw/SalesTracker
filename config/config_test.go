package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
	}{
		{
			name: "default values",
			envVars: map[string]string{},
			want: &Config{
				DatabaseDSN: "host=localhost port=5432 user=postgres password=postgres dbname=analytics sslmode=disable",
				ServerPort:  "8080",
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"DB_HOST":     "customhost",
				"DB_PORT":     "5433",
				"DB_USER":     "customuser",
				"DB_PASSWORD": "custompass",
				"DB_NAME":     "customdb",
				"SERVER_PORT": "9090",
			},
			want: &Config{
				DatabaseDSN: "host=customhost port=5433 user=customuser password=custompass dbname=customdb sslmode=disable",
				ServerPort:  "9090",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Clean up
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			got, err := Load()
			if err != nil {
				t.Errorf("Load() error = %v", err)
				return
			}
			if got.DatabaseDSN != tt.want.DatabaseDSN {
				t.Errorf("Load() DatabaseDSN = %v, want %v", got.DatabaseDSN, tt.want.DatabaseDSN)
			}
			if got.ServerPort != tt.want.ServerPort {
				t.Errorf("Load() ServerPort = %v, want %v", got.ServerPort, tt.want.ServerPort)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "use default when env not set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
		{
			name:         "use env value when set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
