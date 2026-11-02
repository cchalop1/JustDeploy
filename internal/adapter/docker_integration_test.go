package adapter

import (
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
	"context"
	"github.com/moby/moby/client"
	"os"
	"testing"
)

// Opt-in: only resources with this test's random ID are created or removed.
func TestDockerIntegrationPersistenceAndRollback(t *testing.T) {
	if os.Getenv("JUSTDEPLOY_DOCKER_TESTS") != "1" {
		t.Skip("set JUSTDEPLOY_DOCKER_TESTS=1 to test an isolated set of real Docker resources")
	}
	d := NewDockerAdapter()
	if err := d.ConnectClient(); err != nil {
		t.Fatal(err)
	}
	defer d.client.Close()
	id := utils.GenerateRandomPassword(16)
	volume := "justdeploy-test-" + id
	s := domain.Service{Id: id, Name: "audit-test", Type: "database", ImageName: "alpine:3.22", Cmd: []string{"sh", "-c", "echo seeded > /data/probe; exec sleep 300"}, ServiceRuntime: domain.ServiceRuntime{ProjectID: id, Volumes: []domain.ServiceVolume{{Name: volume, Target: "/data"}}}}
	t.Cleanup(func() {
		if err := d.Delete(s.GetDockerName()); err != nil {
			t.Error(err)
		}
		if err := d.Delete(s.GetDockerName() + "-previous"); err != nil {
			t.Error(err)
		}
		if _, err := d.client.VolumeRemove(context.Background(), volume, client.VolumeRemoveOptions{}); err != nil {
			t.Error(err)
		}
		if _, err := d.client.NetworkRemove(context.Background(), ProjectNetwork(s), client.NetworkRemoveOptions{}); err != nil {
			t.Error(err)
		}
	})
	if err := d.EnsureImage(s.ImageName); err != nil {
		t.Fatal(err)
	}
	if err := d.RunImageWithTLS(s, "", false); err != nil {
		t.Fatal(err)
	}
	first, err := d.inspect(context.Background(), s.GetDockerName())
	if err != nil {
		t.Fatal(err)
	}
	s.Cmd = []string{"sh", "-c", "test \"$(cat /data/probe)\" = seeded; exec sleep 300"}
	if err := d.RunImageWithTLS(s, "", false); err != nil {
		t.Fatalf("data not retained across replacement: %v", err)
	}
	second, err := d.inspect(context.Background(), s.GetDockerName())
	if err != nil {
		t.Fatal(err)
	}
	if first.ID == second.ID {
		t.Fatal("test did not replace the container")
	}
	s.Cmd = []string{"sh", "-c", "exit 1"}
	if err := d.RunImageWithTLS(s, "", false); err == nil {
		t.Fatal("crashing replacement accepted")
	}
	restored, err := d.inspect(context.Background(), s.GetDockerName())
	if err != nil {
		t.Fatal(err)
	}
	if restored.ID != second.ID || !restored.State.Running {
		t.Fatal("previous container was not restored")
	}
}
