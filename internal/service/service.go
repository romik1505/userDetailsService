package service

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/romik1505/userDetailsService/internal/cache"
	"github.com/romik1505/userDetailsService/internal/client/statistics"
	"github.com/romik1505/userDetailsService/internal/common"
	"github.com/romik1505/userDetailsService/internal/domain"
	"github.com/romik1505/userDetailsService/internal/queue"
	"github.com/romik1505/userDetailsService/internal/repository"
	"log"
)

type PersonService struct {
	Repo             repository.Persons
	Cache            cache.Cache
	MQ               *queue.Queue
	StatisticsClient statistics.IClient
}

type Persons interface {
	CreatePerson(ctx context.Context, req common.CreatePersonRequest) error
	CreatePersonPushMessage(ctx context.Context, person common.CreatePersonRequest) error
	UpdatePerson(ctx context.Context, person *common.Person) error
	ListPersons(ctx context.Context, filter domain.ListPersonsFilter) (common.Pagination[common.Person], error)
	DeletePerson(ctx context.Context, id int64) error
	EnrichPerson(ctx context.Context, person *common.Person) error
}

func NewPersonService(r repository.Persons, c cache.Cache, sc statistics.IClient, mq *queue.Queue) *PersonService {
	ps := &PersonService{
		Repo:             r,
		Cache:            c,
		StatisticsClient: sc,
		MQ:               mq,
	}
	ps.StartConsume()

	return ps
}

func (p *PersonService) StartConsume() {
	messages := p.MQ.Consumer.Messages()

	for i := 0; i < 3; i++ {
		go Worker(i+1, messages, p)
	}

}

func Worker(id int, mchan <-chan *sarama.ConsumerMessage, p *PersonService) {
	for msg := range mchan {
		log.Printf("worker[%d] receive msg", id+1)
		person, err := UnmarshalPersonMessage(msg)
		if err != nil {
			log.Printf("Error unmarhal person message: %v", err)
		}
		if err := p.CreatePerson(context.Background(), person.Person); err != nil {
			log.Panicf("worker[%d]: error create person: %v", id+1, err)
		}
		log.Printf("worker[%d] finished perform", id+1)
	}
}

const (
	FIO_TOPIC        = "FIO"
	FIO_FAILED_TOPIC = "FIO_FAILED"
)
