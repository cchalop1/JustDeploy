package handlers

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/utils"
	"crypto/sha256"
	"github.com/labstack/echo/v4"
	"sync"
	"time"
)

type githubState struct {
	session [32]byte
	expires time.Time
}

var githubStates = struct {
	sync.Mutex
	values map[string]githubState
}{values: make(map[string]githubState)}

func GithubStateHandler(_ *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		githubStates.Lock()
		defer githubStates.Unlock()
		for key, state := range githubStates.values {
			if time.Now().After(state.expires) {
				delete(githubStates.values, key)
			}
		}
		if len(githubStates.values) >= 100 {
			return echo.NewHTTPError(429, "too many pending GitHub connections")
		}
		value := utils.GenerateRandomPassword(48)
		githubStates.values[value] = githubState{sha256.Sum256([]byte(c.Request().Header.Get("Authorization"))), time.Now().Add(10 * time.Minute)}
		return c.JSON(200, map[string]string{"state": value})
	}
}
func consumeGithubState(c echo.Context) error {
	var req struct {
		State string `json:"state"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "invalid GitHub state")
	}
	githubStates.Lock()
	defer githubStates.Unlock()
	state, ok := githubStates.values[req.State]
	if !ok || time.Now().After(state.expires) || state.session != sha256.Sum256([]byte(c.Request().Header.Get("Authorization"))) {
		return echo.NewHTTPError(403, "invalid or expired GitHub state")
	}
	delete(githubStates.values, req.State)
	return nil
}
