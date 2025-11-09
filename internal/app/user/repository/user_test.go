package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"simple-blog-system/config/db"
	"simple-blog-system/internal/app/user/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	var (
		sqlDB *sql.DB
		err   error
	)

	sqlDB, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	assert.NoError(suite.T(), err)

	gormDB := &db.GormDB{DB: suite.db}
	suite.repository = repository{db: gormDB}
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	sqlDB, err := suite.db.DB()
	assert.NoError(suite.T(), err)
	sqlDB.Close()
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestInsertUser_Success() {
	ctx := context.Background()
	now := time.Now()
	userID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	user := model.AuthUserModel{
		ID:        userID,
		Username:  "testuser",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedBy: "admin",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "auth_user"`)).
		WithArgs(
			user.Username,
			user.Password,
			user.IsActive,
			sqlmock.AnyArg(), // LastLogin
			user.CreatedBy,
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			user.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "updated_by", "deleted_at"}).
			AddRow(user.ID, nil, nil))
	suite.mock.ExpectCommit()

	result, err := suite.repository.InsertUser(ctx, user)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Username, result.Username)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestInsertUser_Error() {
	ctx := context.Background()
	now := time.Now()
	userID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	user := model.AuthUserModel{
		ID:        userID,
		Username:  "testuser",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedBy: "admin",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "auth_user"`)).
		WithArgs(
			user.Username,
			user.Password,
			user.IsActive,
			sqlmock.AnyArg(),
			user.CreatedBy,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			user.ID,
		).
		WillReturnError(gorm.ErrInvalidData)
	suite.mock.ExpectRollback()

	_, err := suite.repository.InsertUser(ctx, user)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidData, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserByUsername_Success() {
	ctx := context.Background()
	username := "testuser"
	now := time.Now()
	userID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	expectedUser := model.AuthUserModel{
		ID:        userID,
		Username:  username,
		CreatedAt: now,
		UpdatedAt: now,
	}

	rows := sqlmock.NewRows([]string{"id", "username", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, created_at, updated_at FROM "auth_user" WHERE username = $1`)).
		WithArgs(username).
		WillReturnRows(rows)

	result, err := suite.repository.GetUserByUsername(ctx, username)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), username, result[0].Username)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserByUsername_NotFound() {
	ctx := context.Background()
	username := "nonexistent"

	rows := sqlmock.NewRows([]string{"id", "username", "created_at", "updated_at"})

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, created_at, updated_at FROM "auth_user" WHERE username = $1`)).
		WithArgs(username).
		WillReturnRows(rows)

	result, err := suite.repository.GetUserByUsername(ctx, username)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserByUsername_Error() {
	ctx := context.Background()
	username := "testuser"

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, created_at, updated_at FROM "auth_user" WHERE username = $1`)).
		WithArgs(username).
		WillReturnError(gorm.ErrInvalidDB)

	result, err := suite.repository.GetUserByUsername(ctx, username)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetPasswordByUsername_Success() {
	ctx := context.Background()
	username := "testuser"
	now := time.Now()
	userID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	expectedUser := model.AuthUserModel{
		ID:        userID,
		Username:  username,
		Password:  "hashedpassword",
		CreatedAt: now,
		UpdatedAt: now,
	}

	rows := sqlmock.NewRows([]string{"id", "password", "username", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Password, expectedUser.Username, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, username, created_at, updated_at FROM "auth_user" WHERE username = $1`)).
		WithArgs(username).
		WillReturnRows(rows)

	result, err := suite.repository.GetPasswordByUsername(ctx, username)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), username, result[0].Username)
	assert.Equal(suite.T(), "hashedpassword", result[0].Password)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetPasswordByUsername_NotFound() {
	ctx := context.Background()
	username := "nonexistent"

	rows := sqlmock.NewRows([]string{"id", "password", "username", "created_at", "updated_at"})

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, username, created_at, updated_at FROM "auth_user" WHERE username = $1`)).
		WithArgs(username).
		WillReturnRows(rows)

	result, err := suite.repository.GetPasswordByUsername(ctx, username)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetPasswordByUsername_Error() {
	ctx := context.Background()
	username := "testuser"

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, username, created_at, updated_at FROM "auth_user" WHERE username = $1`)).
		WithArgs(username).
		WillReturnError(gorm.ErrInvalidDB)

	result, err := suite.repository.GetPasswordByUsername(ctx, username)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdateLastLogin_Success() {
	ctx := context.Background()
	username := "testuser"
	lastLogin := time.Now()

	user := model.AuthUserModel{
		Username:  username,
		LastLogin: lastLogin,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth_user" SET "last_login"=$1,"updated_at"=$2 WHERE username = $3`)).
		WithArgs(lastLogin, sqlmock.AnyArg(), username).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repository.UpdateLastLogin(ctx, user)

	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdateLastLogin_Error() {
	ctx := context.Background()
	username := "testuser"
	lastLogin := time.Now()

	user := model.AuthUserModel{
		Username:  username,
		LastLogin: lastLogin,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth_user" SET "last_login"=$1,"updated_at"=$2 WHERE username = $3`)).
		WithArgs(lastLogin, sqlmock.AnyArg(), username).
		WillReturnError(gorm.ErrInvalidDB)
	suite.mock.ExpectRollback()

	err := suite.repository.UpdateLastLogin(ctx, user)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdateLastLogin_NoRowsAffected() {
	ctx := context.Background()
	username := "nonexistent"
	lastLogin := time.Now()

	user := model.AuthUserModel{
		Username:  username,
		LastLogin: lastLogin,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth_user" SET "last_login"=$1,"updated_at"=$2 WHERE username = $3`)).
		WithArgs(lastLogin, sqlmock.AnyArg(), username).
		WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := suite.repository.UpdateLastLogin(ctx, user)

	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}
