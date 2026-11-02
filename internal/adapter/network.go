package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type NetworkAdapter struct{}

func NewNetworkAdapter() *NetworkAdapter {
	return &NetworkAdapter{}
}

// GetServerURL renvoie l'URL complète incluant le port.
func (n *NetworkAdapter) GetServerURL(port string) (string, error) {
	if url := os.Getenv("JUSTDEPLOY_PUBLIC_URL"); url != "" {
		return url, nil
	}
	return "http://127.0.0.1:" + port, nil
}

// ipifyResponse correspond au JSON retourné par ipify.org
type ipifyResponse struct {
	IP string `json:"ip"`
}

// GetCurrentIP fait une requête à ipify pour récupérer l'IP publique au format JSON.
func (n *NetworkAdapter) GetCurrentIP() (string, error) {
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", fmt.Errorf("échec de la requête ipify: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return "", fmt.Errorf("lecture de la réponse ipify: %w", err)
	}

	var result ipifyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("décodage JSON ipify: %w", err)
	}
	if result.IP == "" {
		return "", fmt.Errorf("ipify a renvoyé un champ vide")
	}

	return result.IP, nil
}
