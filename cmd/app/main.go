package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/api"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
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

	mode := "development"
	if os.Getenv("GIN_MODE") == "release" {
		mode = "production"
	}

	gotenv.Load(fmt.Sprintf(".env.%v", mode))
	gotenv.Load(fmt.Sprintf(".env.%v.local", mode))

	err := db.InitDatabaseClient(os.Getenv("GOOGLE_CLOUD_PROJECT_ID"))
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

	projectApi := api.NewProjectApi(projectUseCase)
	router.POST("/api/projects/list", projectApi.HandleList)
	router.POST("/api/projects/create", projectApi.HandleCreate)
	router.POST("/api/projects/find", projectApi.HandleFind)
	router.POST("/api/projects/update", projectApi.HandleUpdate)
	router.POST("/api/projects/delete", projectApi.HandleDelete)

	chapterApi := api.NewChapterApi(chapterUseCase)
	router.POST("/api/chapters/list", chapterApi.HandleList)
	router.POST("/api/chapters/create", chapterApi.HandleCreate)
	router.POST("/api/chapters/update", chapterApi.HandleUpdate)
	router.POST("/api/chapters/delete", chapterApi.HandleDelete)

	paperApi := api.NewPaperApi(paperUseCase)
	router.POST("/api/papers/find", paperApi.HandleFind)
	router.POST("/api/papers/update", paperApi.HandleUpdate)

	graphApi := api.NewGraphApi(graphUseCase)
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
