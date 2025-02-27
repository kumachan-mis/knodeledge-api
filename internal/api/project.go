package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/middleware"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type ProjectApi interface {
	HandleList(c *gin.Context)
	HandleFind(c *gin.Context)
	HandleCreate(c *gin.Context)
	HandleUpdate(c *gin.Context)
	HandleDelete(c *gin.Context)
}

type projectApi struct {
	verifier middleware.UserVerifier
	usecase  usecase.ProjectUseCase
}

func NewProjectApi(verifier middleware.UserVerifier, usecase usecase.ProjectUseCase) ProjectApi {
	return projectApi{verifier: verifier, usecase: usecase}
}

func (api projectApi) HandleList(c *gin.Context) {
	var request model.ProjectListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectListErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.User.Id)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	res, ucErr := api.usecase.ListProjects(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectListErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (api projectApi) HandleFind(c *gin.Context) {
	var request model.ProjectFindRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectFindErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.User.Id)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	res, ucErr := api.usecase.FindProject(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectFindErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ProjectFindErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (api projectApi) HandleCreate(c *gin.Context) {
	var request model.ProjectCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectCreateErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.User.Id)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	res, ucErr := api.usecase.CreateProject(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectCreateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (api projectApi) HandleUpdate(c *gin.Context) {
	var request model.ProjectUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectUpdateErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.User.Id)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	res, ucErr := api.usecase.UpdateProject(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ProjectUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (api projectApi) HandleDelete(c *gin.Context) {
	var request model.ProjectDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectDeleteErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.User.Id)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	ucErr := api.usecase.DeleteProject(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ProjectDeleteErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ProjectDeleteErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
