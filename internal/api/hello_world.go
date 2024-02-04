package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type HelloWorldApi interface {
	HandleHelloWorld(cxt *gin.Context)
}

type helloWorldApi struct {
	usecase usecase.HelloWorldUseCase
}

func NewHelloWorldApi(usecase usecase.HelloWorldUseCase) HelloWorldApi {
	return helloWorldApi{usecase: usecase}
}

func (a helloWorldApi) HandleHelloWorld(cxt *gin.Context) {
	var request model.HelloWorldRequest
	if err := cxt.ShouldBindJSON(&request); err != nil {
		cxt.JSON(400, model.ApplicationErrorResponse{
			Message: err.Error(),
		})
		return
	}

	message, err := a.usecase.UseHelloWorld(request.Name)

	if err != nil {
		cxt.JSON(500, model.ApplicationErrorResponse{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(200, model.HelloWorldResponse{
		Message: message,
	})
}
