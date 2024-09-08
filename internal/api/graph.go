package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type GraphApi interface {
	HandleFind(c *gin.Context)
	HandleSectionalize(c *gin.Context)
}

type graphApi struct {
	usecase usecase.GraphUseCase
}

func NewGraphApi(usecase usecase.GraphUseCase) GraphApi {
	return graphApi{usecase: usecase}
}

func (api graphApi) HandleFind(c *gin.Context) {
	var request model.GraphFindRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.GraphFindErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, ucErr := api.usecase.FindGraph(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.GraphFindErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
			Section: resErr.Section,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.JSON(http.StatusNotFound, model.GraphFindErrorResponse{
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

func (api graphApi) HandleSectionalize(c *gin.Context) {
	var request model.GraphSectionalizeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.GraphSectionalizeErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, ucErr := api.usecase.SectionalizeGraph(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.GraphSectionalizeErrorResponse{
			Message:  UseCaseErrorToMessage(ucErr),
			User:     resErr.User,
			Project:  resErr.Project,
			Chapter:  resErr.Chapter,
			Sections: resErr.Sections,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		c.JSON(http.StatusBadRequest, model.GraphSectionalizeErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.JSON(http.StatusNotFound, model.GraphSectionalizeErrorResponse{
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

	c.JSON(http.StatusCreated, res)
}
