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
	NotificationHandler *handler.NotificationHandler
	ProjectHandler      *handler.ProjectHandler
	// SprintHandler       *handler.SprintHandler
	TaskHandler *handler.TaskHandler
	UserHandler *handler.UserHandler
	mongoDB     *mongo.Client
}

func NewApplication() (*Application, error) {

	mongoDB, err := repository.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// our repositories will go here
	userRepository := repository.NewMongoUserRepository(mongoDB)
	commentRepository := repository.NewMongoCommentRepository(mongoDB)
	projectRepository := repository.NewMongoProjectRepository(mongoDB)
	taskRepository := repository.NewMongoTaskRepository(mongoDB)
	activityRepository := repository.NewMongoActivityRepository(mongoDB)
	notificationRepository := repository.NewMongoNotificationRepository(mongoDB)

	// our services will go here
	userService := service.NewUserService(userRepository)
	commentService := service.NewCommentService(commentRepository, activityRepository)
	projectService := service.NewProjectService(projectRepository, nil)
	taskService := service.NewTaskService(taskRepository)
	activityService := service.NewActivityService(activityRepository)
	notificationService := service.NewNotificationService(notificationRepository)

	// // our handlers will go here
	// activityHandler := handler.NewActivityHandler()
	// commentHandler := handler.NewCommentHandler()
	// labelHandler := handler.NewLabelHandler()
	// notificationHandler := handler.NewNotificationHandler()
	// projectHandler := handler.NewProjectHandler()
	// sprintHandler := handler.NewSprintHandler()
	// taskHandler := handler.NewTaskHandler()
	middlewareHandler := middleware.UserMiddleware{UserService: *userService}
	commentHandler := handler.NewCommentHandler(commentService, activityService, logger)
	projectHandler := handler.NewProjectHandler(projectService, taskService, logger)
	taskHandler := handler.NewTaskHandler(taskService, projectService, logger)
	userHandler := handler.NewUserHandler(userService, logger)
	activityHandler := handler.NewActivityHandler(activityService, logger)
	notificationHandler := handler.NewNotificationHandler(notificationService, logger)

	app := &Application{
		Logger:              logger,
		UserHandler:         userHandler,
		Middleware:          middlewareHandler,
		ActivityHandler:     activityHandler,
		CommentHandler:      commentHandler,
		ProjectHandler:      projectHandler,
		TaskHandler:         taskHandler,
		NotificationHandler: notificationHandler,
		mongoDB:             mongoDB,
	}

	return app, nil
}
