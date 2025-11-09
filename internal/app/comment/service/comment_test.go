package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"simple-blog-system/internal/app/comment/model"
	"simple-blog-system/internal/app/comment/payload"
	postModel "simple-blog-system/internal/app/post/model"
	userModel "simple-blog-system/internal/app/user/model"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// Mock for ICommentRepository
type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) InsertComment(ctx context.Context, comment model.CommentModel) (model.CommentModel, error) {
	args := m.Called(ctx, comment)
	return args.Get(0).(model.CommentModel), args.Error(1)
}

func (m *MockCommentRepository) UpdateComment(ctx context.Context, comment model.CommentModel) (model.CommentModel, error) {
	args := m.Called(ctx, comment)
	return args.Get(0).(model.CommentModel), args.Error(1)
}

func (m *MockCommentRepository) GetCommentById(ctx context.Context, id string) (*model.CommentModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommentModel), args.Error(1)
}

func (m *MockCommentRepository) DeleteComment(ctx context.Context, comment model.CommentModel) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) GetAllComment(ctx context.Context, page int, limit int) ([]model.CommentModel, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CommentModel), args.Error(1)
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

// Mock for IPostRepository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) InsertPost(ctx context.Context, post postModel.PostModel) (postModel.PostModel, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(postModel.PostModel), args.Error(1)
}

func (m *MockPostRepository) UpdatePost(ctx context.Context, post postModel.PostModel) (postModel.PostModel, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(postModel.PostModel), args.Error(1)
}

func (m *MockPostRepository) GetPostById(ctx context.Context, id string) (*postModel.PostModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*postModel.PostModel), args.Error(1)
}

func (m *MockPostRepository) DeletePost(ctx context.Context, post postModel.PostModel) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) GetAllPost(ctx context.Context, page int, limit int) ([]postModel.PostModel, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postModel.PostModel), args.Error(1)
}

// Test Suite
type CommentServiceTestSuite struct {
	suite.Suite
	service     *service
	commentRepo *MockCommentRepository
	userRepo    *MockUserRepository
	postRepo    *MockPostRepository
	ctx         context.Context
}

func (suite *CommentServiceTestSuite) SetupTest() {
	suite.commentRepo = new(MockCommentRepository)
	suite.userRepo = new(MockUserRepository)
	suite.postRepo = new(MockPostRepository)
	suite.service = &service{
		commentRepo: suite.commentRepo,
		userRepo:    suite.userRepo,
		postRepo:    suite.postRepo,
	}
	suite.ctx = context.Background()
}

func TestCommentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CommentServiceTestSuite))
}

