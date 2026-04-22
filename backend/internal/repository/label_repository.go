package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LabelRepository interface {
	FindByProject(ctx context.Context, projectID string) ([]model.Label, error)
	FindByID(ctx context.Context, id string) (*model.Label, error)
	Create(ctx context.Context, label *model.Label) error
	Update(ctx context.Context, label *model.Label) error
	Delete(ctx context.Context, id string) error
}

type MongoLabelRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewMongoLabelRepository(db *mongo.Client, logger *slog.Logger) *MongoLabelRepository {
	return &MongoLabelRepository{
		collection: db.Database("NoSQL").Collection("labels"),
		logger:     logger,
	}
}

func (m *MongoLabelRepository) FindByProject(ctx context.Context, projectID string) ([]model.Label, error) {
	cursor, err := m.collection.Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		m.logger.Error("repo: find labels by project", "error", err, "project_id", projectID)
		return nil, err
	}
	defer cursor.Close(ctx)

	var labels []model.Label
	err = cursor.All(ctx, &labels)
	if err != nil {
		return nil, err
	}
	if labels == nil {
		labels = []model.Label{}
	}
	return labels, nil
}

func (m *MongoLabelRepository) FindByID(ctx context.Context, id string) (*model.Label, error) {
	var label model.Label
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&label)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		m.logger.Error("repo: find label by id", "error", err, "label_id", id)
		return nil, err
	}
	return &label, nil
}

func (m *MongoLabelRepository) Create(ctx context.Context, label *model.Label) error {
	_, err := m.collection.InsertOne(ctx, label)
	if err != nil {
		m.logger.Error("repo: create label", "error", err, "label_id", label.ID)
	}
	return err
}

func (m *MongoLabelRepository) Update(ctx context.Context, label *model.Label) error {
	filter := bson.M{"_id": label.ID}
	update := bson.M{
		"$set": bson.M{
			"name":  label.Name,
			"color": label.Color,
		},
	}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		m.logger.Error("repo: update label", "error", err, "label_id", label.ID)
	}
	return err
}

func (m *MongoLabelRepository) Delete(ctx context.Context, id string) error {
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		m.logger.Error("repo: delete label", "error", err, "label_id", id)
	}
	return err
}
