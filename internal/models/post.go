package models

import "time"

type Post struct {
	ID              int       `json:"id" db:"id"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	Name            string    `json:"name" db:"name"`
	Author          string    `json:"author" db:"author"`
	Content         string    `json:"content" db:"content"`
	CommentsAllowed bool      `json:"commentsAllowed" db:"comments_allowed"`
	//Comments        []*Comment `json:"comments,omitempty"`
}

func (p InputPost) FromInput() Post {
	return Post{
		Name:            p.Name,
		Author:          p.Author,
		Content:         p.Content,
		CommentsAllowed: p.CommentsAllowed,
	}
}

func (p Post) ToGraph() *PostGraph {
	return &PostGraph{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
		Author:    p.Author,
		Content:   p.Content,
	}
}

func ToPostGraph(posts []Post) []*PostGraph {
	res := make([]*PostGraph, 0, len(posts))

	for _, post := range posts {
		res = append(res, post.ToGraph())
	}

	return res
}
