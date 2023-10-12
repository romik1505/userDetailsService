package server

import (
	"context"
	http2 "github.com/romik1505/userDetailsService/internal/controller/http"
	"github.com/romik1505/userDetailsService/internal/service"
	"log"
	"os"
	"os/signal"
	"time"

	"net/http"
)

type Server struct {
	httpServer    *http.Server
	personService service.Persons
}

func NewServer(ps service.Persons) Server {
	r := http2.NewHandler(ps).NewRouter()

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r,
	}

	return Server{
		httpServer:    srv,
		personService: ps,
	}
}

func (s *Server) Run() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	go func(ch chan os.Signal) {
		if err := s.httpServer.ListenAndServe(); err != nil {

			log.Println(err.Error())
			done <- os.Interrupt
			return
		}
	}(done)

	log.Printf("Server started on %s port", ":8080")

	<-done
	defer close(done)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	log.Println("Server gracefully closed")

	return s.httpServer.Shutdown(ctx)
}
