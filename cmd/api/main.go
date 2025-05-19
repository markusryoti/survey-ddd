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
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

func main() {
	db, err := sql.Open("postgres", "postgres://survey:secret@db:5432/surveydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	surveyRepo := postgres.NewPostgresRepository[*surveys.Survey](db, "surveys")
	_ = postgres.NewPostgresRepository[*surveys.SurveyResponse](db, "responses")

	surveyCommandHandler := command.NewSurveyCommandHandler[*surveys.Survey](surveyRepo)

	surveyHandler := rest.SurveyHandler{
		CommandHandler: surveyCommandHandler,
	}

	r := chi.NewRouter()
	surveyHandler.RegisterRoutes(r)

	http.ListenAndServe(":8080", r)
}
