# ddazo Portfolio — Quality Assurance Report

## Perbaikan lanjutan pada build ini
- Memperbaiki deteksi tab WhatsApp: sebelumnya `window.open()` memakai fitur `noopener,noreferrer`, yang dapat membuat browser mengembalikan `null` walaupun tab sebenarnya berhasil dibuka. Sekarang tab dibuka secara sinkron saat submit, lalu referensi `opener` diputus setelah berhasil.
- Memperjelas pesan ketika browser memblokir tab WhatsApp.
- Memperketat validasi data proyek dari API: kategori dan stack harus tersedia, dan URL repository harus HTTPS pada host `github.com` sebelum dibuat menjadi tautan.
- Menambahkan `rel="noopener noreferrer"` pada tautan eksternal yang memakai tab baru.
- Memperbaiki status navigasi aktif agar kembali ke Home saat pengguna scroll ke bagian paling atas.
- Memperketat validasi proyek di halaman utama dan halaman case study, termasuk format slug, daftar teknologi, konten wajib, dan URL repository GitHub.
- Memperbarui metadata Open Graph homepage menjadi URL absolut dan menambahkan metadata case study yang diperbarui sesuai proyek yang dimuat.
- Memperbarui dokumentasi agar sesuai dengan penggunaan chart kontribusi GitHub eksternal.

## Automated checks completed
- `go test ./...` — PASS
- `go test -race ./...` — PASS
- `go vet ./...` — PASS
- `node --check frontend/assets/js/app.js` — PASS
- `node --check frontend/assets/js/project.js` — PASS
- `sh -n run-local.sh` — PASS
- HTML structural checks — PASS (no duplicate IDs, missing image alt attributes, unlabeled form controls, or unsafe target-blank links)
- Local API smoke check — `/api/health` returned HTTP 200 with JSON status `ok` during the local run.
- Homepage and case study HTML routes — HTTP 200 during the local run.
- Static assets (CSS, JavaScript, CV PDF, profile image, `robots.txt`, and `sitemap.xml`) — HTTP 200 during the local run.
- Project API — returned two valid featured project records; unknown slug returned HTTP 404.
- Invalid trailing health route `/api/health/` returned HTTP 404 as expected.

## Not verified in this environment
- Physical Android phone and tablet rendering, including portrait/landscape.
- Real interaction with installed WhatsApp and Email applications.
- Lighthouse scores for Performance, Accessibility, Best Practices, and SEO. Lighthouse is not installed in this environment.
- Firefox and Safari rendering. Chromium is installed, but automated headless page capture did not complete within the environment's execution limit; therefore no browser screenshot or visual pass is claimed.
- External social destinations were not individually opened and confirmed in a real browser session.

## Important interpretation
The source-level and Go test suite pass, but these checks do not prove that every interaction works on a physical device. The contact form prepares a message in WhatsApp or the visitor's email application; it does not transmit or store the message. Review the prefilled message before sending.

## Final manual QA required
1. Open the deployed site on Android phone and tablet in portrait and landscape.
2. Test all navbar items, project details, CV download, GitHub, LinkedIn, Instagram, and Email.
3. Submit a valid test message using both WhatsApp and Email; verify the recipient and prefilled text before sending.
4. Run Chrome DevTools Lighthouse for all four categories in both Mobile and Desktop modes and save the reports.
5. Compare the same pages in Chrome, Firefox, and Safari.

Do not treat the previous provisional score of 88/100 as a final verified score until the manual checks and Lighthouse reports are available.

## Layout refinement — 2026-11-03

Applied a focused responsive-layout pass while preserving the existing Warm Dark Premium visual identity:

- Compact laptop hero (901–1100px): reduced heading scale and text width, tightened supporting-copy spacing, and repositioned the background visual to improve text/visual balance.
- Project cards: standardized artwork cropping and aspect ratios across desktop, tablet, and narrow-phone layouts; cards use a consistent image frame rather than content-dependent image heights.
- Journey timeline: refined date-column width, inter-column spacing, vertical rhythm, and mobile stacking.
- Contact section: reduced decorative-art height on tablet and narrow phones so the form remains the primary focus.
- 360–390px phones: added targeted adjustments for hero typography, gutters, buttons, specialty labels, section headings, project padding, contact title, and social links.

Validation performed after this pass:
- `go test ./...` — PASS
- `go test -race ./...` — PASS
- `go vet ./...` — PASS
- `node --check frontend/assets/js/app.js` — PASS
- `node --check frontend/assets/js/project.js` — PASS
- `sh -n run-local.sh` — PASS
- CSS parsed with tinycss2 — PASS

Physical-device screenshots, Lighthouse scoring, and Firefox/Safari visual comparison remain unverified and are not claimed as passed.

