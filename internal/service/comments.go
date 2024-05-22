package service

import (
	"github.com/ykkssyaa/Posts_Service/internal/gateway"
	"github.com/ykkssyaa/Posts_Service/pkg/logger"
)

type CommentsService struct {
	repo   gateway.Comments
	logger *logger.Logger
}

func NewCommentsService(repo gateway.Comments, logger *logger.Logger) *CommentsService {
	return &CommentsService{repo: repo, logger: logger}
}
