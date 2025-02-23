package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/testutil"
	mock_repository "github.com/kumachan-mis/knodeledge-api/mock/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListChaptersValidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	maxLengthChapterName := testutil.RandomString(100)
	maxLengthSectionName := testutil.RandomString(100)

	r := mock_repository.NewMockChapterRepository(ctrl)
	r.EXPECT().
		FetchChapters(testutil.ReadOnlyUserId(), "0000000000000001").
		Return(map[string]record.ChapterEntry{
			"1000000000000003": {
				Name:   "Chapter 3",
				Number: 3,
				Sections: []record.SectionEntry{
					{
						Id:        "2000000000000003",
						Name:      "Section 1",
						UserId:    testutil.ReadOnlyUserId(),
						CreatedAt: testutil.Date().Add(-3 * time.Hour),
						UpdatedAt: testutil.Date().Add(-3 * time.Hour),
					},
				},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-3 * time.Hour),
				UpdatedAt: testutil.Date().Add(-3 * time.Hour),
			},
			"1000000000000001": {
				Name:   "Chapter 1",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "2000000000000001",
						Name:      "Section 1",
						UserId:    testutil.ReadOnlyUserId(),
						CreatedAt: testutil.Date().Add(-2 * time.Hour),
						UpdatedAt: testutil.Date().Add(-2 * time.Hour),
					},
					{
						Id:        "2000000000000002",
						Name:      "Section 2",
						UserId:    testutil.ReadOnlyUserId(),
						CreatedAt: testutil.Date().Add(-2 * time.Hour),
						UpdatedAt: testutil.Date().Add(-2 * time.Hour),
					},
				},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-2 * time.Hour),
				UpdatedAt: testutil.Date().Add(-2 * time.Hour),
			},
			"1000000000000004": {
				Name:   maxLengthChapterName,
				Number: 4,
				Sections: []record.SectionEntry{
					{
						Id:        "2000000000000004",
						Name:      maxLengthSectionName,
						UserId:    testutil.ReadOnlyUserId(),
						CreatedAt: testutil.Date().Add(-4 * time.Hour),
						UpdatedAt: testutil.Date().Add(-4 * time.Hour),
					},
				},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-4 * time.Hour),
				UpdatedAt: testutil.Date().Add(-4 * time.Hour),
			},
			"1000000000000002": {
				Name:      "Chapter 2",
				Number:    2,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-1 * time.Hour),
				UpdatedAt: testutil.Date().Add(-1 * time.Hour),
			},
		}, nil)

	pr := mock_repository.NewMockPaperRepository(ctrl)

	s := service.NewChapterService(r, pr)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.Nil(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.Nil(t, err)

	chapters, err := s.ListChapters(*userId, *projectId)
	assert.Nil(t, err)

	assert.Len(t, chapters, 4)

	chapter := chapters[0]
	sections := chapter.Sections()
	assert.Equal(t, "1000000000000001", chapter.Id().Value())
	assert.Equal(t, "Chapter 1", chapter.Name().Value())
	assert.Equal(t, 1, chapter.Number().Value())
	assert.Len(t, sections, 2)
	assert.Equal(t, "2000000000000001", sections[0].Id().Value())
	assert.Equal(t, "Section 1", sections[0].Name().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), sections[0].CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), sections[0].UpdatedAt().Value())
	assert.Equal(t, "2000000000000002", sections[1].Id().Value())
	assert.Equal(t, "Section 2", sections[1].Name().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), sections[1].CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), sections[1].UpdatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), chapter.UpdatedAt().Value())

	chapter = chapters[1]
	sections = chapter.Sections()
	assert.Equal(t, "1000000000000002", chapter.Id().Value())
	assert.Equal(t, "Chapter 2", chapter.Name().Value())
	assert.Equal(t, 2, chapter.Number().Value())
	assert.Len(t, sections, 0)
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), chapter.UpdatedAt().Value())

	chapter = chapters[2]
	sections = chapter.Sections()
	assert.Equal(t, "1000000000000003", chapter.Id().Value())
	assert.Equal(t, "Chapter 3", chapter.Name().Value())
	assert.Equal(t, 3, chapter.Number().Value())
	assert.Len(t, sections, 1)
	assert.Equal(t, "2000000000000003", sections[0].Id().Value())
	assert.Equal(t, "Section 1", sections[0].Name().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), sections[0].CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), sections[0].UpdatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), chapter.UpdatedAt().Value())

	chapter = chapters[3]
	sections = chapter.Sections()
	assert.Equal(t, "1000000000000004", chapter.Id().Value())
	assert.Equal(t, maxLengthChapterName, chapter.Name().Value())
	assert.Equal(t, 4, chapter.Number().Value())
	assert.Len(t, sections, 1)
	assert.Equal(t, "2000000000000004", sections[0].Id().Value())
	assert.Equal(t, maxLengthSectionName, sections[0].Name().Value())
	assert.Equal(t, testutil.Date().Add(-4*time.Hour), sections[0].CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-4*time.Hour), sections[0].UpdatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-4*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-4*time.Hour), chapter.UpdatedAt().Value())
}

func TestListChaptersNoEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockChapterRepository(ctrl)
	r.EXPECT().
		FetchChapters(testutil.ReadOnlyUserId(), "0000000000000001").
		Return(map[string]record.ChapterEntry{}, nil)

	pr := mock_repository.NewMockPaperRepository(ctrl)

	s := service.NewChapterService(r, pr)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.Nil(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.Nil(t, err)

	chapters, sErr := s.ListChapters(*userId, *projectId)
	assert.Nil(t, sErr)

	assert.Len(t, chapters, 0)
}

func TestListChaptersInvalidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tooLongChapterName := testutil.RandomString(101)
	tooLongSectionName := testutil.RandomString(101)

	tt := []struct {
		name          string
		chapterId     string
		chapter       record.ChapterEntry
		expectedError string
	}{
		{
			name:      "should return error when chapter id is empty",
			chapterId: "",
			chapter: record.ChapterEntry{
				Name:      "Chapter 1",
				Number:    1,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (id): chapter id is required, but got ''",
		},
		{
			name:      "should return error when chapter name is empty",
			chapterId: "1000000000000001",
			chapter: record.ChapterEntry{
				Name:      "",
				Number:    1,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): chapter name is required, but got ''",
		},
		{
			name:      "should return error when chapter name is too long",
			chapterId: "1000000000000001",
			chapter: record.ChapterEntry{
				Name:      tooLongChapterName,
				Number:    1,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): " + fmt.Sprintf(
				"chapter name cannot be longer than 100 characters, but got '%v'",
				tooLongChapterName,
			),
		},
		{
			name:      "should return error when chapter number is zero",
			chapterId: "1000000000000001",
			chapter: record.ChapterEntry{
				Name:      "Chapter 1",
				Number:    0,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (number): chapter number must be greater than 0, but got 0",
		},
		{
			name:      "should return error when section id is empty",
			chapterId: "1000000000000001",
			chapter: record.ChapterEntry{
				Name:   "Chapter 1",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "",
						Name:      "Section 1",
						UserId:    testutil.ReadOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (section): failed to convert entry to entity (id): " +
				"section id is required, but got ''",
		},
		{
			name:      "should return error when section name is empty",
			chapterId: "1000000000000001",
			chapter: record.ChapterEntry{
				Name:   "Chapter 1",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "2000000000000001",
						Name:      "",
						UserId:    testutil.ReadOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (section): failed to convert entry to entity (name): " +
				"section name is required, but got ''",
		},
		{
			name:      "should return error when section name is too long",
			chapterId: "1000000000000001",
			chapter: record.ChapterEntry{
				Name:   "Chapter 1",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "2000000000000001",
						Name:      tooLongSectionName,
						UserId:    testutil.ReadOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				},
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (section): failed to convert entry to entity (name): " +
				fmt.Sprintf("section name cannot be longer than 100 characters, but got '%v'", tooLongSectionName),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				FetchChapters(testutil.ReadOnlyUserId(), "0000000000000001").
				Return(map[string]record.ChapterEntry{
					tc.chapterId: tc.chapter,
				}, nil)

			pr := mock_repository.NewMockPaperRepository(ctrl)

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.Nil(t, err)

			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)

			chapters, sErr := s.ListChapters(*userId, *projectId)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %v", tc.expectedError), sErr.Error())
			assert.Nil(t, chapters)
		})
	}
}

