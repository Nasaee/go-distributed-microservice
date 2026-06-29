package data

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var client *mongo.Client

func New(m *mongo.Client) Models {
	client = m

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}
