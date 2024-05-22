package service

import (
	"github.com/ykkssyaa/Posts_Service/internal/gateway"
	"github.com/ykkssyaa/Posts_Service/internal/models"
	"github.com/ykkssyaa/Posts_Service/pkg/logger"
)

type Services struct {
	Posts
	Comments
}

func NewServices(gateways *gateway.Gateways, logger *logger.Logger) *Services {
	return &Services{
		Posts:    NewPostsService(gateways.Posts, logger),
		Comments: NewCommentsService(gateways.Comments, logger),
	}
}

type Posts interface {
	CreatePost(post models.Post) (models.Post, error)
	GetPostById(id int) (models.Post, error)
	GetAllPosts(page, pageSize *int) ([]models.Post, error)
}

type Comments interface {
}
