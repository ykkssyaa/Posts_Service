package gateway

import "github.com/ykkssyaa/Posts_Service/internal/models"

type Gateways struct {
	Posts
	Comments
}

func NewGateways(posts Posts, comments Comments) *Gateways {
	return &Gateways{
		Posts:    posts,
		Comments: comments,
	}
}

type Posts interface {
	CreatePost(post models.Post) (models.Post, error)
	GetPostById(id int) (models.Post, error)
	GetAllPosts(limit, offset int) ([]models.Post, error)
}

type Comments interface {
}
