package service

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/romik1505/userDetailsService/internal/cache"
	"github.com/romik1505/userDetailsService/internal/common"
	"github.com/romik1505/userDetailsService/internal/domain"
	mock_statistics "github.com/romik1505/userDetailsService/pkg/mocks/client/statistics"
	mock_repository "github.com/romik1505/userDetailsService/pkg/mocks/repository/person"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPersonService_CreatePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := cache.NewCache()
	stats := mock_statistics.NewMockIClient(ctrl)
	repo := mock_repository.NewMockPersons(ctrl)

	personService := PersonService{repo, c, nil, stats}

	t.Run("create person", func(t *testing.T) {
		input := common.CreatePersonRequest{
			Name:       "Ivan",
			Surname:    "Markov",
			Patronymic: "Dmitrievich",
		}

		stats.EXPECT().GetMostLikelyAge(gomock.Any(), "Ivan").Return(33, nil).AnyTimes()

		stats.EXPECT().GetMostLikelyGender(gomock.Any(), "Ivan").Return("male", nil).AnyTimes()

		stats.EXPECT().GetMostLikelyNationality(gomock.Any(), "Ivan").Return("RU", nil).AnyTimes()

		repo.EXPECT().Create(gomock.Any(), &domain.Person{
			ID:          domain.NewInt64(0),
			Name:        domain.NewNullString("Ivan"),
			Surname:     domain.NewNullString("Markov"),
			Patronymic:  domain.NewNullString("Dmitrievich"),
			Age:         domain.NewInt32(33),
			Gender:      domain.NewNullString("male"),
			Nationality: domain.NewNullString("RU"),
		}).Return(nil)

		err := personService.CreatePerson(context.Background(), input)
		require.Nil(t, err)
	})
}

func TestPersonService_UpdatePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockPersons(ctrl)
	p := &PersonService{
		Repo: repo,
	}

	tests := []struct {
		name    string
		update  *common.Person
		wantErr bool
		err     error
		hook    func()
	}{
		{
			name: "update all fields",
			update: &common.Person{
				ID:          44,
				Name:        "Makar",
				Surname:     "Markov",
				Patronymic:  "Genadievich",
				Age:         55,
				Gender:      "male",
				Nationality: "UA",
			},
			wantErr: false,
			hook: func() {
				repo.EXPECT().Update(gomock.Any(), &domain.Person{
					ID:          domain.NewInt64(44),
					Name:        domain.NewNullString("Makar"),
					Surname:     domain.NewNullString("Markov"),
					Patronymic:  domain.NewNullString("Genadievich"),
					Age:         domain.NewInt32(55),
					Gender:      domain.NewNullString("male"),
					Nationality: domain.NewNullString("UA"),
				}, "name", "surname", "surname_index_col", "patronymic", "age", "gender", "nationality").
					Do(func(ctx context.Context, person *domain.Person, fields ...string) {
						person.CreatedAt = domain.NewNullTime(time.Now().Add(-time.Minute))
						person.UpdatedAt = domain.NewNullTime(time.Now())
					}).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hook != nil {
				tt.hook()
			}

			err := p.UpdatePerson(context.Background(), tt.update)

			require.ErrorIs(t, err, tt.err)
			require.NotEmpty(t, tt.update.CreatedAt)
			require.NotEmpty(t, tt.update.UpdatedAt)
		})
	}
}

func TestPersonService_ListPersons(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockPersons(ctrl)

	ps := PersonService{
		Repo: repo,
	}

	t.Run("ok case", func(t *testing.T) {
		filter := domain.ListPersonsFilter{
			IDs:           []int64{11, 22, 33},
			Surname:       "Ivanov",
			GenderIn:      []string{"male"},
			NationalityIn: []string{"RU"},
			Page:          1,
			Limit:         100,
		}
		repo.EXPECT().List(gomock.Any(), filter).Return([]domain.Person{
			{
				ID:          domain.NewInt64(11),
				Name:        domain.NewNullString("Vasya"),
				Surname:     domain.NewNullString("Ivanov"),
				Patronymic:  domain.NewNullString("Ilyich"),
				Age:         domain.NewInt32(42),
				Gender:      domain.NewNullString("male"),
				Nationality: domain.NewNullString("RU"),
				TotalItems:  domain.NewInt64(3),
			}, {
				ID:          domain.NewInt64(22),
				Name:        domain.NewNullString("Misha"),
				Surname:     domain.NewNullString("Ivanov"),
				Patronymic:  domain.NewNullString("Yrievich"),
				Age:         domain.NewInt32(23),
				Gender:      domain.NewNullString("male"),
				Nationality: domain.NewNullString("RU"),
				TotalItems:  domain.NewInt64(3),
			}, {
				ID:          domain.NewInt64(33),
				Name:        domain.NewNullString("Anatolii"),
				Surname:     domain.NewNullString("Ivanov"),
				Patronymic:  domain.NewNullString("Mihailovich"),
				Age:         domain.NewInt32(55),
				Gender:      domain.NewNullString("male"),
				Nationality: domain.NewNullString("RU"),
				TotalItems:  domain.NewInt64(3),
			},
		}, int64(3), nil)
		pagination, err := ps.ListPersons(context.Background(), filter)
		require.Nil(t, err)
		out := common.PersonPagination{
			Items: []common.Person{
				{
					ID:          11,
					Name:        "Vasya",
					Surname:     "Ivanov",
					Patronymic:  "Ilyich",
					Age:         42,
					Gender:      "male",
					Nationality: "RU",
				},
				{
					ID:          22,
					Name:        "Misha",
					Surname:     "Ivanov",
					Patronymic:  "Yrievich",
					Age:         23,
					Gender:      "male",
					Nationality: "RU",
				},
				{
					ID:          33,
					Name:        "Anatolii",
					Surname:     "Ivanov",
					Patronymic:  "Mihailovich",
					Age:         55,
					Gender:      "male",
					Nationality: "RU",
				},
			},
			TotalItems: 3,
		}
		require.Equal(t, common.PersonPagination(pagination), out)
	})
}

func TestPersonService_DeletePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockPersons(ctrl)
	p := PersonService{Repo: repo}
	t.Run("delete with empty id", func(t *testing.T) {
		err := p.DeletePerson(context.Background(), 0)
		require.ErrorIs(t, ErrBadRequest, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo.EXPECT().Delete(gomock.Any(), int64(1)).Return(0, sql.ErrNoRows)
		err := p.DeletePerson(context.Background(), 1)
		require.ErrorIs(t, ErrNotFound, err)
	})

	t.Run("success delete", func(t *testing.T) {
		repo.EXPECT().Delete(gomock.Any(), int64(1)).Return(1, nil)
		err := p.DeletePerson(context.Background(), 1)
		require.Nil(t, err)
	})
}
