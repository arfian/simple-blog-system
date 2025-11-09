package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"simple-blog-system/config/db"
	"simple-blog-system/internal/app/post/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository
}

func (suite *PostRepositoryTestSuite) SetupTest() {
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

func (suite *PostRepositoryTestSuite) TearDownTest() {
	sqlDB, err := suite.db.DB()
	assert.NoError(suite.T(), err)
	sqlDB.Close()
}

func TestPostRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostRepositoryTestSuite))
}

func (suite *PostRepositoryTestSuite) TestInsertPost_Success() {
	ctx := context.Background()
	now := time.Now()
	postID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	post := model.PostModel{
		ID:        postID,
		Username:  "testuser",
		Title:     "Test Post",
		Body:      "This is a test post body",
		Status:    "published",
		CreatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "posts"`)).
		WithArgs(
			post.Username,
			post.Title,
			post.Body,
			post.Status,
			post.CreatedBy,
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			post.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "updated_by", "deleted_at"}).
			AddRow(post.ID, nil, nil))
	suite.mock.ExpectCommit()

	result, err := suite.repository.InsertPost(ctx, post)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), post.Title, result.Title)
	assert.Equal(suite.T(), post.Username, result.Username)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestInsertPost_Error() {
	ctx := context.Background()
	now := time.Now()
	postID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	post := model.PostModel{
		ID:        postID,
		Username:  "testuser",
		Title:     "Test Post",
		Body:      "This is a test post body",
		Status:    "published",
		CreatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "posts"`)).
		WithArgs(
			post.Username,
			post.Title,
			post.Body,
			post.Status,
			post.CreatedBy,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			post.ID,
		).
		WillReturnError(gorm.ErrInvalidData)
	suite.mock.ExpectRollback()

	_, err := suite.repository.InsertPost(ctx, post)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidData, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestUpdatePost_Success() {
	ctx := context.Background()
	now := time.Now()
	postID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	post := model.PostModel{
		ID:        postID,
		Username:  "testuser",
		Title:     "Updated Post",
		Body:      "This is an updated post body",
		Status:    "published",
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts"`)).
		WithArgs(
			post.Username,
			post.Title,
			post.Body,
			post.Status,
			post.CreatedBy,
			post.UpdatedBy,
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			post.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	result, err := suite.repository.UpdatePost(ctx, post)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), post.Title, result.Title)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestUpdatePost_Error() {
	ctx := context.Background()
	now := time.Now()
	postID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	post := model.PostModel{
		ID:        postID,
		Username:  "testuser",
		Title:     "Updated Post",
		Body:      "This is an updated post body",
		Status:    "published",
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts"`)).
		WithArgs(
			post.Username,
			post.Title,
			post.Body,
			post.Status,
			post.CreatedBy,
			post.UpdatedBy,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			post.ID,
		).
		WillReturnError(gorm.ErrInvalidDB)
	suite.mock.ExpectRollback()

	_, err := suite.repository.UpdatePost(ctx, post)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestGetPostById_Success() {
	ctx := context.Background()
	postID := "123e4567-e89b-12d3-a456-426614174000"
	now := time.Now()

	expectedPost := model.PostModel{
		ID:        strfmt.UUID4(postID),
		Username:  "testuser",
		Title:     "Test Post",
		Body:      "This is a test post body",
		Status:    "published",
		CreatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	rows := sqlmock.NewRows([]string{"id", "username", "title", "body", "status", "created_by", "updated_by", "created_at", "updated_at"}).
		AddRow(expectedPost.ID, expectedPost.Username, expectedPost.Title, expectedPost.Body, expectedPost.Status, expectedPost.CreatedBy, nil, expectedPost.CreatedAt, expectedPost.UpdatedAt)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE id = $1 AND "posts"."deleted_at" IS NULL ORDER BY "posts"."id" LIMIT`)).
		WithArgs(postID, 1).
		WillReturnRows(rows)

	result, err := suite.repository.GetPostById(ctx, postID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedPost.Title, result.Title)
	assert.Equal(suite.T(), expectedPost.Username, result.Username)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestGetPostById_NotFound() {
	ctx := context.Background()
	postID := "nonexistent-id"

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE id = $1 AND "posts"."deleted_at" IS NULL ORDER BY "posts"."id" LIMIT`)).
		WithArgs(postID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := suite.repository.GetPostById(ctx, postID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
	assert.NotNil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestGetPostById_Error() {
	ctx := context.Background()
	postID := "123e4567-e89b-12d3-a456-426614174000"

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE id = $1 AND "posts"."deleted_at" IS NULL ORDER BY "posts"."id" LIMIT`)).
		WithArgs(postID, 1).
		WillReturnError(gorm.ErrInvalidDB)

	result, err := suite.repository.GetPostById(ctx, postID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.NotNil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestDeletePost_Success() {
	ctx := context.Background()
	postID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	post := model.PostModel{
		ID: postID,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts" SET "deleted_at"=$1 WHERE "posts"."id" = $2 AND "posts"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), post.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repository.DeletePost(ctx, post)

	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestDeletePost_Error() {
	ctx := context.Background()
	postID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	post := model.PostModel{
		ID: postID,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts" SET "deleted_at"=$1 WHERE "posts"."id" = $2 AND "posts"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), post.ID).
		WillReturnError(gorm.ErrInvalidDB)
	suite.mock.ExpectRollback()

	err := suite.repository.DeletePost(ctx, post)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestDeletePost_NotFound() {
	ctx := context.Background()
	postID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	post := model.PostModel{
		ID: postID,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts" SET "deleted_at"=$1 WHERE "posts"."id" = $2 AND "posts"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), post.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := suite.repository.DeletePost(ctx, post)

	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestGetAllPost_Success() {
	ctx := context.Background()
	page := 1
	limit := 10
	now := time.Now()

	post1 := model.PostModel{
		ID:        strfmt.UUID4("123e4567-e89b-12d3-a456-426614174001"),
		Username:  "user1",
		Title:     "Post 1",
		Body:      "Body 1",
		Status:    "published",
		CreatedBy: "user1",
		CreatedAt: now,
		UpdatedAt: now,
	}

	post2 := model.PostModel{
		ID:        strfmt.UUID4("123e4567-e89b-12d3-a456-426614174002"),
		Username:  "user2",
		Title:     "Post 2",
		Body:      "Body 2",
		Status:    "draft",
		CreatedBy: "user2",
		CreatedAt: now,
		UpdatedAt: now,
	}

	rows := sqlmock.NewRows([]string{"id", "username", "title", "body", "status", "created_by", "updated_by", "created_at", "updated_at"}).
		AddRow(post1.ID, post1.Username, post1.Title, post1.Body, post1.Status, post1.CreatedBy, nil, post1.CreatedAt, post1.UpdatedAt).
		AddRow(post2.ID, post2.Username, post2.Title, post2.Body, post2.Status, post2.CreatedBy, nil, post2.CreatedAt, post2.UpdatedAt)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE "posts"."deleted_at" IS NULL LIMIT`)).
		WithArgs(limit).
		WillReturnRows(rows)

	result, err := suite.repository.GetAllPost(ctx, page, limit)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), post1.Title, result[0].Title)
	assert.Equal(suite.T(), post2.Title, result[1].Title)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestGetAllPost_EmptyResult() {
	ctx := context.Background()
	page := 1
	limit := 10

	rows := sqlmock.NewRows([]string{"id", "username", "title", "body", "status", "created_by", "updated_by", "created_at", "updated_at"})

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE "posts"."deleted_at" IS NULL LIMIT`)).
		WithArgs(limit).
		WillReturnRows(rows)

	result, err := suite.repository.GetAllPost(ctx, page, limit)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestGetAllPost_Error() {
	ctx := context.Background()
	page := 1
	limit := 10

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE "posts"."deleted_at" IS NULL LIMIT`)).
		WithArgs(limit).
		WillReturnError(gorm.ErrInvalidDB)

	result, err := suite.repository.GetAllPost(ctx, page, limit)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostRepositoryTestSuite) TestGetAllPost_Pagination() {
	ctx := context.Background()
	page := 2
	limit := 5
	offset := 5

	rows := sqlmock.NewRows([]string{"id", "username", "title", "body", "status", "created_by", "updated_by", "created_at", "updated_at"})

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE "posts"."deleted_at" IS NULL LIMIT $1 OFFSET $2`)).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	result, err := suite.repository.GetAllPost(ctx, page, limit)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}
