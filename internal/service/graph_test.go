package service_test

import (
	"fmt"
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/testutil"
	mock_repository "github.com/kumachan-mis/knodeledge-api/mock/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindGraphValidEntry(t *testing.T) {
	maxLengthGraphName := testutil.RandomString(100)
	maxLengthGraphParagraph := testutil.RandomString(40000)

	tt := []struct {
		name  string
		entry record.GraphEntry
	}{
		{
			name: "should return graph with valid entry",
			entry: record.GraphEntry{
				Name:      "Section",
				Paragraph: "This is section content. This is section content.",
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
		},
		{
			name: "should return graph with max-length valid entry",
			entry: record.GraphEntry{
				Name:      maxLengthGraphName,
				Paragraph: maxLengthGraphParagraph,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockGraphRepository(ctrl)
			r.EXPECT().
				FetchGraph(testutil.ReadOnlyUserId(), "0000000000000001", "1000000000000001", "2000000000000001").
				Return(&tc.entry, nil)

			cr := mock_repository.NewMockChapterRepository(ctrl)

			s := service.NewGraphService(r, cr)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)
			sectionId, err := domain.NewSectionIdObject("2000000000000001")
			assert.Nil(t, err)

			graph, sErr := s.FindGraph(*userId, *projectId, *chapterId, *sectionId)
			assert.Nil(t, sErr)

			assert.Equal(t, tc.entry.Paragraph, graph.Paragraph().Value())
			assert.Equal(t, tc.entry.CreatedAt, graph.CreatedAt().Value())
			assert.Equal(t, tc.entry.UpdatedAt, graph.UpdatedAt().Value())
		})
	}
}

func TestFindGraphInvalidEntry(t *testing.T) {
	tooLongGraphName := testutil.RandomString(101)
	tooLongGraphParagraph := testutil.RandomString(40001)

	tt := []struct {
		name          string
		entry         record.GraphEntry
		expectedError string
	}{
		{
			name: "should return error when name is empty",
			entry: record.GraphEntry{
				Name:      "",
				Paragraph: "This is section content. This is section content.",
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): " +
				"graph name is required, but got ''",
		},
		{
			name: "should return error when name is too long",
			entry: record.GraphEntry{
				Name:      tooLongGraphName,
				Paragraph: "This is section content. This is section content.",
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: fmt.Sprintf("failed to convert entry to entity (name): "+
				"graph name cannot be longer than 100 characters, but got '%v'", tooLongGraphName),
		},
		{
			name: "should return error when paragraph is too long",
			entry: record.GraphEntry{
				Name:      "Section",
				Paragraph: tooLongGraphParagraph,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (paragraph): " +
				"graph paragraph must be less than or equal to 40000 bytes, but got 40001 bytes",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockGraphRepository(ctrl)
			r.EXPECT().
				FetchGraph(testutil.ReadOnlyUserId(), "0000000000000001", "1000000000000001", "2000000000000001").
				Return(&tc.entry, nil)

			cr := mock_repository.NewMockChapterRepository(ctrl)

			s := service.NewGraphService(r, cr)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)
			sectionId, err := domain.NewSectionIdObject("2000000000000001")
			assert.Nil(t, err)

			graph, sErr := s.FindGraph(*userId, *projectId, *chapterId, *sectionId)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %v", tc.expectedError), sErr.Error())
			assert.Nil(t, graph)
		})
	}
}

func TestFindGraphRepositoryError(t *testing.T) {
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
			errorMessage:  "graph not found",
			expectedError: "failed to find graph: graph not found",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns read failure error",
			errorCode:     repository.ReadFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to fetch graph: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockGraphRepository(ctrl)
			r.EXPECT().
				FetchGraph(testutil.ReadOnlyUserId(), "0000000000000001", "1000000000000001", "2000000000000001").
				Return(nil, repository.Errorf(tc.errorCode, tc.errorMessage))

			cr := mock_repository.NewMockChapterRepository(ctrl)

			s := service.NewGraphService(r, cr)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)
			sectionId, err := domain.NewSectionIdObject("2000000000000001")
			assert.Nil(t, err)

			graph, sErr := s.FindGraph(*userId, *projectId, *chapterId, *sectionId)
			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, graph)
		})
	}
}

func TestSectionalizeIntoGraphsValidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sectonName := "Section"
	sectionContent := "This is section content. This is section content."
	maxLengthSectonName := testutil.RandomString(100)
	maxLengthSectionContent := testutil.RandomString(40000)

	r := mock_repository.NewMockGraphRepository(ctrl)
	r.EXPECT().
		GraphExists(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
		Return(false, nil)
	r.EXPECT().
		InsertGraphs(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001",
			[]record.GraphWithoutAutofieldEntry{
				{
					Name:      sectonName,
					Paragraph: sectionContent,
				},
				{
					Name:      maxLengthSectonName,
					Paragraph: maxLengthSectionContent,
				},
			}).
		Return([]string{"2000000000000001", "2000000000000002"}, []record.GraphEntry{
			{
				Name:      sectonName,
				Paragraph: sectionContent,
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			{
				Name:      maxLengthSectonName,
				Paragraph: maxLengthSectionContent,
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
		}, nil)

	chapterRepository := mock_repository.NewMockChapterRepository(ctrl)
	chapterRepository.EXPECT().
		UpdateChapterSections(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001",
			[]record.SectionWithoutAutofieldEntry{
				{
					Id:   "2000000000000001",
					Name: sectonName,
				},
				{
					Id:   "2000000000000002",
					Name: maxLengthSectonName,
				},
			}).
		Return([]record.SectionEntry{
			{
				Id:        "2000000000000001",
				Name:      sectonName,
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			{
				Id:        "2000000000000002",
				Name:      maxLengthSectonName,
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
		}, nil)
	s := service.NewGraphService(r, chapterRepository)

	userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
	assert.NoError(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.NoError(t, err)

	chapterId, err := domain.NewChapterIdObject("1000000000000001")
	assert.NoError(t, err)

	name, err := domain.NewSectionNameObject(sectonName)
	assert.NoError(t, err)
	content, err := domain.NewSectionContentObject(sectionContent)
	assert.NoError(t, err)
	section1 := domain.NewSectionWithoutAutofieldEntity(*name, *content)

	name, err = domain.NewSectionNameObject(maxLengthSectonName)
	assert.NoError(t, err)
	content, err = domain.NewSectionContentObject(maxLengthSectionContent)
	assert.NoError(t, err)
	section2 := domain.NewSectionWithoutAutofieldEntity(*name, *content)

	sections, err := domain.NewSectionWithoutAutofieldEntityList([]domain.SectionWithoutAutofieldEntity{*section1, *section2})
	assert.NoError(t, err)

	insertedGraphs, sErr := s.SectionalizeIntoGraphs(*userId, *projectId, *chapterId, *sections)
	assert.Nil(t, sErr)

	assert.Len(t, insertedGraphs, 2)

	graph := insertedGraphs[0]
	assert.Equal(t, "2000000000000001", graph.Id().Value())
	assert.Equal(t, sectonName, graph.Name().Value())
	assert.Equal(t, sectionContent, graph.Paragraph().Value())
	assert.Equal(t, testutil.Date(), graph.CreatedAt().Value())
	assert.Equal(t, testutil.Date(), graph.UpdatedAt().Value())

	graph = insertedGraphs[1]
	assert.Equal(t, "2000000000000002", graph.Id().Value())
	assert.Equal(t, maxLengthSectonName, graph.Name().Value())
	assert.Equal(t, maxLengthSectionContent, graph.Paragraph().Value())
	assert.Equal(t, testutil.Date(), graph.CreatedAt().Value())
	assert.Equal(t, testutil.Date(), graph.UpdatedAt().Value())
}

func TestSectionalizeIntoGraphsGraphExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sectonName := "Section"
	sectionContent := "This is section content. This is section content."

	r := mock_repository.NewMockGraphRepository(ctrl)
	r.EXPECT().
		GraphExists(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
		Return(true, nil)

	chapterRepository := mock_repository.NewMockChapterRepository(ctrl)

	s := service.NewGraphService(r, chapterRepository)

	userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
	assert.NoError(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.NoError(t, err)

	chapterId, err := domain.NewChapterIdObject("1000000000000001")
	assert.NoError(t, err)

	name, err := domain.NewSectionNameObject(sectonName)
	assert.NoError(t, err)
	content, err := domain.NewSectionContentObject(sectionContent)
	assert.NoError(t, err)
	section := domain.NewSectionWithoutAutofieldEntity(*name, *content)

	sections, err := domain.NewSectionWithoutAutofieldEntityList([]domain.SectionWithoutAutofieldEntity{*section})
	assert.NoError(t, err)

	insertedGraphs, sErr := s.SectionalizeIntoGraphs(*userId, *projectId, *chapterId, *sections)

	assert.NotNil(t, sErr)
	assert.Equal(t, service.InvalidArgument, sErr.Code())
	assert.Equal(t, "invalid argument: failed to sectionalize into graphs: graph already exists", sErr.Error())
	assert.Nil(t, insertedGraphs)
}

func TestSectionalizeIntoGraphsInvalidInsertedGraph(t *testing.T) {
	tooLongGraphName := testutil.RandomString(101)
	tooLongGraphParagraph := testutil.RandomString(40001)

	tt := []struct {
		name          string
		insertedGraph record.GraphEntry
		expectedError string
	}{
		{
			name: "should return error when graph name is empty",
			insertedGraph: record.GraphEntry{
				Name:      "",
				Paragraph: "This is section content. This is section content.",
			},
			expectedError: "failed to convert entry to entity (name): " +
				"graph name is required, but got ''",
		},
		{
			name: "should return error when graph name is too long",
			insertedGraph: record.GraphEntry{
				Name:      tooLongGraphName,
				Paragraph: "This is section content. This is section content.",
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: fmt.Sprintf("failed to convert entry to entity (name): "+
				"graph name cannot be longer than 100 characters, but got '%v'", tooLongGraphName),
		},
		{
			name: "should return error when graph paragraph is too long",
			insertedGraph: record.GraphEntry{
				Name:      "Section",
				Paragraph: tooLongGraphParagraph,
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (paragraph): " +
				"graph paragraph must be less than or equal to 40000 bytes, but got 40001 bytes",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sectonName := "Section"
			sectionContent := "This is section content. This is section content."

			r := mock_repository.NewMockGraphRepository(ctrl)
			r.EXPECT().
				GraphExists(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(false, nil)
			r.EXPECT().
				InsertGraphs(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001",
					[]record.GraphWithoutAutofieldEntry{
						{
							Name:      sectonName,
							Paragraph: sectionContent,
						},
					}).
				Return([]string{"2000000000000001"}, []record.GraphEntry{tc.insertedGraph}, nil)

			chapterRepository := mock_repository.NewMockChapterRepository(ctrl)

			s := service.NewGraphService(r, chapterRepository)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.NoError(t, err)

			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.NoError(t, err)

			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.NoError(t, err)

			name, err := domain.NewSectionNameObject(sectonName)
			assert.NoError(t, err)
			content, err := domain.NewSectionContentObject(sectionContent)
			assert.NoError(t, err)

			section := domain.NewSectionWithoutAutofieldEntity(*name, *content)

			sections, err := domain.NewSectionWithoutAutofieldEntityList([]domain.SectionWithoutAutofieldEntity{*section})
			assert.NoError(t, err)

			insertedGraphs, sErr := s.SectionalizeIntoGraphs(*userId, *projectId, *chapterId, *sections)

			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %v", tc.expectedError), sErr.Error())
			assert.Nil(t, insertedGraphs)
		})
	}
}

func TestSectionalizeIntoGraphsRepositoryGraphExistsError(t *testing.T) {
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
			errorMessage:  "chapter not found",
			expectedError: "failed to sectionalize into graphs: chapter not found",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns read failure error",
			errorCode:     repository.ReadFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to check graph existence: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sectonName := "Section"
			sectionContent := "This is section content. This is section content."

			r := mock_repository.NewMockGraphRepository(ctrl)
			r.EXPECT().
				GraphExists(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(false, repository.Errorf(tc.errorCode, tc.errorMessage))

			chapterRepository := mock_repository.NewMockChapterRepository(ctrl)

			s := service.NewGraphService(r, chapterRepository)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.NoError(t, err)

			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.NoError(t, err)

			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.NoError(t, err)

			name, err := domain.NewSectionNameObject(sectonName)
			assert.NoError(t, err)
			content, err := domain.NewSectionContentObject(sectionContent)
			assert.NoError(t, err)
			section := domain.NewSectionWithoutAutofieldEntity(*name, *content)

			sections, err := domain.NewSectionWithoutAutofieldEntityList([]domain.SectionWithoutAutofieldEntity{*section})
			assert.NoError(t, err)

			insertedGraphs, sErr := s.SectionalizeIntoGraphs(*userId, *projectId, *chapterId, *sections)

			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, insertedGraphs)
		})
	}
}

func TestSectionalizeIntoGraphsRepotoryInsertGraphsError(t *testing.T) {
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
			errorMessage:  "chapter not found",
			expectedError: "failed to sectionalize into graphs: chapter not found",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns write failure error",
			errorCode:     repository.WriteFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to insert graphs: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sectonName := "Section"
			sectionContent := "This is section content. This is section content."

			r := mock_repository.NewMockGraphRepository(ctrl)
			r.EXPECT().
				GraphExists(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(false, nil)
			r.EXPECT().
				InsertGraphs(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001",
					[]record.GraphWithoutAutofieldEntry{
						{
							Name:      sectonName,
							Paragraph: sectionContent,
						},
					}).
				Return(nil, nil, repository.Errorf(tc.errorCode, tc.errorMessage))

			chapterRepository := mock_repository.NewMockChapterRepository(ctrl)

			s := service.NewGraphService(r, chapterRepository)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.NoError(t, err)

			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.NoError(t, err)

			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.NoError(t, err)

			name, err := domain.NewSectionNameObject(sectonName)
			assert.NoError(t, err)
			content, err := domain.NewSectionContentObject(sectionContent)
			assert.NoError(t, err)
			section := domain.NewSectionWithoutAutofieldEntity(*name, *content)

			sections, err := domain.NewSectionWithoutAutofieldEntityList([]domain.SectionWithoutAutofieldEntity{*section})
			assert.NoError(t, err)

			insertedGraphs, sErr := s.SectionalizeIntoGraphs(*userId, *projectId, *chapterId, *sections)

			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, insertedGraphs)
		})
	}
}

func TestSectionalizeIntoGraphsChapterRepositoryError(t *testing.T) {
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
			errorMessage:  "chapter not found",
			expectedError: "failed to sectionalize into graphs: chapter not found",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns write failure error",
			errorCode:     repository.WriteFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to update sections of chapter: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sectonName := "Section"
			sectionContent := "This is section content. This is section content."

			r := mock_repository.NewMockGraphRepository(ctrl)
			r.EXPECT().
				GraphExists(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(false, nil)
			r.EXPECT().
				InsertGraphs(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001",
					[]record.GraphWithoutAutofieldEntry{
						{
							Name:      sectonName,
							Paragraph: sectionContent,
						},
					}).
				Return([]string{"2000000000000001"}, []record.GraphEntry{
					{
						Name:      sectonName,
						Paragraph: sectionContent,
						UserId:    testutil.ModifyOnlyUserId(),
						CreatedAt: testutil.Date(),
						UpdatedAt: testutil.Date(),
					},
				}, nil)

			chapterRepository := mock_repository.NewMockChapterRepository(ctrl)
			chapterRepository.EXPECT().
				UpdateChapterSections(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001",
					[]record.SectionWithoutAutofieldEntry{
						{
							Id:   "2000000000000001",
							Name: sectonName,
						},
					}).
				Return(nil, repository.Errorf(tc.errorCode, tc.errorMessage))

			s := service.NewGraphService(r, chapterRepository)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.NoError(t, err)

			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.NoError(t, err)

			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.NoError(t, err)

			name, err := domain.NewSectionNameObject(sectonName)
			assert.NoError(t, err)
			content, err := domain.NewSectionContentObject(sectionContent)
			assert.NoError(t, err)
			section := domain.NewSectionWithoutAutofieldEntity(*name, *content)

			sections, err := domain.NewSectionWithoutAutofieldEntityList([]domain.SectionWithoutAutofieldEntity{*section})
			assert.NoError(t, err)

			insertedGraphs, sErr := s.SectionalizeIntoGraphs(*userId, *projectId, *chapterId, *sections)

			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, insertedGraphs)
		})
	}
}
