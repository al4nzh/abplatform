package main

import (
	"context"
	"log"
	//"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"abplatform/internal/handler"
	"abplatform/internal/repository"
	"abplatform/internal/service"
	"abplatform/internal/kafka"
)

func main() {
	conn, err := pgx.Connect(context.Background(),
		"postgres://postgres:postgres@localhost:5433/ab_platform?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	r := gin.Default()

	// repos
	experimentRepo := repository.NewExperimentRepository(conn)
	assignmentRepo := repository.NewAssignmentRepository(conn)

	// services
	experimentService := service.NewExperimentService(experimentRepo)
	assignmentService := service.NewAssignmentService(assignmentRepo)

	// handlers
	experimentHandler := handler.NewExperimentHandler(experimentService)
	assignmentHandler := handler.NewAssignmentHandler(assignmentService)
	resultsHandler := handler.NewResultsHandler(conn)

	// kafka
	producer := kafka.NewProducer("localhost:29092")
	eventHandler := handler.NewEventHandler(producer)
	
	// routes
	r.POST("/experiments", experimentHandler.CreateExperiment)
	r.GET("/assign", assignmentHandler.Assign)
	r.POST("/events", eventHandler.TrackEvent)
	r.GET("/results", resultsHandler.GetResults)

	log.Println("API running on :8080")
	r.Run(":8080")
}