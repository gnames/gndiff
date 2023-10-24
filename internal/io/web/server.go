package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	gndiff "github.com/gnames/gndiff/pkg"
	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

var apiPath = "/api/v0/"

func Run(gnd gndiff.GNdiff, port int) {
	var err error
	e := echo.New()

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	e.GET("/api", apiInfo)
	e.GET("/api/", apiInfo)
	e.GET(strings.TrimRight(apiPath, "/"), apiInfo)
	e.GET(apiPath, apiInfo)
	e.GET(apiPath+"ping", apiPing)
	e.GET(apiPath+"version", apiVersion(gnd))
	e.POST(apiPath+"diff", apiDiffPOST(gnd))

	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	err = e.StartServer(s)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

}

func apiPing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func apiVersion(gnd gndiff.GNdiff) func(echo.Context) error {
	return func(c echo.Context) error {
		res := gnd.GetVersion()
		return c.JSON(http.StatusOK, res)
	}
}
