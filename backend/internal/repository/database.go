package repository

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		panic("No Database URI Found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		slog.Error("could not connect to MongoDB", "error", err)
		return nil, err
	}

	slog.Info("connected to MongoDB")
	return client, nil
}
