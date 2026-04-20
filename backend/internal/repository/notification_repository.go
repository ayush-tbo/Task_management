package repository

import (
	"context"

	"github.com/floqast/task-management/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationRepository interface {
	FindByUser(ctx context.Context, userID string) ([]model.Notification, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	Create(ctx context.Context, notification *model.Notification) error
}

// MongoDB Implementations

type MongoNotificationRepository struct {
	collection *mongo.Collection
}

func NewMongoNotificationRepository(db *mongo.Client) *MongoNotificationRepository {
	return &MongoNotificationRepository{
		collection: db.Database("NoSQL").Collection("notifications"),
	}
}

func (m *MongoNotificationRepository) FindByUser(ctx context.Context, userID string) ([]model.Notification, error) {
	var notifications []model.Notification

	cursor, err := m.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &notifications)
	if err != nil {
		return nil, err
	}

	if notifications == nil {
		notifications = []model.Notification{}
	}

	return notifications, nil
}

func (m *MongoNotificationRepository) MarkRead(ctx context.Context, id string) error {
	_, err := m.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_read": true}})
	return err
}

func (m *MongoNotificationRepository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := m.collection.UpdateMany(ctx, bson.M{"user_id": userID, "is_read": false}, bson.M{"$set": bson.M{"is_read": true}})
	return err
}

func (m *MongoNotificationRepository) Create(ctx context.Context, n *model.Notification) error {
	_, err := m.collection.InsertOne(ctx, n)
	return err
}
