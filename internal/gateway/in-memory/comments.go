package gateway

import "github.com/jmoiron/sqlx"

type CommentsInMemory struct {
	db *sqlx.DB
}
