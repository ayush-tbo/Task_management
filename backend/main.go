package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/floqast/task-management/backend/internal/handler"
	"github.com/floqast/task-management/backend/internal/repository"
	"github.com/joho/godotenv"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {

	// Load the .env file from the current directory
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message":"Hello from Go backend!"}`)
	})

	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.yaml"),
	))
	mux.HandleFunc("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	userH := handler.NewUserHandler()
	projectH := handler.NewProjectHandler()
	taskH := handler.NewTaskHandler()
	commentH := handler.NewCommentHandler()
	labelH := handler.NewLabelHandler()
	activityH := handler.NewActivityHandler()
	sprintH := handler.NewSprintHandler()
	notifH := handler.NewNotificationHandler()

	mux.HandleFunc("GET /api/users/me", userH.GetCurrentUser)
	mux.HandleFunc("PUT /api/users/me", userH.UpdateCurrentUser)
	mux.HandleFunc("GET /api/users", userH.ListUsers)
	mux.HandleFunc("POST /api/users/invite", userH.InviteUser)

	mux.HandleFunc("GET /api/projects", projectH.ListProjects)
	mux.HandleFunc("POST /api/projects", projectH.CreateProject)
	mux.HandleFunc("GET /api/projects/{projectId}", projectH.GetProject)
	mux.HandleFunc("PUT /api/projects/{projectId}", projectH.UpdateProject)
	mux.HandleFunc("DELETE /api/projects/{projectId}", projectH.DeleteProject)
	mux.HandleFunc("GET /api/projects/{projectId}/members", projectH.ListProjectMembers)
	mux.HandleFunc("POST /api/projects/{projectId}/members", projectH.AddProjectMember)
	mux.HandleFunc("DELETE /api/projects/{projectId}/members/{userId}", projectH.RemoveProjectMember)
	mux.HandleFunc("GET /api/projects/{projectId}/charts/status", projectH.GetStatusChart)
	mux.HandleFunc("GET /api/projects/{projectId}/charts/priority", projectH.GetPriorityChart)

	mux.HandleFunc("GET /api/projects/{projectId}/tasks", taskH.ListTasks)
	mux.HandleFunc("POST /api/projects/{projectId}/tasks", taskH.CreateTask)
	mux.HandleFunc("GET /api/tasks/my", taskH.GetMyTasks)
	mux.HandleFunc("GET /api/tasks/{taskId}", taskH.GetTask)
	mux.HandleFunc("PUT /api/tasks/{taskId}", taskH.UpdateTask)
	mux.HandleFunc("DELETE /api/tasks/{taskId}", taskH.DeleteTask)
	mux.HandleFunc("PUT /api/tasks/{taskId}/assign", taskH.AssignTask)
	mux.HandleFunc("PUT /api/tasks/{taskId}/status", taskH.UpdateTaskStatus)
	mux.HandleFunc("GET /api/tasks/{taskId}/time", taskH.GetTaskTimeTracking)
	mux.HandleFunc("PUT /api/tasks/{taskId}/time", taskH.LogTaskTime)

	mux.HandleFunc("GET /api/tasks/{taskId}/comments", commentH.ListComments)
	mux.HandleFunc("POST /api/tasks/{taskId}/comments", commentH.CreateComment)
	mux.HandleFunc("PUT /api/comments/{commentId}", commentH.UpdateComment)
	mux.HandleFunc("DELETE /api/comments/{commentId}", commentH.DeleteComment)

	mux.HandleFunc("GET /api/projects/{projectId}/labels", labelH.ListLabels)
	mux.HandleFunc("POST /api/projects/{projectId}/labels", labelH.CreateLabel)
	mux.HandleFunc("PUT /api/labels/{labelId}", labelH.UpdateLabel)
	mux.HandleFunc("DELETE /api/labels/{labelId}", labelH.DeleteLabel)

	mux.HandleFunc("GET /api/projects/{projectId}/activity", activityH.GetProjectActivity)
	mux.HandleFunc("GET /api/tasks/{taskId}/activity", activityH.GetTaskActivity)

	mux.HandleFunc("GET /api/projects/{projectId}/sprints", sprintH.ListSprints)
	mux.HandleFunc("POST /api/projects/{projectId}/sprints", sprintH.CreateSprint)
	mux.HandleFunc("GET /api/sprints/{sprintId}", sprintH.GetSprint)
	mux.HandleFunc("PUT /api/sprints/{sprintId}", sprintH.UpdateSprint)
	mux.HandleFunc("DELETE /api/sprints/{sprintId}", sprintH.DeleteSprint)
	mux.HandleFunc("POST /api/sprints/{sprintId}/tasks", sprintH.AddTaskToSprint)
	mux.HandleFunc("DELETE /api/sprints/{sprintId}/tasks", sprintH.RemoveTaskFromSprint)

	mux.HandleFunc("GET /api/notifications", notifH.ListNotifications)
	mux.HandleFunc("PUT /api/notifications/{notificationId}/read", notifH.MarkNotificationRead)
	mux.HandleFunc("PUT /api/notifications/read-all", notifH.MarkAllNotificationsRead)

	log.Printf("Server starting on :%s", port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/", port)

	_ = repository.ConnectDB()

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