func TestListChaptersRepositoryError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     repository.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  service.ErrorCode
	}{
		{
			name:          "should return error when repository returns not found error",
			errorCode:     repository.NotFoundError,
			errorMessage:  "failed to fetch project",
			expectedError: "failed to list chapters: failed to fetch project",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns read failure error",
			errorCode:     repository.ReadFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to fetch chapters: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				FetchChapters(testutil.ReadOnlyUserId(), "0000000000000001").
				Return(nil, repository.Errorf(tc.errorCode, "%s", tc.errorMessage))

			pr := mock_repository.NewMockPaperRepository(ctrl)

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.Nil(t, err)

			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)

			chapters, sErr := s.ListChapters(*userId, *projectId)
			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, chapters)
		})
	}
}

func TestCreateChapterValidEntry(t *testing.T) {
	maxLengthChapterName := testutil.RandomString(100)

	tt := []struct {
		name    string
		chapter record.ChapterWithoutAutofieldEntry
	}{
		{
			name: "should return chapter with valid entry",
			chapter: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter One",
				Number: 1,
			},
		},
		{
			name: "should return chapter with max-length valid entry",
			chapter: record.ChapterWithoutAutofieldEntry{
				Name:   maxLengthChapterName,
				Number: 1,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				InsertChapter(testutil.ModifyOnlyUserId(), "0000000000000001", tc.chapter).
				Return("1000000000000001", &record.ChapterEntry{
					Name:      tc.chapter.Name,
					Number:    tc.chapter.Number,
					Sections:  []record.SectionEntry{},
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				}, nil)

			paper := record.PaperWithoutAutofieldEntry{Content: ""}
			pr := mock_repository.NewMockPaperRepository(ctrl)
			pr.EXPECT().
				InsertPaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", paper).
				Return("1000000000000001", &record.PaperEntry{
					Content:   paper.Content,
					UserId:    testutil.ModifyOnlyUserId(),
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				}, nil)

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)

			name, err := domain.NewChapterNameObject(tc.chapter.Name)
			assert.Nil(t, err)
			number, err := domain.NewChapterNumberObject(tc.chapter.Number)
			assert.Nil(t, err)

			chapter := domain.NewChapterWithoutAutofieldEntity(*name, *number)

			createdChapter, sErr := s.CreateChapter(*userId, *projectId, *chapter)
			assert.Nil(t, sErr)

			assert.Equal(t, "1000000000000001", createdChapter.Id().Value())
			assert.Equal(t, tc.chapter.Name, createdChapter.Name().Value())
			assert.Equal(t, tc.chapter.Number, createdChapter.Number().Value())
			assert.Len(t, createdChapter.Sections(), 0)
			assert.Equal(t, testutil.Date(), createdChapter.CreatedAt().Value())
			assert.Equal(t, testutil.Date(), createdChapter.UpdatedAt().Value())
		})
	}
}

