package handlers

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"regexp"
)

var deliveryPattern = regexp.MustCompile(`^[a-zA-Z0-9-]{1,100}$`)

func VerifyGithubSignature(body []byte, signature, secret string) bool {
	if len(secret) < 16 || len(signature) != 71 || signature[:7] != "sha256=" {
		return false
	}
	supplied, err := hex.DecodeString(signature[7:])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hmac.Equal(supplied, mac.Sum(nil))
}
func PostGithubEvent(ds *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := io.ReadAll(io.LimitReader(c.Request().Body, 1024*1024+1))
		if err != nil || len(body) > 1024*1024 {
			return echo.NewHTTPError(413, "invalid payload")
		}
		settings, err := ds.DatabaseAdapter.ReadSettings()
		if err != nil {
			return echo.NewHTTPError(503, "settings unavailable")
		}
		if !VerifyGithubSignature(body, c.Request().Header.Get("X-Hub-Signature-256"), settings.GithubAppSettings.WebhookSecret) {
			return echo.NewHTTPError(401, "invalid webhook signature")
		}
		if c.Request().Header.Get("X-GitHub-Event") == "ping" {
			return c.NoContent(204)
		}
		if c.Request().Header.Get("X-GitHub-Event") != "push" {
			return c.NoContent(204)
		}
		var event application.GithubPushEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return echo.NewHTTPError(400, "invalid payload")
		}
		if !application.ValidRepo(event.Repository.FullName) || event.Deleted {
			return c.NoContent(204)
		}
		if !ds.Operation.TryLock() {
			return echo.NewHTTPError(409, "another operation is in progress; retry delivery later")
		}
		defer ds.Operation.Unlock()
		delivery := c.Request().Header.Get("X-GitHub-Delivery")
		if !deliveryPattern.MatchString(delivery) {
			return echo.NewHTTPError(400, "invalid delivery ID")
		}
		claimed, err := ds.DatabaseAdapter.ClaimDelivery(delivery)
		if err != nil {
			return err
		}
		if !claimed {
			return c.NoContent(204)
		}
		if err := application.ManageVerifiedGithubEvent(ds, event); err != nil {
			if releaseErr := ds.DatabaseAdapter.ReleaseDelivery(delivery); releaseErr != nil {
				return releaseErr
			}
			return echo.NewHTTPError(502, "deployment failed; delivery may be retried")
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Processed"})
	}
}
