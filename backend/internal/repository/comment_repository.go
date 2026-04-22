package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CommentRepository interface {
	FindByTask(ctx context.Context, taskID string) ([]model.Comment, error)
	FindByID(ctx context.Context, id string) (*model.Comment, error)
	Create(ctx context.Context, comment *model.Comment) error
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, taskID string) error
}

// MongoDB Implementations

type MongoCommentRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewMongoCommentRepository(db *mongo.Client, logger *slog.Logger) *MongoCommentRepository {
	return &MongoCommentRepository{
		collection: db.Database("NoSQL").Collection("comments"),
		logger:     logger,
	}
}

func (m *MongoCommentRepository) FindByTask(ctx context.Context, taskID string) ([]model.Comment, error) {
	var comments []model.Comment

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"task_id": taskID}}},
		{{Key: "$sort", Value: bson.D{{Key: "updated_at", Value: -1}}}},
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
		m.logger.Error("repo: find comments by task aggregate", "error", err, "task_id", taskID)
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &comments)
	if err != nil {
		m.logger.Error("repo: find comments by task decode", "error", err, "task_id", taskID)
		return nil, err
	}

	if comments == nil {
		comments = []model.Comment{}
	}

	return comments, nil
}

func (m *MongoCommentRepository) FindByID(ctx context.Context, id string) (*model.Comment, error) {
	var comment model.Comment

	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Comment Not Found
		}
		m.logger.Error("repo: find comment by id", "error", err, "comment_id", id)
		return nil, err
	}
	return &comment, nil
}

func (m *MongoCommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	_, err := m.collection.InsertOne(ctx, comment)
	if err != nil {
		m.logger.Error("repo: create comment", "error", err, "comment_id", comment.ID, "task_id", comment.TaskID)
	}
	return err
}

func (m *MongoCommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	filter := bson.M{"_id": comment.ID}

	update := bson.M{
		"$set": bson.M{
			"content":    comment.Content,
			"updated_at": comment.UpdatedAt,
		},
	}

	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		m.logger.Error("repo: update comment", "error", err, "comment_id", comment.ID)
	}
	return err
}

func (m *MongoCommentRepository) Delete(ctx context.Context, id string) error {
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		m.logger.Error("repo: delete comment", "error", err, "comment_id", id)
	}
	return err
}

func (m *MongoCommentRepository) DeleteAll(ctx context.Context, taskID string) error {
	_, err := m.collection.DeleteMany(ctx, bson.M{"task_id": taskID})
	if err != nil {
		m.logger.Error("repo: delete all comments", "error", err, "task_id", taskID)
	}
	return err
}
