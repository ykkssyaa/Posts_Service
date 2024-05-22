package models

import "time"

type Comment struct {
	ID        int        `json:"id" db:"id"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	Author    string     `json:"author" db:"author"`
	Content   string     `json:"content" db:"content"`
	Post      int        `json:"post" db:"post"`
	Replies   []*Comment `json:"replies,omitempty"`
	ReplyTo   *int       `json:"replyTo,omitempty" db:"reply_to"`
}

func (c InputComment) FromInput() Comment {
	return Comment{
		Author:  c.Author,
		Content: c.Content,
		Post:    c.Post,
		ReplyTo: c.ReplyTo,
	}
}
