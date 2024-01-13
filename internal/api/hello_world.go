package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
)

func HelloWorldHandler(cxt *gin.Context) {
	var request model.HelloWorldRequest
	if err := cxt.ShouldBindJSON(&request); err != nil {
		cxt.JSON(400, model.ApplicationErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if request.Name == "" {
		cxt.JSON(200, model.HelloWorldResponse{
			Message: "Hello World!",
		})
		return
	}

	message, err := db.GetHelloWorld(request.Name)
	if err != nil {
		message, err = db.SetHelloWorld(request.Name, fmt.Sprintf("Hello, %s!", request.Name))
	}

	if err != nil {
		cxt.JSON(500, model.ApplicationErrorResponse{
			Message: err.Error(),
		})
	}

	cxt.JSON(200, model.HelloWorldResponse{
		Message: message,
	})
}
