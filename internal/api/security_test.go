package api

import (
	"bytes"
	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/handlers"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/domain"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func testApplication(t *testing.T) (*Application, *service.DeployService) {
	t.Helper()
	dir := t.TempDir()
	oldDir, oldDB := internal.JUSTDEPLOY_FOLDER, internal.DATABASE_SQLITE_PATH
	internal.JUSTDEPLOY_FOLDER = dir
	internal.DATABASE_SQLITE_PATH = filepath.Join(dir, "database.db")
	db := adapter.NewDatabaseAdapter()
	db.Init()
	if err := db.SaveServer(domain.Server{Id: "server"}); err != nil {
		t.Fatal(err)
	}
	ds := &service.DeployService{DatabaseAdapter: db}
	app := NewApplication("5915")
	CreateRoutes(app, ds)
	t.Cleanup(func() { db.Close(); internal.JUSTDEPLOY_FOLDER = oldDir; internal.DATABASE_SQLITE_PATH = oldDB })
	return app, ds
}
func request(app *Application, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	app.Echo.ServeHTTP(rec, req)
	return rec
}
func TestSetupRequiresValidCredentialsAndCannotOverwriteAdmin(t *testing.T) {
	app, _ := testApplication(t)
	body := map[string]string{"email": "admin@example.com", "password": "secure-password-123", "domain": "example.com"}
	body["email"] = ""
	if rec := request(app, "/api/setup", body, nil); rec.Code != 400 {
		t.Fatalf("empty email accepted: %d", rec.Code)
	}
	body["email"] = "admin@example.com"
	body["password"] = "short"
	if rec := request(app, "/api/setup", body, nil); rec.Code != 400 {
		t.Fatalf("short password accepted: %d", rec.Code)
	}
	body["password"] = "secure-password-123"
	if rec := request(app, "/api/setup", body, nil); rec.Code != 200 {
		t.Fatalf("setup: %d %s", rec.Code, rec.Body.String())
	}
	if rec := request(app, "/api/setup", body, nil); rec.Code != 409 {
		t.Fatalf("admin overwritten: %d", rec.Code)
	}
	if rec := request(app, "/api/deploy", nil, nil); rec.Code != 401 {
		t.Fatalf("deploy route unprotected: %d", rec.Code)
	}
}
func TestConcurrentSetupOnlyOneRequestCreatesAdmin(t *testing.T) {
	app, ds := testApplication(t)
	var wg sync.WaitGroup
	results := make(chan int, 2)
	for _, email := range []string{"one@example.com", "two@example.com"} {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			results <- request(app, "/api/setup", map[string]string{
				"email": email, "password": "secure-password-123",
			}, nil).Code
		}(email)
	}
	wg.Wait()
	close(results)
	statuses := map[int]int{}
	for code := range results {
		statuses[code]++
	}
	if statuses[http.StatusOK] != 1 || statuses[http.StatusConflict] != 1 {
		t.Fatalf("expected one successful setup and one conflict: %v", statuses)
	}
	if ds.DatabaseAdapter.GetSettings().AdminEmail == "" {
		t.Fatal("admin was not persisted")
	}
}

