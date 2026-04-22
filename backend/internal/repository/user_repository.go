package repository

import (
	"context"
	"crypto/sha256"
	"errors"
	"log/slog"
	"time"

	"github.com/floqast/task-management/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByName(ctx context.Context, name string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetToken(ctx context.Context, scope string, tokenPlainText string) (*model.Token, error) // can move to token repo
	CreateToken(ctx context.Context, token *model.Token) error
}

// MongoDB Implementations

type MongoUserRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewMongoUserRepository(db *mongo.Client, logger *slog.Logger) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Database("NoSQL").Collection("users"),
		logger:     logger,
	}
}

func (m *MongoUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User Not Found
		}
		m.logger.Error("repo: find user by id", "error", err, "user_id", id)
		return nil, err
	}
	return &user, nil
}

func (m *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := m.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User Not Found
		}
		m.logger.Error("repo: find user by email", "error", err, "email", email)
		return nil, err
	}
	return &user, nil
}

func (m *MongoUserRepository) FindByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := m.collection.FindOne(ctx, bson.M{"name": name}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User Not Found
		}
		m.logger.Error("repo: find user by name", "error", err, "name", name)
		return nil, err
	}
	return &user, nil
}

func (m *MongoUserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User

	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		m.logger.Error("repo: get all users", "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	if users == nil {
		users = []model.User{}
	}

	return users, nil
}

func (m *MongoUserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := m.collection.InsertOne(ctx, user)
	if err != nil {
		m.logger.Error("repo: create user", "error", err, "user_id", user.ID)
	}
	return err
}

func (m *MongoUserRepository) Update(ctx context.Context, user *model.User) error {
	filter := bson.M{"_id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"name":          user.Name,
			"avatar_url":    user.AvatarURL,
			"updated_at":    user.UpdatedAt,
			"email":         user.Email,
			"password_hash": user.PasswordHash,
		},
	}

	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		m.logger.Error("repo: update user", "error", err, "user_id", user.ID)
	}
	return err
}

func (m *MongoUserRepository) GetToken(ctx context.Context, scope string, plainTextPassword string) (*model.Token, error) {
	tokenHash := sha256.Sum256([]byte(plainTextPassword))

	tokenCollection := m.collection.Database().Collection("tokens")

	var token model.Token

	tokenFilter := bson.M{
		"hash":   tokenHash[:],
		"scope":  scope,
		"expiry": bson.M{"$gt": time.Now()},
	}

	err := tokenCollection.FindOne(ctx, tokenFilter).Decode(&token)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Token not found or expired
		}
		m.logger.Error("repo: get token", "error", err, "scope", scope)
		return nil, err
	}

	return &token, err
}

func (m *MongoUserRepository) CreateToken(ctx context.Context, token *model.Token) error {
	tokenCollection := m.collection.Database().Collection("tokens")
	_, err := tokenCollection.InsertOne(ctx, token)
	if err != nil {
		m.logger.Error("repo: create token", "error", err, "user_id", token.UserID)
	}
	return err
}
