package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/romik1505/userDetailsService/internal/cache"
	"github.com/romik1505/userDetailsService/internal/common"
	mock_statistics "github.com/romik1505/userDetailsService/pkg/mocks/client/statistics"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPersonService_enrichLikelyAge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := cache.NewCache()
	stats := mock_statistics.NewMockIClient(ctrl)
	p := PersonService{
		Cache:            c,
		StatisticsClient: stats,
	}
	t.Run("persons age found in cache", func(t *testing.T) {
		c.Del(context.Background(), "person:Oleg:age")

		input := &common.Person{Name: "Oleg"}
		c.Set(context.Background(), "person:Oleg:age", 11, time.Second*5)
		if err := p.enrichLikelyAge(context.Background(), input); err != nil {
			require.Nil(t, err)
		}
		require.Equal(t, int32(11), input.Age)

		c.Del(context.Background(), "person:Oleg:age")
	})
	t.Run("persons age not found in cache", func(t *testing.T) {
		c.Del(context.Background(), "person:Valera:age")

		input := &common.Person{Name: "Valera"}
		stats.EXPECT().GetMostLikelyAge(gomock.Any(), "Valera").Return(69, nil)
		if err := p.enrichLikelyAge(context.Background(), input); err != nil {
			require.Nil(t, err)
		}
		require.Equal(t, int32(69), input.Age)

		age, err := c.Get(context.Background(), "person:Valera:age").Int()
		require.Nil(t, err)
		require.Equal(t, 69, age)

		c.Del(context.Background(), "person:Valera:age")
	})
	t.Run("statistics client error", func(t *testing.T) {
		input := &common.Person{Name: "Afanasii"}
		stats.EXPECT().
			GetMostLikelyAge(gomock.Any(), "Afanasii").
			Return(0, fmt.Errorf("client error"))

		err := p.enrichLikelyAge(context.Background(), input)
		require.Error(t, fmt.Errorf("client error"), err)
	})
}

func TestPersonService_enrichLikelyGender(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := cache.NewCache()
	stats := mock_statistics.NewMockIClient(ctrl)
	p := PersonService{
		Cache:            c,
		StatisticsClient: stats,
	}
	t.Run("persons gender found in cache", func(t *testing.T) {
		c.Del(context.Background(), "person:Oleg:gender")

		input := &common.Person{Name: "Oleg"}
		c.Set(context.Background(), "person:Oleg:gender", "male", time.Second*5)
		if err := p.enrichLikelyGender(context.Background(), input); err != nil {
			require.Nil(t, err)
		}
		require.Equal(t, "male", input.Gender)

		c.Del(context.Background(), "person:Oleg:gender")
	})
	t.Run("persons gender not found in cache", func(t *testing.T) {
		c.Del(context.Background(), "person:Valera:gender")

		input := &common.Person{Name: "Valera"}
		stats.EXPECT().GetMostLikelyGender(gomock.Any(), "Valera").Return("male", nil)
		if err := p.enrichLikelyGender(context.Background(), input); err != nil {
			require.Nil(t, err)
		}
		require.Equal(t, "male", input.Gender)

		gender, err := c.Get(context.Background(), "person:Valera:gender").Result()
		require.Nil(t, err)
		require.Equal(t, "male", gender)

		c.Del(context.Background(), "person:Valera:gender")
	})
	t.Run("statistics client error", func(t *testing.T) {
		input := &common.Person{Name: "Afanasii"}
		stats.EXPECT().
			GetMostLikelyGender(gomock.Any(), "Afanasii").
			Return("", fmt.Errorf("client error"))

		err := p.enrichLikelyGender(context.Background(), input)
		require.Error(t, fmt.Errorf("client error"), err)
	})
}

func TestPersonService_enrichLikelyNationality(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := cache.NewCache()
	stats := mock_statistics.NewMockIClient(ctrl)
	p := PersonService{
		Cache:            c,
		StatisticsClient: stats,
	}
	t.Run("persons nationality found in cache", func(t *testing.T) {
		c.Del(context.Background(), "person:Oleg:nationality")

		input := &common.Person{Name: "Oleg"}
		c.Set(context.Background(), "person:Oleg:nationality", "UA", time.Second*5)
		if err := p.enrichLikelyNationality(context.Background(), input); err != nil {
			require.Nil(t, err)
		}
		require.Equal(t, input.Nationality, "UA")

		c.Del(context.Background(), "person:Oleg:nationality")
	})
	t.Run("persons nationality not found in cache", func(t *testing.T) {
		c.Del(context.Background(), "person:Valera:nationality")

		input := &common.Person{Name: "Valera"}
		stats.EXPECT().GetMostLikelyNationality(gomock.Any(), "Valera").Return("RU", nil)
		if err := p.enrichLikelyNationality(context.Background(), input); err != nil {
			require.Nil(t, err)
		}
		require.Equal(t, "RU", input.Nationality)

		nationality, err := c.Get(context.Background(), "person:Valera:nationality").Result()
		require.Nil(t, err)
		require.Equal(t, "RU", nationality)

		c.Del(context.Background(), "person:Valera:nationality")
	})
	t.Run("statistics client error", func(t *testing.T) {
		input := &common.Person{Name: "Afanasii"}
		stats.EXPECT().
			GetMostLikelyNationality(gomock.Any(), "Afanasii").
			Return("", fmt.Errorf("client error"))
		err := p.enrichLikelyNationality(context.Background(), input)
		require.Error(t, fmt.Errorf("client error"), err)
	})
}