func TestCreateChapterInvalidCreatedEntry(t *testing.T) {
	tooLongChapterName := testutil.RandomString(101)
	tooLongSectionName := testutil.RandomString(101)

	tt := []struct {
		name           string
		createdChapter record.ChapterEntry
		expectedError  string
	}{
		{
			name: "should return error when chapter name is empty",
			createdChapter: record.ChapterEntry{
				Name:      "",
				Number:    1,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): chapter name is required, but got ''",
		},
		{
			name: "should return error when chapter name is too long",
			createdChapter: record.ChapterEntry{
				Name:      tooLongChapterName,
				Number:    1,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): " + fmt.Sprintf(
				"chapter name cannot be longer than 100 characters, but got '%v'",
				tooLongChapterName,
			),
		},
		{
			name: "should return error when chapter number is zero",
			createdChapter: record.ChapterEntry{
				Name:      "Chapter One",
				Number:    0,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (number): chapter number must be greater than 0, but got 0",
		},
		{
			name: "should return error when section id is empty",
			createdChapter: record.ChapterEntry{
				Name:   "Chapter One",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "",
						Name:      "Section One",
						UserId:    testutil.ModifyOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (section): failed to convert entry to entity (id): " +
				"section id is required, but got ''",
		},
		{
			name: "should return error when section name is empty",
			createdChapter: record.ChapterEntry{
				Name:   "Chapter One",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "2000000000000001",
						Name:      "",
						UserId:    testutil.ModifyOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (section): failed to convert entry to entity (name): " +
				"section name is required, but got ''",
		},
		{
			name: "should return error when section name is too long",
			createdChapter: record.ChapterEntry{
				Name:   "Chapter One",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "2000000000000001",
						Name:      tooLongSectionName,
						UserId:    testutil.ModifyOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (section): failed to convert entry to entity (name): " +
				fmt.Sprintf("section name cannot be longer than 100 characters, but got '%v'", tooLongSectionName),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				InsertChapter(testutil.ModifyOnlyUserId(), "0000000000000001", record.ChapterWithoutAutofieldEntry{
					Name:   "Chapter One",
					Number: 1,
				}).
				Return("1000000000000001", &tc.createdChapter, nil)

			paper := record.PaperWithoutAutofieldEntry{Content: ""}
			pr := mock_repository.NewMockPaperRepository(ctrl)
			pr.EXPECT().
				InsertPaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", paper).
				Return("1000000000000001", &record.PaperEntry{
					Content:   paper.Content,
					UserId:    testutil.ModifyOnlyUserId(),
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				}, nil)

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)

			name, err := domain.NewChapterNameObject("Chapter One")
			assert.Nil(t, err)
			number, err := domain.NewChapterNumberObject(1)
			assert.Nil(t, err)

			chapter := domain.NewChapterWithoutAutofieldEntity(*name, *number)

			createdChapter, sErr := s.CreateChapter(*userId, *projectId, *chapter)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())

			assert.Equal(t, fmt.Sprintf("domain failure: %v", tc.expectedError), sErr.Error())
			assert.Nil(t, createdChapter)
		})
	}
}

func TestCreateChapterRepositoryError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     repository.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  service.ErrorCode
	}{
		{
			name:          "should return error when repository returns invalid argument error",
			errorCode:     repository.InvalidArgumentError,
			errorMessage:  "chapter number is too large",
			expectedError: "failed to create chapter: chapter number is too large",
			expectedCode:  service.InvalidArgumentError,
		},
		{
			name:          "should return error when repository returns not found error",
			errorCode:     repository.NotFoundError,
			errorMessage:  "failed to fetch project",
			expectedError: "failed to create chapter: failed to fetch project",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns write failure error",
			errorCode:     repository.WriteFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to create chapter: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				InsertChapter(testutil.ModifyOnlyUserId(), "0000000000000001", record.ChapterWithoutAutofieldEntry{
					Name:   "Chapter One",
					Number: 1,
				}).
				Return("", nil, repository.Errorf(tc.errorCode, "%s", tc.errorMessage))

			pr := mock_repository.NewMockPaperRepository(ctrl)

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)

			name, err := domain.NewChapterNameObject("Chapter One")
			assert.Nil(t, err)
			number, err := domain.NewChapterNumberObject(1)
			assert.Nil(t, err)

			chapter := domain.NewChapterWithoutAutofieldEntity(*name, *number)

			createdChapter, sErr := s.CreateChapter(*userId, *projectId, *chapter)
			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, createdChapter)
		})
	}
}

