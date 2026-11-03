package application

import (
	"cchalop1.com/deploy/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"testing"
	"time"
)

func TestJWTRequiresConfiguredSecretIdentityAndExpiry(t *testing.T) {
	secret := strings.Repeat("s", 64)
	token, err := GenerateJWT("admin@example.com", secret)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateJWT(token, secret); err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateJWT(token, ""); err == nil {
		t.Fatal("empty key accepted")
	}
	for _, claims := range []jwt.MapClaims{{"email": "admin@example.com"}, {"email": "admin@example.com", "exp": time.Now().Add(-time.Hour).Unix()}, {"email": "", "exp": time.Now().Add(time.Hour).Unix()}} {
		token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
		if _, err := ValidateJWT(token, secret); err == nil {
			t.Fatalf("invalid claims accepted: %v", claims)
		}
	}
}
func TestPrivateServiceDomainAndInputValidation(t *testing.T) {
	s := domain.Service{Id: "abc", Name: "app", Type: "github_repo"}
	if domain := buildServiceDomain(s, "example.com"); domain != "" {
		t.Fatal("private service has a route")
	}
	for _, repo := range []string{"invalid", "../repo", "owner/../../etc", "https://github.com/owner/repo", "owner/repo?token=secret"} {
		if ValidRepo(repo) {
			t.Fatalf("invalid repository accepted: %s", repo)
		}
	}
	if !ValidRepo("owner/repo-name") {
		t.Fatal("valid repository rejected")
	}
}
