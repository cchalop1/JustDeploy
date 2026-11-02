package domain

import "encoding/json"

type GithubAppSettings struct {
	InstallationID string `json:"installation_id,omitempty"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	WebhookSecret  string `json:"webhook"`
	Pem            string `json:"pem"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
}

type Settings struct {
	GithubAppSettings GithubAppSettings `json:"github_app_settings"`
	GithubToken       string            `json:"github_token"`
	AdminEmail        string            `json:"admin_email"`
	AdminPasswordHash string            `json:"admin_password_hash"`
	JwtSecret         string            `json:"jwt_secret"`
}

func (s *GithubAppSettings) UnmarshalJSON(data []byte) error {
	type plain GithubAppSettings
	var decoded plain
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	var legacy struct {
		ClientSecret  string `json:"ClientSecret"`
		WebhookSecret string `json:"WebhookSecret"`
	}
	if err := json.Unmarshal(data, &legacy); err != nil {
		return err
	}
	if decoded.ClientSecret == "" {
		decoded.ClientSecret = legacy.ClientSecret
	}
	if decoded.WebhookSecret == "" {
		decoded.WebhookSecret = legacy.WebhookSecret
	}
	*s = GithubAppSettings(decoded)
	return nil
}
