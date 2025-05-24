package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/markusryoti/survey-ddd/internal/adapters/postgres"
	"github.com/markusryoti/survey-ddd/internal/adapters/rest"
	"github.com/markusryoti/survey-ddd/internal/application/command"
	"github.com/markusryoti/survey-ddd/internal/application/query"
)

func main() {
	db, err := sql.Open("postgres", "postgres://survey:secret@db:5432/surveydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	transactional := postgres.NewPostgresTransactionalProvider(db)

	surveyCommandHandler := command.NewCommandHandler(transactional)
	queryHandler := query.NewQueryHandler(transactional)

	surveyHandler := rest.SurveyHandler{
		CommandHandler: surveyCommandHandler,
		QueryHandler:   queryHandler,
	}

	r := chi.NewRouter()
	surveyHandler.RegisterRoutes(r)

	http.ListenAndServe(":8080", r)
}
