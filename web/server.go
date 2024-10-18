package web

import (
	"net/http"
	"os"

	"codeberg.org/polarhive/pasta/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var router *echo.Echo

func Serve(db *util.DB) {
	router = echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	registerCrudRoutes(router, db)

	router.GET("/robots.txt", func(c echo.Context) error {
		robotsTxt := `User-agent: *
Disallow: /`
		return c.Blob(http.StatusOK, "text/plain", []byte(robotsTxt))

	})
	router.Logger.Fatal(router.Start(":" + os.Getenv("SERVER_PORT")))
}
