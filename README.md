# ddazo — Golang Backend Developer Portfolio

Responsive portfolio website for Daffa Hari Mulyono, built with a Go backend organized using Clean Architecture and a lightweight HTML/CSS/JavaScript frontend.

## Requirements
- Go 1.22+
- No third-party Go dependencies

## Run locally
```bash
go run ./cmd/server
```
Open http://localhost:8080

## Project structure
```text
ddazo-portfolio/
├── cmd/server/main.go
├── internal/
│   ├── domain/                 # Enterprise entities and repository contracts
│   ├── usecase/                # Application business rules
│   ├── repository/memory/      # Infrastructure implementation
│   └── delivery/http/          # HTTP handlers and routing
├── web/
│   ├── index.html
│   ├── assets/css/styles.css
│   ├── assets/js/main.js
│   └── assets/images/          # Add your own profile/project images here
├── go.mod
└── README.md
```

## Clean Architecture notes
- `domain` contains core entities and repository interfaces.
- `usecase` contains application rules and depends only on domain contracts.
- `repository/memory` implements the repository contract.
- `delivery/http` adapts HTTP requests to use cases.
- `cmd/server` is the composition root and wires dependencies.

The contact endpoint currently validates and stores messages in memory for demonstration. Configure a persistent repository and email provider before production use.

## Personalization checklist
- Add your CV as `web/assets/files/daffa-hari-mulyono-cv.pdf`.
- Add your circular profile image as `web/assets/images/profile.jpg`.
- Replace project preview placeholders with real screenshots.
- Verify project descriptions, repository URLs, and deployment URLs before publishing.
- The displayed skill chips are technology categories, not proficiency percentages.
