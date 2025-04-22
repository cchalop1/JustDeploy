package dto

type SaveGithubToken struct {
	Token string `json:"token"`
}

type GithubIsConnected struct {
	IsConnected bool `json:"isConnected"`
}

type GithubEvent struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}
