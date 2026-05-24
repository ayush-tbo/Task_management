package repository

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func EnsureIndexes(client *mongo.Client, logger *slog.Logger) error {
	db := client.Database("NoSQL")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	type collectionIndex struct {
		collection string
		models     []mongo.IndexModel
	}

	indexes := []collectionIndex{
		// users: lookup by email (unique) and name
		{
			collection: "users",
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "email", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{{Key: "name", Value: 1}},
				},
			},
		},
		// tokens: compound lookup by hash + scope + expiry
		{
			collection: "tokens",
			models: []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "hash", Value: 1},
						{Key: "scope", Value: 1},
						{Key: "expiry", Value: 1},
					},
				},
			},
		},
		// project_members: compound unique on project+user, individual indexes for lookups
		{
			collection: "project_members",
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "project_id", Value: 1}, {Key: "user_id", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{{Key: "user_id", Value: 1}},
				},
			},
		},
		// tasks: queried by project, assignee, status, priority, sprint
		{
			collection: "tasks",
			models: []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "created_at", Value: -1}},
				},
				{
					Keys: bson.D{{Key: "assignee_id", Value: 1}, {Key: "created_at", Value: -1}},
				},
				{
					Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "status", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "priority", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "sprint_id", Value: 1}},
				},
			},
		},
		// comments: queried by task_id
		{
			collection: "comments",
			models: []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "task_id", Value: 1}, {Key: "updated_at", Value: -1}},
				},
			},
		},
		// activity: queried by project and task, sorted by created_at
		{
			collection: "activity",
			models: []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "created_at", Value: -1}},
				},
				{
					Keys: bson.D{{Key: "task_id", Value: 1}, {Key: "created_at", Value: -1}},
				},
			},
		},
		// notifications: queried by user, sorted by is_read + created_at
		{
			collection: "notifications",
			models: []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "is_read", Value: 1}, {Key: "created_at", Value: -1}},
				},
			},
		},
		// sprints: queried by project, optionally filtered by is_active
		{
			collection: "sprints",
			models: []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "is_active", Value: 1}},
				},
			},
		},
		// labels: queried by project
		{
			collection: "labels",
			models: []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "project_id", Value: 1}},
				},
			},
		},
	}

	for _, ci := range indexes {
		coll := db.Collection(ci.collection)
		_, err := coll.Indexes().CreateMany(ctx, ci.models)
		if err != nil {
			logger.Error("failed to create indexes", "collection", ci.collection, "error", err)
			return err
		}
		logger.Info("indexes ensured", "collection", ci.collection)
	}

	return nil
}
