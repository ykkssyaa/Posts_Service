// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type Comment struct {
	ID        string     `json:"id"`
	CreatedAt string     `json:"createdAt"`
	Author    string     `json:"author"`
	Content   string     `json:"content"`
	Post      string     `json:"post"`
	Replies   []*Comment `json:"replies,omitempty"`
	ReplyTo   *string    `json:"replyTo,omitempty"`
}

type InputComment struct {
	Author  string  `json:"author"`
	Content string  `json:"content"`
	Post    string  `json:"post"`
	ReplyTo *string `json:"replyTo,omitempty"`
}

type InputPost struct {
	Name            string `json:"name"`
	Content         string `json:"content"`
	Author          string `json:"author"`
	CommentsAllowed bool   `json:"commentsAllowed"`
}

type Mutation struct {
}

type Post struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	Content   string `json:"content"`
}

type PostDetails struct {
	ID              string     `json:"id"`
	CreatedAt       string     `json:"createdAt"`
	Name            string     `json:"name"`
	Author          string     `json:"author"`
	Content         string     `json:"content"`
	CommentsAllowed bool       `json:"commentsAllowed"`
	Comments        []*Comment `json:"comments,omitempty"`
}

type Query struct {
}

type Subscription struct {
}