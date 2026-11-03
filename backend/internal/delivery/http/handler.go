package httpdelivery

import (
	"dazo-portfolio/internal/usecase"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

type Handler struct{ projects *usecase.ProjectService }

func NewHandler(projects *usecase.ProjectService) http.Handler {
	h := &Handler{projects: projects}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.health)
	mux.HandleFunc("GET /api/projects", h.listProjects)
	mux.HandleFunc("GET /api/projects/{slug}", h.getProject)
	return mux
}
func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "dazo-portfolio-api"})
}
func (h *Handler) listProjects(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.projects.ListFeatured())
}
func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	project, err := h.projects.GetBySlug(strings.TrimSpace(r.PathValue("slug")))
	switch {
	case errors.Is(err, usecase.ErrInvalidSlug):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid project slug"})
	case errors.Is(err, usecase.ErrProjectNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "project not found"})
	case err != nil:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	default:
		writeJSON(w, http.StatusOK, project)
	}
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("encode HTTP JSON response", "error", err, "status", status)
	}
}