func TestConcurrentAdminCreationOnlyOneWins(t *testing.T) {
	_, ds := testApplication(t)
	var wg sync.WaitGroup
	successes := make(chan bool, 2)
	for _, email := range []string{"one@example.com", "two@example.com"} {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			successes <- ds.DatabaseAdapter.InitializeAdmin(email, "hash", "secret", "example.com") == nil
		}(email)
	}
	wg.Wait()
	close(successes)
	count := 0
	for success := range successes {
		if success {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("%d admin creations succeeded", count)
	}
}
func TestWebhookRejectsForgeryAndDeduplicates(t *testing.T) {
	app, ds := testApplication(t)
	secret := strings.Repeat("s", 32)
	settings := ds.DatabaseAdapter.GetSettings()
	settings.GithubAppSettings.WebhookSecret = secret
	if err := ds.DatabaseAdapter.SaveSettings(settings); err != nil {
		t.Fatal(err)
	}
	body := map[string]any{"ref": "refs/heads/main", "repository": map[string]string{"full_name": "owner/unknown"}}
	if rec := request(app, "/api/github/events", body, nil); rec.Code != 401 {
		t.Fatalf("forged webhook: %d", rec.Code)
	}
	data, _ := json.Marshal(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(data)
	headers := map[string]string{"X-Hub-Signature-256": "sha256=" + hex.EncodeToString(mac.Sum(nil)), "X-GitHub-Event": "push", "X-GitHub-Delivery": "delivery-1"}
	if rec := request(app, "/api/github/events", body, headers); rec.Code != 200 {
		t.Fatalf("verified event: %d %s", rec.Code, rec.Body.String())
	}
	if rec := request(app, "/api/github/events", body, headers); rec.Code != 204 {
		t.Fatalf("duplicate delivery: %d", rec.Code)
	}
	ds.Operation.Lock()
	if rec := request(app, "/api/github/events", body, headers); rec.Code != 409 {
		t.Fatalf("concurrent job accepted: %d", rec.Code)
	}
	ds.Operation.Unlock()
	if handlers.VerifyGithubSignature(append(data, ' '), headers["X-Hub-Signature-256"], secret) {
		t.Fatal("modified payload authenticated")
	}
}
func TestRemoteHTTPRejectedAndForwardingNotBlindlyTrusted(t *testing.T) {
	app, _ := testApplication(t)
	for _, forwarded := range []string{"", "https"} {
		req := httptest.NewRequest("GET", "http://example.com/api/info", nil)
		req.RemoteAddr = "192.0.2.5:1234"
		req.Header.Set("X-Forwarded-Proto", forwarded)
		rec := httptest.NewRecorder()
		app.Echo.ServeHTTP(rec, req)
		if rec.Code != 403 {
			t.Fatalf("remote HTTP accepted: %d", rec.Code)
		}
	}
}

func TestServiceUpdatesPreserveIdentityAndRejectStaleDrafts(t *testing.T) {
	app, ds := testApplication(t)
	settings := ds.DatabaseAdapter.GetSettings()
	settings.AdminEmail = "admin@example.com"
	settings.JwtSecret = strings.Repeat("s", 64)
	if err := ds.DatabaseAdapter.SaveSettings(settings); err != nil {
		t.Fatal(err)
	}
	token, err := application.GenerateJWT(settings.AdminEmail, settings.JwtSecret)
	if err != nil {
		t.Fatal(err)
	}
	s := domain.Service{Id: "service", Name: "app", Type: "github_repo", CurrentPath: "/trusted/workspace", ImageName: "trusted", ServiceRuntime: domain.ServiceRuntime{ProjectID: "group"}}
	if err := ds.DatabaseAdapter.SaveService(s); err != nil {
		t.Fatal(err)
	}
	s.CurrentPath = "/etc"
	s.ImageName = "attacker/image"
	s.Name = "changed"
	s.Envs = nil
	update := func() *httptest.ResponseRecorder {
		body, _ := json.Marshal(s)
		req := httptest.NewRequest(http.MethodPut, "/api/service", bytes.NewReader(body))
		req.RemoteAddr = "127.0.0.1:1234"
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		app.Echo.ServeHTTP(rec, req)
		return rec
	}
	if rec := update(); rec.Code != 200 {
		t.Fatalf("update: %d %s", rec.Code, rec.Body.String())
	}
	saved, err := ds.DatabaseAdapter.GetServiceById(s.Id)
	if err != nil {
		t.Fatal(err)
	}
	if saved.CurrentPath != "/trusted/workspace" || saved.ImageName != "trusted" || saved.Name != "app" {
		t.Fatal("immutable fields overwritten")
	}
	if rec := update(); rec.Code != 409 {
		t.Fatalf("stale update accepted: %d", rec.Code)
	}
}
