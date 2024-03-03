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

	res, ucErr := api.usecase.ListProjects(request)

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.ProjectListErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
		})
		return
	}

	if ucErr != nil {
		c.JSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
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

	res, ucErr := api.usecase.CreateProject(request)

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.ProjectCreateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if ucErr != nil {
		c.JSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}
