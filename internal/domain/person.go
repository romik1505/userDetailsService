package domain

import (
	"database/sql"
)

const PersonTableName = "persons"

const (
	PersonNameField        = "name"
	PersonSurnameField     = "surname"
	PersonPatronymicField  = "patronymic"
	PersonAgeField         = "age"
	PersonGenderField      = "gender"
	PersonNationalityField = "nationality"
	PersonSurnameIndexCol  = "surname_index_col"
)

var PersonFields = []string{PersonNameField, PersonSurnameField, PersonPatronymicField,
	PersonAgeField, PersonGenderField, PersonNationalityField}

type Person struct {
	ID          sql.NullInt64  `db:"id"`
	Name        sql.NullString `db:"name"`
	Surname     sql.NullString `db:"surname"`
	Patronymic  sql.NullString `db:"patronymic"`
	Age         sql.NullInt32  `db:"age"`
	Gender      sql.NullString `db:"gender"`
	Nationality sql.NullString `db:"nationality"`

	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`

	TotalItems sql.NullInt64 `db:"total_items"`
}

type ListPersonsFilter struct {
	IDs           []int64  `json:"ids,omitempty" form:"-"`
	Surname       string   `json:"surname,omitempty" form:"surname"`
	Name          string   `json:"name_like,omitempty" form:"name"`
	AgeLtOrEq     int      `json:"age_lt_or_eq" form:"age_lt_or_eq"`
	AgeEq         int      `json:"age_eq" form:"age_eq"`
	AgeGtOrEq     int      `json:"age_gt_or_eq" form:"age_gt_or_eq"`
	GenderIn      []string `json:"gender_in" form:"-"`
	NationalityIn []string `json:"nationality_in" form:"-"`

	Page  int `json:"page,omitempty" form:"page"`
	Limit int `json:"limit,omitempty" form:"limit"`
}
