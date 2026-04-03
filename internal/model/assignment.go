package model

import "time"

type Assignment struct {
	ID           int       `json:"id"`
	ExperimentID int       `json:"experiment_id"`
	UserID       string    `json:"user_id"`
	Variant      string    `json:"variant"`
	AssignedAt   time.Time `json:"assigned_at"`
}
