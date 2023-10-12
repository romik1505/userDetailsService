package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/romik1505/userDetailsService/internal/config"
	"github.com/romik1505/userDetailsService/internal/domain"
	"log"
	"strings"
	"time"
)

type PersonsRepo struct {
	*sqlx.DB
}

const (
	DefaultLimit = 100
)

func NewPersonsRepo() *PersonsRepo {
	connect, err := sqlx.Connect(config.Config.DBConfig.Driver, config.Config.DBConfig.ConnectionString())
	if err != nil {
		panic(err)
	}

	if config.Config.AppLevel != "test" {
		if err := goose.SetDialect(config.Config.DBConfig.Driver); err != nil {
			panic(err)
			log.Fatalf("goose set dialect error: %v", err)
		}

		if err := goose.Up(connect.DB, "migrations"); err != nil {
			log.Fatalf("goose up :%s", err.Error())
		}
	}

	return &PersonsRepo{
		DB: connect,
	}
}

func (p PersonsRepo) Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(p.DB)
}

func (p PersonsRepo) Create(ctx context.Context, person *domain.Person) error {
	q := p.Builder().Insert("persons").SetMap(map[string]interface{}{
		"name":              person.Name,
		"surname":           person.Surname,
		"patronymic":        person.Patronymic,
		"age":               person.Age,
		"gender":            person.Gender,
		"nationality":       person.Nationality,
		"surname_index_col": sq.Expr("to_tsvector(?)", person.Surname),
	}).Suffix("RETURNING id, created_at")

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}

	r := p.QueryRowxContext(ctx, query, args...)

	return r.StructScan(person)
}

func (p PersonsRepo) Delete(ctx context.Context, id int64) (int, error) {
	q := p.Builder().Delete("persons").Where(sq.Eq{"id": id})

	res, err := q.ExecContext(ctx)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(affected), nil
}

var col = map[string]func(*domain.Person) interface{}{
	domain.PersonNameField:        func(p *domain.Person) interface{} { return p.Name },
	domain.PersonSurnameField:     func(p *domain.Person) interface{} { return p.Surname },
	domain.PersonPatronymicField:  func(p *domain.Person) interface{} { return p.Patronymic },
	domain.PersonAgeField:         func(p *domain.Person) interface{} { return p.Age },
	domain.PersonGenderField:      func(p *domain.Person) interface{} { return p.Gender },
	domain.PersonNationalityField: func(p *domain.Person) interface{} { return p.Nationality },
	domain.PersonSurnameIndexCol:  func(p *domain.Person) interface{} { return sq.Expr("to_tsvector(?)", p.Surname) },
}

func (p PersonsRepo) Update(ctx context.Context, person *domain.Person, fields ...string) error {
	query := p.Builder().Update("persons").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": person.ID}).
		Suffix("RETURNING id, name, surname, patronymic, age, gender, nationality, created_at, updated_at")

	for _, fieldToUpdate := range fields {
		getter, ok := col[fieldToUpdate]
		if !ok {
			return fmt.Errorf("field not found")
		}

		query = query.Set(fieldToUpdate, getter(person))
	}

	q, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = p.QueryRowxContext(ctx, q, args...).StructScan(person)
	if err != nil {
		return err
	}

	return nil
}

func applyListPersonsFilter(s sq.SelectBuilder, f domain.ListPersonsFilter) sq.SelectBuilder {
	if len(f.IDs) > 0 {
		s = s.Where(sq.Eq{"id": f.IDs})
	}

	if f.Surname != "" {
		s = s.Where("surname_index_col @@ to_tsquery(?)", f.Surname)
	}

	if f.Name != "" {
		s = s.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", f.Name)})
	}

	if f.AgeLtOrEq != 0 {
		s = s.Where(sq.LtOrEq{"age": f.AgeLtOrEq})
	}

	if f.AgeEq != 0 {
		s = s.Where(sq.Eq{"age": f.AgeEq})
	}

	if f.AgeGtOrEq != 0 {
		s = s.Where(sq.GtOrEq{"age": f.AgeGtOrEq})
	}

	if len(f.GenderIn) > 0 {
		s = s.Where(sq.Eq{"gender": f.GenderIn})
	}

	if len(f.NationalityIn) > 0 {
		s = s.Where(sq.Eq{"nationality": f.NationalityIn})
	}

	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 || f.Limit > 10_000 {
		f.Limit = DefaultLimit
	}

	s = s.Limit(uint64(f.Limit)).Offset(uint64((f.Page - 1) * f.Limit))
	return s
}

func (p *PersonsRepo) List(ctx context.Context, filter domain.ListPersonsFilter) ([]domain.Person, int64, error) {
	q := p.Builder().
		Select("id", strings.Join(domain.PersonFields, " ,"), "updated_at", "created_at",
			"COUNT(*) OVER() as total_items").
		From("persons").OrderBy("id")

	q = applyListPersonsFilter(q, filter)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	persons := make([]domain.Person, 0, 10)

	for rows.Next() {
		var buff domain.Person
		err := rows.StructScan(&buff)
		if err != nil {
			return nil, 0, err
		}

		persons = append(persons, buff)
	}
	var totalItems int64

	if len(persons) != 0 {
		totalItems = persons[0].TotalItems.Int64
	}

	return persons, totalItems, nil
}
