package api

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/domain"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Application struct{ Echo *echo.Echo }

func NewApplication(port string) *Application {
	e := echo.New()
	e.HideBanner = true
	// Never trust arbitrary forwarded headers for client identity or TLS.
	e.IPExtractor = echo.ExtractIPDirect()
	e.Server.ReadHeaderTimeout = 10 * time.Second
	e.Server.ReadTimeout = 30 * time.Second
	e.Server.WriteTimeout = 40 * time.Minute
	e.Server.IdleTimeout = 60 * time.Second
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code, message := 500, "internal server error"
		if errors.Is(err, domain.ErrNotFound) {
			code = 404
			message = domain.ErrNotFound.Error()
		}
		if errors.Is(err, domain.ErrConflict) {
			code = 409
			message = err.Error()
		}
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
			if s, ok := he.Message.(string); ok {
				message = s
			}
		}
		if !c.Response().Committed {
			_ = c.JSON(code, dto.ResponseApi{Message: message})
		}
	}
	e.Use(middleware.Recover(), middleware.BodyLimit("1M"), middleware.Secure())
	e.Use(requireSecureTransport)
	origins := []string{"http://localhost:" + port, "http://localhost:5173", "http://127.0.0.1:5173"}
	if origin := os.Getenv("JUSTDEPLOY_PUBLIC_URL"); origin != "" {
		origins = append(origins, origin)
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: origins, AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE}, AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization}}))
	return &Application{Echo: e}
}
func requireSecureTransport(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		host, _, err := net.SplitHostPort(c.Request().RemoteAddr)
		if err != nil {
			return echo.NewHTTPError(403, "invalid peer")
		}
		ip := net.ParseIP(host)
		secure := c.Request().TLS != nil
		if !secure && c.Request().Header.Get("X-Forwarded-Proto") == "https" {
			for _, cidr := range strings.Split(os.Getenv("JUSTDEPLOY_TRUSTED_PROXIES"), ",") {
				_, network, err := net.ParseCIDR(strings.TrimSpace(cidr))
				if err == nil && network.Contains(ip) {
					secure = true
					break
				}
			}
		}
		if !secure && (ip == nil || !ip.IsLoopback()) {
			return echo.NewHTTPError(403, "HTTPS is required; configure TLS or a trusted HTTPS reverse proxy")
		}
		return next(c)
	}
}
func (app *Application) StartServer(port string) {
	addr := os.Getenv("JUSTDEPLOY_LISTEN_ADDR")
	if addr == "" {
		addr = "127.0.0.1"
	}
	address := net.JoinHostPort(addr, port)
	cert, key := os.Getenv("JUSTDEPLOY_TLS_CERT"), os.Getenv("JUSTDEPLOY_TLS_KEY")
	var err error
	if cert != "" && key != "" {
		err = app.Echo.StartTLS(address, cert, key)
	} else {
		err = app.Echo.Start(address)
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
