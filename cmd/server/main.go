package main

import (
	"context"
	"fmt"
	"github.com/captainlettuce/sales-stats/internal/repository"
	"github.com/captainlettuce/sales-stats/internal/server"
	"github.com/captainlettuce/sales-stats/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func createServer(ctx context.Context, bind string) *http.Server {
	db, err := sqlx.Open("postgres", "postgres://postgres:superDevPassword@db:5432/test-db?sslmode=disable")
	if err != nil {
		slog.With(slog.Any("err", err)).Error("Could not connect to DB")
		os.Exit(-1)
	}

	timeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err = db.PingContext(timeout); err != nil {
		slog.With(slog.Any("err", err)).Error("Could not connect to DB")
		os.Exit(-1)
	}

	r := repository.New(db)
	s := service.New(r)
	serv := server.New(s)

	return serv.CreateHTTPServer(bind)
}

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	rootContext, cancel := context.WithCancel(context.Background())

	serv := createServer(rootContext, ":"+port)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGINT)

		<-stop

		timeoutCtx, innerCancel := context.WithTimeout(rootContext, time.Second)
		defer innerCancel()
		err := serv.Shutdown(timeoutCtx)
		if err != nil {
			slog.With(slog.Any("err", err)).Error("Server shutdown failed")
			os.Exit(-1)
		}
		cancel()
	}()

	slog.Info(fmt.Sprintf("Listening on http://0.0.0.0:%s", port))
	err := serv.ListenAndServe()
	if err != nil {
		slog.With(slog.Any("err", err)).Error("Http server crashed")
	}
}
