package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/middleware"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type chaptersApi struct {
	verifier middleware.UserVerifier
	usecase  usecase.ChapterUseCase
}

func NewChaptersApi(verifier middleware.UserVerifier, usecase usecase.ChapterUseCase) openapi.ChaptersAPI {
	return chaptersApi{verifier: verifier, usecase: usecase}
}

func (api chaptersApi) ChaptersList(c *gin.Context) {
	var request openapi.ChapterListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterListErrorResponse{
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

	res, ucErr := api.usecase.ListChapters(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterListErrorResponse{
			Message:   UseCaseErrorToMessage(ucErr),
			UserId:    resErr.UserId,
			ProjectId: resErr.ProjectId,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterListErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.ChapterListErrorResponse{
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

func (api chaptersApi) ChaptersCreate(c *gin.Context) {
	var request openapi.ChapterCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterCreateErrorResponse{
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

	res, ucErr := api.usecase.CreateChapter(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterCreateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterCreateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.ChapterCreateErrorResponse{
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

	c.JSON(http.StatusCreated, res)
}

func (api chaptersApi) ChaptersUpdate(c *gin.Context) {
	var request openapi.ChapterUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterUpdateErrorResponse{
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

	res, ucErr := api.usecase.UpdateChapter(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.ChapterUpdateErrorResponse{
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

func (api chaptersApi) ChaptersDelete(c *gin.Context) {
	var request openapi.ChapterDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterDeleteErrorResponse{
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

	ucErr := api.usecase.DeleteChapter(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.ChapterDeleteErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.ChapterDeleteErrorResponse{
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
