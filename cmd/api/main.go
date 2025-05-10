package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/markusryoti/survey-ddd/internal/adapters/postgres"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

func main() {
	ctx := context.TODO()

	db, err := sql.Open("postgres", "postgres://survey:secret@db:5432/surveydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	repo := postgres.NewPostgresRepository(db, "surveys", func() *surveys.Survey {
		return &surveys.Survey{}
	})

	repo.Save(ctx, &surveys.Survey{})

	// conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ch, err := conn.Channel()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// publisher := &rabbitmq.Publisher{Channel: ch, Exchange: "survey_events"}

	log.Println("Survey service running...")
	select {}
}
