package persistent

import (
	"context"

	"github.com/ne-ray/tcp-inbox/config"
	"github.com/ne-ray/tcp-inbox/internal/entity"
)

// WordOfWisdomRepo -.
type WordOfWisdomRepo struct {
	storageCfg config.Storage
}

// New -.
func New(storageCfg config.Storage) *WordOfWisdomRepo {
	return &WordOfWisdomRepo{storageCfg: storageCfg}
}

func (*WordOfWisdomRepo) Save(context.Context, entity.WordOfWisdom) (entity.WordOfWisdom, error) {
	//FIXME: реализовать сохранение на диск
	return entity.WordOfWisdom{}, nil
}
