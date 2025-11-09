package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"simple-blog-system/config"
	"simple-blog-system/internal/app/user/model"
	"simple-blog-system/pkg/encrypt"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock for IUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) InsertUser(ctx context.Context, user model.AuthUserModel) (model.AuthUserModel, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.AuthUserModel), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) ([]model.AuthUserModel, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AuthUserModel), args.Error(1)
}

func (m *MockUserRepository) GetPasswordByUsername(ctx context.Context, username string) ([]model.AuthUserModel, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AuthUserModel), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, user model.AuthUserModel) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// TestMain initializes the config before running tests
func TestMain(m *testing.M) {
	// Set required environment variables for testing
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("SIGNING_KEY", "test-secret-key-for-unit-tests")
	os.Setenv("JWT_SIGNING_KEY", "test-secret-key-for-unit-tests")
	os.Setenv("DB_DSN", "test-db-dsn")
	os.Setenv("DB_MAX_OPEN_CONN", "10")
	os.Setenv("DB_MAX_IDLE_CONN", "5")
	os.Setenv("DB_MAX_LIFETIME_CONN", "60")
	os.Setenv("DB_MAX_IDLETIME_CONN", "30")
	
	// Initialize config
	config.InitConfig()
	
	// Run tests
	code := m.Run()
	os.Exit(code)
}

// Test Suite
type UserServiceTestSuite struct {
	suite.Suite
	service  *service
	userRepo *MockUserRepository
	ctx      context.Context
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.userRepo = new(MockUserRepository)
	suite.service = &service{
		userRepo: suite.userRepo,
	}
	suite.ctx = context.Background()
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestRegister_Success() {
	username := "newuser"
	password := "password123"

	user := model.AuthUserModel{
		Username: username,
		Password: password,
	}

	// Mock: User doesn't exist yet
	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]model.AuthUserModel{}, nil)

	// Mock: User insertion succeeds
	suite.userRepo.On("InsertUser", suite.ctx, mock.MatchedBy(func(u model.AuthUserModel) bool {
		return u.Username == username && u.CreatedBy == username && u.Password != password
	})).Return(model.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}, nil)

	token, err := suite.service.Register(suite.ctx, user)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestRegister_UserAlreadyExists() {
	username := "existinguser"
	password := "password123"

	user := model.AuthUserModel{
		Username: username,
		Password: password,
	}

	existingUser := model.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
	}

	// Mock: User already exists
	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]model.AuthUserModel{existingUser}, nil)

	token, err := suite.service.Register(suite.ctx, user)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Equal(suite.T(), "user already exists", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestRegister_GetUserError() {
	username := "newuser"
	password := "password123"

	user := model.AuthUserModel{
		Username: username,
		Password: password,
	}

	// Mock: Database error when checking user
	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]model.AuthUserModel{}, errors.New("database error"))

	token, err := suite.service.Register(suite.ctx, user)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Equal(suite.T(), "database error", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestRegister_InsertUserError() {
	username := "newuser"
	password := "password123"

	user := model.AuthUserModel{
		Username: username,
		Password: password,
	}

	// Mock: User doesn't exist
	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]model.AuthUserModel{}, nil)

	// Mock: Insert fails
	suite.userRepo.On("InsertUser", suite.ctx, mock.Anything).Return(model.AuthUserModel{}, errors.New("insert error"))

	token, err := suite.service.Register(suite.ctx, user)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Equal(suite.T(), "insert error", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestLogin_Success() {
	username := "testuser"
	password := "password123"
	hashedPassword, _ := encrypt.HashPassword(password)

	user := model.AuthUserModel{
		Username: username,
		Password: password,
	}

	existingUser := model.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
		Password: hashedPassword,
	}

	// Mock: Get user with password
	suite.userRepo.On("GetPasswordByUsername", suite.ctx, username).Return([]model.AuthUserModel{existingUser}, nil)

	// Mock: Update last login
	suite.userRepo.On("UpdateLastLogin", suite.ctx, mock.MatchedBy(func(u model.AuthUserModel) bool {
		return u.Username == username && u.UpdatedBy == username
	})).Return(nil)

	token, err := suite.service.Login(suite.ctx, user)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestLogin_UserNotFound() {
	username := "nonexistent"
	password := "password123"

	user := model.AuthUserModel{
		Username: username,
		Password: password,
	}

	// Mock: User not found
	suite.userRepo.On("GetPasswordByUsername", suite.ctx, username).Return([]model.AuthUserModel{}, nil)

	token, err := suite.service.Login(suite.ctx, user)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Equal(suite.T(), "incorrect username or password", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestLogin_DatabaseError() {
	username := "testuser"
	password := "password123"

	user := model.AuthUserModel{
		Username: username,
		Password: password,
	}

	// Mock: Database error
	suite.userRepo.On("GetPasswordByUsername", suite.ctx, username).Return([]model.AuthUserModel{}, errors.New("database error"))

	token, err := suite.service.Login(suite.ctx, user)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Equal(suite.T(), "incorrect username or password", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestLogin_IncorrectPassword() {
	username := "testuser"
	password := "wrongpassword"
	correctPassword := "password123"
	hashedPassword, _ := encrypt.HashPassword(correctPassword)

	user := model.AuthUserModel{
		Username: username,
		Password: password,
	}

	existingUser := model.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
		Password: hashedPassword,
	}

	// Mock: Get user with password
	suite.userRepo.On("GetPasswordByUsername", suite.ctx, username).Return([]model.AuthUserModel{existingUser}, nil)

	token, err := suite.service.Login(suite.ctx, user)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Equal(suite.T(), "incorrect username or password", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestLogin_UpdateLastLoginError() {
	username := "testuser"
	password := "password123"
	hashedPassword, _ := encrypt.HashPassword(password)

	user := model.AuthUserModel{
		Username: username,
		Password: password,
	}

	existingUser := model.AuthUserModel{
		ID:       strfmt.UUID4("user-123"),
		Username: username,
		Password: hashedPassword,
	}

	// Mock: Get user with password
	suite.userRepo.On("GetPasswordByUsername", suite.ctx, username).Return([]model.AuthUserModel{existingUser}, nil)

	// Mock: Update last login fails
	suite.userRepo.On("UpdateLastLogin", suite.ctx, mock.Anything).Return(errors.New("update error"))

	token, err := suite.service.Login(suite.ctx, user)

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Equal(suite.T(), "update error", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetUser_Success() {
	username := "testuser"
	now := time.Now()

	existingUser := model.AuthUserModel{
		ID:        strfmt.UUID4("user-123"),
		Username:  username,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Mock: Get user by username
	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]model.AuthUserModel{existingUser}, nil)

	result, err := suite.service.GetUser(suite.ctx, username)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), username, result.User.Username)
	assert.Equal(suite.T(), existingUser.ID, result.User.ID)
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetUser_UserNotFound() {
	username := "nonexistent"

	// Mock: User not found
	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]model.AuthUserModel{}, nil)

	result, err := suite.service.GetUser(suite.ctx, username)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetUser_DatabaseError() {
	username := "testuser"

	// Mock: Database error
	suite.userRepo.On("GetUserByUsername", suite.ctx, username).Return([]model.AuthUserModel{}, errors.New("database error"))

	result, err := suite.service.GetUser(suite.ctx, username)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.userRepo.AssertExpectations(suite.T())
}

