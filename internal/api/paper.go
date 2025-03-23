package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/middleware"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type papersApi struct {
	verifier middleware.UserVerifier
	usecase  usecase.PaperUseCase
}

func NewPapersApi(verifier middleware.UserVerifier, usecase usecase.PaperUseCase) openapi.PapersAPI {
	return papersApi{verifier: verifier, usecase: usecase}
}

func (api papersApi) PapersFind(c *gin.Context) {
	var request openapi.PaperFindRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.PaperFindErrorResponse{
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

	res, ucErr := api.usecase.FindPaper(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.PaperFindErrorResponse{
			Message:   UseCaseErrorToMessage(ucErr),
			UserId:    resErr.UserId,
			ProjectId: resErr.ProjectId,
			ChapterId: resErr.ChapterId,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.PaperFindErrorResponse{
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

func (api papersApi) PapersUpdate(c *gin.Context) {
	var request openapi.PaperUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.PaperUpdateErrorResponse{
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

	res, ucErr := api.usecase.UpdatePaper(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, openapi.PaperUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Paper:   resErr.Paper,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.PaperUpdateErrorResponse{
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
