package repository
import (
	"context"
	"github.com/jackc/pgx/v5")
import "abplatform/internal/model"

type AssignmentRepository struct{
	conn *pgx.Conn
}
func NewAssignmentRepository(conn * pgx.Conn) *AssignmentRepository {
	return &AssignmentRepository{conn: conn}
}
func (r *AssignmentRepository) GetByExperimentAndUser(
	ctx context.Context,
	experimentID int,
	userID string,
) (*model.Assignment, error) {
	var a model.Assignment
	query := `
		SELECT id, experiment_id, user_id, variant, assigned_at
		FROM assignments
		WHERE experiment_id = $1 AND user_id = $2
	`
	err := r.conn.QueryRow(ctx, query, experimentID, userID).Scan(
		&a.ID,
		&a.ExperimentID,
		&a.UserID,
		&a.Variant,
		&a.AssignedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}
func (r *AssignmentRepository) Create(
	ctx context.Context,
	experimentID int,
	userID string,
	variant string,
) (*model.Assignment, error) {
	var a model.Assignment

	query := `
		INSERT INTO assignments (experiment_id, user_id, variant)
		VALUES ($1, $2, $3)
		RETURNING id, experiment_id, user_id, variant, assigned_at
	`

	err := r.conn.QueryRow(ctx, query, experimentID, userID, variant).Scan(
		&a.ID,
		&a.ExperimentID,
		&a.UserID,
		&a.Variant,
		&a.AssignedAt,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}