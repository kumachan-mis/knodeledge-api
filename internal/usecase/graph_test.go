package usecase_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	mock_service "github.com/kumachan-mis/knodeledge-api/mock/service"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSectionalizeGraphValidEntity(t *testing.T) {
	maxLengthSectionName := testutil.RandomString(100)
	maxLengthSectionContent := testutil.RandomString(40000)

	maxLengthSectionIds := make([]string, 20)
	for i := 0; i < 20; i++ {
		maxLengthSectionIds[i] = testutil.RandomString(16)
	}
	maxLengthSections := make([]model.SectionWithoutAutofield, 20)
	for i := 0; i < 20; i++ {
		maxLengthSections[i] = model.SectionWithoutAutofield{
			Name:    testutil.RandomString(10),
			Content: testutil.RandomString(100),
		}
	}

	tt := []struct {
		name       string
		userId     string
		projectId  string
		chapterId  string
		sectionIds []string
		sections   []model.SectionWithoutAutofield
	}{
		{
			name:      "should sectionalize into graphs",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionIds: []string{
				"2000000000000001",
				"2000000000000002",
			},
			sections: []model.SectionWithoutAutofield{
				{
					Name:    "Section 1",
					Content: "This is the content of section 1. This is the content of section 1.",
				},
				{
					Name:    "Section 2",
					Content: "This is the content of section 2. This is the content of section 2.",
				},
			},
		},
		{
			name:      "should sectionalize into graphs with max length name",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionIds: []string{
				"2000000000000001",
			},
			sections: []model.SectionWithoutAutofield{
				{
					Name:    maxLengthSectionName,
					Content: "This is the content of section. This is the content of section.",
				},
			},
		},
		{
			name:      "should sectionalize into graphs with max length content",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionIds: []string{
				"2000000000000001",
			},
			sections: []model.SectionWithoutAutofield{
				{
					Name:    "Section",
					Content: maxLengthSectionContent,
				},
			},
		},
		{
			name:       "should sectionalize into graphs with max length sections",
			userId:     testutil.ModifyOnlyUserId(),
			projectId:  "0000000000000001",
			chapterId:  "1000000000000001",
			sectionIds: maxLengthSectionIds,
			sections:   maxLengthSections,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockGraphService(ctrl)

			grapths := make([]domain.GraphEntity, len(tc.sections))
			for i, section := range tc.sections {
				id, err := domain.NewGraphIdObject(tc.sectionIds[i])
				assert.Nil(t, err)
				name, err := domain.NewGraphNameObject(section.Name)
				assert.Nil(t, err)
				paragraph, err := domain.NewGraphParagraphObject(section.Content)
				assert.Nil(t, err)
				createdAt, err := domain.NewCreatedAtObject(testutil.Date())
				assert.NoError(t, err)
				updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
				assert.NoError(t, err)
				grapths[i] = *domain.NewGraphEntity(*id, *name, *paragraph, *createdAt, *updatedAt)
			}
			s.EXPECT().
				SectionalizeIntoGraphs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(
					userId domain.UserIdObject,
					projectId domain.ProjectIdObject,
					chapterId domain.ChapterIdObject,
					sections domain.SectionWithoutAutofieldEntityList) {
					assert.Equal(t, tc.userId, userId.Value())
					assert.Equal(t, tc.projectId, projectId.Value())
					assert.Equal(t, tc.chapterId, chapterId.Value())
					assert.Equal(t, len(tc.sections), sections.Len())
					for i, section := range tc.sections {
						assert.Equal(t, section.Name, sections.Value()[i].Name().Value())
						assert.Equal(t, section.Content, sections.Value()[i].Content().Value())
					}
				}).
				Return(grapths, nil)

			uc := usecase.NewGraphUseCase(s)

			res, err := uc.SectionalizeGraph(model.GraphSectionalizeRequest{
				User:     model.UserOnlyId{Id: tc.userId},
				Project:  model.ProjectOnlyId{Id: tc.projectId},
				Chapter:  model.ChapterOnlyId{Id: tc.chapterId},
				Sections: tc.sections,
			})

			assert.Nil(t, err)

			assert.Equal(t, len(tc.sections), len(res.Graphs))
			for i, graph := range res.Graphs {
				assert.Equal(t, tc.sectionIds[i], graph.Id)
				assert.Equal(t, tc.sections[i].Name, graph.Name)
				assert.Equal(t, tc.sections[i].Content, graph.Paragraph)
			}
		})
	}
}

