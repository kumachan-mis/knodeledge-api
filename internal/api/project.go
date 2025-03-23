package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/middleware"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type projectsApi struct {
	verifier middleware.UserVerifier
	usecase  usecase.ProjectUseCase
}

func NewProjectsApi(verifier middleware.UserVerifier, usecase usecase.ProjectUseCase) openapi.ProjectsAPI {
	return projectsApi{verifier: verifier, usecase: usecase}
}

func (api projectsApi) ProjectsList(c *gin.Context) {
	var request openapi.ProjectListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectListErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.UserId)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	res, ucErr := api.usecase.ListProjects(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectListErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			UserId:  resErr.UserId,
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (api projectsApi) ProjectsFind(c *gin.Context) {
	var request openapi.ProjectFindRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectFindErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.UserId)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	res, ucErr := api.usecase.FindProject(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectFindErrorResponse{
			Message:   UseCaseErrorToMessage(ucErr),
			UserId:    resErr.UserId,
			ProjectId: resErr.ProjectId,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.ProjectFindErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (api projectsApi) ProjectsCreate(c *gin.Context) {
	var request openapi.ProjectCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectCreateErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.User.Id)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	res, ucErr := api.usecase.CreateProject(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectCreateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (api projectsApi) ProjectsUpdate(c *gin.Context) {
	var request openapi.ProjectUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectUpdateErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.User.Id)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	res, ucErr := api.usecase.UpdateProject(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.ProjectUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (api projectsApi) ProjectsDelete(c *gin.Context) {
	var request openapi.ProjectDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectDeleteErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	vErr := api.verifier.Verify(c.Request.Context(), request.User.Id)

	if vErr != nil && vErr.Code() == middleware.AuthorizationError {
		c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	if vErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: MiddlewareErrorToMessage(vErr),
		})
		return
	}

	ucErr := api.usecase.DeleteProject(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ProjectDeleteErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.ProjectDeleteErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.ApplicationErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
