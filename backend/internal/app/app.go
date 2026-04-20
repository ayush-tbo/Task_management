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
	userRepository := repository.NewMongoUserRepository(mongoDB)
	commentRepository := repository.NewMongoCommentRepository(mongoDB)
	projectRepository := repository.NewMongoProjectRepository(mongoDB)
	taskRepository := repository.NewMongoTaskRepository(mongoDB)
	activityRepository := repository.NewMongoActivityRepository(mongoDB)
	notificationRepository := repository.NewMongoNotificationRepository(mongoDB)
	sprintRepository := repository.NewMongoSprintRepository(mongoDB)
	labelRepository := repository.NewMongoLabelRepository(mongoDB)

	// our services will go here
	userService := service.NewUserService(userRepository)
	commentService := service.NewCommentService(commentRepository, activityRepository)
	projectService := service.NewProjectService(projectRepository, activityRepository)
	taskService := service.NewTaskService(taskRepository)
	activityService := service.NewActivityService(activityRepository)
	notificationService := service.NewNotificationService(notificationRepository)
	sprintService := service.NewSprintService(sprintRepository, activityRepository)
	labelService := service.NewLabelService(labelRepository, activityRepository)

	//handlers will go here
	middlewareHandler := middleware.UserMiddleware{UserService: *userService}
	commentHandler := handler.NewCommentHandler(commentService, activityService, logger)
	projectHandler := handler.NewProjectHandler(projectService, taskService, activityService, logger)
	taskHandler := handler.NewTaskHandler(taskService, projectService, activityService, logger)
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
