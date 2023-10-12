package repository

import (
	"context"
	"github.com/romik1505/userDetailsService/internal/domain"
)

type Persons interface {
	Create(ctx context.Context, person *domain.Person) error
	Update(ctx context.Context, person *domain.Person, fields ...string) error
	List(ctx context.Context, filter domain.ListPersonsFilter) ([]domain.Person, int64, error)
	Delete(ctx context.Context, id int64) (int, error)
}
