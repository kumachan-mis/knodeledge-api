package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type ChapterUseCase interface {
	ListChapters(req model.ChapterListRequest) (
		*model.ChapterListResponse, *Error[model.ChapterListErrorResponse])
}

type chapterUseCase struct {
	service service.ChapterService
}

func NewChapterUseCase(service service.ChapterService) ChapterUseCase {
	return chapterUseCase{service: service}
}

func (uc chapterUseCase) ListChapters(req model.ChapterListRequest) (
	*model.ChapterListResponse, *Error[model.ChapterListErrorResponse]) {
	uid, uidErr := domain.NewUserIdObject(req.User.Id)
	pid, pidErr := domain.NewProjectIdObject(req.Project.Id)

	uidMsg := ""
	if uidErr != nil {
		uidMsg = uidErr.Error()
	}
	pidMsg := ""
	if pidErr != nil {
		pidMsg = pidErr.Error()
	}

	if uidErr != nil || pidErr != nil {
		return nil, NewModelBasedError(
			InvalidArgumentError,
			model.ChapterListErrorResponse{
				User:    model.UserOnlyIdError{Id: uidMsg},
				Project: model.ProjectOnlyIdError{Id: pidMsg},
			},
		)
	}

	entities, sErr := uc.service.ListChapters(*uid, *pid)
	if sErr != nil {
		return nil, NewMessageBasedError[model.ChapterListErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	chapters := make([]model.Chapter, len(entities))
	i := 0
	for _, entity := range entities {
		chapters[i] = model.Chapter{
			Id:       entity.Id().Value(),
			Name:     entity.Name().Value(),
			Number:   int32(entity.Number().Value()),
			Sections: []model.Section{},
		}
		i++
	}

	return &model.ChapterListResponse{Chapters: chapters}, nil
}
