package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/domain"
)

func TestOnboardingWithoutDomain(t *testing.T) {
	app, ds := testApplication(t)
	rec := request(app, "/api/setup", map[string]string{
		"email": "admin@example.com", "password": "secure-password-123",
	}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("domain-free setup failed: %d %s", rec.Code, rec.Body.String())
	}
	var session struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatal(err)
	}
	if session.Token == "" {
		t.Fatal("setup did not return a session")
	}
	if ds.DatabaseAdapter.GetServer().Domain != "" {
		t.Fatal("setup unexpectedly configured a domain")
	}
	if rec := request(app, "/api/deploy/service/missing", nil, nil); rec.Code != http.StatusUnauthorized {
		t.Fatalf("unprotected deployment: %d", rec.Code)
	}
	headers := map[string]string{"Authorization": "Bearer " + session.Token}
	if rec := request(app, "/api/deploy/service/missing", nil, headers); rec.Code != http.StatusNotFound {
		t.Fatalf("missing service: %d %s", rec.Code, rec.Body.String())
	}
	ds.Operation.Lock()
	rec = request(app, "/api/deploy/service/missing", nil, headers)
	ds.Operation.Unlock()
	if rec.Code != http.StatusConflict {
		t.Fatalf("concurrent deployment accepted: %d", rec.Code)
	}
}

func TestOnboardingDeploymentKeepsOtherGroupsAndOrdersDependencies(t *testing.T) {
	for _, failDependency := range []bool{false, true} {
		t.Run(map[bool]string{false: "running services retained", true: "dependency failure stops application"}[failDependency], func(t *testing.T) {
			app, ds := testApplication(t)
			if err := ds.DatabaseAdapter.InitializeAdmin("admin@example.com", "hash", strings.Repeat("s", 64), ""); err != nil {
				t.Fatal(err)
			}
			token, err := application.GenerateJWT("admin@example.com", strings.Repeat("s", 64))
			if err != nil {
				t.Fatal(err)
			}
			// Save the app first: deployment must still inspect its database first.
			services := []domain.Service{
				{Id: "app", Name: "web", Type: "github_repo", Status: "Running", ServiceRuntime: domain.ServiceRuntime{ProjectID: "project"}},
				{Id: "other", Name: "unrelated", Type: "database", Status: "ready_to_deploy", ServiceRuntime: domain.ServiceRuntime{ProjectID: "other"}},
				{Id: "db", Name: "postgres", Type: "database", Status: "Running", ServiceRuntime: domain.ServiceRuntime{ProjectID: "project"}},
			}
			if err := ds.DatabaseAdapter.SaveServices(services); err != nil {
				t.Fatal(err)
			}
			var visited []string
			docker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/_ping" {
					w.Header().Set("API-Version", "1.44")
					w.Write([]byte("OK"))
					return
				}
				path := strings.TrimPrefix(r.URL.Path, "/v1.44")
				visited = append(visited, r.Method+" "+path)
				if failDependency && strings.Contains(path, "postgres-db") {
					w.WriteHeader(500)
					w.Write([]byte(`{"message":"dependency unavailable"}`))
					return
				}
				if r.Method != "GET" || !strings.HasSuffix(path, "/json") {
					t.Errorf("unexpected Docker mutation: %s %s", r.Method, path)
					w.WriteHeader(500)
					return
				}
				w.Write([]byte(`{"Id":"existing","State":{"Running":true}}`))
			}))
			defer docker.Close()
			t.Setenv("DOCKER_HOST", docker.URL)
			t.Setenv("DOCKER_API_VERSION", "1.44")
			ds.DockerAdapter = adapter.NewDockerAdapter()
			if err := ds.DockerAdapter.ConnectClient(); err != nil {
				t.Fatal(err)
			}
			rec := request(app, "/api/deploy/service/app", nil, map[string]string{"Authorization": "Bearer " + token})
			want := []string{"GET /containers/postgres-db/json", "GET /containers/web-app/json"}
			if failDependency {
				want = want[:1]
				if rec.Code < 400 {
					t.Fatalf("dependency error ignored: %d", rec.Code)
				}
			} else if rec.Code != http.StatusOK {
				t.Fatalf("deploy failed: %d %s", rec.Code, rec.Body.String())
			}
			if !reflect.DeepEqual(visited, want) {
				t.Fatalf("unexpected deployment order/scope: got %v, want %v", visited, want)
			}
		})
	}
}

