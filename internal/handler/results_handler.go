package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"abplatform/internal/service"
)

type ResultsHandler struct {
	conn *pgx.Conn
}

func NewResultsHandler(conn *pgx.Conn) *ResultsHandler {
	return &ResultsHandler{conn: conn}
}

func (h *ResultsHandler) GetResults(c *gin.Context) {
	expIDStr := c.Query("experiment_id")

	expID, err := strconv.Atoi(expIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid experiment_id"})
		return
	}

	var imprA, imprB, convA, convB int

	err = h.conn.QueryRow(c.Request.Context(),
		`SELECT impressions_a, impressions_b, conversions_a, conversions_b 
		 FROM metrics WHERE experiment_id=$1`,
		expID,
	).Scan(&imprA, &imprB, &convA, &convB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch metrics", "details": err.Error()})
		return
	}

	stats := service.CalculateStats(imprA, convA, imprB, convB)

	c.JSON(http.StatusOK, gin.H{
		"experiment_id": expID,
		"A": gin.H{
			"impressions": imprA,
			"conversions": convA,
			"cr":          stats.CR_A,
		},
		"B": gin.H{
			"impressions": imprB,
			"conversions": convB,
			"cr":          stats.CR_B,
		},
		"uplift":      stats.Uplift,
		"z_score":     stats.ZScore,
		"p_value":     stats.PValue,
		"significant": stats.Significant,
	})
}
