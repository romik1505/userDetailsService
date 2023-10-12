package resolver

import "github.com/romik1505/userDetailsService/internal/service"

// This file will not be regenerated automatically.
//go:generate go run github.com/99designs/gqlgen generate
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PersonService service.Persons
}

func NewResolver(ps service.Persons) *Resolver {
	return &Resolver{
		PersonService: ps,
	}
}