func TestSectionalizeGraphDomainValidationError(t *testing.T) {
	tooLongSectionName := testutil.RandomString(101)
	tooLongSectionContent := testutil.RandomString(40001)

	tooLongSectionIds := make([]string, 21)
	for i := 0; i < 21; i++ {
		tooLongSectionIds[i] = testutil.RandomString(16)
	}
	tooLongSections := make([]model.SectionWithoutAutofield, 21)
	for i := 0; i < 21; i++ {
		tooLongSections[i] = model.SectionWithoutAutofield{
			Name:    testutil.RandomString(10),
			Content: testutil.RandomString(100),
		}
	}
	tooLongSectionItemsError := make([]model.SectionWithoutAutofieldError, 21)
	for i := 0; i < 21; i++ {
		tooLongSectionItemsError[i] = model.SectionWithoutAutofieldError{}
	}

	tt := []struct {
		name       string
		userId     string
		projectId  string
		chapterId  string
		sectionIds []string
		sections   []model.SectionWithoutAutofield
		expected   model.GraphSectionalizeErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionIds: []string{
				"2000000000000001",
			},
			sections: []model.SectionWithoutAutofield{
				{
					Name:    "Section",
					Content: "This is the content of section. This is the content of section.",
				},
			},
			expected: model.GraphSectionalizeErrorResponse{
				User: model.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
				Sections: model.SectionWithoutAutofieldListError{
					Items: []model.SectionWithoutAutofieldError{
						{
							Name:    "",
							Content: "",
						},
					},
				},
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "",
			chapterId: "1000000000000001",
			sectionIds: []string{
				"2000000000000001",
			},
			sections: []model.SectionWithoutAutofield{
				{
					Name:    "Section",
					Content: "This is the content of section. This is the content of section.",
				},
			},
			expected: model.GraphSectionalizeErrorResponse{
				Project: model.ProjectOnlyIdError{
					Id: "project id is required, but got ''",
				},
				Sections: model.SectionWithoutAutofieldListError{
					Items: []model.SectionWithoutAutofieldError{
						{
							Name:    "",
							Content: "",
						},
					},
				},
			},
		},
		{
			name:      "should return error when chapter id is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "",
			sectionIds: []string{
				"2000000000000001",
			},
			sections: []model.SectionWithoutAutofield{
				{
					Name:    "Section",
					Content: "This is the content of section. This is the content of section.",
				},
			},
			expected: model.GraphSectionalizeErrorResponse{
				Chapter: model.ChapterOnlyIdError{
					Id: "chapter id is required, but got ''",
				},
				Sections: model.SectionWithoutAutofieldListError{
					Items: []model.SectionWithoutAutofieldError{
						{
							Name:    "",
							Content: "",
						},
					},
				},
			},
		},
		{
			name:       "should return error when section is empty",
			userId:     testutil.ModifyOnlyUserId(),
			projectId:  "0000000000000001",
			chapterId:  "1000000000000001",
			sectionIds: []string{},
			sections:   []model.SectionWithoutAutofield{},
			expected: model.GraphSectionalizeErrorResponse{
				Sections: model.SectionWithoutAutofieldListError{
					Message: "sections are required, but got []",
					Items:   []model.SectionWithoutAutofieldError{},
				},
			},
		},
		{
			name:       "should return error when section is too long",
			userId:     testutil.ModifyOnlyUserId(),
			projectId:  "0000000000000001",
			chapterId:  "1000000000000001",
			sectionIds: tooLongSectionIds,
			sections:   tooLongSections,
			expected: model.GraphSectionalizeErrorResponse{
				Sections: model.SectionWithoutAutofieldListError{
					Message: "sections length must be less than or equal to 20, but got 21",
					Items:   tooLongSectionItemsError,
				},
			},
		},
		{
			name:      "should return error when section name is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionIds: []string{
				"2000000000000001",
				"2000000000000002",
			},
			sections: []model.SectionWithoutAutofield{
				{
					Name:    "Section 1",
					Content: "This is the content of section 1. This is the content of section 1.",
				},
				{
					Name:    "",
					Content: "This is the content of section 2. This is the content of section 2.",
				},
			},
			expected: model.GraphSectionalizeErrorResponse{
				Sections: model.SectionWithoutAutofieldListError{
					Items: []model.SectionWithoutAutofieldError{
						{
							Name:    "",
							Content: "",
						},
						{
							Name:    "section name is required, but got ''",
							Content: "",
						},
					},
				},
			},
		},
		{
			name:      "should return error when section name is too long",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionIds: []string{
				"2000000000000001",
				"2000000000000002",
			},
			sections: []model.SectionWithoutAutofield{
				{
					Name:    "Section 1",
					Content: "This is the content of section 1. This is the content of section 1.",
				},
				{
					Name:    tooLongSectionName,
					Content: "This is the content of section 2. This is the content of section 2.",
				},
			},
			expected: model.GraphSectionalizeErrorResponse{
				Sections: model.SectionWithoutAutofieldListError{
					Items: []model.SectionWithoutAutofieldError{
						{
							Name:    "",
							Content: "",
						},
						{
							Name: fmt.Sprintf("section name cannot be longer than 100 characters, but got '%v'",
								tooLongSectionName),
							Content: "",
						},
					},
				},
			},
		},
		{
			name:      "should return error when section content is too long",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionIds: []string{
				"2000000000000001",
				"2000000000000002",
			},
			sections: []model.SectionWithoutAutofield{
				{
					Name:    "Section 1",
					Content: tooLongSectionContent,
				},
				{
					Name:    "Section 2",
					Content: "This is the content of section 2. This is the content of section 2.",
				},
			},
			expected: model.GraphSectionalizeErrorResponse{
				Sections: model.SectionWithoutAutofieldListError{
					Items: []model.SectionWithoutAutofieldError{
						{
							Name:    "",
							Content: "section content must be less than or equal to 40000 bytes, but got 40001 bytes",
						},
						{
							Name:    "",
							Content: "",
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockGraphService(ctrl)

			uc := usecase.NewGraphUseCase(s)

			res, ucErr := uc.SectionalizeGraph(model.GraphSectionalizeRequest{
				User:     model.UserOnlyId{Id: tc.userId},
				Project:  model.ProjectOnlyId{Id: tc.projectId},
				Chapter:  model.ChapterOnlyId{Id: tc.chapterId},
				Sections: tc.sections,
			})
			assert.NotNil(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestSectionalizeGraphServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when repository returns invalid argument error",
			errorCode:     service.InvalidArgument,
			errorMessage:  "graph already exists",
			expectedError: "invalid argument: graph already exists",
			expectedCode:  usecase.InvalidArgumentError,
		},
		{
			name:          "should return error when chapter not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to find chapter",
			expectedError: "not found: failed to find chapter",
			expectedCode:  usecase.NotFoundError,
		},
		{
			name:          "should return error when repository failure",
			errorCode:     service.RepositoryFailurePanic,
			errorMessage:  "service error",
			expectedError: "internal error: service error",
			expectedCode:  usecase.InternalErrorPanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockGraphService(ctrl)
			s.EXPECT().
				SectionalizeIntoGraphs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, tc.errorMessage))

			uc := usecase.NewGraphUseCase(s)

			res, ucErr := uc.SectionalizeGraph(model.GraphSectionalizeRequest{
				User:    model.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: model.ProjectOnlyId{Id: "0000000000000001"},
				Chapter: model.ChapterOnlyId{Id: "1000000000000001"},
				Sections: []model.SectionWithoutAutofield{
					{
						Name:    "Section",
						Content: "This is the content of section. This is the content of section.",
					},
				},
			})
			assert.NotNil(t, ucErr)

			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())

			assert.Nil(t, res)
		})
	}
}