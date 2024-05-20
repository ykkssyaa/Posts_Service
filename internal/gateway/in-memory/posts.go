package gateway

import "github.com/jmoiron/sqlx"

type PostsInMemory struct {
	db *sqlx.DB
}
