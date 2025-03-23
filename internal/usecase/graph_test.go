package usecase_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/testutil"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	mock_service "github.com/kumachan-mis/knodeledge-api/mock/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindGraphValidEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockGraphService(ctrl)

	id, err := domain.NewGraphIdObject("2000000000000001")
	assert.Nil(t, err)
	name, err := domain.NewGraphNameObject("Section")
	assert.Nil(t, err)
	paragraph, err := domain.NewGraphParagraphObject("This is graph paragraph")
	assert.Nil(t, err)

	emptyChildren, err := domain.NewGraphChildrenEntity([]domain.GraphChildEntity{})
	assert.Nil(t, err)

	grandChildName, err := domain.NewGraphNameObject("GrandChild")
	assert.Nil(t, err)
	grandChildRelation, err := domain.NewGraphRelationObject("grandchild relation")
	assert.Nil(t, err)
	grandChildDescription, err := domain.NewGraphDescriptionObject("grandchild description")
	assert.Nil(t, err)
	grandChild := domain.NewGraphChildEntity(*grandChildName, *grandChildRelation, *grandChildDescription, *emptyChildren)
	grandChildren, err := domain.NewGraphChildrenEntity([]domain.GraphChildEntity{*grandChild})
	assert.Nil(t, err)

	childName, err := domain.NewGraphNameObject("Child")
	assert.Nil(t, err)
	childRelation, err := domain.NewGraphRelationObject("child relation")
	assert.Nil(t, err)
	childDescription, err := domain.NewGraphDescriptionObject("child description")
	assert.Nil(t, err)
	child := domain.NewGraphChildEntity(*childName, *childRelation, *childDescription, *grandChildren)
	children, err := domain.NewGraphChildrenEntity([]domain.GraphChildEntity{*child})
	assert.Nil(t, err)

	createdAt, err := domain.NewCreatedAtObject(testutil.Date())
	assert.Nil(t, err)
	updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
	assert.Nil(t, err)

	graph := domain.NewGraphEntity(*id, *name, *paragraph, *children, *createdAt, *updatedAt)

	s.EXPECT().
		FindGraph(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(
			userId domain.UserIdObject,
			projectId domain.ProjectIdObject,
			chapterId domain.ChapterIdObject,
			sectionId domain.SectionIdObject,
		) {
			assert.Equal(t, testutil.ReadOnlyUserId(), userId.Value())
			assert.Equal(t, "0000000000000001", projectId.Value())
			assert.Equal(t, "1000000000000001", chapterId.Value())
			assert.Equal(t, "2000000000000001", sectionId.Value())
		}).
		Return(graph, nil)

	uc := usecase.NewGraphUseCase(s)

	res, ucErr := uc.FindGraph(openapi.GraphFindRequest{
		UserId:    testutil.ReadOnlyUserId(),
		ProjectId: "0000000000000001",
		ChapterId: "1000000000000001",
		SectionId: "2000000000000001",
	})

	assert.Nil(t, ucErr)

	assert.Equal(t, "2000000000000001", res.Graph.Id)
	assert.Equal(t, "Section", res.Graph.Name)
	assert.Equal(t, "This is graph paragraph", res.Graph.Paragraph)

	assert.Equal(t, 1, len(res.Graph.Children))
	assert.Equal(t, "Child", res.Graph.Children[0].Name)
	assert.Equal(t, "child relation", res.Graph.Children[0].Relation)
	assert.Equal(t, "child description", res.Graph.Children[0].Description)

	assert.Equal(t, 1, len(res.Graph.Children[0].Children))
	assert.Equal(t, "GrandChild", res.Graph.Children[0].Children[0].Name)
	assert.Equal(t, "grandchild relation", res.Graph.Children[0].Children[0].Relation)
	assert.Equal(t, "grandchild description", res.Graph.Children[0].Children[0].Description)
	assert.Equal(t, 0, len(res.Graph.Children[0].Children[0].Children))
}

