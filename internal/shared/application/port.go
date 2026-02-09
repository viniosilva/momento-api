package application

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient interface {
	Ping(ctx context.Context, readPreference *readpref.ReadPref) error
}

type HealthService interface {
	HealthCheck(ctx context.Context) HealthOutput
}
