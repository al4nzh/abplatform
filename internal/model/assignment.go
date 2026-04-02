package model
import "time"
type Assignment struct {
	ID 		 int `json:"id"`
	ExperimentID int `json:"experiment_id"`
	UserID 	 int `json:"user_id"`
	GroupID 	 int `json:"group_id"`
	AssignedAt time.Time `json:"assigned_at"`
}