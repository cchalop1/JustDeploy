package adapter

import (
	"archive/tar"
	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/domain"

	"encoding/json"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func testDB(t *testing.T) *DatabaseAdapter {
	t.Helper()
	dir := t.TempDir()
	oldDir, oldDB := internal.JUSTDEPLOY_FOLDER, internal.DATABASE_SQLITE_PATH
	internal.JUSTDEPLOY_FOLDER = dir
	internal.DATABASE_SQLITE_PATH = filepath.Join(dir, "database.db")
	db := NewDatabaseAdapter()
	db.Init()
	t.Cleanup(func() { db.Close(); internal.JUSTDEPLOY_FOLDER = oldDir; internal.DATABASE_SQLITE_PATH = oldDB })
	return db
}
func TestSQLiteMigrationAndServiceRoundTrip(t *testing.T) {
	db := testDB(t)
	s := domain.Service{Id: "abc", Name: "app", Type: "github_repo", ServiceRuntime: domain.ServiceRuntime{ProjectID: "group", Branch: "trunk", Volumes: []domain.ServiceVolume{{Name: "storage", Target: "/data"}}}}
	if err := db.SaveService(s); err != nil {
		t.Fatal(err)
	}
	saved, err := db.GetServiceById(s.Id)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s.ServiceRuntime, saved.ServiceRuntime) {
		t.Fatalf("runtime metadata changed: %+v", saved)
	}
	// ALTER migration is idempotent and retains old rows.
	db.Close()
	db.Init()
	if _, err := db.GetServiceById(s.Id); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(internal.DATABASE_SQLITE_PATH)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatal("database must be private")
	}
	settings, err := db.ReadSettings()
	if err != nil || settings.AdminEmail != "" {
		t.Fatalf("initial settings: %v", err)
	}
	claimed, err := db.ClaimDelivery("delivery")
	if err != nil || !claimed {
		t.Fatal(err)
	}
	claimed, err = db.ClaimDelivery("delivery")
	if err != nil || claimed {
		t.Fatal("duplicate delivery accepted")
	}
}
func TestMountsPreserveExistingData(t *testing.T) {
	s := domain.Service{Id: "abc", Name: "postgres", Type: "database", ImageName: "postgres:17"}
	old := container.InspectResponse{ID: "old", Mounts: []container.MountPoint{{Type: mount.TypeVolume, Name: "anonymous-existing", Destination: "/var/lib/postgresql/data", RW: true}}}
	mounts, err := serviceMounts(s, old)
	if err != nil {
		t.Fatal(err)
	}
	if len(mounts) != 1 || mounts[0].Source != "anonymous-existing" {
		t.Fatalf("lost old volume: %+v", mounts)
	}
	old.Mounts = nil
	if _, err := serviceMounts(s, old); err == nil {
		t.Fatal("must refuse to replace unmounted stateful data")
	}
	mounts, err = serviceMounts(s, container.InspectResponse{})
	if err != nil {
		t.Fatal(err)
	}
	if len(mounts) != 1 || !strings.HasPrefix(mounts[0].Source, "justdeploy-abc-") {
		t.Fatalf("missing stable named volume: %+v", mounts)
	}
}
func TestBuildContextExcludesCredentials(t *testing.T) {
	dir := t.TempDir()
	for name, data := range map[string]string{"Dockerfile": "FROM scratch", ".env": "PASSWORD=secret", ".git/config": "token", ".dockerignore": "!.env\n!.git\n"} {
		path := filepath.Join(dir, name)
		os.MkdirAll(filepath.Dir(path), 0700)
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	stream, err := makeTar(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	reader := tar.NewReader(stream)
	found := false
	for {
		h, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if h.Name == "Dockerfile" {
			found = true
		}
		if h.Name == ".env" || strings.HasPrefix(h.Name, ".git/") {
			t.Fatalf("credentials in context: %s", h.Name)
		}
	}
	if !found {
		t.Fatal("Dockerfile missing")
	}
}
func TestDockerStreamPropagatesPullFailure(t *testing.T) {
	if err := consumeDockerStream(strings.NewReader(`{"error":"registry denied"}`), nil); err == nil {
		t.Fatal("pull failure ignored")
	}
}
func TestPrivateServiceHasNoTraefikRoute(t *testing.T) {
	d := NewDockerAdapter()
	s := domain.Service{Id: "abc", Name: "app", Type: "github_repo"}
	c := d.ConfigContainer(s)
	if c.Labels["traefik.enable"] != "false" {
		t.Fatal("private service exposed")
	}
	d.ExposeContainer(&c, ExposeContainerParams{IsTls: true, Domain: "app.example.com", Port: "3000"})
	if c.Labels["traefik.http.routers.app-abc.tls"] != "true" {
		t.Fatal("TLS route missing")
	}
}
func TestFailedReplacementRollsBack(t *testing.T) {
	events := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1.44")
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasPrefix(path, "/networks/"):
			io.WriteString(w, `{}`)
		case path == "/containers/app-abc/json":
			json.NewEncoder(w).Encode(container.InspectResponse{ID: "old", State: &container.State{Running: true}, Config: &container.Config{}})
		case path == "/containers/new/json":
			json.NewEncoder(w).Encode(container.InspectResponse{ID: "new", State: &container.State{Running: false}})
		case path == "/containers/create":
			events = append(events, "create")
			w.WriteHeader(201)
			io.WriteString(w, `{"Id":"new"}`)
		case r.Method == "POST" || r.Method == "DELETE":
			events = append(events, r.Method+" "+path+" "+r.URL.Query().Get("name"))
			w.WriteHeader(204)
		default:
			t.Errorf("unexpected Docker request %s %s", r.Method, path)
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	cli, err := client.New(client.WithHost(server.URL), client.WithAPIVersion("1.44"))
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	d := &DockerAdapter{client: cli}
	err = d.RunImageWithTLS(domain.Service{Id: "abc", Name: "app", Type: "github_repo"}, "", false)
	if err == nil {
		t.Fatal("crashed container accepted")
	}
	joined := strings.Join(events, "\n")
	for _, event := range []string{"POST /containers/old/rename app-abc-previous", "create", "POST /containers/old/rename app-abc", "POST /containers/old/start "} {
		if !strings.Contains(joined, event) {
			t.Fatalf("rollback incomplete, missing %q:\n%s", event, joined)
		}
	}
	if strings.Contains(joined, "DELETE /containers/old ") {
		t.Fatal("old version removed on failure")
	}
	// Ensure the fake client remains usable, without any real Docker connection.

}
