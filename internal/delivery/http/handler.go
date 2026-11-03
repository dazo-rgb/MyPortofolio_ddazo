package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/dazo-rgb/ddazo-portfolio/internal/domain"
	"github.com/dazo-rgb/ddazo-portfolio/internal/usecase"
)

type Handler struct {
	contact *usecase.ContactService
	webDir  string
	mux     *http.ServeMux
}

func NewHandler(contact *usecase.ContactService, webDir string) http.Handler {
	h := &Handler{contact: contact, webDir: webDir, mux: http.NewServeMux()}
	h.routes()
	return h
}

func (h *Handler) routes() {
	h.mux.HandleFunc("POST /api/contact", h.submitContact)
	h.mux.Handle("/", http.FileServer(http.Dir(filepath.Clean(h.webDir))))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) submitContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8*1024)
	defer r.Body.Close()

	var message domain.ContactMessage
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&message); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.contact.Submit(r.Context(), message); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, domain.ErrMessageTooLong) {
			status = http.StatusRequestEntityTooLarge
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"message": "Your message has been received by this demo portfolio.",
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
