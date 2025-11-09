package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"simple-blog-system/internal/app/post/model"
	"simple-blog-system/internal/app/post/payload"
	userModel "simple-blog-system/internal/app/user/model"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// Mock for IPostRepository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) InsertPost(ctx context.Context, post model.PostModel) (model.PostModel, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(model.PostModel), args.Error(1)
}

func (m *MockPostRepository) UpdatePost(ctx context.Context, post model.PostModel) (model.PostModel, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(model.PostModel), args.Error(1)
}

func (m *MockPostRepository) GetPostById(ctx context.Context, id string) (*model.PostModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PostModel), args.Error(1)
}

func (m *MockPostRepository) DeletePost(ctx context.Context, post model.PostModel) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) GetAllPost(ctx context.Context, page int, limit int) ([]model.PostModel, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PostModel), args.Error(1)
}

// Mock for IUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) InsertUser(ctx context.Context, user userModel.AuthUserModel) (userModel.AuthUserModel, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(userModel.AuthUserModel), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) ([]userModel.AuthUserModel, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userModel.AuthUserModel), args.Error(1)
}

func (m *MockUserRepository) GetPasswordByUsername(ctx context.Context, username string) ([]userModel.AuthUserModel, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userModel.AuthUserModel), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, user userModel.AuthUserModel) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Test Suite
type PostServiceTestSuite struct {
	suite.Suite
	service  *service
	postRepo *MockPostRepository
	userRepo *MockUserRepository
	ctx      context.Context
}

func (suite *PostServiceTestSuite) SetupTest() {
	suite.postRepo = new(MockPostRepository)
	suite.userRepo = new(MockUserRepository)
	suite.service = &service{
		postRepo: suite.postRepo,
		userRepo: suite.userRepo,
	}
	suite.ctx = context.Background()
}

func TestPostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceTestSuite))
}

