package main

import (
	_ "github.com/romik1505/userDetailsService/docs"
	"github.com/romik1505/userDetailsService/internal/cache"
	"github.com/romik1505/userDetailsService/internal/client/statistics"
	_ "github.com/romik1505/userDetailsService/internal/config"
	"github.com/romik1505/userDetailsService/internal/queue"
	"github.com/romik1505/userDetailsService/internal/queue/consumer"
	"github.com/romik1505/userDetailsService/internal/queue/producer"
	"github.com/romik1505/userDetailsService/internal/repository"
	"github.com/romik1505/userDetailsService/internal/server"
	"github.com/romik1505/userDetailsService/internal/service"
	_ "github.com/romik1505/userDetailsService/migrations"
)

// @title Person API
// @version         1.0
// @host localhost:8080
// @BasePath /api/v1
func main() {
	repo := repository.NewPersonsRepo()
	ch := cache.NewCache()
	stats := statistics.NewClient()

	prod := producer.NewProducer()
	cons := consumer.NewConsumer()
	qu := queue.NewMessageQueue(prod, cons)

	ps := service.NewPersonService(repo, ch, stats, qu)

	srv := server.NewServer(ps)
	srv.Run()
}
