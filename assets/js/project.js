(() => {
  "use strict";
  const root = document.querySelector("#case-content");
  const year = document.querySelector("#current-year");
  if (year) year.textContent = String(new Date().getFullYear());
  if (!root) return;

  const element = (tag, className, text) => {
    const node = document.createElement(tag);
    if (className) node.className = className;
    if (text !== undefined) node.textContent = text;
    return node;
  };
  const slug = new URLSearchParams(window.location.search).get("slug");
  if (!slug || !/^[a-z0-9-]+$/.test(slug)) {
    const error = element("div", "case-error");
    error.append(element("h1", "", "Project not found"), element("p", "", "Choose a project from the portfolio homepage."));
    const back = element("a", "button button-dark", "Back to projects");
    back.href = "/#work";
    error.append(back);
    root.replaceChildren(error);
    return;
  }

  const addListSection = (container, title, items) => {
    const section = element("section");
    section.append(element("h2", "", title));
    const list = element("ul", "case-list");
    items.forEach((item) => list.append(element("li", "", item)));
    section.append(list);
    container.append(section);
  };

  fetch("/data/projects.json", {
    headers: { Accept: "application/json" }
  })
    .then((response) => {
      if (!response.ok) throw new Error("Project details are temporarily unavailable.");
      return response.json();
    })
    .then((projects) => {
      const project = projects.find((item) => item.slug === slug);
      if (!project) throw new Error("This project could not be found.");
      const isNonEmptyString = (value) => typeof value === "string" && value.trim().length > 0;
      const isStringList = (value) => Array.isArray(value) && value.every(isNonEmptyString);
      let repositoryIsValid = false;
      try {
        const repositoryURL = new URL(project.repository);
        repositoryIsValid = repositoryURL.protocol === "https:" && repositoryURL.hostname === "github.com" && repositoryURL.pathname.split("/").filter(Boolean).length >= 2;
      } catch {
        repositoryIsValid = false;
      }
      if (!project || !isNonEmptyString(project.title) || !isNonEmptyString(project.category) ||
          !isNonEmptyString(project.date) || !isNonEmptyString(project.description) ||
          !isNonEmptyString(project.overview) || !isNonEmptyString(project.role) ||
          !isStringList(project.stack) || project.stack.length === 0 ||
          !isStringList(project.features) || !isStringList(project.technicalWork) ||
          !isStringList(project.engineeringNotes) || !repositoryIsValid) {
        throw new Error("Project details are incomplete or contain invalid data.");
      }
      document.title = `${project.title} — ddazo.`;
      const pageDescription = `${project.description.slice(0, 150).trim()}${project.description.length > 150 ? "…" : ""}`;
      const descriptionMeta = document.querySelector('meta[name="description"]');
      if (descriptionMeta) descriptionMeta.setAttribute("content", pageDescription);
      const canonicalURL = `${window.location.origin}/project.html?slug=${encodeURIComponent(slug)}`;
      const canonicalMeta = document.querySelector('link[rel="canonical"]');
      if (canonicalMeta) canonicalMeta.setAttribute("href", canonicalURL);
      const openGraphTitle = document.querySelector('meta[property="og:title"]');
      if (openGraphTitle) openGraphTitle.setAttribute("content", `${project.title} — ddazo.`);
      const openGraphDescription = document.querySelector('meta[property="og:description"]');
      if (openGraphDescription) openGraphDescription.setAttribute("content", pageDescription);
      const openGraphURL = document.querySelector('meta[property="og:url"]');
      if (openGraphURL) openGraphURL.setAttribute("content", canonicalURL);
      const hero = element("header", "case-hero");
      hero.append(element("p", "eyebrow", `CASE STUDY · ${project.category}`));
      hero.append(element("h1", "", project.title));
      hero.append(element("p", "case-lead", project.overview));
      const meta = element("div", "case-meta");
      meta.append(element("span", "", project.date), element("span", "", project.role));
      hero.append(meta);

      const grid = element("div", "case-grid");
      const content = element("div", "case-content");
      const overview = element("section");
      overview.append(element("h2", "", "Project overview"), element("p", "", project.description));
      content.append(overview);
      addListSection(content, "Key features", project.features || []);
      addListSection(content, "Technical implementation", project.technicalWork || []);
      addListSection(content, "Engineering considerations", project.engineeringNotes || []);

      const sidebar = element("aside", "case-sidebar");
      const stackSection = element("section");
      stackSection.append(element("h2", "", "Technology stack"));
      const stack = element("div", "case-stack");
      (project.stack || []).forEach((technology) => stack.append(element("span", "", technology)));
      stackSection.append(stack);
      sidebar.append(stackSection);
      const links = element("div", "case-actions");
      const repository = element("a", "button button-dark", "View GitHub ↗");
      repository.href = project.repository;
      repository.target = "_blank";
      repository.rel = "noopener noreferrer";
      links.append(repository);
      sidebar.append(links);

      grid.append(content, sidebar);
      root.replaceChildren(hero, grid);
    })
    .catch((error) => {
      const message = element("p", "", error.message);
      const back = element("a", "button button-dark", "Back to projects");
      back.href = "/#work";
      root.replaceChildren(element("div", "case-error"));
      root.firstElementChild.append(element("h1", "", "Unable to load project"), message, back);
    });
})();
