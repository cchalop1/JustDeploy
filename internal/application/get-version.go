package application

import "cchalop1.com/deploy/internal/api/dto"

var Version string

func GetVersion() dto.VersionDto {
	return dto.VersionDto{
		Version:   Version,
		GithubUrl: "https://github.com/cchalop1/JustDeploy/releases/tag/" + Version,
	}
}