func TestCreateChapterPaperRepositoryError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     repository.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  service.ErrorCode
	}{
		{
			name:          "should return error when repository returns not found error",
			errorCode:     repository.NotFoundError,
			errorMessage:  "failed to fetch project",
			expectedError: "failed to create paper: failed to fetch project",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns write failure error",
			errorCode:     repository.WriteFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to create paper: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				InsertChapter(testutil.ModifyOnlyUserId(), "0000000000000001", record.ChapterWithoutAutofieldEntry{
					Name:   "Chapter One",
					Number: 1,
				}).
				Return("1000000000000001", &record.ChapterEntry{
					Name:     "Chapter One",
					Number:   1,
					Sections: []record.SectionEntry{},
					UserId:   testutil.ModifyOnlyUserId(),
				}, nil)

			paper := record.PaperWithoutAutofieldEntry{Content: ""}
			pr := mock_repository.NewMockPaperRepository(ctrl)
			pr.EXPECT().
				InsertPaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", paper).
				Return("", nil, repository.Errorf(tc.errorCode, "%s", tc.errorMessage))

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)

			name, err := domain.NewChapterNameObject("Chapter One")
			assert.Nil(t, err)
			number, err := domain.NewChapterNumberObject(1)
			assert.Nil(t, err)

			chapter := domain.NewChapterWithoutAutofieldEntity(*name, *number)

			createdChapter, sErr := s.CreateChapter(*userId, *projectId, *chapter)
			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, createdChapter)
		})
	}
}

func TestUpdateChapterValidEntry(t *testing.T) {
	maxLengthChapterName := testutil.RandomString(100)
	maxLengthSectionName := testutil.RandomString(100)

	tt := []struct {
		name     string
		chapter  record.ChapterWithoutAutofieldEntry
		sections []record.SectionWithoutAutofieldEntry
	}{
		{
			name: "should return chapter with valid entry",
			chapter: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter One",
				Number: 1,
			},
			sections: []record.SectionWithoutAutofieldEntry{
				{
					Id:   "2000000000000001",
					Name: "Section One",
				},
			},
		},
		{
			name: "should return chapter with max-length valid entry",
			chapter: record.ChapterWithoutAutofieldEntry{
				Name:   maxLengthChapterName,
				Number: 1,
			},
			sections: []record.SectionWithoutAutofieldEntry{
				{
					Id:   "2000000000000001",
					Name: maxLengthSectionName,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sectionEntries := make([]record.SectionEntry, len(tc.sections))
			for i, section := range tc.sections {
				sectionEntries[i] = record.SectionEntry{
					Id:        section.Id,
					Name:      section.Name,
					UserId:    testutil.ModifyOnlyUserId(),
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				}
			}

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				UpdateChapter(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", tc.chapter).
				Return(&record.ChapterEntry{
					Name:      tc.chapter.Name,
					Number:    tc.chapter.Number,
					Sections:  sectionEntries,
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				}, nil)

			pr := mock_repository.NewMockPaperRepository(ctrl)

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)

			name, err := domain.NewChapterNameObject(tc.chapter.Name)
			assert.Nil(t, err)
			number, err := domain.NewChapterNumberObject(tc.chapter.Number)
			assert.Nil(t, err)

			sections := make([]domain.SectionOfChapterWithoutAutofieldEntity, len(tc.sections))
			for i, section := range tc.sections {
				sectionId, err := domain.NewSectionIdObject(section.Id)
				assert.Nil(t, err)
				sectionName, err := domain.NewSectionNameObject(section.Name)
				assert.Nil(t, err)

				sections[i] = *domain.NewSectionOfChapterWithoutAutofieldEntity(*sectionId, *sectionName)
			}

			chapter := domain.NewChapterWithoutAutofieldEntity(*name, *number)

			updatedChapter, sErr := s.UpdateChapter(*userId, *projectId, *chapterId, *chapter)
			updatedSections := updatedChapter.Sections()
			assert.Nil(t, sErr)

			assert.Equal(t, "1000000000000001", updatedChapter.Id().Value())
			assert.Equal(t, tc.chapter.Name, updatedChapter.Name().Value())
			assert.Equal(t, tc.chapter.Number, updatedChapter.Number().Value())
			assert.Len(t, updatedSections, len(tc.sections))
			for i, updatedSection := range updatedSections {
				assert.Equal(t, tc.sections[i].Id, updatedSection.Id().Value())
				assert.Equal(t, tc.sections[i].Name, updatedSection.Name().Value())
				assert.Equal(t, testutil.Date(), updatedSection.CreatedAt().Value())
				assert.Equal(t, testutil.Date(), updatedSection.UpdatedAt().Value())
			}
			assert.Equal(t, testutil.Date(), updatedChapter.CreatedAt().Value())
			assert.Equal(t, testutil.Date(), updatedChapter.UpdatedAt().Value())
		})
	}
}

