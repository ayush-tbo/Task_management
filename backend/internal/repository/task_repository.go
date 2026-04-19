package repository

import (
	"context"
	"errors"
	"time"

	"github.com/floqast/task-management/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TaskRepository interface {
	FindByID(ctx context.Context, id string) (*model.Task, error)
	FindByProject(ctx context.Context, projectID string, filters TaskFilters, page, pageSize int) ([]model.Task, int, error)
	FindByAssignee(ctx context.Context, userID string, filters TaskFilters, page, pageSize int) ([]model.Task, int, error)
	Create(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
	CountByStatus(ctx context.Context, projectID string) ([]model.StatusChartEntry, error)
	CountByPriority(ctx context.Context, projectID string) ([]model.PriorityChartEntry, error)
}

type TaskFilters struct {
	Status     *model.TaskStatus
	Priority   *model.Priority
	AssigneeID *string
	ReporterID *string
	LabelIDs   []string
	SprintID   *string
	DueBefore  *string
	DueAfter   *string
	IsPastDue  *bool
	SortBy     string
	SortOrder  string
}

type MongoTaskRepository struct {
	collection *mongo.Collection
}

func NewMongoTaskRepository(db *mongo.Client) *MongoTaskRepository {
	return &MongoTaskRepository{
		collection: db.Database("NoSQL").Collection("tasks"),
	}
}

func (m *MongoTaskRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	var task model.Task
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	if task.DueDate != nil && task.DueDate.Before(time.Now()) && task.Status != model.StatusDone {
		task.IsPastDue = true
	}
	return &task, nil
}

func (m *MongoTaskRepository) Create(ctx context.Context, task *model.Task) error {
	_, err := m.collection.InsertOne(ctx, task)
	return err
}

func (m *MongoTaskRepository) Update(ctx context.Context, task *model.Task) error {
	filter := bson.M{"_id": task.ID}
	update := bson.M{
		"$set": bson.M{
			"title":         task.Title,
			"description":   task.Description,
			"status":        task.Status,
			"priority":      task.Priority,
			"due_date":      task.DueDate,
			"assignee_id":   task.AssigneeID,
			"label_ids":     task.LabelIDs,
			"sprint_id":     task.SprintID,
			"comment_count": task.CommentCount,
			"time_tracking": task.TimeTracking,
			"updated_at":    task.UpdatedAt,
		},
	}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}

func (m *MongoTaskRepository) Delete(ctx context.Context, id string) error {
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func buildTaskFilter(baseFilter bson.M, filters TaskFilters) bson.M {
	if filters.Status != nil {
		baseFilter["status"] = *filters.Status
	}
	if filters.Priority != nil {
		baseFilter["priority"] = *filters.Priority
	}
	if filters.AssigneeID != nil {
		baseFilter["assignee_id"] = *filters.AssigneeID
	}
	if filters.ReporterID != nil {
		baseFilter["reporter_id"] = *filters.ReporterID
	}
	if len(filters.LabelIDs) > 0 {
		baseFilter["label_ids"] = bson.M{"$in": filters.LabelIDs}
	}
	if filters.SprintID != nil {
		baseFilter["sprint_id"] = *filters.SprintID
	}
	return baseFilter
}

func buildSortOption(filters TaskFilters) bson.D {
	sortBy := "created_at"
	sortOrder := -1
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	if filters.SortOrder == "asc" {
		sortOrder = 1
	}
	return bson.D{{Key: sortBy, Value: sortOrder}}
}

func (m *MongoTaskRepository) FindByProject(ctx context.Context, projectID string, filters TaskFilters, page, pageSize int) ([]model.Task, int, error) {
	filter := buildTaskFilter(bson.M{"project_id": projectID}, filters)

	total, err := m.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(buildSortOption(filters))

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, 0, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	now := time.Now()
	for i := range tasks {
		if tasks[i].DueDate != nil && tasks[i].DueDate.Before(now) && tasks[i].Status != model.StatusDone {
			tasks[i].IsPastDue = true
		}
	}

	return tasks, int(total), nil
}

func (m *MongoTaskRepository) FindByAssignee(ctx context.Context, userID string, filters TaskFilters, page, pageSize int) ([]model.Task, int, error) {
	filter := buildTaskFilter(bson.M{"assignee_id": userID}, filters)

	total, err := m.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(buildSortOption(filters))

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, 0, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	now := time.Now()
	for i := range tasks {
		if tasks[i].DueDate != nil && tasks[i].DueDate.Before(now) && tasks[i].Status != model.StatusDone {
			tasks[i].IsPastDue = true
		}
	}

	return tasks, int(total), nil
}

func (m *MongoTaskRepository) CountByStatus(ctx context.Context, projectID string) ([]model.StatusChartEntry, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"project_id": projectID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Status model.TaskStatus `bson:"_id"`
		Count  int              `bson:"count"`
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	entries := make([]model.StatusChartEntry, len(results))
	for i, r := range results {
		entries[i] = model.StatusChartEntry{Status: r.Status, Count: r.Count}
	}
	return entries, nil
}

func (m *MongoTaskRepository) CountByPriority(ctx context.Context, projectID string) ([]model.PriorityChartEntry, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"project_id": projectID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$priority",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Priority model.Priority `bson:"_id"`
		Count    int            `bson:"count"`
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	entries := make([]model.PriorityChartEntry, len(results))
	for i, r := range results {
		entries[i] = model.PriorityChartEntry{Priority: r.Priority, Count: r.Count}
	}
	return entries, nil
}
