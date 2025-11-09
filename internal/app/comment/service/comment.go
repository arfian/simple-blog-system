package service

import (
	"context"
	"errors"

	"simple-blog-system/internal/app/comment/model"
	"simple-blog-system/internal/app/comment/payload"
	"simple-blog-system/internal/app/comment/port"
	postPort "simple-blog-system/internal/app/post/port"
	userPort "simple-blog-system/internal/app/user/port"

	"github.com/go-openapi/strfmt"
)

type service struct {
	commentRepo port.ICommentRepository
	userRepo    userPort.IUserRepository
	postRepo    postPort.IPostRepository
}

func New(commentRepo port.ICommentRepository, userRepo userPort.IUserRepository, postRepo postPort.IPostRepository) port.ICommentService {
	return &service{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		postRepo:    postRepo,
	}
}

func (s *service) AddComment(ctx context.Context, username string, param payload.CommentRequest) (res *model.CommentModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	comment := model.CommentModel{
		Username:  users[0].Username,
		Comment:   param.Comment,
		PostId:    param.PostId,
		CreatedBy: username,
	}
	comment, qerr = s.commentRepo.InsertComment(ctx, comment)
	if qerr != nil {
		return nil, qerr
	}

	post, qerr := s.postRepo.GetPostById(ctx, comment.PostId)
	if qerr != nil {
		return nil, qerr
	}
	comment.Post = *post

	return &comment, nil
}

func (s *service) UpdateComment(ctx context.Context, username string, id string, param payload.CommentRequest) (res *model.CommentModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	comment := model.CommentModel{
		ID:        strfmt.UUID4(id),
		Username:  users[0].Username,
		Comment:   param.Comment,
		PostId:    param.PostId,
		CreatedBy: username,
	}
	comment, qerr = s.commentRepo.UpdateComment(ctx, comment)
	if qerr != nil {
		return nil, qerr
	}

	post, qerr := s.postRepo.GetPostById(ctx, comment.PostId)
	if qerr != nil {
		return nil, qerr
	}
	comment.Post = *post

	return &comment, nil
}

func (s *service) DeleteComment(ctx context.Context, username string, id string) (res *model.CommentModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	comment, err := s.commentRepo.GetCommentById(ctx, id)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	err = s.commentRepo.DeleteComment(ctx, *comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *service) GetAllComment(ctx context.Context, username string, page int, limit int) (res []model.CommentModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	post, err := s.commentRepo.GetAllComment(ctx, page, limit)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	return post, nil
}

func (s *service) GetCommentById(ctx context.Context, username string, id string) (res *model.CommentModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	comment, err := s.commentRepo.GetCommentById(ctx, id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return comment, nil
}
