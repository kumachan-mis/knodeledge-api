package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type PaperApi interface {
	HandleFind(c *gin.Context)
	HandleUpdate(c *gin.Context)
}

type paperApi struct {
	usecase usecase.PaperUseCase
}

func NewPaperApi(usecase usecase.PaperUseCase) PaperApi {
	return paperApi{usecase: usecase}
}

func (api paperApi) HandleFind(c *gin.Context) {
	var request model.PaperFindRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.PaperFindErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, ucErr := api.usecase.FindPaper(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.PaperFindErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.JSON(http.StatusNotFound, model.PaperFindErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
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

func (api paperApi) HandleUpdate(c *gin.Context) {
	var request model.PaperUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.PaperUpdateErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, ucErr := api.usecase.UpdatePaper(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.PaperUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Paper:   resErr.Paper,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.JSON(http.StatusNotFound, model.PaperUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
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