func (suite *PostServiceTestSuite) TestAddPost_Success() {
	username := "testuser"
	now := time.Now()

	param := payload.PostRequest{
		Title:  "Test Post",
		Body:   "This is a test post body",
		Status: "PUBLISH",
	}

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	post := model.PostModel{
		ID:        strfmt.UUID4("post-123"),
		Username:  username,
		Title:     param.Title,
		Body:      param.Body,
		Status:    param.Status,
		CreatedBy: username,
		CreatedAt: now,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("InsertPost", suite.ctx, mock.MatchedBy(func(p model.PostModel) bool {
		return p.Username == username && p.Title == param.Title && p.Body == param.Body && p.Status == param.Status
	})).Return(post, nil)

	result, err := suite.service.AddPost(suite.ctx, username, param)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), post.Title, result.Title)
	assert.Equal(suite.T(), post.Body, result.Body)
	suite.postRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestAddPost_UserNotFound() {
	username := "nonexistent"
	param := payload.PostRequest{
		Title:  "Test Post",
		Body:   "This is a test post body",
		Status: "PUBLISH",
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.AddPost(suite.ctx, username, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestAddPost_UserRepoError() {
	username := "testuser"
	param := payload.PostRequest{
		Title:  "Test Post",
		Body:   "This is a test post body",
		Status: "PUBLISH",
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, errors.New("database error"))

	result, err := suite.service.AddPost(suite.ctx, username, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestAddPost_InsertPostError() {
	username := "testuser"
	param := payload.PostRequest{
		Title:  "Test Post",
		Body:   "This is a test post body",
		Status: "PUBLISH",
	}

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("InsertPost", suite.ctx, mock.Anything).Return(model.PostModel{}, errors.New("insert error"))

	result, err := suite.service.AddPost(suite.ctx, username, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "insert error", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestUpdatePost_Success() {
	username := "testuser"
	postID := "post-123"
	now := time.Now()

	param := payload.PostRequest{
		Title:  "Updated Post",
		Body:   "This is an updated post body",
		Status: "DRAFT",
	}

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	post := model.PostModel{
		ID:        strfmt.UUID4(postID),
		Username:  username,
		Title:     param.Title,
		Body:      param.Body,
		Status:    param.Status,
		CreatedBy: username,
		UpdatedAt: now,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("UpdatePost", suite.ctx, mock.MatchedBy(func(p model.PostModel) bool {
		return p.Username == username && p.Title == param.Title && p.Body == param.Body && p.Status == param.Status
	})).Return(post, nil)

	result, err := suite.service.UpdatePost(suite.ctx, username, postID, param)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), post.Title, result.Title)
	assert.Equal(suite.T(), post.Body, result.Body)
	suite.postRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestUpdatePost_UserNotFound() {
	username := "nonexistent"
	postID := "post-123"
	param := payload.PostRequest{
		Title:  "Updated Post",
		Body:   "This is an updated post body",
		Status: "DRAFT",
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.UpdatePost(suite.ctx, username, postID, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestUpdatePost_UpdateError() {
	username := "testuser"
	postID := "post-123"
	param := payload.PostRequest{
		Title:  "Updated Post",
		Body:   "This is an updated post body",
		Status: "DRAFT",
	}

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("UpdatePost", suite.ctx, mock.Anything).Return(model.PostModel{}, errors.New("update error"))

	result, err := suite.service.UpdatePost(suite.ctx, username, postID, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "update error", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestDeletePost_Success() {
	username := "testuser"
	postID := "post-123"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	post := model.PostModel{
		ID:        strfmt.UUID4(postID),
		Username:  username,
		Title:     "Test Post",
		Body:      "Test Body",
		Status:    "PUBLISH",
		CreatedBy: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("GetPostById", suite.ctx, postID).Return(&post, nil)
	suite.postRepo.On("DeletePost", suite.ctx, post).Return(nil)

	result, err := suite.service.DeletePost(suite.ctx, username, postID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), post.ID, result.ID)
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestDeletePost_UserNotFound() {
	username := "nonexistent"
	postID := "post-123"

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.DeletePost(suite.ctx, username, postID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestDeletePost_PostNotFound() {
	username := "testuser"
	postID := "nonexistent"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("GetPostById", suite.ctx, postID).Return(nil, gorm.ErrRecordNotFound)

	result, err := suite.service.DeletePost(suite.ctx, username, postID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "post not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestDeletePost_DeleteError() {
	username := "testuser"
	postID := "post-123"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	post := model.PostModel{
		ID:        strfmt.UUID4(postID),
		Username:  username,
		Title:     "Test Post",
		Body:      "Test Body",
		Status:    "PUBLISH",
		CreatedBy: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("GetPostById", suite.ctx, postID).Return(&post, nil)
	suite.postRepo.On("DeletePost", suite.ctx, post).Return(errors.New("delete error"))

	result, err := suite.service.DeletePost(suite.ctx, username, postID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "delete error", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetAllPost_Success() {
	username := "testuser"
	page := 1
	limit := 10

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	posts := []model.PostModel{
		{
			ID:       strfmt.UUID4("post-1"),
			Username: username,
			Title:    "Post 1",
			Body:     "Body 1",
			Status:   "PUBLISH",
		},
		{
			ID:       strfmt.UUID4("post-2"),
			Username: username,
			Title:    "Post 2",
			Body:     "Body 2",
			Status:   "DRAFT",
		},
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("GetAllPost", suite.ctx, page, limit).Return(posts, nil)

	result, err := suite.service.GetAllPost(suite.ctx, username, page, limit)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), posts[0].Title, result[0].Title)
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetAllPost_UserNotFound() {
	username := "nonexistent"
	page := 1
	limit := 10

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.GetAllPost(suite.ctx, username, page, limit)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetAllPost_Error() {
	username := "testuser"
	page := 1
	limit := 10

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("GetAllPost", suite.ctx, page, limit).Return(nil, errors.New("database error"))

	result, err := suite.service.GetAllPost(suite.ctx, username, page, limit)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "post not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetById_Success() {
	username := "testuser"
	postID := "post-123"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	post := model.PostModel{
		ID:        strfmt.UUID4(postID),
		Username:  username,
		Title:     "Test Post",
		Body:      "Test Body",
		Status:    "PUBLISH",
		CreatedBy: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("GetPostById", suite.ctx, postID).Return(&post, nil)

	result, err := suite.service.GetById(suite.ctx, username, postID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), post.Title, result.Title)
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetById_UserNotFound() {
	username := "nonexistent"
	postID := "post-123"

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.GetById(suite.ctx, username, postID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetById_PostNotFound() {
	username := "testuser"
	postID := "nonexistent"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.postRepo.On("GetPostById", suite.ctx, postID).Return(nil, gorm.ErrRecordNotFound)

	result, err := suite.service.GetById(suite.ctx, username, postID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "post not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}
