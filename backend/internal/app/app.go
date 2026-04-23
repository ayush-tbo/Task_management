package app

import (
	"log/slog"
	"os"

	"github.com/floqast/task-management/backend/internal/handler"
	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/repository"
	"github.com/floqast/task-management/backend/internal/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	Logger              *slog.Logger
	Middleware          middleware.UserMiddleware
	ActivityHandler     *handler.ActivityHandler
	CommentHandler      *handler.CommentHandler
	LabelHandler        *handler.LabelHandler
	NotificationHandler *handler.NotificationHandler
	ProjectHandler      *handler.ProjectHandler
	SprintHandler       *handler.SprintHandler
	TaskHandler         *handler.TaskHandler
	UserHandler         *handler.UserHandler
	mongoDB             *mongo.Client
}

func NewApplication() (*Application, error) {

	mongoDB, err := repository.ConnectDB()
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// our repositories will go here
	userRepository := repository.NewMongoUserRepository(mongoDB, logger)
	commentRepository := repository.NewMongoCommentRepository(mongoDB, logger)
	projectRepository := repository.NewMongoProjectRepository(mongoDB, logger)
	taskRepository := repository.NewMongoTaskRepository(mongoDB, logger)
	activityRepository := repository.NewMongoActivityRepository(mongoDB, logger)
	notificationRepository := repository.NewMongoNotificationRepository(mongoDB, logger)
	sprintRepository := repository.NewMongoSprintRepository(mongoDB, logger)
	labelRepository := repository.NewMongoLabelRepository(mongoDB, logger)

	// our services will go here
	userService := service.NewUserService(userRepository, logger)
	commentService := service.NewCommentService(commentRepository, activityRepository, logger)
	projectService := service.NewProjectService(projectRepository, activityRepository, logger)
	taskService := service.NewTaskService(taskRepository, logger)
	activityService := service.NewActivityService(activityRepository, logger)
	notificationService := service.NewNotificationService(notificationRepository, logger)
	sprintService := service.NewSprintService(sprintRepository, activityRepository, logger)
	labelService := service.NewLabelService(labelRepository, activityRepository, logger)

	//handlers will go here
	middlewareHandler := middleware.UserMiddleware{UserService: *userService}
	commentHandler := handler.NewCommentHandler(commentService, taskService, activityService, notificationService, logger)
	projectHandler := handler.NewProjectHandler(projectService, taskService, activityService, userService, notificationService, logger)
	taskHandler := handler.NewTaskHandler(taskService, projectService, commentService, activityService, notificationService, logger)
	userHandler := handler.NewUserHandler(userService, logger)
	activityHandler := handler.NewActivityHandler(activityService, logger)
	notificationHandler := handler.NewNotificationHandler(notificationService, logger)
	sprintHandler := handler.NewSprintHandler(sprintService, taskService, logger)
	labelHandler := handler.NewLabelHandler(labelService, logger)

	app := &Application{
		Logger:              logger,
		UserHandler:         userHandler,
		Middleware:          middlewareHandler,
		ActivityHandler:     activityHandler,
		CommentHandler:      commentHandler,
		LabelHandler:        labelHandler,
		NotificationHandler: notificationHandler,
		ProjectHandler:      projectHandler,
		SprintHandler:       sprintHandler,
		TaskHandler:         taskHandler,
		mongoDB:             mongoDB,
	}

	return app, nil
}