func TestFindGraphDomainValidationError(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		sectionId string
		expected  openapi.GraphFindErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionId: "2000000000000001",
			expected: openapi.GraphFindErrorResponse{
				UserId:    "user id is required, but got ''",
				ProjectId: "",
				ChapterId: "",
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			chapterId: "1000000000000001",
			sectionId: "2000000000000001",
			expected: openapi.GraphFindErrorResponse{
				UserId:    "",
				ProjectId: "project id is required, but got ''",
				ChapterId: "",
			},
		},
		{
			name:      "should return error when chapter id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "",
			sectionId: "2000000000000001",
			expected: openapi.GraphFindErrorResponse{
				UserId:    "",
				ProjectId: "",
				ChapterId: "chapter id is required, but got ''",
			},
		},
		{
			name:      "should return error when section id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionId: "",
			expected: openapi.GraphFindErrorResponse{
				UserId:    "",
				ProjectId: "",
				ChapterId: "",
				SectionId: "section id is required, but got ''",
			},
		},
		{
			name:      "should return error when all fields are empty",
			userId:    "",
			projectId: "",
			chapterId: "",
			sectionId: "",
			expected: openapi.GraphFindErrorResponse{
				UserId:    "user id is required, but got ''",
				ProjectId: "project id is required, but got ''",
				ChapterId: "chapter id is required, but got ''",
				SectionId: "section id is required, but got ''",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockGraphService(ctrl)

			uc := usecase.NewGraphUseCase(s)

			res, ucErr := uc.FindGraph(openapi.GraphFindRequest{
				UserId:    tc.userId,
				ProjectId: tc.projectId,
				ChapterId: tc.chapterId,
				SectionId: tc.sectionId,
			})

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestFindGraphServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when project not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to find project",
			expectedError: "not found: failed to find project",
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

			uc := usecase.NewGraphUseCase(s)

			s.EXPECT().
				FindGraph(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			res, ucErr := uc.FindGraph(openapi.GraphFindRequest{
				UserId:    testutil.ReadOnlyUserId(),
				ProjectId: "0000000000000001",
				ChapterId: "1000000000000001",
				SectionId: "2000000000000001",
			})

			assert.Nil(t, res)
			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
		})
	}
}

