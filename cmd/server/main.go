package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"workshop/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"workshop/internal/handlers"
	"workshop/internal/users"
)

const databaseDSN = "host=postgres port=5432 user=db_user password=db_pass dbname=workshop sslmode=disable"

func main() {
	fmt.Println("starting")
	defer fmt.Println("shutdown")

	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		log.Fatalln(err)
	}

	repo := storage.NewStorage(db)
	us := users.NewService(repo)

	uh := handlers.NewUsers(us, repo)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", uh.Create)
		r.Get("/{userId}", uh.Get)
	})

	s := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	fmt.Println(s.ListenAndServe())
}
