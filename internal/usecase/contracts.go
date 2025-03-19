// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/ne-ray/tcp-inbox/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// WordOfWisdom -.
	WordOfWisdom interface {
		Post(context.Context, entity.WordOfWisdom) (entity.WordOfWisdom, error)
	}
)
