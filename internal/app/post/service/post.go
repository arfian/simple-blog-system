package service

import (
	"context"
	"errors"

	"simple-blog-system/internal/app/post/model"
	"simple-blog-system/internal/app/post/payload"
	"simple-blog-system/internal/app/post/port"
	userPort "simple-blog-system/internal/app/user/port"

	"github.com/go-openapi/strfmt"
)

type service struct {
	postRepo port.IPostRepository
	userRepo userPort.IUserRepository
}

func New(postRepo port.IPostRepository, userRepo userPort.IUserRepository) port.IPostService {
	return &service{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (s *service) AddPost(ctx context.Context, username string, param payload.PostRequest) (res *model.PostModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	post := model.PostModel{
		Username:  users[0].Username,
		Title:     param.Title,
		Body:      param.Body,
		Status:    param.Status,
		CreatedBy: username,
	}
	post, qerr = s.postRepo.InsertPost(ctx, post)
	if qerr != nil {
		return nil, qerr
	}

	return &post, nil
}

func (s *service) UpdatePost(ctx context.Context, username string, id string, param payload.PostRequest) (res *model.PostModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	post := model.PostModel{
		ID:        strfmt.UUID4(id),
		Username:  users[0].Username,
		Title:     param.Title,
		Body:      param.Body,
		Status:    param.Status,
		CreatedBy: username,
	}
	post, qerr = s.postRepo.UpdatePost(ctx, post)
	if qerr != nil {
		return nil, qerr
	}

	return &post, nil
}

func (s *service) DeletePost(ctx context.Context, username string, id string) (res *model.PostModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	post, err := s.postRepo.GetPostById(ctx, id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	err = s.postRepo.DeletePost(ctx, *post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) GetAllPost(ctx context.Context, username string, page int, limit int) (res []model.PostModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	post, err := s.postRepo.GetAllPost(ctx, page, limit)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (s *service) GetById(ctx context.Context, username string, id string) (res *model.PostModel, err error) {
	users, qerr := s.userRepo.GetUserByUsername(ctx, username)
	if len(users) == 0 || qerr != nil {
		return nil, errors.New("user not found")
	}

	post, err := s.postRepo.GetPostById(ctx, id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return post, nil
}