func TestUpdateChapterInvalidUpdatedEntry(t *testing.T) {
	tooLongChapterName := testutil.RandomString(101)
	tooLongSectionName := testutil.RandomString(101)

	tt := []struct {
		name           string
		updatedChapter record.ChapterEntry
		expectedError  string
	}{
		{
			name: "should return error when chapter name is empty",
			updatedChapter: record.ChapterEntry{
				Name:      "",
				Number:    1,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): chapter name is required, but got ''",
		},
		{
			name: "should return error when chapter name is too long",
			updatedChapter: record.ChapterEntry{
				Name:      tooLongChapterName,
				Number:    1,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): " + fmt.Sprintf(
				"chapter name cannot be longer than 100 characters, but got '%v'",
				tooLongChapterName,
			),
		},
		{
			name: "should return error when chapter number is zero",
			updatedChapter: record.ChapterEntry{
				Name:      "Chapter One",
				Number:    0,
				Sections:  []record.SectionEntry{},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (number): chapter number must be greater than 0, but got 0",
		},
		{
			name: "should return error when section id is empty",
			updatedChapter: record.ChapterEntry{
				Name:   "Chapter One",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "",
						Name:      "Section One",
						UserId:    testutil.ModifyOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (section): failed to convert entry to entity (id): " +
				"section id is required, but got ''",
		},
		{
			name: "should return error when section name is empty",
			updatedChapter: record.ChapterEntry{
				Name:   "Chapter One",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "2000000000000001",
						Name:      "",
						UserId:    testutil.ModifyOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (section): failed to convert entry to entity (name): " +
				"section name is required, but got ''",
		},
		{
			name: "should return error when section name is too long",
			updatedChapter: record.ChapterEntry{
				Name:   "Chapter One",
				Number: 1,
				Sections: []record.SectionEntry{
					{
						Id:        "2000000000000001",
						Name:      tooLongSectionName,
						UserId:    testutil.ModifyOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				},
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (section): failed to convert entry to entity (name): " +
				fmt.Sprintf("section name cannot be longer than 100 characters, but got '%v'", tooLongSectionName),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				UpdateChapter(
					testutil.ModifyOnlyUserId(),
					"0000000000000001",
					"1000000000000001",
					record.ChapterWithoutAutofieldEntry{
						Name:   "Chapter One",
						Number: 1,
					},
				).
				Return(&tc.updatedChapter, nil)

			pr := mock_repository.NewMockPaperRepository(ctrl)

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)

			name, err := domain.NewChapterNameObject("Chapter One")
			assert.Nil(t, err)
			number, err := domain.NewChapterNumberObject(1)
			assert.Nil(t, err)

			chapter := domain.NewChapterWithoutAutofieldEntity(*name, *number)

			updatedChapter, sErr := s.UpdateChapter(*userId, *projectId, *chapterId, *chapter)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())

			assert.Equal(t, fmt.Sprintf("domain failure: %v", tc.expectedError), sErr.Error())
			assert.Nil(t, updatedChapter)
		})
	}
}

func TestUpdateChapterRepositoryError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     repository.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  service.ErrorCode
	}{
		{
			name:          "should return error when repository returns not found error",
			errorCode:     repository.NotFoundError,
			errorMessage:  "failed to update chapter",
			expectedError: "failed to update chapter: failed to update chapter",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns invalid argument error",
			errorCode:     repository.InvalidArgumentError,
			errorMessage:  "chapter number is too large",
			expectedError: "failed to update chapter: chapter number is too large",
			expectedCode:  service.InvalidArgumentError,
		},
		{
			name:          "should return error when repository returns write failure error",
			errorCode:     repository.WriteFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to update chapter: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				UpdateChapter(
					testutil.ModifyOnlyUserId(),
					"0000000000000001",
					"1000000000000001",
					record.ChapterWithoutAutofieldEntry{
						Name:   "Chapter One",
						Number: 1,
					},
				).
				Return(nil, repository.Errorf(tc.errorCode, "%s", tc.errorMessage))

			pr := mock_repository.NewMockPaperRepository(ctrl)

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)

			name, err := domain.NewChapterNameObject("Chapter One")
			assert.Nil(t, err)
			number, err := domain.NewChapterNumberObject(1)
			assert.Nil(t, err)

			chapter := domain.NewChapterWithoutAutofieldEntity(*name, *number)

			updatedChapter, sErr := s.UpdateChapter(*userId, *projectId, *chapterId, *chapter)

			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, updatedChapter)
		})
	}
}

func TestDeleteChapterValidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockChapterRepository(ctrl)
	r.EXPECT().
		DeleteChapter(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
		Return(nil)

	pr := mock_repository.NewMockPaperRepository(ctrl)
	pr.EXPECT().
		DeletePaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
		Return(nil)

	s := service.NewChapterService(r, pr)

	userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
	assert.Nil(t, err)
	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.Nil(t, err)
	chapterId, err := domain.NewChapterIdObject("1000000000000001")
	assert.Nil(t, err)

	sErr := s.DeleteChapter(*userId, *projectId, *chapterId)
	assert.Nil(t, sErr)
}

func TestDeleteChapterRepositoryError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     repository.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  service.ErrorCode
	}{
		{
			name:          "should return error when repository returns not found error",
			errorCode:     repository.NotFoundError,
			errorMessage:  "failed to delete chapter",
			expectedError: "failed to delete chapter: failed to delete chapter",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns write failure error",
			errorCode:     repository.WriteFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to delete chapter: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				DeleteChapter(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(repository.Errorf(tc.errorCode, "%s", tc.errorMessage))

			pr := mock_repository.NewMockPaperRepository(ctrl)
			pr.EXPECT().
				DeletePaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(nil)

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)

			sErr := s.DeleteChapter(*userId, *projectId, *chapterId)
			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
		})
	}
}

func TestDeletePaperRepositoryError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     repository.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  service.ErrorCode
	}{
		{
			name:          "should return error when repository returns not found error",
			errorCode:     repository.NotFoundError,
			errorMessage:  "failed to delete paper",
			expectedError: "failed to delete paper: failed to delete paper",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns write failure error",
			errorCode:     repository.WriteFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to delete paper: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)

			pr := mock_repository.NewMockPaperRepository(ctrl)
			pr.EXPECT().
				DeletePaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(repository.Errorf(tc.errorCode, "%s", tc.errorMessage))

			s := service.NewChapterService(r, pr)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)

			sErr := s.DeleteChapter(*userId, *projectId, *chapterId)
			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
		})
	}
}
