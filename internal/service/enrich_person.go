package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/romik1505/userDetailsService/internal/common"
	"log"
	"sync"
	"time"
)

type EnrichesFunc func(context.Context, *common.Person) error

const (
	PersonAgeKeyMask         = "person:%s:age"
	PersonGenderKeyMask      = "person:%s:gender"
	PersonNationalityKeyMask = "person:%s:nationality"
)

func (p *PersonService) enrichLikelyAge(ctx context.Context, person *common.Person) error {
	age, err := p.Cache.Get(ctx, fmt.Sprintf(PersonAgeKeyMask, person.Name)).Int()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return err
		}

		log.Printf("error cache get person age: %v", err)

		clientAge, clientErr := p.StatisticsClient.GetMostLikelyAge(ctx, person.Name)
		if clientErr != nil {
			return clientErr
		}
		age = clientAge

		setCmd := p.Cache.Set(ctx, fmt.Sprintf(PersonAgeKeyMask, person.Name), clientAge, time.Hour*24)
		if setCmd.Err() != nil {
			return setCmd.Err()
		}
	}

	person.Age = int32(age)
	return nil
}

func (p *PersonService) enrichLikelyGender(ctx context.Context, person *common.Person) error {
	gender, err := p.Cache.Get(ctx, fmt.Sprintf(PersonGenderKeyMask, person.Name)).Result()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return err
		}

		log.Printf("error cache get person gender: %v", err)

		clientGender, clientErr := p.StatisticsClient.GetMostLikelyGender(ctx, person.Name)
		if clientErr != nil {
			return clientErr
		}

		gender = clientGender

		setCmd := p.Cache.Set(ctx, fmt.Sprintf(PersonGenderKeyMask, person.Name), clientGender, time.Hour*24)
		if setCmd.Err() != nil {
			return setCmd.Err()
		}
	}

	person.Gender = gender
	return nil
}

func (p *PersonService) enrichLikelyNationality(ctx context.Context, person *common.Person) error {
	nationality, err := p.Cache.Get(ctx, fmt.Sprintf(PersonNationalityKeyMask, person.Name)).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return err
		}

		log.Printf("error cache get person nationality: %v", err)

		clientNationality, clientErr := p.StatisticsClient.GetMostLikelyNationality(ctx, person.Name)
		if clientErr != nil {
			return clientErr
		}
		nationality = clientNationality

		setCmd := p.Cache.Set(ctx, fmt.Sprintf(PersonNationalityKeyMask, person.Name), clientNationality, time.Hour*24)
		if setCmd.Err() != nil {
			return setCmd.Err()
		}

	}

	person.Nationality = nationality
	return nil
}

func (p *PersonService) enrich(ctx context.Context, person *common.Person, fns ...EnrichesFunc) error {
	wg := sync.WaitGroup{}

	wg.Add(len(fns))
	var someError error
	for _, fn := range fns {
		go func(fn EnrichesFunc) {
			defer wg.Done()
			if err := fn(ctx, person); err != nil {
				someError = err
			}
		}(fn)
	}
	wg.Wait()

	return someError
}

func (p *PersonService) EnrichPerson(ctx context.Context, person *common.Person) error {
	enrichFuncs := []EnrichesFunc{
		p.enrichLikelyAge,
		p.enrichLikelyGender,
		p.enrichLikelyNationality,
	}

	if err := p.enrich(ctx, person, enrichFuncs...); err != nil {
		return err
	}

	return nil
}
