package app

import (
	"log"
	"os"

	"github.com/floqast/task-management/backend/internal/handler"
	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/repository"
	"github.com/floqast/task-management/backend/internal/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	Logger          *log.Logger
	Middleware      middleware.UserMiddleware
	ActivityHandler *handler.ActivityHandler
	CommentHandler  *handler.CommentHandler
	// LabelHandler        *handler.LabelHandler
	// NotificationHandler *handler.NotificationHandler
	// ProjectHandler      *handler.ProjectHandler
	// SprintHandler       *handler.SprintHandler
	// TaskHandler         *handler.TaskHandler
	UserHandler *handler.UserHandler
	mongoDB     *mongo.Client
}

func NewApplication() (*Application, error) {

	mongoDB, err := repository.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Our Repositories will go here
	userRepository := repository.NewMongoUserRepository(mongoDB)
	commentRepository := repository.NewMongoCommentRepository(mongoDB)
	activityRepository := repository.NewMongoActivityRepository(mongoDB)

	// Our Services will go here
	userService := service.NewUserService(userRepository)
	commentService := service.NewCommentService(commentRepository, activityRepository)
	activityService := service.NewActivityService(activityRepository)

	// // Our Handlers will go here
	// activityHandler := handler.NewActivityHandler()
	// commentHandler := handler.NewCommentHandler()
	// labelHandler := handler.NewLabelHandler()
	// notificationHandler := handler.NewNotificationHandler()
	// projectHandler := handler.NewProjectHandler()
	// sprintHandler := handler.NewSprintHandler()
	// taskHandler := handler.NewTaskHandler()
	middlewareHandler := middleware.UserMiddleware{UserService: *userService}
	userHandler := handler.NewUserHandler(userService, logger)
	commentHandler := handler.NewCommentHandler(commentService, activityService, logger)
	activityHandler := handler.NewActivityHandler(activityService, logger)

	app := &Application{
		Logger:          logger,
		Middleware:      middlewareHandler,
		ActivityHandler: activityHandler,
		UserHandler:     userHandler,
		CommentHandler:  commentHandler,
		mongoDB:         mongoDB,
	}

	return app, nil
}
