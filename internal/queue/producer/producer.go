package producer

import (
	"github.com/IBM/sarama"
	"github.com/romik1505/userDetailsService/internal/config"
	"log"
	"strings"
)

type Producer struct {
	sarama.SyncProducer
}

func NewProducer() *Producer {
	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = 5
	conf.Producer.Return.Errors = true
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(config.Config.KafkaConfig.Brokers, ","), conf)
	if err != nil {
		log.Panicf("Error creating producer: %v", err)
	}

	return &Producer{producer}
}
