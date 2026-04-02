package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"

	"abplatform/internal/handler"
	"abplatform/internal/repository"
	"abplatform/internal/service"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5433/ab_platform?sslmode=disable")
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}
	defer conn.Close(context.Background())

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal("failed to ping postgres:", err)
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:29092"),
		Topic:    "experiment-events",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	experimentRepo := repository.NewExperimentRepository(conn)
	experimentService := service.NewExperimentService(experimentRepo)
	experimentHandler := handler.NewExperimentHandler(experimentService)

	assignmentRepo := repository.NewAssignmentRepository(conn)
    assignmentService := service.NewAssignmentService(assignmentRepo)
    assignmentHandler := handler.NewAssignmentHandler(assignmentService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"db":     "connected",
			"kafka":  "configured",
		})
	})

	r.GET("/kafka-test", func(c *gin.Context) {
		err := writer.WriteMessages(c.Request.Context(),
			kafka.Message{
				Key:   []byte("test"),
				Value: []byte(`{"message":"hello kafka"}`),
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "message sent",
		})
	})

	r.POST("/experiments", experimentHandler.CreateExperiment)
	r.GET("/assign", assignmentHandler.Assign)

	log.Println("API running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}