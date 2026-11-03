package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	httpdelivery "dazo-portfolio/internal/delivery/http"
	"dazo-portfolio/internal/domain"
	"dazo-portfolio/internal/repository"
	"dazo-portfolio/internal/usecase"
)

type Project = domain.Project

var _ usecase.ProjectRepository = (*repository.InMemoryProjectRepository)(nil)

const contentSecurityPolicy = "default-src 'self'; base-uri 'self'; object-src 'none'; frame-ancestors 'none'; form-action 'self'; script-src 'self' 'sha256-4Bc1lM0SKDNmGuiyfLFEAsFTdoPeR5K+53J2KUfi1C4='; style-src 'self' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com data:; img-src 'self' data: https://ghchart.rshah.org; connect-src 'self'; upgrade-insecure-requests"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	if err := run(logger); err != nil {
		logger.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	if _, err := net.LookupPort("tcp", port); err != nil {
		return fmt.Errorf("invalid PORT %q: %w", port, err)
	}

	frontendDir, err := findFrontendDir()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           requestLogger(logger, newHandler(frontendDir)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", server.Addr, err)
	}
	logger.Info("server started", "address", listener.Addr().String(), "port", port, "frontend_dir", frontendDir)

	signalContext, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()
	return serveUntilSignal(server, listener, signalContext.Done(), logger, 10*time.Second)
}

func serveUntilSignal(server *http.Server, listener net.Listener, shutdownSignal <-chan struct{}, logger *slog.Logger, shutdownTimeout time.Duration) error {
	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- server.Serve(listener)
	}()

	select {
	case err := <-serveErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	case <-shutdownSignal:
		logger.Info("shutdown signal received")
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := shutdownHTTPServer(shutdownContext, server, logger); err != nil {
		return err
	}

	if err := <-serveErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP during shutdown: %w", err)
	}
	return nil
}

func shutdownHTTPServer(ctx context.Context, server *http.Server, logger *slog.Logger) error {
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	logger.Info("server stopped gracefully")
	return nil
}

func newHandler(frontendDir string) http.Handler {
	projectRepository := repository.NewInMemoryProjectRepository(defaultProjects())
	projectService := usecase.NewProjectService(projectRepository)
	apiHandler := httpdelivery.NewHandler(projectService)

	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	if frontendDir != "" {
		fileServer := http.FileServer(http.Dir(frontendDir))
		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "..") {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			w.Header().Set("Cache-Control", "no-cache, must-revalidate")
			fileServer.ServeHTTP(w, r)
		}))
	}
	return securityHeaders(mux)
}

func defaultProjects() []domain.Project {
	return []domain.Project{
		{
			Slug: "dazo-musicyt", Title: "Dazo - Musicyt", Category: "Audio · Engineering", Date: "October 2026",
			Description:      "A lightweight, responsive web-based music platform for searching and playing YouTube music through a focused audio player.",
			Overview:         "Dazo - Musicyt is a web-based music platform that brings search, playback, audio customization, lyrics, and personal music management into one responsive experience.",
			Role:             "Personal project · Backend and application development",
			Stack:            []string{"Go", "Rust", "WebAssembly", "HTML5", "CSS3", "JavaScript", "Bootstrap 5", "Web Audio API", "yt-dlp", "Docker"},
			Features:         []string{"YouTube music search and playback", "Go REST API endpoints for music search, audio stream extraction, and audio proxy handling", "Customizable audio equalizer using the Web Audio API", "Synchronized lyrics and dynamic album artwork", "Playback queue, shuffle, repeat, and playlist management", "Liked songs and playback history using LocalStorage"},
			TechnicalWork:    []string{"Implemented backend endpoints in Go to support music search and audio delivery workflows", "Integrated yt-dlp to retrieve direct audio streams", "Built interactive playback controls and music management with JavaScript", "Connected the equalizer interface to the Web Audio API", "Included Rust and WebAssembly in the audio-engine technology stack"},
			EngineeringNotes: []string{"Coordinating search, stream retrieval, and browser playback across separate application responsibilities", "Keeping playback state, queue actions, and user-managed music data consistent", "Providing responsive controls while supporting configurable audio processing"},
			Repository:       "https://github.com/dazo-rgb/dazo-musicyt", Featured: true,
		},
		{
			Slug: "jatinusa", Title: "Jatinusa - Website", Category: "E-commerce · Web", Date: "September 2026",
			Description:      "A full-stack furniture commerce platform designed for teak furniture businesses to showcase products, manage inventory, and provide product price estimates.",
			Overview:         "Jatinusa is a furniture e-commerce website concept for teak furniture businesses. It combines a responsive product catalog with customer-facing shopping features and product management capabilities.",
			Role:             "Personal project · Full-stack web development",
			Stack:            []string{"HTML5", "CSS3", "JavaScript ES6+", "Bootstrap 5", "Node.js", "Express.js", "MySQL", "LocalStorage", "Git", "GitHub"},
			Features:         []string{"Responsive furniture catalog and product browsing", "Detailed product information pages", "User authentication", "Shopping cart functionality", "Product CRUD and inventory management", "Furniture price estimation"},
			TechnicalWork:    []string{"Developed responsive catalog and product interfaces using HTML, CSS, and JavaScript", "Built backend functionality with Node.js and Express.js", "Integrated MySQL for product data and database operations", "Implemented authentication, product CRUD, cart, and price estimation features"},
			EngineeringNotes: []string{"Connecting product presentation with persistent relational product data", "Supporting both customer shopping flows and product-management operations", "Keeping catalog and estimation interactions understandable across screen sizes"},
			Repository:       "https://github.com/dazo-rgb/JATINUSA-WEB", Featured: true,
		},
	}
}

func findFrontendDir() (string, error) {
	candidates := []string{filepath.Join(".", "frontend"), filepath.Join("..", "frontend")}
	if _, sourceFile, _, ok := runtime.Caller(0); ok {
		projectDir := filepath.Dir(filepath.Dir(sourceFile))
		candidates = append(candidates, filepath.Join(projectDir, "frontend"))
	}
	seen := map[string]bool{}
	for _, candidate := range candidates {
		absolute, err := filepath.Abs(candidate)
		if err != nil || seen[absolute] {
			continue
		}
		seen[absolute] = true
		if info, err := os.Stat(filepath.Join(absolute, "index.html")); err == nil && !info.IsDir() {
			return absolute, nil
		}
	}
	return "", errors.New("frontend directory not found: expected frontend/index.html")
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("encode JSON response", "error", err, "status", status)
	}
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", contentSecurityPolicy)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

func (recorder *statusRecorder) WriteHeader(status int) {
	if recorder.wroteHeader {
		return
	}
	recorder.status = status
	recorder.wroteHeader = true
	recorder.ResponseWriter.WriteHeader(status)
}

func (recorder *statusRecorder) Write(body []byte) (int, error) {
	if !recorder.wroteHeader {
		recorder.WriteHeader(http.StatusOK)
	}
	n, err := recorder.ResponseWriter.Write(body)
	recorder.bytes += n
	return n, err
}

func (recorder *statusRecorder) Unwrap() http.ResponseWriter {
	return recorder.ResponseWriter
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.status,
			"bytes", recorder.bytes,
			"duration_ms", time.Since(started).Milliseconds(),
		)
	})
}
