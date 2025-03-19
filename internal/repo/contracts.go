package repo

import (
	"context"

	"github.com/ne-ray/tcp-inbox/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// WordOfWisdom -.
	WordOfWisdom interface {
		Save(context.Context, entity.WordOfWisdom) (entity.WordOfWisdom, error)
	}
)
