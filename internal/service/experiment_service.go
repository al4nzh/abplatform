package service

import (
	"context"

	"abplatform/internal/model"
	"abplatform/internal/repository"
)

type ExperimentService struct {
	repo *repository.ExperimentRepository
}

func NewExperimentService(repo *repository.ExperimentRepository) *ExperimentService {
	return &ExperimentService{repo: repo}
}

func (s *ExperimentService) CreateExperiment(ctx context.Context, name string) (*model.Experiment, error) {
	return s.repo.Create(ctx, name)
}