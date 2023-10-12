package repository

import (
	"context"
	"database/sql"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/romik1505/userDetailsService/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPersonsRepo_Create(t *testing.T) {
	t.Run("create new person", func(t *testing.T) {
		newPerson := &domain.Person{
			Name:        domain.NewNullString("Irina"),
			Surname:     domain.NewNullString("Ivanova"),
			Patronymic:  domain.NewNullString("Sergeevna"),
			Age:         domain.NewInt32(44),
			Gender:      domain.NewNullString("female"),
			Nationality: domain.NewNullString("RU"),
		}
		personsRepo.DB.MustExec("DELETE FROM persons;")

		err := personsRepo.Create(context.Background(), newPerson)
		require.Nil(t, err)

		outputPerson := &domain.Person{}
		err = personsRepo.QueryRowx("Select id, name, surname, patronymic, age, gender, nationality,"+
			"created_at, updated_at from persons where id=$1", newPerson.ID).StructScan(outputPerson)
		require.Nil(t, err)
		require.Equal(t, newPerson, outputPerson)
		require.NotEmpty(t, newPerson.CreatedAt.Time)
	})
}

func TestPersonsRepo_Delete(t *testing.T) {

	t.Run("delete exist person", func(t *testing.T) {
		var id int64

		err := personsRepo.DB.QueryRow("INSERT INTO persons(name, surname, patronymic, age, gender, nationality, surname_index_col) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
			"Egor", "Medvedev", "Olegovich", 25, "male", "RU", "'medvedev':1").Scan(&id)

		require.Nil(t, err)

		output := &domain.Person{}
		err = personsRepo.DB.QueryRowx("SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at FROM persons WHERE id=$1", id).StructScan(output)
		require.Nil(t, err)
		require.Equal(t, id, output.ID.Int64)
		require.NotEmpty(t, output.CreatedAt)

		affected, err := personsRepo.Delete(context.Background(), id)
		require.Equal(t, 1, affected)
		require.Nil(t, err)

		err = personsRepo.DB.QueryRowx("SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at FROM persons WHERE id=$1", id).StructScan(output)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestPersonsRepo_List(t *testing.T) {
	personsRepo.DB.MustExec("DELETE FROM persons")

	inputPersons := []domain.Person{
		{
			Name:        domain.NewNullString("Ivan"),
			Surname:     domain.NewNullString("Ivanov"),
			Patronymic:  domain.NewNullString("Ivanovich"),
			Age:         domain.NewInt32(33),
			Gender:      domain.NewNullString("male"),
			Nationality: domain.NewNullString("RU"),
		},
		{
			Name:        domain.NewNullString("Rostislav"),
			Surname:     domain.NewNullString("Moromishko"),
			Patronymic:  domain.NewNullString("Valerievich"),
			Age:         domain.NewInt32(44),
			Gender:      domain.NewNullString("male"),
			Nationality: domain.NewNullString("UA"),
		},
		{
			Name:        domain.NewNullString("Nastya"),
			Surname:     domain.NewNullString("Datsenko"),
			Patronymic:  domain.NewNullString("Olegovna"),
			Age:         domain.NewInt32(25),
			Gender:      domain.NewNullString("female"),
			Nationality: domain.NewNullString("BY"),
		},
		{
			Name:        domain.NewNullString("Nastya"),
			Surname:     domain.NewNullString("Ivanova"),
			Patronymic:  domain.NewNullString("Markovich"),
			Age:         domain.NewInt32(27),
			Gender:      domain.NewNullString("female"),
			Nationality: domain.NewNullString("RU"),
		},
		{
			Name:        domain.NewNullString("Vyacheslav"),
			Surname:     domain.NewNullString("Surnamovich"),
			Patronymic:  domain.NewNullString("Sergeevich"),
			Age:         domain.NewInt32(44),
			Gender:      domain.NewNullString("male"),
			Nationality: domain.NewNullString("RU"),
		},
	}

	for i := range inputPersons {
		personsRepo.Create(context.Background(), &inputPersons[i])
	}

	tests := []struct {
		name      string
		filter    domain.ListPersonsFilter
		want      []domain.Person
		wantTotal int64
		wantErr   bool
		err       error
	}{
		{
			name:      "empty filter",
			want:      inputPersons,
			wantTotal: int64(len(inputPersons)),
		},
		{
			name: "id filet",
			filter: domain.ListPersonsFilter{
				IDs: []int64{inputPersons[0].ID.Int64, inputPersons[1].ID.Int64},
			},
			want:      inputPersons[:2],
			wantTotal: 2,
		},
		{
			name: "surname filter",
			filter: domain.ListPersonsFilter{
				Surname: "Moromishko",
			},
			want:      inputPersons[1:2],
			wantTotal: 1,
		},
		{
			name: "name filter",
			filter: domain.ListPersonsFilter{
				Name: "Nastya",
			},
			want:      inputPersons[2:4],
			wantTotal: 2,
		},
		{
			name: "gender filter",
			filter: domain.ListPersonsFilter{
				GenderIn: []string{"male"},
			},
			want:      []domain.Person{inputPersons[0], inputPersons[1], inputPersons[4]},
			wantTotal: 3,
		},
		{
			name: "nationality filter",
			filter: domain.ListPersonsFilter{
				NationalityIn: []string{"RU", "BY"},
			},
			want:      []domain.Person{inputPersons[0], inputPersons[2], inputPersons[3], inputPersons[4]},
			wantTotal: 4,
		},
		{
			name: "pagination filter",
			filter: domain.ListPersonsFilter{
				Page:  1,
				Limit: 2,
			},
			want:      inputPersons[:2],
			wantTotal: 5,
		},
		{
			name: "age interval filter",
			filter: domain.ListPersonsFilter{
				AgeGtOrEq: 18,
				AgeLtOrEq: 30,
			},
			want:      inputPersons[2:4],
			wantTotal: 2,
		},
		{
			name: "age equals filter",
			filter: domain.ListPersonsFilter{
				AgeEq: 25,
			},
			want:      inputPersons[2:3],
			wantTotal: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			persons, totalItems, err := personsRepo.List(context.Background(), tt.filter)
			if tt.wantErr {
				require.Equal(t, tt.wantErr, err)
			}
			require.Equal(t, tt.wantTotal, totalItems)

			require.Empty(t, cmp.Diff(tt.want, persons, cmpopts.IgnoreFields(domain.Person{}, "CreatedAt", "UpdatedAt", "TotalItems")))
		})
	}
}

func TestPersonsRepo_Update(t *testing.T) {

	tests := []struct {
		name         string
		input        *domain.Person
		updatePerson *domain.Person
		updateFields []string
		exceptPerson *domain.Person
		hookBefore   func(t *testing.T) int64
		wantErr      bool
		err          error
	}{
		{
			name: "updated exist person",
			input: &domain.Person{
				Name:        domain.NewNullString("Kirill"),
				Surname:     domain.NewNullString("Vitalievich"),
				Patronymic:  domain.NewNullString("Patronymicovich"),
				Age:         domain.NewInt32(44),
				Gender:      domain.NewNullString("male"),
				Nationality: domain.NewNullString("BY"),
			},
			updatePerson: &domain.Person{
				Surname:     domain.NewNullString("NewSurnamovich"),
				Age:         domain.NewInt32(55),
				Gender:      domain.NewNullString("female"),
				Nationality: domain.NewNullString("RU"),
			},
			updateFields: []string{domain.PersonSurnameField, domain.PersonGenderField, domain.PersonAgeField, domain.PersonNationalityField},
			exceptPerson: &domain.Person{
				Name:        domain.NewNullString("Kirill"),
				Surname:     domain.NewNullString("NewSurnamovich"),
				Patronymic:  domain.NewNullString("Patronymicovich"),
				Age:         domain.NewInt32(55),
				Gender:      domain.NewNullString("female"),
				Nationality: domain.NewNullString("RU"),
			},
			hookBefore: func(t *testing.T) int64 {
				var id int64
				err := personsRepo.DB.QueryRow("INSERT INTO persons(name, surname, patronymic, age, gender, nationality, surname_index_col) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
					"Kirill", "Vitalievich", "Patronymicovich", 44, "male", "BY", "'vitalievich':1").Scan(&id)
				require.Nil(t, err)
				return id
			},
		},
		{
			name: "update not exist person",
			input: &domain.Person{
				Name:        domain.NewNullString("Kirill"),
				Surname:     domain.NewNullString("Vitalievich"),
				Patronymic:  domain.NewNullString("Patronymicovich"),
				Age:         domain.NewInt32(44),
				Gender:      domain.NewNullString("male"),
				Nationality: domain.NewNullString("BY"),
			},
			updatePerson: &domain.Person{
				Surname:     domain.NewNullString("NewSurnamovich"),
				Age:         domain.NewInt32(55),
				Gender:      domain.NewNullString("female"),
				Nationality: domain.NewNullString("RU"),
			},
			updateFields: []string{domain.PersonSurnameField, domain.PersonGenderField, domain.PersonAgeField, domain.PersonNationalityField},
			wantErr:      true,
			err:          sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id int64
			if tt.hookBefore != nil {
				id = tt.hookBefore(t)
			}

			tt.updatePerson.ID = domain.NewInt64(id)

			if tt.exceptPerson != nil {
				tt.exceptPerson.ID = domain.NewInt64(id)
			}

			if err := personsRepo.Update(context.Background(), tt.updatePerson, tt.updateFields...); err != nil {
				require.ErrorIs(t, err, tt.err)
			}

			if !tt.wantErr {
				out := &domain.Person{}
				err := personsRepo.DB.QueryRowx("SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at FROM persons WHERE id=$1", id).StructScan(out)
				require.Nil(t, err)

				require.Empty(t, cmp.Diff(tt.exceptPerson, out, cmpopts.IgnoreFields(domain.Person{}, "CreatedAt", "UpdatedAt")))
			}

			personsRepo.DB.MustExec("DELETE FROM persons;")
		})
	}
}
