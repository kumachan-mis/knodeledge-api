package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
)

type ChapterApi interface {
	HandleList(c *gin.Context)
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
		c.JSON(400, model.ChapterListErrorResponse{
			Message: JsonBindErrorToMessage(err),
		})
		return
	}

	res, ucErr := api.usecase.ListChapters(request)

	if ucErr != nil && ucErr.Code() == usecase.InvalidArgumentError {
		resErr := UseCaseErrorToResponse(ucErr)
		c.JSON(http.StatusBadRequest, model.ChapterListErrorResponse{
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

	c.JSON(http.StatusOK, res)
}