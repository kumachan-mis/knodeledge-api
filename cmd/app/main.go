package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/api"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/subosito/gotenv"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	mode := "development"
	if os.Getenv("GIN_MODE") == "release" {
		mode = "production"
	}

	gotenv.Load(fmt.Sprintf(".env.%v", mode))
	gotenv.Load(fmt.Sprintf(".env.%v.local", mode))

	err := db.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	router := gin.Default()
	router.SetTrustedProxies([]string{os.Getenv("TRUSTED_PROXY")})

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("ALLOW_ORIGIN")}
	router.Use(cors.New(config))

	router.GET("/", func(cxt *gin.Context) {
		cxt.JSON(200, gin.H{
			"environment": env,
			"ginVersion":  gin.Version,
			"ginMode":     gin.Mode(),
		})
	})

	router.POST("/api/hello-world", api.HelloWorldHandler)

	router.Run(":8080")
}
