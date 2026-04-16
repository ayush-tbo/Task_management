package app

import (
	"log"
	"os"

	"github.com/floqast/task-management/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	Logger *log.Logger
	// ActivityHandler     *handler.ActivityHandler
	// CommentHandler      *handler.CommentHandler
	// LabelHandler        *handler.LabelHandler
	// NotificationHandler *handler.NotificationHandler
	// ProjectHandler      *handler.ProjectHandler
	// SprintHandler       *handler.SprintHandler
	// TaskHandler         *handler.TaskHandler
	// UserHandler         *handler.UserHandler
	mongoDB *mongo.Client
}

func NewApplication() (*Application, error) {

	mongoDB, err := repository.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Our Repositories will go here

	// Our Services will go here

	// // Our Handlers will go here
	// activityHandler := handler.NewActivityHandler()
	// commentHandler := handler.NewCommentHandler()
	// labelHandler := handler.NewLabelHandler()
	// notificationHandler := handler.NewNotificationHandler()
	// projectHandler := handler.NewProjectHandler()
	// sprintHandler := handler.NewSprintHandler()
	// taskHandler := handler.NewTaskHandler()
	// userHandler := handler.NewUserHandler()

	app := &Application{
		Logger:  logger,
		mongoDB: mongoDB,
	}

	return app, nil
}
