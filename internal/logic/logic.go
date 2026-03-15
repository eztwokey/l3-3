package logic

import (
	"github.com/wb-go/wbf/logger"

	"github.com/eztwokey/l3-3/internal/storage"
)

type Logic struct {
	store  *storage.Storage
	logger logger.Logger
}

func New(store *storage.Storage, logger logger.Logger) *Logic {
	return &Logic{store: store, logger: logger}
}
