package application

import (
	"encoding/json"
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
)

var Version string

func GetVersion() dto.VersionDto {
	return dto.VersionDto{
		Version:   Version,
		GithubUrl: "https://github.com/cchalop1/JustDeploy/releases/tag/" + Version,
	}
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func CheckIsNewVersion() bool {
	resp, err := http.Get("https://api.github.com/repos/cchalop1/JustDeploy/releases/latest")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false
	}
	return Version != release.TagName
}
