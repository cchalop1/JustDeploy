package adapter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComposeFilenameVariantsAndErrors(t *testing.T) {
	for _, name := range []string{"compose.yml", "compose.yaml", "docker-compose.yml", "docker-compose.yaml"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, name), []byte("services:\n  db:\n    image: postgres:17\n    environment:\n      POSTGRES_PASSWORD: example\n"), 0600); err != nil {
				t.Fatal(err)
			}
			configs, err := NewFilesystemAdapter().GetComposeConfigOfDeploy(dir)
			if err != nil {
				t.Fatal(err)
			}
			if configs["db"].Image != "postgres:17" {
				t.Fatal("service missing")
			}
		})
	}
	for _, yaml := range []string{"services: [", "services:\n  db:\n    image: postgres:17\n    privileged: true\n", "services:\n  db:\n    image: postgres:17\n    environment:\n      PASSWORD: ${PASSWORD}\n"} {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(yaml), 0600)
		if _, err := NewFilesystemAdapter().GetComposeConfigOfDeploy(dir); err == nil {
			t.Fatal("invalid/unsupported Compose accepted")
		}
	}
}
