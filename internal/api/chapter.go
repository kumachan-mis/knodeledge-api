package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type ChapterApi interface {
	HandleList(c *gin.Context)
	HandleCreate(c *gin.Context)
	HandleUpdate(c *gin.Context)
	HandleDelete(c *gin.Context)
}

type chapterApi struct {
	usecase usecase.ChapterUseCase
}

func NewChapterApi(usecase usecase.ChapterUseCase) ChapterApi {
	return chapterApi{usecase: usecase}
}

func (api chapterApi) HandleList(c *gin.Context) {
	var request model.ChapterListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ChapterListErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, ucErr := api.usecase.ListChapters(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.ChapterListErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		c.JSON(http.StatusBadRequest, model.ChapterListErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.JSON(http.StatusNotFound, model.ChapterListErrorResponse{
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

func (api chapterApi) HandleCreate(c *gin.Context) {
	var request model.ChapterCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ChapterCreateErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, ucErr := api.usecase.CreateChapter(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.ChapterCreateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		c.JSON(http.StatusBadRequest, model.ChapterCreateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.JSON(http.StatusNotFound, model.ChapterCreateErrorResponse{
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

func (api chapterApi) HandleUpdate(c *gin.Context) {
	var request model.ChapterUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ChapterUpdateErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, ucErr := api.usecase.UpdateChapter(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.ChapterUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		c.JSON(http.StatusBadRequest, model.ChapterUpdateErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.JSON(http.StatusNotFound, model.ChapterUpdateErrorResponse{
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

func (api chapterApi) HandleDelete(c *gin.Context) {
	var request model.ChapterDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ChapterDeleteErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	ucErr := api.usecase.DeleteChapter(request)

	if ucErr != nil && ucErr.Code() == usecase.DomainValidationError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.ChapterDeleteErrorResponse{
			Message: UseCaseErrorToMessage(ucErr),
			User:    resErr.User,
			Project: resErr.Project,
			Chapter: resErr.Chapter,
		})
		return
	}

	if ucErr != nil && ucErr.Code() == usecase.NotFoundError {
		c.JSON(http.StatusNotFound, model.ChapterDeleteErrorResponse{
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

	c.JSON(http.StatusNoContent, nil)
}
