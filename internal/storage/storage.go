package storage

import (
	"errors"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Storage {
	return &Storage{db: db}
}
