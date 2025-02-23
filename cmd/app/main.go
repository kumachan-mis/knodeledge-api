package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/api"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/middleware"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	"github.com/subosito/gotenv"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	gotenv.Load(fmt.Sprintf(".env.%v", env))
	gotenv.Load(fmt.Sprintf(".env.%v.local", env))

	err := db.InitDatabaseClient(os.Getenv("GOOGLE_CLOUD_PROJECT_ID"))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	router := gin.Default()
	router.SetTrustedProxies([]string{os.Getenv("TRUSTED_PROXY")})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{os.Getenv("ALLOW_ORIGIN")}
	router.Use(cors.New(corsConfig))

	auth0Config := middleware.Auth0JwTConfig{
		Domain:   os.Getenv("AUTH0_DOMAIN"),
		Audience: os.Getenv("AUTH0_AUDIENCE"),
	}
	router.Use(middleware.Auth0JWT(auth0Config))

	router.GET("/", func(cxt *gin.Context) {
		cxt.JSON(http.StatusOK, gin.H{
			"environment": env,
			"ginVersion":  gin.Version,
			"ginMode":     gin.Mode(),
		})
	})

	client := db.FirestoreClient()
	if client == nil {
		log.Fatalf("Failed to get firestore client")
	}

	projectRepository := repository.NewProjectRepository(*client)
	chapterRepository := repository.NewChapterRepository(*client)
	paperRepository := repository.NewPaperRepository(*client)
	graphRepository := repository.NewGraphRepository(*client)

	projectService := service.NewProjectService(projectRepository)
	chapterService := service.NewChapterService(chapterRepository, paperRepository)
	paperService := service.NewPaperService(paperRepository)
	graphService := service.NewGraphService(graphRepository, chapterRepository)

	projectUseCase := usecase.NewProjectUseCase(projectService)
	chapterUseCase := usecase.NewChapterUseCase(chapterService)
	paperUseCase := usecase.NewPaperUseCase(paperService)
	graphUseCase := usecase.NewGraphUseCase(graphService)

	userVerifier := middleware.NewUserVerifier()

	projectApi := api.NewProjectApi(userVerifier, projectUseCase)
	router.POST("/api/projects/list", projectApi.HandleList)
	router.POST("/api/projects/create", projectApi.HandleCreate)
	router.POST("/api/projects/find", projectApi.HandleFind)
	router.POST("/api/projects/update", projectApi.HandleUpdate)
	router.POST("/api/projects/delete", projectApi.HandleDelete)

	chapterApi := api.NewChapterApi(userVerifier, chapterUseCase)
	router.POST("/api/chapters/list", chapterApi.HandleList)
	router.POST("/api/chapters/create", chapterApi.HandleCreate)
	router.POST("/api/chapters/update", chapterApi.HandleUpdate)
	router.POST("/api/chapters/delete", chapterApi.HandleDelete)

	paperApi := api.NewPaperApi(userVerifier, paperUseCase)
	router.POST("/api/papers/find", paperApi.HandleFind)
	router.POST("/api/papers/update", paperApi.HandleUpdate)

	graphApi := api.NewGraphApi(userVerifier, graphUseCase)
	router.POST("/api/graphs/find", graphApi.HandleFind)
	router.POST("/api/graphs/update", graphApi.HandleUpdate)
	router.POST("/api/graphs/delete", graphApi.HandleDelete)
	router.POST("/api/graphs/sectionalize", graphApi.HandleSectionalize)

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to run gin server: %v", err)
	}

	err = db.FinalizeDatabaseClient()
	if err != nil {
		log.Fatalf("Failed to finalize database: %v", err)
	}
}
