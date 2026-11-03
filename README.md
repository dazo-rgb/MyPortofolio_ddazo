# ddazo. — Warm Dark Premium Portfolio

Portfolio website for Daffa Hari Mulyono, focused on backend development with Go, REST APIs, and database-driven applications.

## Technology
- Frontend: semantic HTML5, CSS3, vanilla JavaScript
- Backend: Go standard library (`net/http`)
- API: health check, featured project collection, and project case-study detail
- Tests: Go `httptest` unit tests, race detector, and coverage reporting
- Runtime: graceful SIGINT/SIGTERM shutdown and JSON structured request logs
- Security: Content Security Policy plus standard browser security headers
- Assets: locally stored SVG artwork, profile photo, and CV PDF
- No frontend framework and no third-party Go dependencies

## Project structure
```text
dazo-warm-dark-portfolio/
├── backend/
│   ├── go.mod
│   ├── main.go
│   └── main_test.go
├── frontend/
│   ├── index.html
│   ├── project.html
│   ├── robots.txt
│   ├── sitemap.xml
│   └── assets/
│       ├── css/style.css
│       ├── js/app.js
│       ├── js/project.js
│       ├── images/
│       └── docs/Daffa-Hari-Mulyono-CV.pdf
├── run-local.sh
└── README.md
```

## Run locally
Requires Go 1.22 or newer.

From the project root:
```bash
./run-local.sh
```

Or:
```bash
cd backend
go run .
```

Open `http://localhost:8080`. Do not open the HTML file directly or use Live Server for the full-stack preview; the Go server serves both the frontend and API from one origin.

## API endpoints
| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/health` | Service health status |
| GET | `/api/projects` | Featured project collection |
| GET | `/api/projects/{slug}` | Full case-study data for one project |

Example:
```bash
curl http://localhost:8080/api/health
curl http://localhost:8080/api/projects
curl http://localhost:8080/api/projects/dazo-musicyt
```

The homepage fetches `/api/projects` and renders project cards from the returned JSON. Its built-in HTML cards remain as a graceful fallback if the API is unavailable. Each card links to `project.html?slug=...`, which fetches the matching case study from the Go detail endpoint.

## Tests and quality checks
Run from `backend/`:
```bash
go test ./...
go test -race ./...
go vet ./...
```

The Go test suite covers the health endpoint, project list and detail payloads, unknown projects, unsupported methods, API response headers, CSP and security headers, static frontend serving, JSON structured request logging, and graceful shutdown behavior. Current automated statement coverage target is at least 75%; generate a fresh report with `go test -cover ./...`.

Check frontend JavaScript syntax:
```bash
node --check ../frontend/assets/js/app.js
node --check ../frontend/assets/js/project.js
```

The frontend includes an 8-second API request timeout with a built-in project-card fallback, accessible loading/error status messages, and client-side contact validation. The contact form opens WhatsApp or the visitor's email application; it does not send or store messages on the server. The WhatsApp handoff checks whether a new tab was blocked and clears the new tab's opener reference after it opens. Complete a final visual check on real phones/tablets and keyboard-only navigation before production deployment.

## Reliability, accessibility, and SEO
- Semantic page landmarks, skip link, descriptive image alternatives, form labels, keyboard-accessible navigation, and visible focus styles.
- Mobile menu supports Escape and closes after navigation or outside click.
- Reduced-motion preference is respected.
- Contact form validates required fields and provides a status message; it opens WhatsApp or the visitor's email application and does not store messages.
- Security response headers include a restrictive Content Security Policy, `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, and `Permissions-Policy`. The CSP explicitly allows the Google Fonts stylesheet/font hosts and the GitHub contribution chart image host; the homepage JSON-LD script is authorized by an exact SHA-256 hash.
- Open Graph metadata, page descriptions, canonical URL, favicon, `robots.txt`, and `sitemap.xml` are included.
- Hero, project, contact, favicon, and social-preview artwork are local SVG files. Google Fonts remain an optional external enhancement; system font fallbacks are defined.
- A GitHub contribution chart is embedded from an external chart service; visitors can also open the original GitHub profile directly.
- The Go server logs startup, HTTP method/path/status/response bytes/duration, shutdown signals, and shutdown failures as JSON using `log/slog`. It attempts a 10-second graceful shutdown on SIGINT/SIGTERM.
- Project metadata intentionally remains in a small in-memory Go slice: the portfolio currently exposes two read-only case studies and has no content-management workflow. A database/CMS migration is deferred until project editing, authentication, or a larger catalog is required.

## Personal links
- Email: daffaharimulyono@gmail.com
- WhatsApp: +62 897-7400-705
- GitHub: https://github.com/dazo-rgb
- LinkedIn: https://www.linkedin.com/in/daffa-hari-mulyono-8982313a3/
- Instagram: https://www.instagram.com/ddaffa_harii/

## Content accuracy
- Education: Universitas Islam 45 Bekasi, 2020–2024, GPA 3.43/4.00.
- Work experience dates follow the supplied CV.
- The supplied CV lists the Hacktiv8 course as September 2026–January 2026, a reversed range. The public timeline therefore shows `2026 · TRAINING` without inventing a month range. Confirm the correct dates before adding them.
- Project descriptions and technology lists are based on the supplied CV and existing portfolio source. No user counts, performance metrics, or business outcomes have been invented.
- Confirm all contact links, repository visibility, project dates, and CV details before deployment.

## Deployment notes
The current application serves frontend and API from one Go server, so same-origin requests need no CORS configuration. For Vercel or another static-only host, deploy the Go backend as a serverless/API service or use a Go-capable host; static hosting alone will not run `main.go`. If frontend and API are hosted on separate origins, configure an explicit API base URL and a restricted CORS policy rather than using a wildcard in production.

## Backend architecture

The Go backend is organized using Clean Architecture boundaries:

- `backend/internal/domain`: pure project entity.
- `backend/internal/usecase`: application operations and the `ProjectRepository` output port.
- `backend/internal/repository`: in-memory infrastructure adapter implementing the application port, with defensive deep copies.
- `backend/internal/delivery/http`: HTTP adapter depending only on use cases.
- `backend/main.go`: composition root, static seed data, dependency wiring, server lifecycle, middleware, security headers, and structured request logging.

Dependencies point inward: delivery depends on use cases; use cases depend only on domain and define the repository port; the in-memory adapter implements that port; `main.go` wires the concrete adapter. Slugs are validated before repository access. Project data remains in memory for the current read-only scope.
