package main

import (
	"apollo-sample-tracker/api"
	"apollo-sample-tracker/lib/database"
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"strings"

	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed database/schema.sql
var db_schema string

func main() {
	r := gin.Default()
	godotenv.Load()
	_, devmode := os.LookupEnv("DEV")
	if devmode {
		println("Running frontend proxy")
		r.Any("/app/*proxypath", proxy)
	} else {
		println("Serving frontend static files")
		r.LoadHTMLGlob("./client/*.html")
		r.GET("/app/*path", func(c *gin.Context) {
			path := c.Param("path")
			// Check if the first segment is "assets"
			segments := strings.SplitN(path, "/", 3)
			if len(segments) > 1 && segments[1] == "assets" {
				c.File("./client" + path)
				return
			}
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})
	}
	database.Connect(context.Background(), db_schema)
	api.Routes(r)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/app")
	})

	if runtime.GOOS == "windows" && devmode {
		r.Run("localhost:80")
	} else {
		r.Run("0.0.0.0:80")
	}
}

func proxy(c *gin.Context) {
	remote, err := url.Parse("http://" + c.Request.Host + ":5173/")
	if err != nil {
		println("Could not resolve Proxy URL")
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(c.Writer, c.Request)
}
