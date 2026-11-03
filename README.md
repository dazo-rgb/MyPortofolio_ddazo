# ddazo. — Vercel-ready portfolio

Static portfolio for Daffa Hari Mulyono. It is prepared for direct deployment to Vercel with **no build command, no framework, no environment variables, and no backend server**.

## Deploy to Vercel

1. Extract this ZIP on your computer.
2. Open Vercel and choose **Add New → Project**.
3. Import the GitHub repository containing these files, or upload/deploy the extracted project using Vercel CLI.
4. In Project Settings → Build and Development Settings, use:
   - Framework Preset: **Other**
   - Build Command: **leave empty**
   - Output Directory: **`.`** (the project root)
   - Install Command: **leave empty**
5. Deploy.

The `index.html`, `project.html`, `assets/`, `data/`, and `vercel.json` files must be at the project root (not inside an extra nested folder).

## Run locally

From this folder, run:

```bash
python3 -m http.server 8000
```

Then open http://localhost:8000.

## Project structure

```text
.
├── index.html
├── project.html
├── vercel.json
├── robots.txt
├── sitemap.xml
├── data/
│   └── projects.json
└── assets/
    ├── css/
    ├── docs/
    ├── images/
    └── js/
```

## What changed for static deployment

- Removed the Go backend and its tests because Vercel static hosting does not run a persistent Go HTTP server.
- Moved the existing read-only project case-study data into `data/projects.json`.
- Updated the homepage and case-study page to load project data from that local JSON file.
- Kept the existing visual design, assets, CV, contact links, navigation, and project case studies.
- No database or server environment variables are required. The contact form prepares a WhatsApp or email message; it does not store or send messages itself.

## Before publishing

Check that the CV is current and that the personal links, contact details, project dates, and GitHub repositories are correct.
