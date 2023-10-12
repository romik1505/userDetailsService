package consumer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/romik1505/userDetailsService/internal/config"
	"log"
	"strings"
)

type Consumer struct {
	ready chan bool
	Group sarama.ConsumerGroup
	MsgCh chan *sarama.ConsumerMessage
}

const (
	FIO_TOPIC = "FIO"
)

func NewConsumer() *Consumer {
	conf := sarama.NewConfig()
	conf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	consumer := &Consumer{
		ready: make(chan bool),
		MsgCh: make(chan *sarama.ConsumerMessage, 2),
	}
	group, err := sarama.NewConsumerGroup(
		strings.Split(config.Config.KafkaConfig.Brokers, ","),
		config.Config.KafkaConfig.GroupID,
		conf,
	)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	consumer.Group = group

	cgHandler := CGHandler{
		consumer: consumer,
	}

	go consumer.startConsume(context.Background(), &cgHandler)

	<-consumer.ready
	return consumer
}

func (c Consumer) Messages() <-chan *sarama.ConsumerMessage {
	return c.MsgCh
}

func (c Consumer) startConsume(ctx context.Context, handler *CGHandler) {
	for {
		if err := c.Group.Consume(ctx, []string{FIO_TOPIC}, handler); err != nil {
			log.Printf("Error consume message: %v", err)
		}
	}
}
