package domain

type GithubAppSettings struct {
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret`
	WebhookSecret string `json:"webhook`
	Pem           string `json:"pem"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
}

type Settings struct {
	GithubAppSettings GithubAppSettings `json:"github_app_settings"`
	GithubToken       string            `json:"github_token"`
	ApiKey            string            `json:"api_key"`
}
