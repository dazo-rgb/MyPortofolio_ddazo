package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testHandler() http.Handler { return newHandler("") }

func TestHealthEndpoint(t *testing.T) {
	recorder := httptest.NewRecorder()
	testHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status payload = %q, want ok", body["status"])
	}
}

func TestProjectsEndpoint(t *testing.T) {
	recorder := httptest.NewRecorder()
	testHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/projects", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body []Project
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 2 {
		t.Fatalf("projects = %d, want 2", len(body))
	}
	for _, project := range body {
		if project.Slug == "" || project.Repository == "" || len(project.Stack) == 0 {
			t.Errorf("incomplete project payload: %+v", project)
		}
	}
}

func TestProjectDetailEndpoint(t *testing.T) {
	recorder := httptest.NewRecorder()
	testHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/projects/dazo-musicyt", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body Project
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Slug != "dazo-musicyt" || len(body.Features) == 0 || len(body.TechnicalWork) == 0 {
		t.Fatalf("unexpected detail payload: %+v", body)
	}
}

func TestUnknownProjectReturns404(t *testing.T) {
	recorder := httptest.NewRecorder()
	testHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/projects/not-real", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	recorder := httptest.NewRecorder()
	testHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/projects", nil))
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}

func TestSecurityHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	testHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	for _, header := range []string{"X-Content-Type-Options", "Referrer-Policy", "X-Frame-Options", "Permissions-Policy"} {
		if recorder.Header().Get(header) == "" {
			t.Errorf("missing security header %s", header)
		}
	}
}

func TestStaticFrontendServed(t *testing.T) {
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "index.html"), []byte("<!doctype html><title>portfolio test</title>"), 0o600); err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	newHandler(directory).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), "portfolio test") {
		t.Fatal("static index response did not contain expected content")
	}
}

func TestAPIResponsesAreNotCachedAndUseJSON(t *testing.T) {
	for _, path := range []string{"/api/health", "/api/projects", "/api/projects/jatinusa"} {
		recorder := httptest.NewRecorder()
		testHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK {
			t.Errorf("%s status = %d, want %d", path, recorder.Code, http.StatusOK)
		}
		if got := recorder.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
			t.Errorf("%s Content-Type = %q, want application/json", path, got)
		}
		if got := recorder.Header().Get("Cache-Control"); got != "no-store" {
			t.Errorf("%s Cache-Control = %q, want no-store", path, got)
		}
	}
}

func TestProjectDetailPayloadContainsRequiredFields(t *testing.T) {
	for _, slug := range []string{"dazo-musicyt", "jatinusa"} {
		recorder := httptest.NewRecorder()
		testHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/projects/"+slug, nil))
		var project Project
		if err := json.Unmarshal(recorder.Body.Bytes(), &project); err != nil {
			t.Fatal(err)
		}
		if project.Slug != slug || project.Title == "" || project.Description == "" || project.Overview == "" {
			t.Errorf("incomplete project detail for %s: %+v", slug, project)
		}
		if len(project.Features) == 0 || len(project.TechnicalWork) == 0 {
			t.Errorf("missing project details for %s", slug)
		}
	}
}

func TestContentSecurityPolicyAndSecurityHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	csp := recorder.Header().Get("Content-Security-Policy")
	for _, directive := range []string{
		"default-src 'self'",
		"object-src 'none'",
		"frame-ancestors 'none'",
		"https://fonts.googleapis.com",
		"https://fonts.gstatic.com",
		"https://ghchart.rshah.org",
	} {
		if !strings.Contains(csp, directive) {
			t.Errorf("CSP missing directive %q: %s", directive, csp)
		}
	}
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}

func TestContentSecurityPolicyAllowsExactJSONLDHash(t *testing.T) {
	frontendDir, err := findFrontendDir()
	if err != nil {
		t.Fatal(err)
	}
	html, err := os.ReadFile(filepath.Join(frontendDir, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	const opening = `<script type="application/ld+json">`
	const closing = `</script>`
	start := strings.Index(string(html), opening)
	if start < 0 {
		t.Fatal("JSON-LD script opening tag not found")
	}
	start += len(opening)
	end := strings.Index(string(html)[start:], closing)
	if end < 0 {
		t.Fatal("JSON-LD script closing tag not found")
	}
	jsonLD := string(html)[start : start+end]
	digest := sha256.Sum256([]byte(jsonLD))
	expectedSource := "'sha256-" + base64.StdEncoding.EncodeToString(digest[:]) + "'"
	if !strings.Contains(contentSecurityPolicy, expectedSource) {
		t.Fatalf("CSP does not allow the current JSON-LD hash %s", expectedSource)
	}
}

func TestRequestLoggerWritesStructuredFields(t *testing.T) {
	var output strings.Builder
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	handler := requestLogger(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created"))
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/test-path", nil))

	var entry map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(output.String())), &entry); err != nil {
		t.Fatalf("log is not valid JSON: %v; output=%q", err, output.String())
	}
	if entry["msg"] != "http request" || entry["method"] != http.MethodPost || entry["path"] != "/test-path" {
		t.Fatalf("unexpected structured log fields: %#v", entry)
	}
	if entry["status"] != float64(http.StatusCreated) || entry["bytes"] != float64(len("created")) {
		t.Fatalf("unexpected status/bytes in log: %#v", entry)
	}
}

func TestRequestLoggerUsesImplicitOKStatus(t *testing.T) {
	var output strings.Builder
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	handler := requestLogger(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	var entry map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(output.String())), &entry); err != nil {
		t.Fatal(err)
	}
	if entry["status"] != float64(http.StatusOK) {
		t.Fatalf("status = %#v, want 200", entry["status"])
	}
}

