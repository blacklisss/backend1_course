package main

import (
	"context"
	"gb/backend1_course/internal/infrastructure/api/handlers"
	"gb/backend1_course/internal/infrastructure/api/routeropenapi"
	"gb/backend1_course/internal/infrastructure/db/pgstore"
	"gb/backend1_course/internal/infrastructure/server"
	"gb/backend1_course/internal/usecases/app/repos/linkrepo"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	lst, err := pgstore.NewLinks(os.Getenv("PG_DSN")) // "postgres://postgres:password@localhost/test?sslmode=disable"
	if err != nil {
		log.Fatal(err)
	}

	ls := linkrepo.NewLinks(lst)
	hs := handlers.NewHandlers(ls)

	h := routeropenapi.Handler(hs)
	srv := server.NewServer(":8000", h)

	srv.Start(ls)
	log.Print("Start")

	<-ctx.Done()

	srv.Stop()
	cancel()
	lst.Close()

	log.Print("Exit")
}
