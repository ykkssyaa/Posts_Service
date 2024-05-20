package service

import "github.com/ykkssyaa/Posts_Service/internal/gateway"

type PostsService struct {
	repo gateway.Posts
}

func NewPostsService(repo gateway.Posts) *PostsService {
	return &PostsService{repo: repo}
}