func TestStatusRecorderWriteHeaderAndUnwrap(t *testing.T) {
	underlying := httptest.NewRecorder()
	recorder := &statusRecorder{ResponseWriter: underlying, status: http.StatusOK}
	recorder.WriteHeader(http.StatusAccepted)
	recorder.WriteHeader(http.StatusInternalServerError)
	n, err := recorder.Write([]byte("accepted"))
	if err != nil {
		t.Fatal(err)
	}
	if n != len("accepted") || recorder.status != http.StatusAccepted || recorder.bytes != len("accepted") {
		t.Fatalf("unexpected recorder state: n=%d status=%d bytes=%d", n, recorder.status, recorder.bytes)
	}
	if recorder.Unwrap() != underlying {
		t.Fatal("Unwrap did not return the underlying ResponseWriter")
	}
}

func TestWriteJSONSetsHeadersAndBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeJSON(recorder, http.StatusAccepted, map[string]string{"result": "queued"})
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusAccepted)
	}
	if !strings.Contains(recorder.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("unexpected Content-Type: %q", recorder.Header().Get("Content-Type"))
	}
	if recorder.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", recorder.Header().Get("Cache-Control"))
	}
	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil || body["result"] != "queued" {
		t.Fatalf("unexpected JSON body: %q, err=%v", recorder.Body.String(), err)
	}
}

type failingResponseWriter struct {
	header http.Header
}

func (writer *failingResponseWriter) Header() http.Header {
	return writer.header
}

func (writer *failingResponseWriter) WriteHeader(int) {}

func (writer *failingResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("simulated write failure")
}

func TestWriteJSONHandlesEncoderWriteFailure(t *testing.T) {
	writer := &failingResponseWriter{header: make(http.Header)}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
	if !strings.Contains(writer.Header().Get("Content-Type"), "application/json") {
		t.Fatal("Content-Type was not set before encoding")
	}
}

func TestFindFrontendDir(t *testing.T) {
	directory, err := findFrontendDir()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(directory, "index.html")); err != nil {
		t.Fatalf("resolved frontend directory does not contain index.html: %v", err)
	}
}

func TestRunRejectsInvalidPort(t *testing.T) {
	t.Setenv("PORT", "not-a-port")
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	if err := run(logger); err == nil || !strings.Contains(err.Error(), "invalid PORT") {
		t.Fatalf("run error = %v, want invalid PORT error", err)
	}
}

func TestServeUntilSignalGracefullyWaitsForActiveRequest(t *testing.T) {
	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(requestStarted)
		<-releaseRequest
		_, _ = w.Write([]byte("finished"))
	})}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	shutdownSignal := make(chan struct{})
	var output strings.Builder
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	done := make(chan error, 1)
	go func() {
		done <- serveUntilSignal(server, listener, shutdownSignal, logger, 2*time.Second)
	}()

	responseDone := make(chan error, 1)
	go func() {
		response, requestErr := http.Get("http://" + listener.Addr().String())
		if requestErr != nil {
			responseDone <- requestErr
			return
		}
		defer response.Body.Close()
		body, readErr := io.ReadAll(response.Body)
		if readErr == nil && string(body) != "finished" {
			readErr = fmt.Errorf("response body = %q, want finished", body)
		}
		responseDone <- readErr
	}()

	select {
	case <-requestStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}
	close(shutdownSignal)

	select {
	case err := <-done:
		t.Fatalf("server shut down before active request completed: %v", err)
	case <-time.After(60 * time.Millisecond):
	}
	close(releaseRequest)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("serveUntilSignal returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down")
	}
	if err := <-responseDone; err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "server stopped gracefully") {
		t.Fatalf("missing graceful shutdown log: %s", output.String())
	}
}

func TestServeUntilSignalReturnsServeError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	server := &http.Server{Handler: http.NotFoundHandler()}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	err = serveUntilSignal(server, listener, make(chan struct{}), logger, time.Second)
	if err == nil || !strings.Contains(err.Error(), "serve HTTP") {
		t.Fatalf("error = %v, want serve HTTP error", err)
	}
}

func TestShutdownHTTPServerHonorsDeadline(t *testing.T) {
	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(requestStarted)
		<-releaseRequest
		_, _ = w.Write([]byte("done"))
	})}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = server.Serve(listener) }()
	responseDone := make(chan struct{})
	go func() {
		defer close(responseDone)
		response, getErr := http.Get("http://" + listener.Addr().String())
		if getErr == nil {
			_ = response.Body.Close()
		}
	}()
	select {
	case <-requestStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	if err := shutdownHTTPServer(ctx, server, logger); err == nil {
		t.Fatal("expected shutdown deadline error")
	}
	close(releaseRequest)
	_ = server.Close()
	select {
	case <-responseDone:
	case <-time.After(2 * time.Second):
		t.Fatal("request goroutine did not finish")
	}
}
