package model
//import "time"
type Event struct {
	ExperimentID int    `json:"experiment_id"`
	UserID       string `json:"user_id"`
	EventName    string `json:"event_name"`
}
