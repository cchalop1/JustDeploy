package dto

type SaveGithubToken struct {
	Token string `json:"token"`
}

type GithubIsConnected struct {
	IsConnected bool `json:"isConnected"`
}
