package distance

import (
	"encoding/json"
	"github.com/gauravbansal74/solution/shared/database"
)

type Entity struct {
	ID            string     `bson:"_id" json:"token"`
	Path          [][]string `bson:"path" json:"path"`
	Message       string     `bson:"message"`
	Status        uint       `bson:"status" json:"status"` // 0 Failed, 1 Success, 2 In Progress
	TotalDistance int        `bson:"total_distance" json:"total_distance"`
	TotalTime     int        `bson:"total_time" json:"total_time"`
}

// New Message entity with default values
func New() (*Entity, error) {
	var err error
	entity := &Entity{}
	// Set the default parameters
	entity.Status = 2
	entity.Message = "In Progress"
	entity.ID, err = database.UUID()
	// If error on UUID generation
	if err != nil {
		return entity, err
	}
	return entity, nil
}

// Parse message string to struct values
func Parse(body string) (*Entity, error) {
	var err error
	entity := &Entity{}
	err = json.Unmarshal([]byte(body), &entity)
	if err != nil {
		return entity, err
	}
	return entity, nil
}
