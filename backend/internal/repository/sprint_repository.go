package repository

import (
	"context"
	"errors"

	"github.com/floqast/task-management/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SprintRepository interface {
	FindByProject(ctx context.Context, projectID string, activeOnly bool) ([]model.Sprint, error)
	FindByID(ctx context.Context, id string) (*model.Sprint, error)
	Create(ctx context.Context, sprint *model.Sprint) error
	Update(ctx context.Context, sprint *model.Sprint) error
	Delete(ctx context.Context, id string) error
	AddTask(ctx context.Context, sprintID, taskID string) error
	RemoveTask(ctx context.Context, sprintID, taskID string) error
}

type MongoSprintRepository struct {
	collection *mongo.Collection
}

func NewMongoSprintRepository(db *mongo.Client) *MongoSprintRepository {
	return &MongoSprintRepository{
		collection: db.Database("NoSQL").Collection("sprints"),
	}
}

func (m *MongoSprintRepository) FindByProject(ctx context.Context, projectID string, activeOnly bool) ([]model.Sprint, error) {
	filter := bson.M{"project_id": projectID}
	if activeOnly {
		filter["is_active"] = true
	}

	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sprints []model.Sprint
	err = cursor.All(ctx, &sprints)
	if err != nil {
		return nil, err
	}
	if sprints == nil {
		sprints = []model.Sprint{}
	}
	return sprints, nil
}

func (m *MongoSprintRepository) FindByID(ctx context.Context, id string) (*model.Sprint, error) {
	var sprint model.Sprint
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sprint)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &sprint, nil
}

func (m *MongoSprintRepository) Create(ctx context.Context, sprint *model.Sprint) error {
	_, err := m.collection.InsertOne(ctx, sprint)
	return err
}

func (m *MongoSprintRepository) Update(ctx context.Context, sprint *model.Sprint) error {
	filter := bson.M{"_id": sprint.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       sprint.Name,
			"label":      sprint.Label,
			"start_date": sprint.StartDate,
			"end_date":   sprint.EndDate,
			"is_active":  sprint.IsActive,
			"task_count": sprint.TaskCount,
			"updated_at": sprint.UpdatedAt,
		},
	}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}

func (m *MongoSprintRepository) Delete(ctx context.Context, id string) error {
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (m *MongoSprintRepository) AddTask(ctx context.Context, sprintID, taskID string) error {
	filter := bson.M{"_id": sprintID}
	update := bson.M{"$inc": bson.M{"task_count": 1}}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}

func (m *MongoSprintRepository) RemoveTask(ctx context.Context, sprintID, taskID string) error {
	filter := bson.M{"_id": sprintID}
	update := bson.M{"$inc": bson.M{"task_count": -1}}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}
