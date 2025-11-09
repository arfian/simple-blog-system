package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"simple-blog-system/config/db"
	"simple-blog-system/internal/app/comment/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CommentRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository
}

func (suite *CommentRepositoryTestSuite) SetupTest() {
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

func (suite *CommentRepositoryTestSuite) TearDownTest() {
	sqlDB, err := suite.db.DB()
	assert.NoError(suite.T(), err)
	sqlDB.Close()
}

func TestCommentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CommentRepositoryTestSuite))
}

func (suite *CommentRepositoryTestSuite) TestInsertComment_Success() {
	ctx := context.Background()
	now := time.Now()
	commentID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	comment := model.CommentModel{
		ID:        commentID,
		Username:  "testuser",
		Comment:   "This is a test comment",
		PostId:    "post-123",
		CreatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "comments"`)).
		WithArgs(
			comment.Username,
			comment.Comment,
			comment.PostId,
			comment.CreatedBy,
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			comment.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "updated_by", "deleted_at"}).
			AddRow(comment.ID, nil, nil))
	suite.mock.ExpectCommit()

	result, err := suite.repository.InsertComment(ctx, comment)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), comment.Comment, result.Comment)
	assert.Equal(suite.T(), comment.Username, result.Username)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestInsertComment_Error() {
	ctx := context.Background()
	now := time.Now()
	commentID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	comment := model.CommentModel{
		ID:        commentID,
		Username:  "testuser",
		Comment:   "This is a test comment",
		PostId:    "post-123",
		CreatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "comments"`)).
		WithArgs(
			comment.Username,
			comment.Comment,
			comment.PostId,
			comment.CreatedBy,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			comment.ID,
		).
		WillReturnError(gorm.ErrInvalidData)
	suite.mock.ExpectRollback()

	_, err := suite.repository.InsertComment(ctx, comment)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidData, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestUpdateComment_Success() {
	ctx := context.Background()
	now := time.Now()
	commentID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	comment := model.CommentModel{
		ID:        commentID,
		Username:  "testuser",
		Comment:   "Updated comment",
		PostId:    "post-123",
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments"`)).
		WithArgs(
			comment.Username,
			comment.Comment,
			comment.PostId,
			comment.CreatedBy,
			comment.UpdatedBy,
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			comment.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	result, err := suite.repository.UpdateComment(ctx, comment)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), comment.Comment, result.Comment)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestUpdateComment_Error() {
	ctx := context.Background()
	now := time.Now()
	commentID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	comment := model.CommentModel{
		ID:        commentID,
		Username:  "testuser",
		Comment:   "Updated comment",
		PostId:    "post-123",
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments"`)).
		WithArgs(
			comment.Username,
			comment.Comment,
			comment.PostId,
			comment.CreatedBy,
			comment.UpdatedBy,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			comment.ID,
		).
		WillReturnError(gorm.ErrInvalidDB)
	suite.mock.ExpectRollback()

	_, err := suite.repository.UpdateComment(ctx, comment)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestGetCommentById_Success() {
	ctx := context.Background()
	commentID := "123e4567-e89b-12d3-a456-426614174000"
	postID := "post-123"
	now := time.Now()

	expectedComment := model.CommentModel{
		ID:        strfmt.UUID4(commentID),
		Username:  "testuser",
		Comment:   "Test comment",
		PostId:    postID,
		CreatedBy: "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Expect the main query for comment
	commentRows := sqlmock.NewRows([]string{"id", "username", "comment", "post_id", "created_by", "updated_by", "created_at", "updated_at"}).
		AddRow(expectedComment.ID, expectedComment.Username, expectedComment.Comment, expectedComment.PostId, expectedComment.CreatedBy, nil, expectedComment.CreatedAt, expectedComment.UpdatedAt)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE id = $1 AND "comments"."deleted_at" IS NULL ORDER BY "comments"."id" LIMIT`)).
		WithArgs(commentID, 1).
		WillReturnRows(commentRows)

	// Expect the preload query for Post
	postRows := sqlmock.NewRows([]string{"id", "username", "title", "body", "status", "created_by", "updated_by", "created_at", "updated_at"}).
		AddRow(postID, "postuser", "Post Title", "Post Body", "published", "postuser", nil, now, now)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE "posts"."id" = $1 AND "posts"."deleted_at" IS NULL`)).
		WithArgs(postID).
		WillReturnRows(postRows)

	result, err := suite.repository.GetCommentById(ctx, commentID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedComment.Comment, result.Comment)
	assert.Equal(suite.T(), expectedComment.Username, result.Username)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestGetCommentById_NotFound() {
	ctx := context.Background()
	commentID := "nonexistent-id"

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE id = $1 AND "comments"."deleted_at" IS NULL ORDER BY "comments"."id" LIMIT`)).
		WithArgs(commentID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := suite.repository.GetCommentById(ctx, commentID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
	assert.NotNil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestGetCommentById_Error() {
	ctx := context.Background()
	commentID := "123e4567-e89b-12d3-a456-426614174000"

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE id = $1 AND "comments"."deleted_at" IS NULL ORDER BY "comments"."id" LIMIT`)).
		WithArgs(commentID, 1).
		WillReturnError(gorm.ErrInvalidDB)

	result, err := suite.repository.GetCommentById(ctx, commentID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.NotNil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestDeleteComment_Success() {
	ctx := context.Background()
	commentID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	comment := model.CommentModel{
		ID: commentID,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET "deleted_at"=$1 WHERE "comments"."id" = $2 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), comment.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repository.DeleteComment(ctx, comment)

	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestDeleteComment_Error() {
	ctx := context.Background()
	commentID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	comment := model.CommentModel{
		ID: commentID,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET "deleted_at"=$1 WHERE "comments"."id" = $2 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), comment.ID).
		WillReturnError(gorm.ErrInvalidDB)
	suite.mock.ExpectRollback()

	err := suite.repository.DeleteComment(ctx, comment)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestDeleteComment_NotFound() {
	ctx := context.Background()
	commentID := strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000")

	comment := model.CommentModel{
		ID: commentID,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET "deleted_at"=$1 WHERE "comments"."id" = $2 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), comment.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := suite.repository.DeleteComment(ctx, comment)

	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestGetAllComment_Success() {
	ctx := context.Background()
	page := 1
	limit := 10
	now := time.Now()

	comment1 := model.CommentModel{
		ID:        strfmt.UUID4("123e4567-e89b-12d3-a456-426614174001"),
		Username:  "user1",
		Comment:   "Comment 1",
		PostId:    "post-1",
		CreatedBy: "user1",
		CreatedAt: now,
		UpdatedAt: now,
	}

	comment2 := model.CommentModel{
		ID:        strfmt.UUID4("123e4567-e89b-12d3-a456-426614174002"),
		Username:  "user2",
		Comment:   "Comment 2",
		PostId:    "post-2",
		CreatedBy: "user2",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Expect the main query for comments
	commentRows := sqlmock.NewRows([]string{"id", "username", "comment", "post_id", "created_by", "updated_by", "created_at", "updated_at"}).
		AddRow(comment1.ID, comment1.Username, comment1.Comment, comment1.PostId, comment1.CreatedBy, nil, comment1.CreatedAt, comment1.UpdatedAt).
		AddRow(comment2.ID, comment2.Username, comment2.Comment, comment2.PostId, comment2.CreatedBy, nil, comment2.CreatedAt, comment2.UpdatedAt)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE "comments"."deleted_at" IS NULL LIMIT`)).
		WithArgs(limit).
		WillReturnRows(commentRows)

	// Expect the preload query for Posts
	postRows := sqlmock.NewRows([]string{"id", "username", "title", "body", "status", "created_by", "updated_by", "created_at", "updated_at"}).
		AddRow("post-1", "postuser1", "Post 1", "Body 1", "published", "postuser1", nil, now, now).
		AddRow("post-2", "postuser2", "Post 2", "Body 2", "published", "postuser2", nil, now, now)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE "posts"."id" IN ($1,$2) AND "posts"."deleted_at" IS NULL`)).
		WithArgs("post-1", "post-2").
		WillReturnRows(postRows)

	result, err := suite.repository.GetAllComment(ctx, page, limit)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), comment1.Comment, result[0].Comment)
	assert.Equal(suite.T(), comment2.Comment, result[1].Comment)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestGetAllComment_EmptyResult() {
	ctx := context.Background()
	page := 1
	limit := 10

	commentRows := sqlmock.NewRows([]string{"id", "username", "comment", "post_id", "created_by", "updated_by", "created_at", "updated_at"})

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE "comments"."deleted_at" IS NULL LIMIT`)).
		WithArgs(limit).
		WillReturnRows(commentRows)

	result, err := suite.repository.GetAllComment(ctx, page, limit)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestGetAllComment_Error() {
	ctx := context.Background()
	page := 1
	limit := 10

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE "comments"."deleted_at" IS NULL LIMIT`)).
		WithArgs(limit).
		WillReturnError(gorm.ErrInvalidDB)

	result, err := suite.repository.GetAllComment(ctx, page, limit)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidDB, err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CommentRepositoryTestSuite) TestGetAllComment_Pagination() {
	ctx := context.Background()
	page := 2
	limit := 5
	offset := 5

	commentRows := sqlmock.NewRows([]string{"id", "username", "comment", "post_id", "created_by", "updated_by", "created_at", "updated_at"})

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE "comments"."deleted_at" IS NULL LIMIT $1 OFFSET $2`)).
		WithArgs(limit, offset).
		WillReturnRows(commentRows)

	result, err := suite.repository.GetAllComment(ctx, page, limit)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}
