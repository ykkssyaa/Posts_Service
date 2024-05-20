package graphql

import "github.com/ykkssyaa/Posts_Service/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service.Posts
	service.Comments
}
