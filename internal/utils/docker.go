package utils

import "strings"

// GetDockerHubUrl returns the Docker Hub URL for a given image name
func GetDockerHubUrl(image string) string {
	// Remove tag if present
	imageParts := strings.Split(image, ":")
	imageName := imageParts[0]

	// Handle official images (no slash)
	if !strings.Contains(imageName, "/") {
		return "https://hub.docker.com/_/" + imageName
	}

	// Handle organization/user images
	return "https://hub.docker.com/r/" + imageName
}
