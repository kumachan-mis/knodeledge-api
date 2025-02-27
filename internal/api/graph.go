package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/middleware"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type GraphApi interface {
	HandleFind(c *gin.Context)
	HandleUpdate(c *gin.Context)
	HandleDelete(c *gin.Context)
	HandleSectionalize(c *gin.Context)
}

type graphApi struct {
	verifier middleware.UserVerifier
	usecase  usecase.GraphUseCase
}

func NewGraphApi(verifier middleware.UserVerifier, usecase usecase.GraphUseCase) GraphApi {
	return graphApi{verifier: verifier, usecase: usecase}
}

func (api graphApi) HandleFind(c *gin.Context) {
	var request model.GraphFindRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.GraphFindErrorResponse{
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

	res, ucErr := api.usecase.FindGraph(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.GraphFindErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
			Section: resErr.Section,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, model.GraphFindErrorResponse{
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

func (api graphApi) HandleUpdate(c *gin.Context) {
	var request model.GraphUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.GraphUpdateErrorResponse{
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

	res, ucErr := api.usecase.UpdateGraph(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.GraphUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
			Graph:   resErr.Graph,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, model.GraphUpdateErrorResponse{
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

func (api graphApi) HandleDelete(c *gin.Context) {
	var request model.GraphDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.GraphDeleteErrorResponse{
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

	ucErr := api.usecase.DeleteGraph(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.GraphDeleteErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
			Section: resErr.Section,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, model.GraphDeleteErrorResponse{
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

func (api graphApi) HandleSectionalize(c *gin.Context) {
	var request model.GraphSectionalizeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.GraphSectionalizeErrorResponse{
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

	res, ucErr := api.usecase.SectionalizeGraph(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.GraphSectionalizeErrorResponse{
			Message:  UseCaseErrorToMessage(ucErr),
			User:     resErr.User,
			Project:  resErr.Project,
			Chapter:  resErr.Chapter,
			Sections: resErr.Sections,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.GraphSectionalizeErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, model.GraphSectionalizeErrorResponse{
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

	c.JSON(http.StatusCreated, res)
}
