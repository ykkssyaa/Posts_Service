package service

import "github.com/ykkssyaa/Posts_Service/internal/gateway"

type Services struct {
	Posts
	Comments
}

func NewServices(gateways *gateway.Gateways) *Services {
	return &Services{
		Posts:    NewPostsService(gateways.Posts),
		Comments: NewCommentsService(gateways.Comments),
	}
}

type Posts interface {
}

type Comments interface {
}
