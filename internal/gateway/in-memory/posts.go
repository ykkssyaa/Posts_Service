package in_memory

import (
	"database/sql"
	"errors"
	"github.com/ykkssyaa/Posts_Service/internal/consts"
	"github.com/ykkssyaa/Posts_Service/internal/models"
	"sync"
	"time"
)

type PostsInMemory struct {
	postCounter int
	posts       []models.Post

	mu sync.RWMutex
}

func NewPostsInMemory(count int) *PostsInMemory {
	return &PostsInMemory{
		postCounter: 0,
		posts:       make([]models.Post, 0, count),
	}
}

func (p *PostsInMemory) CreatePost(post models.Post) (models.Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.postCounter++
	t := time.Now()

	post.ID = p.postCounter
	post.CreatedAt = t

	p.posts = append(p.posts, post)

	return post, nil
}

func (p *PostsInMemory) GetPostById(id int) (models.Post, error) {

	p.mu.RLock()
	defer p.mu.RUnlock()

	if id > p.postCounter || id <= 0 {
		return models.Post{}, sql.ErrNoRows
	}

	return p.posts[id-1], nil
}

func (p *PostsInMemory) GetAllPosts(limit, offset int) ([]models.Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if offset > p.postCounter {
		return nil, nil
	}

	if offset+limit > p.postCounter || limit == -1 {
		return p.posts[offset:], nil
	}

	if offset < 0 || limit < 0 {
		return nil, errors.New(consts.WrongLimitOffsetError)
	}

	return p.posts[offset : offset+limit], nil
}
