package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/romik1505/userDetailsService/internal/common"
	"github.com/romik1505/userDetailsService/internal/domain"
	"log"
)

func (p *PersonService) CreatePersonPushMessage(ctx context.Context, person common.CreatePersonRequest) error {
	err := person.Validate()
	msg := PersonMessage{
		Person: person,
	}
	if err != nil {
		msg.ErrorReason = err.Error()
	}

	if err := p.sendPersonToMQ(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (p *PersonService) sendPersonToMQ(ctx context.Context, msg PersonMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	var topic string
	switch {
	case msg.ErrorReason != "":
		topic = FIO_FAILED_TOPIC
	default:
		topic = FIO_TOPIC
	}

	if _, _, err := p.MQ.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}); err != nil {
		return err
	}

	return nil
}

func (p *PersonService) CreatePerson(ctx context.Context, req common.CreatePersonRequest) error {
	person := req.ToPerson()
	if err := p.EnrichPerson(ctx, &person); err != nil {
		return err
	}

	model := common.ConvertPerson(person)
	err := p.Repo.Create(ctx, &model)
	if err != nil {
		log.Printf("db: error create person: %v", err)
		return err
	}

	return nil
}

func (p *PersonService) UpdatePerson(ctx context.Context, person *common.Person) error {
	if person.ID == 0 {
		return ErrBadRequest
	}

	model := common.ConvertPerson(*person)
	fields := person.EditableFields()

	err := p.Repo.Update(ctx, &model, fields...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}

		return ErrInternalError
	}
	out := common.ConvertModelPerson(model)
	*person = out
	return nil
}

func (p *PersonService) ListPersons(ctx context.Context, filter domain.ListPersonsFilter) (common.Pagination[common.Person], error) {
	persons, totalItems, err := p.Repo.List(ctx, filter)
	if err != nil {
		return common.Pagination[common.Person]{}, ErrInternalError
	}
	if len(persons) == 0 {
		return common.Pagination[common.Person]{}, ErrEmptyResponse
	}
	return common.NewPagination(common.ConvertModelPersons(persons), totalItems), err
}

func (p *PersonService) DeletePerson(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrBadRequest
	}

	affected, err := p.Repo.Delete(ctx, id)
	if affected == 0 {
		return ErrNotFound
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}
