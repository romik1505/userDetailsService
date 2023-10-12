package common

import (
	"github.com/romik1505/userDetailsService/internal/domain"
	"time"
)

type Person struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Surname     string `json:"surname,omitempty"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int32  `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`

	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (p Person) EditableFields() []string {
	fields := make([]string, 0)
	if p.Name != "" {
		fields = append(fields, domain.PersonNameField)
	}
	if p.Surname != "" {
		fields = append(fields, domain.PersonSurnameField, domain.PersonSurnameIndexCol)
	}
	if p.Patronymic != "" {
		fields = append(fields, domain.PersonPatronymicField)
	}
	if p.Age != 0 {
		fields = append(fields, domain.PersonAgeField)
	}
	if p.Gender != "" {
		fields = append(fields, domain.PersonGenderField)
	}
	if p.Nationality != "" {
		fields = append(fields, domain.PersonNationalityField)
	}

	return fields
}

func ConvertModelPerson(person domain.Person) Person {

	return Person{
		ID:          person.ID.Int64,
		Name:        person.Name.String,
		Surname:     person.Surname.String,
		Patronymic:  person.Patronymic.String,
		Age:         person.Age.Int32,
		Gender:      person.Gender.String,
		Nationality: person.Nationality.String,
		CreatedAt:   person.CreatedAt.Time,
		UpdatedAt:   domain.ConvertNullTime(person.UpdatedAt),
	}
}

func ConvertPerson(person Person) domain.Person {
	ret := domain.Person{
		ID:          domain.NewInt64(person.ID),
		Name:        domain.NewNullString(person.Name),
		Surname:     domain.NewNullString(person.Surname),
		Patronymic:  domain.NewNullString(person.Patronymic),
		Age:         domain.NewInt32(person.Age),
		Gender:      domain.NewNullString(person.Gender),
		Nationality: domain.NewNullString(person.Nationality),
	}
	return ret
}

func ConvertModelPersons(persons []domain.Person) []Person {
	res := make([]Person, len(persons))
	for i, person := range persons {
		res[i] = ConvertModelPerson(person)
	}
	return res
}
