package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type ProjectApi interface {
	HandleList(c *gin.Context)
}

type projectApi struct {
	usecase usecase.ProjectUseCase
}

func NewProjectApi(usecase usecase.ProjectUseCase) ProjectApi {
	return projectApi{usecase: usecase}
}

func (api projectApi) HandleList(c *gin.Context) {
	var request model.ProjectListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ProjectListErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, err := api.usecase.ListProjects(request)

	if err != nil && err.Code() == usecase.InvalidArgumentError {
		c.JSON(http.StatusBadRequest, model.ProjectListErrorResponse{
			Message: UseCaseErrorToMessage(err),
			User:    UseCaseErrorToResponse(err).User,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
