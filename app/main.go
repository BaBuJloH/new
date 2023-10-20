package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"new/repositories/orders_repo"
	"new/server"
	"new/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	ctx := context.Background()

	fmt.Println("Start apipost server\n")
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	cfg.Print()

	conn, err := pgxpool.New(ctx, cfg.PgConnectUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	repos := service.Repositories{
		Orders: orders_repo.New(ctx, conn),
	}

	srvc := service.New(repos)

	srv := &http.Server{
		Addr:    cfg.Listen,
		Handler: server.New(cfg.Listen, srvc),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown, err: %s", err.Error())
		os.Exit(0)
	}
	log.Print("Shutting down server...")
}
