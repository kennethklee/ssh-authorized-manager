package routes

import (
	"fmt"
	"io/fs"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

func Register(router *echo.Echo) {
	// Serve static files
	if os.Getenv("APP_ENV") == "production" {
		fileSystem := echo.MustSubFS(router.Filesystem, "static")
		router.GET("/*", spaFS(fileSystem))
	} else {
		fmt.Println("  - Development mode -- proxy frontend to http://web:5173")
		router.GET("/*", reverseProxy("http://web:5173"))
	}

	router.GET("/api/me", func(c echo.Context) error {
		if c.Get(apis.ContextAdminKey) != nil {
			return c.JSON(200, c.Get("admin"))
		} else {
			return c.JSON(200, c.Get("user"))
		}
	})
}

func spaFS(fileSystem fs.FS) echo.HandlerFunc {
	staticHandler := echo.StaticDirectoryHandler(fileSystem, true)
	return func(c echo.Context) error {
		// If err is echo.ErrNotFound, serve root
		err := staticHandler(c)
		if err == echo.ErrNotFound {
			return c.FileFS(".", fileSystem) // SPA time!
		}

		return err
	}
}

// Simple reverse proxy for development mode (don't use in production)
func reverseProxy(host string) echo.HandlerFunc {
	remote, _ := url.Parse(host)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	return func(c echo.Context) error {
		proxy.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
