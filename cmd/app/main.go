package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/interal/api"
	"github.com/subosito/gotenv"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	gotenv.Load(fmt.Sprintf(".env.%s", env))

	router := gin.Default()
	router.GET("/", func(cxt *gin.Context) {
		cxt.JSON(200, gin.H{
			"environment": env,
			"ginVersion":  gin.Version,
			"ginMode":     gin.Mode(),
		})
	})

	router.POST("/api/hello-world", api.HelloWorldHandler)

	router.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
