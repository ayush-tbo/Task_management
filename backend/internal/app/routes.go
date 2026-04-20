package app

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *Application) *chi.Mux {

	r := chi.NewRouter()

	r.Post("/api/users/register", app.UserHandler.RegisterUser)
	r.Post("/api/users/login", app.UserHandler.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.Get("/api/users/{id}", app.Middleware.RequireUser(app.UserHandler.GetUserById))
		r.Patch("/api/users/{id}", app.Middleware.RequireUser(app.UserHandler.UpdateUserById))
		r.Get("/api/users", app.Middleware.RequireUser(app.UserHandler.AllUsers))

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

		r.Get("/api/notifications", app.Middleware.RequireUser(app.NotificationHandler.ListNotifications))
		r.Put("/api/notifications/{id}/read", app.Middleware.RequireUser(app.NotificationHandler.MarkNotificationRead))
		r.Put("/api/notifications/read-all", app.Middleware.RequireUser(app.NotificationHandler.MarkAllNotificationsRead))
	})

	// mux.HandleFunc("GET /api/projects/{projectId}/labels", labelH.ListLabels)
	// mux.HandleFunc("POST /api/projects/{projectId}/labels", labelH.CreateLabel)
	// mux.HandleFunc("PUT /api/labels/{labelId}", labelH.UpdateLabel)
	// mux.HandleFunc("DELETE /api/labels/{labelId}", labelH.DeleteLabel)

	// mux.HandleFunc("GET /api/projects/{projectId}/sprints", sprintH.ListSprints)
	// mux.HandleFunc("POST /api/projects/{projectId}/sprints", sprintH.CreateSprint)
	// mux.HandleFunc("GET /api/sprints/{sprintId}", sprintH.GetSprint)
	// mux.HandleFunc("PUT /api/sprints/{sprintId}", sprintH.UpdateSprint)
	// mux.HandleFunc("DELETE /api/sprints/{sprintId}", sprintH.DeleteSprint)
	// mux.HandleFunc("POST /api/sprints/{sprintId}/tasks", sprintH.AddTaskToSprint)
	// mux.HandleFunc("DELETE /api/sprints/{sprintId}/tasks", sprintH.RemoveTaskFromSprint)

	return r
}
