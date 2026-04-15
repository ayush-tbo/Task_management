package handler

import "net/http"

type ProjectHandler struct{}

func NewProjectHandler() *ProjectHandler { return &ProjectHandler{} }

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ProjectHandler) ListProjectMembers(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ProjectHandler) AddProjectMember(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ProjectHandler) RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ProjectHandler) GetStatusChart(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *ProjectHandler) GetPriorityChart(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}