func TestUpdateGraphContentValidEntity(t *testing.T) {
	maxLengthParagraph := testutil.RandomString(40000)

	maxLengthChildName := testutil.RandomString(100)
	maxLengthRelation := testutil.RandomString(100)
	maxLengthDescription := testutil.RandomString(400)

	tt := []struct {
		name        string
		paragraph   string
		childName   string
		relation    string
		description string
	}{
		{
			name:        "should update graph content",
			paragraph:   "This is updated graph paragraph",
			childName:   "Child",
			relation:    "child relation",
			description: "child description",
		},
		{
			name:        "should update graph content with max length paragraph",
			paragraph:   maxLengthParagraph,
			childName:   maxLengthChildName,
			relation:    maxLengthRelation,
			description: maxLengthDescription,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockGraphService(ctrl)

			id, err := domain.NewGraphIdObject("2000000000000001")
			assert.Nil(t, err)
			name, err := domain.NewGraphNameObject("Section")
			assert.Nil(t, err)
			paragraph, err := domain.NewGraphParagraphObject(tc.paragraph)
			assert.Nil(t, err)

			emptyChildren, err := domain.NewGraphChildrenEntity([]domain.GraphChildEntity{})
			assert.Nil(t, err)

			grandChildName, err := domain.NewGraphNameObject(tc.childName)
			assert.Nil(t, err)
			grandChildRelation, err := domain.NewGraphRelationObject(tc.relation)
			assert.Nil(t, err)
			grandChildDescription, err := domain.NewGraphDescriptionObject(tc.description)
			assert.Nil(t, err)
			grandChild := domain.NewGraphChildEntity(*grandChildName, *grandChildRelation, *grandChildDescription, *emptyChildren)
			grandChildren, err := domain.NewGraphChildrenEntity([]domain.GraphChildEntity{*grandChild})
			assert.Nil(t, err)

			childName, err := domain.NewGraphNameObject(tc.childName)
			assert.Nil(t, err)
			childRelation, err := domain.NewGraphRelationObject(tc.relation)
			assert.Nil(t, err)
			childDescription, err := domain.NewGraphDescriptionObject(tc.description)
			assert.Nil(t, err)
			child := domain.NewGraphChildEntity(*childName, *childRelation, *childDescription, *grandChildren)
			children, err := domain.NewGraphChildrenEntity([]domain.GraphChildEntity{*child})
			assert.Nil(t, err)

			createdAt, err := domain.NewCreatedAtObject(testutil.Date())
			assert.Nil(t, err)
			updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
			assert.Nil(t, err)

			graph := domain.NewGraphEntity(*id, *name, *paragraph, *children, *createdAt, *updatedAt)

			s.EXPECT().
				UpdateGraphContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(
					userId domain.UserIdObject,
					projectId domain.ProjectIdObject,
					chapterId domain.ChapterIdObject,
					graphId domain.GraphIdObject,
					graph domain.GraphContentEntity,
				) {
					assert.Equal(t, testutil.ModifyOnlyUserId(), userId.Value())
					assert.Equal(t, "0000000000000001", projectId.Value())
					assert.Equal(t, "1000000000000001", chapterId.Value())
					assert.Equal(t, "2000000000000001", graphId.Value())
					assert.Equal(t, tc.paragraph, graph.Paragraph().Value())
				}).
				Return(graph, nil)

			uc := usecase.NewGraphUseCase(s)

			res, ucErr := uc.UpdateGraph(openapi.GraphUpdateRequest{
				User:    openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: openapi.ProjectOnlyId{Id: "0000000000000001"},
				Chapter: openapi.ChapterOnlyId{Id: "1000000000000001"},
				Graph: openapi.GraphContent{
					Id:        "2000000000000001",
					Paragraph: tc.paragraph,
					Children: []openapi.GraphChild{
						{
							Name:        tc.childName,
							Relation:    tc.relation,
							Description: tc.description,
							Children: []openapi.GraphChild{
								{
									Name:        tc.childName,
									Relation:    tc.relation,
									Description: tc.description,
									Children:    []openapi.GraphChild{},
								},
							},
						},
					},
				},
			})

			assert.Nil(t, ucErr)
			assert.Equal(t, "2000000000000001", res.Graph.Id)
			assert.Equal(t, "Section", res.Graph.Name)
			assert.Equal(t, tc.paragraph, res.Graph.Paragraph)
		})
	}
}

