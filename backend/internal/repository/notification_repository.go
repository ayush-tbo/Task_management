package repository

import (
	"context"
	"log/slog"

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
	logger     *slog.Logger
}

func NewMongoNotificationRepository(db *mongo.Client, logger *slog.Logger) *MongoNotificationRepository {
	return &MongoNotificationRepository{
		collection: db.Database("NoSQL").Collection("notifications"),
		logger:     logger,
	}
}

func (m *MongoNotificationRepository) FindByUser(ctx context.Context, userID string) ([]model.Notification, error) {
	var notifications []model.Notification

	cursor, err := m.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		m.logger.Error("repo: find notifications by user", "error", err, "user_id", userID)
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
	if err != nil {
		m.logger.Error("repo: mark notification read", "error", err, "notification_id", id)
	}
	return err
}

func (m *MongoNotificationRepository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := m.collection.UpdateMany(ctx, bson.M{"user_id": userID, "is_read": false}, bson.M{"$set": bson.M{"is_read": true}})
	if err != nil {
		m.logger.Error("repo: mark all notifications read", "error", err, "user_id", userID)
	}
	return err
}

func (m *MongoNotificationRepository) Create(ctx context.Context, n *model.Notification) error {
	_, err := m.collection.InsertOne(ctx, n)
	if err != nil {
		m.logger.Error("repo: create notification", "error", err, "user_id", n.UserID)
	}
	return err
}
