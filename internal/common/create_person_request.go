package common

import "fmt"

type CreatePersonRequest struct {
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
}

func (c CreatePersonRequest) Bind() (Person, error) {
	return c.ToPerson(), c.Validate()
}

func (c CreatePersonRequest) ToPerson() Person {
	return Person{
		Name:       c.Name,
		Surname:    c.Surname,
		Patronymic: c.Patronymic,
	}
}

func (c CreatePersonRequest) Validate() error {
	if c.Surname == "" {
		return fmt.Errorf("person surname not set")
	}
	if c.Name == "" {
		return fmt.Errorf("person name not set")
	}

	return nil
}
