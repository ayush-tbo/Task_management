package repository

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ActivityRepository interface {
	FindByProject(ctx context.Context, projectID string) ([]model.ActivityEntry, error)
	FindByTask(ctx context.Context, taskID string) ([]model.ActivityEntry, error)
	Create(ctx context.Context, entry *model.ActivityEntry) error
}

// MongoDB Implementations

type MongoActivityRepository struct {
	collection *mongo.Collection
}

func NewMongoActivityRepository(db *mongo.Client) *MongoActivityRepository {
	return &MongoActivityRepository{
		collection: db.Database("NoSQL").Collection("activity"),
	}
}

func (m *MongoActivityRepository) FindByProject(ctx context.Context, projectID string) ([]model.ActivityEntry, error) {
	var activities []model.ActivityEntry

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"project_id": projectID}}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}}},
	}

	cursor, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &activities)
	if err != nil {
		return nil, err
	}

	if activities == nil {
		activities = []model.ActivityEntry{}
	}

	return activities, nil
}

func (m *MongoActivityRepository) FindByTask(ctx context.Context, taskID string) ([]model.ActivityEntry, error) {
	var activities []model.ActivityEntry

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"task_id": taskID}}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}}},
	}

	cursor, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &activities)
	if err != nil {
		return nil, err
	}

	if activities == nil {
		activities = []model.ActivityEntry{}
	}

	return activities, nil
}

func (m *MongoActivityRepository) Create(ctx context.Context, activity *model.ActivityEntry) error {
	_, err := m.collection.InsertOne(ctx, activity)
	return err
}
