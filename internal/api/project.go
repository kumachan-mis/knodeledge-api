package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type ProjectApi interface {
	HandleList(c *gin.Context)
	HandleCreate(c *gin.Context)
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
		resErr := UseCaseErrorToResponse(err)
		c.JSON(http.StatusBadRequest, model.ProjectListErrorResponse{
			Message: UseCaseErrorToMessage(err),
			User:    resErr.User,
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

func (api projectApi) HandleCreate(c *gin.Context) {
	var request model.ProjectCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ProjectCreateErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, err := api.usecase.CreateProject(request)

	if err != nil && err.Code() == usecase.InvalidArgumentError {
		resErr := UseCaseErrorToResponse(err)
		c.JSON(http.StatusBadRequest, model.ProjectCreateErrorResponse{
			Message: UseCaseErrorToMessage(err),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}