- Contact section aligned to the same centered max-width container (`--content`) as the other main sections; desktop columns now shrink safely without forcing a 280px minimum image column.


## Backend hardening — 2026-11-03

- Added a restrictive Content Security Policy: `default-src 'self'`, no plugins/objects, no framing, same-origin API connections, explicit Google Fonts stylesheet/font hosts, and the GitHub contribution chart image host. The homepage JSON-LD is allowed through a SHA-256 script hash rather than a blanket inline-script exception.
- Removed the font stylesheet's inline `onload` handler so the policy does not need `unsafe-inline` for script attributes. The Google Fonts stylesheet now loads as a normal stylesheet; local font fallbacks remain available if the external font service is unavailable.
- Added SIGINT/SIGTERM handling and a 10-second graceful HTTP shutdown. Active requests are allowed to complete within the shutdown deadline.
- Replaced plain request logs with JSON structured logging using Go `log/slog`; records include method, path, status, response bytes, and duration.
- Expanded tests for CSP directives and JSON-LD hash integrity, structured logs, response recorder behavior, JSON response handling, frontend directory discovery, invalid port configuration, graceful shutdown, and shutdown deadlines.
- Project metadata remains in-memory by design. A database or CMS is not justified for two read-only portfolio case studies; the README records the current boundary and the trigger for a future migration.

Validation:
- `go test ./...` — PASS
- `go test -race ./...` — PASS
- `go vet ./...` — PASS
- `go test -cover ./...` — PASS (main package 75.0%; HTTP delivery 90.0%; repository 100.0%; use case 100.0%; domain contains data declarations only).
- `gofmt` — PASS
- Frontend JavaScript syntax — PASS
- CSS parsing — PASS
- Real binary smoke check — health and project endpoints returned HTTP 200; CSP was present without `unsafe-inline`; SIGTERM exited with code 0 after logging `server stopped gracefully`.

Still requires real-browser Lighthouse reports and physical/cross-browser testing; those scores are not fabricated or claimed.

## Clean Architecture refactor

- Separated project entities/data into `internal/domain`.
- Added a repository contract and in-memory adapter in `internal/repository`.
- Added project application services in `internal/usecase`.
- Moved API routing and response delivery into `internal/delivery/http`.
- Kept `main.go` as the composition root and server lifecycle layer.
- Updated the experience entry to `SD Jaya Suti Abadi · Tambun Selatan`.

Validation after refactor: `go test -count=1 ./...`, `go test -race -count=1 ./...`, and `go vet ./...`.


## Clean Architecture and data-isolation correction — 2026-11-03

- Moved the repository interface into the use-case layer, so application rules no longer import infrastructure.
- Kept static project data in the composition root, leaving the domain package focused on the entity.
- Added defensive deep copies for project slices when data enters or leaves the in-memory repository.
- Added use-case slug validation and HTTP 400 responses for invalid slugs; unknown valid slugs remain HTTP 404.
- Added tests for dependency inversion, invalid slug rejection, error propagation, and repository data isolation.
- HTTP delivery now depends only on use-case abstractions and services.


## Final verification after architecture correction — 2026-11-03

### Backend test results
- `go test -count=1 ./...` — PASS
- `go test -race -count=1 ./...` — PASS
- `go vet ./...` — PASS
- `go test -count=1 -cover ./...` — PASS
  - Main package: 75.3%
  - HTTP delivery: 89.5%
  - Repository: 100.0%
  - Use case: 100.0%
- `gofmt` — PASS

### Frontend and runtime checks
- JavaScript syntax checks for `app.js` and `project.js` — PASS.
- Shell syntax check for `run-local.sh` — PASS.
- HTML checks for duplicate IDs, missing image alt text, and unsafe new-tab links — PASS.
- Local asset references in homepage and case-study HTML — PASS.
- CSS parsing — PASS.
- Real Go binary smoke test: health 200, project list 200 (2 records), project detail 200, unknown project 404, invalid slug 400, homepage 200 — PASS.
- Static delivery: homepage, case study, CSS, JavaScript, CV PDF, profile image, robots.txt, and sitemap.xml — HTTP 200.
- Updated experience location `SD Jaya Suti Abadi · Tambun Selatan` — verified in served homepage.
- Use-case package imports no concrete repository package — verified.

### Remaining verification limits
- Lighthouse scores were not generated.
- No real-browser screenshots or physical-device tests were completed.
- Firefox and Safari visual checks remain outstanding.
- WhatsApp and email handoff still require a manual test in installed applications.

Scores should therefore distinguish source-level quality from visual/browser verification; no Lighthouse or device-test pass is implied.
