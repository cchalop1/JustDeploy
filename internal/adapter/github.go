package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/golang-jwt/jwt"
)

type GithubRepo struct {
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	Private   bool   `json:"private"`
	ID        int    `json:"id"`
	UpdatedAt string `json:"updated_at"`
}

type GithubReposResponse struct {
	Repositories []GithubRepo `json:"repositories"`
}

type GitHubAppResponse struct {
	ID            int               `json:"id"`
	Slug          string            `json:"slug"`
	NodeID        string            `json:"node_id"`
	Owner         GitHubUser        `json:"owner"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	ExternalURL   string            `json:"external_url"`
	HTMLURL       string            `json:"html_url"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
	Permissions   map[string]string `json:"permissions"`
	Events        []string          `json:"events"`
	ClientID      string            `json:"client_id"`
	ClientSecret  string            `json:"client_secret"`
	WebhookSecret string            `json:"webhook_secret"`
	Pem           string            `json:"pem"`
}

type GitHubUser struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	NodeID    string `json:"node_id"`
	URL       string `json:"url"`
	ReposURL  string `json:"repos_url"`
	EventsURL string `json:"events_url"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
}

type InstallationTokenResponse struct {
	Token string `json:"token"`
}

type GithubAdapter struct{}

func NewGithubAdapter() *GithubAdapter {
	return &GithubAdapter{}
}

func (g *GithubAdapter) GetInstallationToken(appID int, privateKey string, installationID string) (string, error) {
	jwtToken, err := generateJWT(appID, privateKey)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to get installation token, status code: %d", res.StatusCode)
	}

	var response InstallationTokenResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Token, nil
}

func (g *GithubAdapter) GetRepos(token string) ([]GithubRepo, error) {
	url := "https://api.github.com/installation/repositories?per_page=100"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// TODO: regenerate a token when is expired
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repositories, status code: %d", res.StatusCode)
	}

	var response GithubReposResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Sort repositories by updated_at in descending order
	sort.Slice(response.Repositories, func(i, j int) bool {
		return response.Repositories[i].UpdatedAt > response.Repositories[j].UpdatedAt
	})

	return response.Repositories, nil
}

func generateJWT(appID int, privateKey string) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iat": now,
		"exp": now + (10 * 60),
		"iss": appID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func (g *GithubAdapter) FinalizeGithubAppCreation(code string) (GitHubAppResponse, error) {
	apiUrl := "https://api.github.com/app-manifests/" + code + "/conversions"
	res, err := http.Post(apiUrl, "application/vnd.github+json", nil)
	if err != nil {
		return GitHubAppResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return GitHubAppResponse{}, errors.New("failed to finalize GitHub App creation")
	}

	var response GitHubAppResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return GitHubAppResponse{}, err
	}

	return response, nil
}

type Installation struct {
	ID int `json:"id"`
}

func (g *GithubAdapter) getInstallationID(appID int, privateKey string) (int, error) {
	jwtToken, err := generateJWT(appID, privateKey)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("GET", "https://api.github.com/app/installations", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to list installations, status code: %d", res.StatusCode)
	}

	var installations []Installation
	if err := json.NewDecoder(res.Body).Decode(&installations); err != nil {
		return 0, err
	}

	if len(installations) == 0 {
		return 0, fmt.Errorf("no installations found for this app")
	}

	fmt.Println(installations)
	return installations[0].ID, nil
}
