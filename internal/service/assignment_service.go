package service

import (
	"context"
	"crypto/sha1"
	"fmt"

	"abplatform/internal/model"
	"abplatform/internal/repository"
)

type AssignmentService struct {
	repo *repository.AssignmentRepository
}

func NewAssignmentService(repo *repository.AssignmentRepository) *AssignmentService {
	return &AssignmentService{repo: repo}
}

func (s *AssignmentService) AssignUser(
	ctx context.Context,
	experimentID int,
	userID string,
) (*model.Assignment, error) {
	existing, err := s.repo.GetByExperimentAndUser(ctx, experimentID, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	variant := chooseVariant(experimentID, userID)
	created, err := s.repo.Create(ctx, experimentID, userID, variant)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func chooseVariant(experimentID int, userID string) string {
	input := fmt.Sprintf("%d:%s", experimentID, userID)
	hash := sha1.Sum([]byte(input))

	if hash[0]%2 == 0 {
		return "A"
	}
	return "B"
}
