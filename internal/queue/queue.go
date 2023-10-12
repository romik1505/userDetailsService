package queue

import (
	"github.com/romik1505/userDetailsService/internal/queue/consumer"
	"github.com/romik1505/userDetailsService/internal/queue/producer"
)

type Queue struct {
	Producer *producer.Producer
	Consumer *consumer.Consumer
}

func NewMessageQueue(p *producer.Producer, c *consumer.Consumer) *Queue {
	return &Queue{
		Producer: p,
		Consumer: c,
	}
}
