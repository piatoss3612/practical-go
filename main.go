package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"syscall"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	// 1. DB 연결
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB connection established")

	srv := &Server{
		Server: &http.Server{
			Addr: ":8080",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte("Hello World"))
			}),
		},
		db: db,
	}

	log.Println("Starting server...")

	go func() {
		// 2. 서버 시작
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx2, cancel2 := GracefulShutdownCtx(ctx, func() {
		// 3. 서버 종료
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}

		// 4. DB 연결 종료
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}

		log.Println("DB connection closed")

		log.Println("Server shutdown complete")
	}, syscall.SIGINT, syscall.SIGTERM)
	defer cancel2()

	<-ctx2.Done()
}

type Server struct {
	*http.Server
	db *sql.DB
}

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./sql.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
