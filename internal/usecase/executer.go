package usecase

import (
	"context"

	result "github.com/selfshop-dev/lib-result"
)

type Executer[C, R any] interface {
	Execute(ctx context.Context, cmd C) result.Value[R]
}

//go:generate mockgen -typed -destination=../mocks/mock_executer.go -package=mocks github.com/selfshop-dev/ms-catalog/internal/usecase Executer
