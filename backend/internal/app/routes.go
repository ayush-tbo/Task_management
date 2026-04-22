package app

import (
	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *Application) *chi.Mux {

	r := chi.NewRouter()

	// log every request: method, path, user, status, duration
	r.Use(middleware.RequestLogger(app.Logger))

	// Auth routes
	r.Post("/api/users/register", app.UserHandler.RegisterUser)
	r.Post("/api/users/login", app.UserHandler.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		// user routes
		r.Get("/api/users/{id}", app.Middleware.RequireUser(app.UserHandler.GetUserById))
		r.Patch("/api/users/{id}", app.Middleware.RequireUser(app.UserHandler.UpdateUserById))
		r.Get("/api/users", app.Middleware.RequireUser(app.UserHandler.AllUsers))

		// comment routes
		r.Post("/api/tasks/{id}/comments", app.Middleware.RequireUser(app.CommentHandler.CreateComment))
		r.Get("/api/tasks/{id}/comments", app.Middleware.RequireUser(app.CommentHandler.ListComments))
		r.Put("/api/comments/{id}", app.Middleware.RequireUser(app.CommentHandler.UpdateComment))
		r.Delete("/api/comments/{id}", app.Middleware.RequireUser(app.CommentHandler.DeleteComment))

		// project routes
		r.Get("/api/projects", app.Middleware.RequireUser(app.ProjectHandler.ListProjects))
		r.Post("/api/projects", app.Middleware.RequireUser(app.ProjectHandler.CreateProject))
		r.Get("/api/projects/{id}", app.Middleware.RequireUser(app.ProjectHandler.GetProject))
		r.Put("/api/projects/{id}", app.Middleware.RequireUser(app.ProjectHandler.UpdateProject))
		r.Delete("/api/projects/{id}", app.Middleware.RequireUser(app.ProjectHandler.DeleteProject))
		r.Get("/api/projects/{id}/members", app.Middleware.RequireUser(app.ProjectHandler.ListProjectMembers))
		r.Post("/api/projects/{id}/members", app.Middleware.RequireUser(app.ProjectHandler.AddProjectMember))
		r.Delete("/api/projects/{id}/members/{userId}", app.Middleware.RequireUser(app.ProjectHandler.RemoveProjectMember))
		r.Get("/api/projects/{id}/charts/status", app.Middleware.RequireUser(app.ProjectHandler.GetStatusChart))
		r.Get("/api/projects/{id}/charts/priority", app.Middleware.RequireUser(app.ProjectHandler.GetPriorityChart))

		// task routes
		r.Get("/api/projects/{id}/tasks", app.Middleware.RequireUser(app.TaskHandler.ListTasks))
		r.Post("/api/projects/{id}/tasks", app.Middleware.RequireUser(app.TaskHandler.CreateTask))
		r.Get("/api/tasks/my", app.Middleware.RequireUser(app.TaskHandler.GetMyTasks))
		r.Get("/api/tasks/{id}", app.Middleware.RequireUser(app.TaskHandler.GetTask))
		r.Put("/api/tasks/{id}", app.Middleware.RequireUser(app.TaskHandler.UpdateTask))
		r.Delete("/api/tasks/{id}", app.Middleware.RequireUser(app.TaskHandler.DeleteTask))
		r.Put("/api/tasks/{id}/assign", app.Middleware.RequireUser(app.TaskHandler.AssignTask))
		r.Put("/api/tasks/{id}/status", app.Middleware.RequireUser(app.TaskHandler.UpdateTaskStatus))
		r.Get("/api/tasks/{id}/time", app.Middleware.RequireUser(app.TaskHandler.GetTaskTimeTracking))
		r.Put("/api/tasks/{id}/time", app.Middleware.RequireUser(app.TaskHandler.LogTaskTime))
		r.Get("/api/projects/{id}/activity", app.Middleware.RequireUser(app.ActivityHandler.GetProjectActivity))
		r.Get("/api/tasks/{id}/activity", app.Middleware.RequireUser(app.ActivityHandler.GetTaskActivity))

		// notification routes
		r.Get("/api/notifications", app.Middleware.RequireUser(app.NotificationHandler.ListNotifications))
		r.Put("/api/notifications/{id}/read", app.Middleware.RequireUser(app.NotificationHandler.MarkNotificationRead))
		r.Put("/api/notifications/read-all", app.Middleware.RequireUser(app.NotificationHandler.MarkAllNotificationsRead))

		// label routes
		r.Get("/api/projects/{id}/labels", app.Middleware.RequireUser(app.LabelHandler.ListLabels))
		r.Post("/api/projects/{id}/labels", app.Middleware.RequireUser(app.LabelHandler.CreateLabel))
		r.Put("/api/labels/{id}", app.Middleware.RequireUser(app.LabelHandler.UpdateLabel))
		r.Delete("/api/labels/{id}", app.Middleware.RequireUser(app.LabelHandler.DeleteLabel))

		// sprint routes
		r.Get("/api/projects/{id}/sprints", app.Middleware.RequireUser(app.SprintHandler.ListSprints))
		r.Post("/api/projects/{id}/sprints", app.Middleware.RequireUser(app.SprintHandler.CreateSprint))
		r.Get("/api/sprints/{id}", app.Middleware.RequireUser(app.SprintHandler.GetSprint))
		r.Put("/api/sprints/{id}", app.Middleware.RequireUser(app.SprintHandler.UpdateSprint))
		r.Delete("/api/sprints/{id}", app.Middleware.RequireUser(app.SprintHandler.DeleteSprint))
		r.Post("/api/sprints/{id}/tasks", app.Middleware.RequireUser(app.SprintHandler.AddTaskToSprint))
		r.Delete("/api/sprints/{id}/tasks", app.Middleware.RequireUser(app.SprintHandler.RemoveTaskFromSprint))
	})

	return r
}
