package repository
import (
	"context"

	"github.com/jackc/pgx/v5"
	"abplatform/internal/model"
)
type ExperimentRepository struct {
	conn *pgx.Conn
}

func NewExperimentRepository(conn *pgx.Conn) *ExperimentRepository {
	return &ExperimentRepository{conn: conn}
}

func (r *ExperimentRepository) Create(ctx context.Context, name string) (*model.Experiment, error) {
	var exp model.Experiment

	query := `
		INSERT INTO experiments (name)
		VALUES ($1)
		RETURNING id, name, status, created_at
	`

	err := r.conn.QueryRow(ctx, query, name).Scan(
		&exp.ID,
		&exp.Name,
		&exp.Status,
		&exp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &exp, nil
}