func TestFirstPrivateDatabaseCreationAndDeployment(t *testing.T) {
	app, ds := testApplication(t)
	secret := strings.Repeat("s", 64)
	if err := ds.DatabaseAdapter.InitializeAdmin("admin@example.com", "hash", secret, ""); err != nil {
		t.Fatal(err)
	}
	token, err := application.GenerateJWT("admin@example.com", secret)
	if err != nil {
		t.Fatal(err)
	}
	headers := map[string]string{"Authorization": "Bearer " + token}
	createdContainer, startedContainer, createdVolume := false, false, false
	docker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := strings.TrimPrefix(r.URL.Path, "/v1.44")
		switch {
		case path == "/_ping":
			w.Header().Set("API-Version", "1.44")
			w.Write([]byte("OK"))
		case r.Method == "GET" && strings.HasPrefix(path, "/images/"):
			w.Write([]byte(`{"Id":"postgres-image"}`))
		case r.Method == "GET" && strings.HasPrefix(path, "/networks/"):
			w.Write([]byte(`{"Id":"private-network"}`))
		case r.Method == "POST" && path == "/volumes/create":
			createdVolume = true
			w.WriteHeader(201)
			w.Write([]byte(`{"Name":"persistent-volume"}`))
		case r.Method == "POST" && path == "/containers/create":
			var config struct {
				Labels     map[string]string
				HostConfig struct {
					PortBindings map[string]any
					Mounts       []struct{ Type, Target string }
				}
			}
			if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
				t.Error(err)
			}
			if config.Labels["traefik.enable"] != "false" || len(config.HostConfig.PortBindings) != 0 {
				t.Error("private database exposed publicly")
			}
			if len(config.HostConfig.Mounts) != 1 || config.HostConfig.Mounts[0].Type != "volume" || config.HostConfig.Mounts[0].Target != "/var/lib/postgresql/data" {
				t.Errorf("database persistence missing: %+v", config.HostConfig.Mounts)
			}
			createdContainer = true
			w.WriteHeader(201)
			w.Write([]byte(`{"Id":"new-db"}`))
		case r.Method == "POST" && path == "/containers/new-db/start":
			startedContainer = true
			w.WriteHeader(204)
		case r.Method == "GET" && path == "/containers/new-db/json":
			w.Write([]byte(`{"Id":"new-db","State":{"Running":true,"Health":{"Status":"healthy"}}}`))
		case r.Method == "GET" && strings.HasPrefix(path, "/containers/"):
			w.WriteHeader(404)
			w.Write([]byte(`{"message":"No such container"}`))
		default:
			t.Errorf("unexpected Docker call: %s %s", r.Method, path)
			w.WriteHeader(500)
		}
	}))
	defer docker.Close()
	t.Setenv("DOCKER_HOST", docker.URL)
	t.Setenv("DOCKER_API_VERSION", "1.44")
	ds.DockerAdapter = adapter.NewDockerAdapter()
	if err := ds.DockerAdapter.ConnectClient(); err != nil {
		t.Fatal(err)
	}
	rec := request(app, "/api/database/create", map[string]string{"databaseName": "Postgres"}, headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	var created domain.Service
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.Id == "" || created.Status != "ready_to_deploy" || created.ExposeSettings.IsExposed {
		t.Fatalf("invalid created database: %+v", created)
	}
	rec = request(app, "/api/deploy/service/"+created.Id, nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("deploy: %d %s", rec.Code, rec.Body.String())
	}
	saved, err := ds.DatabaseAdapter.GetServiceById(created.Id)
	if err != nil {
		t.Fatal(err)
	}
	if saved.Status != "Running" || saved.Url != "" || !createdContainer || !startedContainer || !createdVolume {
		t.Fatal("private database did not complete deployment with persistent storage")
	}
}
