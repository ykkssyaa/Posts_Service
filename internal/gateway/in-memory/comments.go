package in_memory

import (
	"errors"
	"github.com/ykkssyaa/Posts_Service/internal/consts"
	"github.com/ykkssyaa/Posts_Service/internal/models"
	"sync"
	"time"
)

type CommentsInMemory struct {
	commCounter int
	comments    []models.Comment

	mu sync.RWMutex
}

func NewCommentsInMemory(count int) *CommentsInMemory {
	return &CommentsInMemory{
		commCounter: 0,
		comments:    make([]models.Comment, 0, count),
	}
}

func (c *CommentsInMemory) CreateComment(comment models.Comment) (models.Comment, error) {

	c.mu.Lock()
	defer c.mu.Unlock()

	c.commCounter++
	t := time.Now()

	comment.ID = c.commCounter
	comment.CreatedAt = t

	c.comments = append(c.comments, comment)

	return comment, nil

}

func (c *CommentsInMemory) GetCommentsByPost(postId, limit, offset int) ([]*models.Comment, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	var res []*models.Comment

	for _, comment := range c.comments {
		if comment.ReplyTo == nil && comment.Post == postId {
			com := comment
			res = append(res, &com)
		}
	}

	if offset > len(res) {
		return nil, nil
	}

	if offset+limit > len(res) || limit == -1 {
		return res[offset:], nil
	}

	if offset < 0 || limit < 0 {
		return nil, errors.New(consts.WrongLimitOffsetError)
	}

	return res[offset : offset+limit], nil
}

func (c *CommentsInMemory) GetRepliesOfComment(commentId int) ([]*models.Comment, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	if commentId > c.commCounter {
		return nil, nil
	}

	var res []*models.Comment

	for _, comment := range c.comments {
		if comment.ReplyTo != nil && *comment.ReplyTo == commentId {
			res = append(res, &comment)
		}
	}

	return res, nil
}
