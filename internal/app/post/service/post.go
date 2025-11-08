package service

import (
	"context"
	"errors"

	"simple-blog-system/internal/app/post/model"
	"simple-blog-system/internal/app/post/payload"
	"simple-blog-system/internal/app/post/port"
	userPort "simple-blog-system/internal/app/user/port"
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
