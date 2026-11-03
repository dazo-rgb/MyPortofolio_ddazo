package httpdelivery

import (
	"dazo-portfolio/internal/domain"
	"dazo-portfolio/internal/repository"
	"dazo-portfolio/internal/usecase"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandler() http.Handler {
	repo := repository.NewInMemoryProjectRepository([]domain.Project{{Slug: "demo", Title: "Demo", Featured: true}})
	return NewHandler(usecase.NewProjectService(repo))
}
func TestProjectEndpoints(t *testing.T) {
	h := testHandler()
	cases := []struct {
		name, path string
		want       int
	}{{"health", "/api/health", 200}, {"list", "/api/projects", 200}, {"detail", "/api/projects/demo", 200}, {"unknown", "/api/projects/missing", 404}, {"invalid slug", "/api/projects/%2E%2E%2Fdemo", 400}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if rec.Code != tc.want {
				t.Fatalf("status=%d want=%d body=%s", rec.Code, tc.want, rec.Body.String())
			}
			if tc.want == 200 {
				var v any
				if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
func TestMethodNotAllowed(t *testing.T) {
	rec := httptest.NewRecorder()
	testHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/projects", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d", rec.Code)
	}
}
