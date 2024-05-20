package service

import "github.com/ykkssyaa/Posts_Service/internal/gateway"

type CommentsService struct {
	repo gateway.Comments
}

func NewCommentsService(repo gateway.Comments) *CommentsService {
	return &CommentsService{repo: repo}
}
