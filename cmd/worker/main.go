package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/jackc/pgx/v5"

	"abplatform/internal/model"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5433/ab_platform?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   "experiment-events",
		GroupID: "metrics-worker",
	})

	log.Println("Worker started...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("error reading:", err)
			continue
		}

		var event model.Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("error decoding:", err)
			continue
		}

		log.Println("Received event:", event)

		processEvent(conn, event)
	}
}
func processEvent(conn *pgx.Conn, event model.Event) {
	ctx := context.Background()

	var variant string

	err := conn.QueryRow(ctx,
		`SELECT variant FROM assignments WHERE experiment_id=$1 AND user_id=$2`,
		event.ExperimentID, event.UserID,
	).Scan(&variant)

	if err != nil {
		log.Println("assignment not found:", err)
		return
	}

	// ===== IMPRESSION =====
	if event.EventName == "impression" {

		// пробуем вставить уникально
		cmdTag, err := conn.Exec(ctx,
			`INSERT INTO user_events (experiment_id, user_id, event_name)
			 VALUES ($1, $2, 'impression')
			 ON CONFLICT DO NOTHING`,
			event.ExperimentID, event.UserID,
		)

		if err != nil {
			log.Println("error inserting impression:", err)
			return
		}

		// если не вставилось → уже был
		if cmdTag.RowsAffected() == 0 {
			return
		}

		// увеличиваем metrics
		if variant == "A" {
			conn.Exec(ctx,
				`UPDATE metrics SET impressions_a = impressions_a + 1 WHERE experiment_id=$1`,
				event.ExperimentID)
		} else {
			conn.Exec(ctx,
				`UPDATE metrics SET impressions_b = impressions_b + 1 WHERE experiment_id=$1`,
				event.ExperimentID)
		}
	}

	// ===== CLICK =====
	if event.EventName == "click" {

		// 1. проверяем был ли impression
		var exists bool
		err := conn.QueryRow(ctx,
			`SELECT EXISTS (
				SELECT 1 FROM user_events
				WHERE experiment_id=$1 AND user_id=$2 AND event_name='impression'
			)`,
			event.ExperimentID, event.UserID,
		).Scan(&exists)

		if err != nil {
			log.Println("error checking impression:", err)
			return
		}

		if !exists {
			log.Println("skip click without impression")
			return
		}

		// 2. пробуем вставить click
		cmdTag, err := conn.Exec(ctx,
			`INSERT INTO user_events (experiment_id, user_id, event_name)
			 VALUES ($1, $2, 'click')
			 ON CONFLICT DO NOTHING`,
			event.ExperimentID, event.UserID,
		)

		if err != nil {
			log.Println("error inserting click:", err)
			return
		}

		// если уже был click → игнор
		if cmdTag.RowsAffected() == 0 {
			return
		}

		// 3. увеличиваем metrics
		if variant == "A" {
			conn.Exec(ctx,
				`UPDATE metrics SET conversions_a = conversions_a + 1 WHERE experiment_id=$1`,
				event.ExperimentID)
		} else {
			conn.Exec(ctx,
				`UPDATE metrics SET conversions_b = conversions_b + 1 WHERE experiment_id=$1`,
				event.ExperimentID)
		}
	}
}