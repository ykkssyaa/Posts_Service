package gateway

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
}

type Comments interface {
}
