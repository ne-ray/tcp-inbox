package wordofwisdom

import (
	"context"
	"fmt"

	"github.com/ne-ray/tcp-inbox/internal/entity"
	"github.com/ne-ray/tcp-inbox/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo repo.WordOfWisdom
}

// New -.
func New(r repo.WordOfWisdom) *UseCase {
	return &UseCase{
		repo: r,
	}
}

// Post - post part of word-of-wisdom.
func (uc *UseCase) Post(ctx context.Context, p entity.WordOfWisdom) (entity.WordOfWisdom, error) {
	w, err := uc.repo.Save(ctx, p)
	if err != nil {
		return entity.WordOfWisdom{}, fmt.Errorf("WordOfWisdomUseCase - Post - uc.repo.Save: %w", err)
	}

	return w, nil
}
