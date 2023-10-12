package common

type Pagination[T any] struct {
	Items      []T   `json:"items,omitempty"`
	TotalItems int64 `json:"total_items,omitempty"`
}

type PersonPagination Pagination[Person]

func NewPagination[T any](items []T, total int64) Pagination[T] {
	return Pagination[T]{
		Items:      items,
		TotalItems: total,
	}
}
