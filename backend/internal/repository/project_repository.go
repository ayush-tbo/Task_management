package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/floqast/task-management/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProjectRepository interface {
	FindByID(ctx context.Context, id string) (*model.Project, error)
	FindByUser(ctx context.Context, userID string, page, pageSize int) ([]model.Project, int, error)
	Create(ctx context.Context, project *model.Project) error
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id string) error
	ListMembers(ctx context.Context, projectID string) ([]model.ProjectMember, error)
	AddMember(ctx context.Context, projectID string, member *model.ProjectMember) error
	RemoveMember(ctx context.Context, projectID, userID string) error
	IncrementMemberCount(ctx context.Context, projectID string, delta int) error
	IncrementTaskCount(ctx context.Context, projectID string, delta int) error
}

type MongoProjectRepository struct {
	collection *mongo.Collection
	members    *mongo.Collection
	tasks      *mongo.Collection
	logger     *slog.Logger
}

func NewMongoProjectRepository(db *mongo.Client, logger *slog.Logger) *MongoProjectRepository {
	return &MongoProjectRepository{
		collection: db.Database("NoSQL").Collection("projects"),
		members:    db.Database("NoSQL").Collection("project_members"),
		tasks:      db.Database("NoSQL").Collection("tasks"),
		logger:     logger,
	}
}

func (m *MongoProjectRepository) FindByID(ctx context.Context, id string) (*model.Project, error) {
	var project model.Project
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // project not found
		}
		m.logger.Error("repo: find project by id", "error", err, "project_id", id)
		return nil, err
	}
	return &project, nil
}

func (m *MongoProjectRepository) Create(ctx context.Context, project *model.Project) error {
	_, err := m.collection.InsertOne(ctx, project)
	if err != nil {
		m.logger.Error("repo: create project", "error", err, "project_id", project.ID)
	}
	return err
}

func (m *MongoProjectRepository) Update(ctx context.Context, project *model.Project) error {
	filter := bson.M{"_id": project.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        project.Name,
			"task_count":  project.TaskCount,
			"description": project.Description,
			"updated_at":  project.UpdatedAt,
		},
	}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		m.logger.Error("repo: update project", "error", err, "project_id", project.ID)
	}
	return err
}

func (m *MongoProjectRepository) Delete(ctx context.Context, id string) error {
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		m.logger.Error("repo: delete project", "error", err, "project_id", id)
	}
	return err
}

func (m *MongoProjectRepository) FindByUser(ctx context.Context, userID string, page, pageSize int) ([]model.Project, int, error) {
	var projects []model.Project
	cursor, err := m.members.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var memberDocs []struct {
		ProjectID string `bson:"project_id"`
	}
	err = cursor.All(ctx, &memberDocs)
	if err != nil {
		return nil, 0, err
	}
	if len(memberDocs) == 0 {
		return []model.Project{}, 0, nil
	}
	projectIDs := make([]string, len(memberDocs))
	for i, m := range memberDocs {
		projectIDs[i] = m.ProjectID
	}
	filter := bson.M{"_id": bson.M{"$in": projectIDs}}
	totalCount, err := m.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize))
	cursor2, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor2.Close(ctx)
	err = cursor2.All(ctx, &projects)
	if err != nil {
		return nil, 0, err
	}

	// Compute member and task counts dynamically
	for i, p := range projects {
		mc, _ := m.members.CountDocuments(ctx, bson.M{"project_id": p.ID})
		tc, _ := m.tasks.CountDocuments(ctx, bson.M{"project_id": p.ID})
		projects[i].MemberCount = int(mc)
		projects[i].TaskCount = int(tc)
	}

	return projects, int(totalCount), nil
}

func (m *MongoProjectRepository) ListMembers(ctx context.Context, projectID string) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	cursor, err := m.members.Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &members)
	if err != nil {
		return nil, err
	}
	if members == nil {
		members = []model.ProjectMember{}
	}
	return members, nil
}

func (m *MongoProjectRepository) AddMember(ctx context.Context, projectID string, member *model.ProjectMember) error {
	_, err := m.members.InsertOne(ctx, bson.M{
		"project_id": projectID,
		"user_id":    member.UserID,
		"name":       member.Name,
		"email":      member.Email,
		"avatar_url": member.AvatarURL,
		"role":       member.Role,
		"joined_at":  member.JoinedAt,
	})
	if err != nil {
		m.logger.Error("repo: add project member", "error", err, "project_id", projectID, "user_id", member.UserID)
	}
	return err
}

func (m *MongoProjectRepository) RemoveMember(ctx context.Context, projectID, userID string) error {
	_, err := m.members.DeleteOne(ctx, bson.M{"project_id": projectID, "user_id": userID})
	if err != nil {
		m.logger.Error("repo: remove project member", "error", err, "project_id", projectID, "user_id", userID)
	}
	return err
}

func (m *MongoProjectRepository) IncrementMemberCount(ctx context.Context, projectID string, delta int) error {
	_, err := m.collection.UpdateOne(ctx, bson.M{"_id": projectID}, bson.M{"$inc": bson.M{"member_count": delta}})
	if err != nil {
		m.logger.Error("repo: increment member count", "error", err, "project_id", projectID, "delta", delta)
	}
	return err
}

func (m *MongoProjectRepository) IncrementTaskCount(ctx context.Context, projectID string, delta int) error {
	_, err := m.collection.UpdateOne(ctx, bson.M{"_id": projectID}, bson.M{"$inc": bson.M{"task_count": delta}})
	if err != nil {
		m.logger.Error("repo: increment task count", "error", err, "project_id", projectID, "delta", delta)
	}
	return err
}