func TestUpdateGraphContentDomainValidationError(t *testing.T) {
	tooLongParagraph := testutil.RandomString(40001)

	tooLongChildName := testutil.RandomString(101)
	tooLongRelation := testutil.RandomString(101)
	tooLongDescription := testutil.RandomString(401)

	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		graphId   string
		paragraph string
		children  []openapi.GraphChild
		expected  openapi.GraphUpdateErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			graphId:   "2000000000000001",
			paragraph: "This is updated graph paragraph",
			children:  []openapi.GraphChild{},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Graph: openapi.GraphContentError{
					Id:        "",
					Paragraph: "",
					Children: openapi.GraphChildrenError{
						Message: "",
						Items:   []openapi.GraphChildError{},
					},
				},
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "",
			chapterId: "1000000000000001",
			graphId:   "2000000000000001",
			paragraph: "This is updated graph paragraph",
			children:  []openapi.GraphChild{},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: "project id is required, but got ''"},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Graph: openapi.GraphContentError{
					Id:        "",
					Paragraph: "",
					Children: openapi.GraphChildrenError{
						Message: "",
						Items:   []openapi.GraphChildError{},
					},
				},
			},
		},
		{
			name:      "should return error when chapter id is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "",
			graphId:   "2000000000000001",
			paragraph: "This is updated graph paragraph",
			children:  []openapi.GraphChild{},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: "chapter id is required, but got ''"},
				Graph: openapi.GraphContentError{
					Id:        "",
					Paragraph: "",
					Children: openapi.GraphChildrenError{
						Message: "",
						Items:   []openapi.GraphChildError{},
					},
				},
			},
		},
		{
			name:      "should return error when graph id is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			graphId:   "",
			paragraph: "This is updated graph paragraph",
			children: []openapi.GraphChild{
				{
					Name:        "Child",
					Relation:    "child relation",
					Description: "child description",
					Children:    []openapi.GraphChild{},
				},
			},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Graph: openapi.GraphContentError{
					Id:        "graph id is required, but got ''",
					Paragraph: "",
					Children: openapi.GraphChildrenError{
						Message: "",
						Items: []openapi.GraphChildError{
							{
								Name:        "",
								Relation:    "",
								Description: "",
								Children: openapi.GraphChildrenError{
									Message: "",
									Items:   []openapi.GraphChildError{},
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "should return error when paragraph is too long",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			graphId:   "2000000000000001",
			paragraph: tooLongParagraph,
			children:  []openapi.GraphChild{},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Graph: openapi.GraphContentError{
					Id: "",
					Paragraph: fmt.Sprintf(
						"graph paragraph must be less than or equal to 40000 bytes, but got %d bytes",
						len(tooLongParagraph)),
					Children: openapi.GraphChildrenError{
						Message: "",
						Items:   []openapi.GraphChildError{},
					},
				},
			},
		},
		{
			name:      "should return error when child name is too long",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			graphId:   "2000000000000001",
			paragraph: "This is updated graph paragraph",
			children: []openapi.GraphChild{
				{
					Name:        tooLongChildName,
					Relation:    "child relation",
					Description: "child description",
					Children:    []openapi.GraphChild{},
				},
			},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Graph: openapi.GraphContentError{
					Id:        "",
					Paragraph: "",
					Children: openapi.GraphChildrenError{
						Message: "",
						Items: []openapi.GraphChildError{
							{
								Name: fmt.Sprintf(
									"graph name cannot be longer than 100 characters, but got '%v'",
									tooLongChildName),
								Relation:    "",
								Description: "",
								Children: openapi.GraphChildrenError{
									Message: "",
									Items:   []openapi.GraphChildError{},
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "should return error when child relation is too long",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			graphId:   "2000000000000001",
			paragraph: "This is updated graph paragraph",
			children: []openapi.GraphChild{
				{
					Name:        "Child",
					Relation:    tooLongRelation,
					Description: "child description",
					Children:    []openapi.GraphChild{},
				},
			},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Graph: openapi.GraphContentError{
					Id:        "",
					Paragraph: "",
					Children: openapi.GraphChildrenError{
						Message: "",
						Items: []openapi.GraphChildError{
							{
								Name: "",
								Relation: fmt.Sprintf(
									"graph relation cannot be longer than 100 characters, but got '%v'",
									tooLongRelation),
								Description: "",
								Children: openapi.GraphChildrenError{
									Message: "",
									Items:   []openapi.GraphChildError{},
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "should return error when child description is too long",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			graphId:   "2000000000000001",
			paragraph: "This is updated graph paragraph",
			children: []openapi.GraphChild{
				{
					Name:        "Child",
					Relation:    "child relation",
					Description: tooLongDescription,
					Children:    []openapi.GraphChild{},
				},
			},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Graph: openapi.GraphContentError{
					Id:        "",
					Paragraph: "",
					Children: openapi.GraphChildrenError{
						Message: "",
						Items: []openapi.GraphChildError{
							{
								Name:     "",
								Relation: "",
								Description: fmt.Sprintf(
									"graph description cannot be longer than 400 characters, but got '%v'",
									tooLongDescription),
								Children: openapi.GraphChildrenError{
									Message: "",
									Items:   []openapi.GraphChildError{},
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "should return error when child names are duplicated",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			graphId:   "2000000000000001",
			paragraph: "This is updated graph paragraph",
			children: []openapi.GraphChild{
				{
					Name:        "Child",
					Relation:    "child relation",
					Description: "child description",
					Children:    []openapi.GraphChild{},
				},
				{
					Name:        "Child",
					Relation:    "child relation",
					Description: "child description",
					Children:    []openapi.GraphChild{},
				},
			},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Graph: openapi.GraphContentError{
					Id:        "",
					Paragraph: "",
					Children: openapi.GraphChildrenError{
						Message: "names of children must be unique, but got 'Child' duplicated",
						Items: []openapi.GraphChildError{
							{
								Name:        "",
								Relation:    "",
								Description: "",
								Children: openapi.GraphChildrenError{
									Message: "",
									Items:   []openapi.GraphChildError{},
								},
							},
							{
								Name:        "",
								Relation:    "",
								Description: "",
								Children: openapi.GraphChildrenError{
									Message: "",
									Items:   []openapi.GraphChildError{},
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "should return error when grand child has errors",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			graphId:   "2000000000000001",
			paragraph: "This is updated graph paragraph",
			children: []openapi.GraphChild{
				{
					Name:        "Child",
					Relation:    "child relation",
					Description: "child description",
					Children: []openapi.GraphChild{
						{
							Name:        tooLongChildName,
							Relation:    tooLongRelation,
							Description: tooLongDescription,
							Children:    []openapi.GraphChild{},
						},
					},
				},
			},
			expected: openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Graph: openapi.GraphContentError{
					Id:        "",
					Paragraph: "",
					Children: openapi.GraphChildrenError{
						Message: "",
						Items: []openapi.GraphChildError{
							{
								Name:        "",
								Relation:    "",
								Description: "",
								Children: openapi.GraphChildrenError{
									Message: "",
									Items: []openapi.GraphChildError{
										{
											Name: fmt.Sprintf(
												"graph name cannot be longer than 100 characters, but got '%v'",
												tooLongChildName),
											Relation: fmt.Sprintf(
												"graph relation cannot be longer than 100 characters, but got '%v'",
												tooLongRelation),
											Description: fmt.Sprintf(
												"graph description cannot be longer than 400 characters, but got '%v'",
												tooLongDescription),
											Children: openapi.GraphChildrenError{
												Message: "",
												Items:   []openapi.GraphChildError{},
											},
										},
									},
								},
							},
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

			res, ucErr := uc.UpdateGraph(openapi.GraphUpdateRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter: openapi.ChapterOnlyId{Id: tc.chapterId},
				Graph: openapi.GraphContent{
					Id:        tc.graphId,
					Paragraph: tc.paragraph,
					Children:  tc.children,
				},
			})

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestUpdateGraphContentServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when graph not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to find graph",
			expectedError: "not found: failed to find graph",
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

			uc := usecase.NewGraphUseCase(s)

			s.EXPECT().
				UpdateGraphContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			res, ucErr := uc.UpdateGraph(openapi.GraphUpdateRequest{
				User:    openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: openapi.ProjectOnlyId{Id: "0000000000000001"},
				Chapter: openapi.ChapterOnlyId{Id: "1000000000000001"},
				Graph: openapi.GraphContent{
					Id:        "2000000000000001",
					Paragraph: "This is updated graph paragraph",
				},
			})

			assert.Nil(t, res)
			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
		})
	}
}

func TestDeleteGraphValidEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockGraphService(ctrl)
	s.EXPECT().
		DeleteGraph(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(
			userId domain.UserIdObject,
			projectId domain.ProjectIdObject,
			chapterId domain.ChapterIdObject,
			sectionId domain.SectionIdObject,
		) {
			assert.Equal(t, testutil.ModifyOnlyUserId(), userId.Value())
			assert.Equal(t, "0000000000000001", projectId.Value())
			assert.Equal(t, "1000000000000001", chapterId.Value())
			assert.Equal(t, "2000000000000001", sectionId.Value())
		}).
		Return(nil)

	uc := usecase.NewGraphUseCase(s)

	ucErr := uc.DeleteGraph(openapi.GraphDeleteRequest{
		User:    openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
		Project: openapi.ProjectOnlyId{Id: "0000000000000001"},
		Chapter: openapi.ChapterOnlyId{Id: "1000000000000001"},
		Section: openapi.SectionOnlyId{Id: "2000000000000001"},
	})

	assert.Nil(t, ucErr)
}

func TestDeleteGraphDomainValidationError(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		sectionId string
		expected  openapi.GraphDeleteErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionId: "2000000000000001",
			expected: openapi.GraphDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Section: openapi.SectionOnlyIdError{Id: ""},
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "",
			chapterId: "1000000000000001",
			sectionId: "2000000000000001",
			expected: openapi.GraphDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: "project id is required, but got ''"},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Section: openapi.SectionOnlyIdError{Id: ""},
			},
		},
		{
			name:      "should return error when chapter id is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "",
			sectionId: "2000000000000001",
			expected: openapi.GraphDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: "chapter id is required, but got ''"},
				Section: openapi.SectionOnlyIdError{Id: ""},
			},
		},
		{
			name:      "should return error when section id is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionId: "",
			expected: openapi.GraphDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
				Section: openapi.SectionOnlyIdError{Id: "section id is required, but got ''"},
			},
		},
		{
			name:      "should return error when all fields are empty",
			userId:    "",
			projectId: "",
			chapterId: "",
			sectionId: "",
			expected: openapi.GraphDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: openapi.ProjectOnlyIdError{Id: "project id is required, but got ''"},
				Chapter: openapi.ChapterOnlyIdError{Id: "chapter id is required, but got ''"},
				Section: openapi.SectionOnlyIdError{Id: "section id is required, but got ''"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockGraphService(ctrl)

			uc := usecase.NewGraphUseCase(s)

			ucErr := uc.DeleteGraph(openapi.GraphDeleteRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter: openapi.ChapterOnlyId{Id: tc.chapterId},
				Section: openapi.SectionOnlyId{Id: tc.sectionId},
			})

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())
		})
	}
}

func TestDeleteGraphServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when graph not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to delete graph",
			expectedError: "not found: failed to delete graph",
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

			uc := usecase.NewGraphUseCase(s)

			s.EXPECT().
				DeleteGraph(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			ucErr := uc.DeleteGraph(openapi.GraphDeleteRequest{
				User:    openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: openapi.ProjectOnlyId{Id: "0000000000000001"},
				Chapter: openapi.ChapterOnlyId{Id: "1000000000000001"},
				Section: openapi.SectionOnlyId{Id: "2000000000000001"},
			})

			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
		})
	}
}

func TestSectionalizeGraphValidEntity(t *testing.T) {
	maxLengthSectionName := testutil.RandomString(100)
	maxLengthSectionContent := testutil.RandomString(40000)

	maxLengthSectionIds := make([]string, 20)
	for i := 0; i < 20; i++ {
		maxLengthSectionIds[i] = testutil.RandomString(16)
	}
	maxLengthSections := make([]openapi.SectionWithoutAutofield, 20)
	for i := 0; i < 20; i++ {
		maxLengthSections[i] = openapi.SectionWithoutAutofield{
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
		sections   []openapi.SectionWithoutAutofield
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
			sections: []openapi.SectionWithoutAutofield{
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
			sections: []openapi.SectionWithoutAutofield{
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
			sections: []openapi.SectionWithoutAutofield{
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
				children, err := domain.NewGraphChildrenEntity([]domain.GraphChildEntity{})
				assert.Nil(t, err)
				createdAt, err := domain.NewCreatedAtObject(testutil.Date())
				assert.NoError(t, err)
				updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
				assert.NoError(t, err)
				grapths[i] = *domain.NewGraphEntity(*id, *name, *paragraph, *children, *createdAt, *updatedAt)
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

			res, err := uc.SectionalizeGraph(openapi.GraphSectionalizeRequest{
				User:     openapi.UserOnlyId{Id: tc.userId},
				Project:  openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter:  openapi.ChapterOnlyId{Id: tc.chapterId},
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
	tooLongSections := make([]openapi.SectionWithoutAutofield, 21)
	for i := 0; i < 21; i++ {
		tooLongSections[i] = openapi.SectionWithoutAutofield{
			Name:    testutil.RandomString(10),
			Content: testutil.RandomString(100),
		}
	}
	tooLongSectionItemsError := make([]openapi.SectionWithoutAutofieldError, 21)
	for i := 0; i < 21; i++ {
		tooLongSectionItemsError[i] = openapi.SectionWithoutAutofieldError{}
	}

	tt := []struct {
		name       string
		userId     string
		projectId  string
		chapterId  string
		sectionIds []string
		sections   []openapi.SectionWithoutAutofield
		expected   openapi.GraphSectionalizeErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			sectionIds: []string{
				"2000000000000001",
			},
			sections: []openapi.SectionWithoutAutofield{
				{
					Name:    "Section",
					Content: "This is the content of section. This is the content of section.",
				},
			},
			expected: openapi.GraphSectionalizeErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
				Sections: openapi.SectionWithoutAutofieldListError{
					Items: []openapi.SectionWithoutAutofieldError{
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
			sections: []openapi.SectionWithoutAutofield{
				{
					Name:    "Section",
					Content: "This is the content of section. This is the content of section.",
				},
			},
			expected: openapi.GraphSectionalizeErrorResponse{
				Project: openapi.ProjectOnlyIdError{
					Id: "project id is required, but got ''",
				},
				Sections: openapi.SectionWithoutAutofieldListError{
					Items: []openapi.SectionWithoutAutofieldError{
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
			sections: []openapi.SectionWithoutAutofield{
				{
					Name:    "Section",
					Content: "This is the content of section. This is the content of section.",
				},
			},
			expected: openapi.GraphSectionalizeErrorResponse{
				Chapter: openapi.ChapterOnlyIdError{
					Id: "chapter id is required, but got ''",
				},
				Sections: openapi.SectionWithoutAutofieldListError{
					Items: []openapi.SectionWithoutAutofieldError{
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
			sections:   []openapi.SectionWithoutAutofield{},
			expected: openapi.GraphSectionalizeErrorResponse{
				Sections: openapi.SectionWithoutAutofieldListError{
					Message: "sections are required, but got []",
					Items:   []openapi.SectionWithoutAutofieldError{},
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
			expected: openapi.GraphSectionalizeErrorResponse{
				Sections: openapi.SectionWithoutAutofieldListError{
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
			sections: []openapi.SectionWithoutAutofield{
				{
					Name:    "Section 1",
					Content: "This is the content of section 1. This is the content of section 1.",
				},
				{
					Name:    "",
					Content: "This is the content of section 2. This is the content of section 2.",
				},
			},
			expected: openapi.GraphSectionalizeErrorResponse{
				Sections: openapi.SectionWithoutAutofieldListError{
					Items: []openapi.SectionWithoutAutofieldError{
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
			sections: []openapi.SectionWithoutAutofield{
				{
					Name:    "Section 1",
					Content: "This is the content of section 1. This is the content of section 1.",
				},
				{
					Name:    tooLongSectionName,
					Content: "This is the content of section 2. This is the content of section 2.",
				},
			},
			expected: openapi.GraphSectionalizeErrorResponse{
				Sections: openapi.SectionWithoutAutofieldListError{
					Items: []openapi.SectionWithoutAutofieldError{
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
			sections: []openapi.SectionWithoutAutofield{
				{
					Name:    "Section 1",
					Content: tooLongSectionContent,
				},
				{
					Name:    "Section 2",
					Content: "This is the content of section 2. This is the content of section 2.",
				},
			},
			expected: openapi.GraphSectionalizeErrorResponse{
				Sections: openapi.SectionWithoutAutofieldListError{
					Items: []openapi.SectionWithoutAutofieldError{
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

			res, ucErr := uc.SectionalizeGraph(openapi.GraphSectionalizeRequest{
				User:     openapi.UserOnlyId{Id: tc.userId},
				Project:  openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter:  openapi.ChapterOnlyId{Id: tc.chapterId},
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
			errorCode:     service.InvalidArgumentError,
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
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			uc := usecase.NewGraphUseCase(s)

			res, ucErr := uc.SectionalizeGraph(openapi.GraphSectionalizeRequest{
				User:    openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: openapi.ProjectOnlyId{Id: "0000000000000001"},
				Chapter: openapi.ChapterOnlyId{Id: "1000000000000001"},
				Sections: []openapi.SectionWithoutAutofield{
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