func (suite *CommentServiceTestSuite) TestAddComment_Success() {
	username := "testuser"
	now := time.Now()
	
	param := payload.CommentRequest{
		Comment: "Test comment",
		PostId:  "post-123",
	}

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	comment := model.CommentModel{
		ID:        strfmt.UUID4("comment-123"),
		Username:  username,
		Comment:   param.Comment,
		PostId:    param.PostId,
		CreatedBy: username,
		CreatedAt: now,
	}

	post := postModel.PostModel{
		ID:       strfmt.UUID4("post-123"),
		Username: "postuser",
		Title:    "Test Post",
		Body:     "Test Body",
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("InsertComment", suite.ctx, mock.MatchedBy(func(c model.CommentModel) bool {
		return c.Username == username && c.Comment == param.Comment && c.PostId == param.PostId
	})).Return(comment, nil)
	suite.postRepo.On("GetPostById", suite.ctx, param.PostId).Return(&post, nil)

	result, err := suite.service.AddComment(suite.ctx, username, param)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), comment.Comment, result.Comment)
	assert.Equal(suite.T(), post.Title, result.Post.Title)
	suite.commentRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestAddComment_UserNotFound() {
	username := "nonexistent"
	param := payload.CommentRequest{
		Comment: "Test comment",
		PostId:  "post-123",
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.AddComment(suite.ctx, username, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestAddComment_UserRepoError() {
	username := "testuser"
	param := payload.CommentRequest{
		Comment: "Test comment",
		PostId:  "post-123",
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, errors.New("database error"))

	result, err := suite.service.AddComment(suite.ctx, username, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestAddComment_InsertCommentError() {
	username := "testuser"
	param := payload.CommentRequest{
		Comment: "Test comment",
		PostId:  "post-123",
	}

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("InsertComment", suite.ctx, mock.Anything).Return(model.CommentModel{}, errors.New("insert error"))

	result, err := suite.service.AddComment(suite.ctx, username, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "insert error", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestAddComment_GetPostError() {
	username := "testuser"
	now := time.Now()
	
	param := payload.CommentRequest{
		Comment: "Test comment",
		PostId:  "post-123",
	}

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	comment := model.CommentModel{
		ID:        strfmt.UUID4("comment-123"),
		Username:  username,
		Comment:   param.Comment,
		PostId:    param.PostId,
		CreatedBy: username,
		CreatedAt: now,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("InsertComment", suite.ctx, mock.Anything).Return(comment, nil)
	suite.postRepo.On("GetPostById", suite.ctx, param.PostId).Return(nil, errors.New("post not found"))

	result, err := suite.service.AddComment(suite.ctx, username, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "post not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestUpdateComment_Success() {
	username := "testuser"
	commentID := "comment-123"
	now := time.Now()
	
	param := payload.CommentRequest{
		Comment: "Updated comment",
		PostId:  "post-123",
	}

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	comment := model.CommentModel{
		ID:        strfmt.UUID4(commentID),
		Username:  username,
		Comment:   param.Comment,
		PostId:    param.PostId,
		CreatedBy: username,
		UpdatedAt: now,
	}

	post := postModel.PostModel{
		ID:       strfmt.UUID4("post-123"),
		Username: "postuser",
		Title:    "Test Post",
		Body:     "Test Body",
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("UpdateComment", suite.ctx, mock.MatchedBy(func(c model.CommentModel) bool {
		return c.Username == username && c.Comment == param.Comment && c.PostId == param.PostId
	})).Return(comment, nil)
	suite.postRepo.On("GetPostById", suite.ctx, param.PostId).Return(&post, nil)

	result, err := suite.service.UpdateComment(suite.ctx, username, commentID, param)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), comment.Comment, result.Comment)
	assert.Equal(suite.T(), post.Title, result.Post.Title)
	suite.commentRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
	suite.postRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestUpdateComment_UserNotFound() {
	username := "nonexistent"
	commentID := "comment-123"
	param := payload.CommentRequest{
		Comment: "Updated comment",
		PostId:  "post-123",
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.UpdateComment(suite.ctx, username, commentID, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestUpdateComment_UpdateError() {
	username := "testuser"
	commentID := "comment-123"
	param := payload.CommentRequest{
		Comment: "Updated comment",
		PostId:  "post-123",
	}

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("UpdateComment", suite.ctx, mock.Anything).Return(model.CommentModel{}, errors.New("update error"))

	result, err := suite.service.UpdateComment(suite.ctx, username, commentID, param)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "update error", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestDeleteComment_Success() {
	username := "testuser"
	commentID := "comment-123"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	comment := model.CommentModel{
		ID:        strfmt.UUID4(commentID),
		Username:  username,
		Comment:   "Test comment",
		PostId:    "post-123",
		CreatedBy: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("GetCommentById", suite.ctx, commentID).Return(&comment, nil)
	suite.commentRepo.On("DeleteComment", suite.ctx, comment).Return(nil)

	result, err := suite.service.DeleteComment(suite.ctx, username, commentID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), comment.ID, result.ID)
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestDeleteComment_UserNotFound() {
	username := "nonexistent"
	commentID := "comment-123"

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.DeleteComment(suite.ctx, username, commentID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestDeleteComment_CommentNotFound() {
	username := "testuser"
	commentID := "nonexistent"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("GetCommentById", suite.ctx, commentID).Return(nil, gorm.ErrRecordNotFound)

	result, err := suite.service.DeleteComment(suite.ctx, username, commentID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "comment not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestDeleteComment_DeleteError() {
	username := "testuser"
	commentID := "comment-123"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	comment := model.CommentModel{
		ID:        strfmt.UUID4(commentID),
		Username:  username,
		Comment:   "Test comment",
		PostId:    "post-123",
		CreatedBy: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("GetCommentById", suite.ctx, commentID).Return(&comment, nil)
	suite.commentRepo.On("DeleteComment", suite.ctx, comment).Return(errors.New("delete error"))

	result, err := suite.service.DeleteComment(suite.ctx, username, commentID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "delete error", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestGetAllComment_Success() {
	username := "testuser"
	page := 1
	limit := 10

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	comments := []model.CommentModel{
		{
			ID:       strfmt.UUID4("comment-1"),
			Username: username,
			Comment:  "Comment 1",
			PostId:   "post-1",
		},
		{
			ID:       strfmt.UUID4("comment-2"),
			Username: username,
			Comment:  "Comment 2",
			PostId:   "post-2",
		},
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("GetAllComment", suite.ctx, page, limit).Return(comments, nil)

	result, err := suite.service.GetAllComment(suite.ctx, username, page, limit)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), comments[0].Comment, result[0].Comment)
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestGetAllComment_UserNotFound() {
	username := "nonexistent"
	page := 1
	limit := 10

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.GetAllComment(suite.ctx, username, page, limit)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestGetAllComment_Error() {
	username := "testuser"
	page := 1
	limit := 10

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("GetAllComment", suite.ctx, page, limit).Return(nil, errors.New("database error"))

	result, err := suite.service.GetAllComment(suite.ctx, username, page, limit)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "comment not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestGetCommentById_Success() {
	username := "testuser"
	commentID := "comment-123"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	comment := model.CommentModel{
		ID:        strfmt.UUID4(commentID),
		Username:  username,
		Comment:   "Test comment",
		PostId:    "post-123",
		CreatedBy: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("GetCommentById", suite.ctx, commentID).Return(&comment, nil)

	result, err := suite.service.GetCommentById(suite.ctx, username, commentID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), comment.Comment, result.Comment)
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestGetCommentById_UserNotFound() {
	username := "nonexistent"
	commentID := "comment-123"

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{}, nil)

	result, err := suite.service.GetCommentById(suite.ctx, username, commentID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CommentServiceTestSuite) TestGetCommentById_CommentNotFound() {
	username := "testuser"
	commentID := "nonexistent"

	user := userModel.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]userModel.AuthUserModel{user}, nil)
	suite.commentRepo.On("GetCommentById", suite.ctx, commentID).Return(nil, gorm.ErrRecordNotFound)

	result, err := suite.service.GetCommentById(suite.ctx, username, commentID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "post not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
	suite.commentRepo.AssertExpectations(suite.T())
}
