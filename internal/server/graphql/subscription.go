package graphql

import (
	"errors"
	"github.com/ykkssyaa/Posts_Service/internal/consts"
	"github.com/ykkssyaa/Posts_Service/internal/models"
	"sync"
)

type Observers interface {
	CreateObserver(postId int) (int, chan *models.Comment, error)
	DeleteObserver(postId, chanId int) error
	NotifyObservers(postId int, comment models.Comment) error
}

type CommentsObservers struct {
	chans   map[int][]CommentObserver
	counter int
	mu      sync.Mutex
}

type CommentObserver struct {
	ch chan *models.Comment
	id int
}

func NewCommentsObserver() *CommentsObservers {
	return &CommentsObservers{
		chans:   make(map[int][]CommentObserver),
		counter: 0,
		mu:      sync.Mutex{},
	}
}

func (c *CommentsObservers) CreateObserver(postId int) (int, chan *models.Comment, error) {

	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan *models.Comment)
	c.counter++

	c.chans[postId] = append(c.chans[postId], CommentObserver{ch: ch, id: c.counter})

	return c.counter, ch, nil

}

func (c *CommentsObservers) DeleteObserver(postId, chanId int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	obs := c.chans[postId]
	for i, observer := range obs {
		if observer.id == chanId {
			c.chans[postId] = append(obs[:i], obs[i+1:]...)
		}
	}

	return nil
}

func (c *CommentsObservers) NotifyObservers(postId int, comment models.Comment) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	obs, ok := c.chans[postId]
	if ok {
		for _, observer := range obs {
			observer.ch <- &comment
		}
	} else {
		return errors.New(consts.ThereIsNoObserversError)
	}

	return nil
}
