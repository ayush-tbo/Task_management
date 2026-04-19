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
	Logger *log.Logger
	// ActivityHandler     *handler.ActivityHandler
	CommentHandler *handler.CommentHandler
	// LabelHandler        *handler.LabelHandler
	// NotificationHandler *handler.NotificationHandler
	ProjectHandler *handler.ProjectHandler
	// SprintHandler       *handler.SprintHandler
	TaskHandler *handler.TaskHandler
	UserHandler *handler.UserHandler
	Middleware  middleware.UserMiddleware
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

	// our services will go here
	userService := service.NewUserService(userRepository)
	commentService := service.NewCommentService(commentRepository)
	projectService := service.NewProjectService(projectRepository, nil)
	taskService := service.NewTaskService(taskRepository)

	// // our handlers will go here
	// activityHandler := handler.NewActivityHandler()
	// commentHandler := handler.NewCommentHandler()
	// labelHandler := handler.NewLabelHandler()
	// notificationHandler := handler.NewNotificationHandler()
	// projectHandler := handler.NewProjectHandler()
	// sprintHandler := handler.NewSprintHandler()
	// taskHandler := handler.NewTaskHandler()
	userHandler := handler.NewUserHandler(userService, logger)
	middlewareHandler := middleware.UserMiddleware{UserService: *userService}
	commentHandler := handler.NewCommentHandler(commentService, logger)
	projectHandler := handler.NewProjectHandler(projectService, taskService, logger)
	taskHandler := handler.NewTaskHandler(taskService, projectService, logger)

	app := &Application{
		Logger:         logger,
		UserHandler:    userHandler,
		Middleware:     middlewareHandler,
		CommentHandler: commentHandler,
		ProjectHandler: projectHandler,
		TaskHandler:    taskHandler,
		mongoDB:        mongoDB,
	}

	return app, nil
}
