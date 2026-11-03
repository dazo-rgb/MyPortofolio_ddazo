(() => {
  "use strict";

  const menuButton = document.querySelector(".menu-toggle");
  const nav = document.querySelector("#primary-nav");
  const navLinks = [...document.querySelectorAll(".nav-link")];

  const closeMenu = () => {
    if (!menuButton || !nav) return;
    nav.classList.remove("is-open");
    menuButton.setAttribute("aria-expanded", "false");
    menuButton.setAttribute("aria-label", "Open navigation");
  };

  if (menuButton && nav) {
    menuButton.addEventListener("click", () => {
      const open = menuButton.getAttribute("aria-expanded") !== "true";
      menuButton.setAttribute("aria-expanded", String(open));
      menuButton.setAttribute("aria-label", open ? "Close navigation" : "Open navigation");
      nav.classList.toggle("is-open", open);
    });
    navLinks.forEach((link) => link.addEventListener("click", closeMenu));
    document.addEventListener("keydown", (event) => {
      if (event.key === "Escape") closeMenu();
    });
    document.addEventListener("click", (event) => {
      if (nav.classList.contains("is-open") && !nav.contains(event.target) && !menuButton.contains(event.target)) closeMenu();
    });
  }

  // Route every navigation item to its matching section, including Home -> page top.
  navLinks.forEach((link) => {
    link.addEventListener("click", (event) => {
      const hash = link.getAttribute("href");
      if (!hash || !hash.startsWith("#")) return;

      const targetId = hash.slice(1);
      const target = document.getElementById(targetId);
      if (!target) return;

      event.preventDefault();
      closeMenu();

      const isHome = targetId === "top";
      const header = document.querySelector(".site-header");
      const headerOffset = header ? header.getBoundingClientRect().height : 0;
      const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
      const behavior = reduceMotion ? "auto" : "smooth";

      window.history.pushState(null, "", hash);

      // Wait for the mobile menu to finish closing before calculating the final scroll.
      window.requestAnimationFrame(() => {
        if (isHome) {
          // Home always means the absolute top of the document. Do not assign
          // scrollTop afterwards: that would cancel the native smooth animation.
          window.scrollTo({ top: 0, left: 0, behavior });
          return;
        }

        const finalTop = Math.max(
          0,
          window.scrollY + target.getBoundingClientRect().top - headerOffset - 12
        );
        window.scrollTo({ top: finalTop, left: 0, behavior });
      });
    });
  });

  const year = document.querySelector("#current-year");
  if (year) year.textContent = String(new Date().getFullYear());

  const revealTargets = document.querySelectorAll(
    ".intro-heading, .intro-copy, .skills-preview, .project-card, .expertise-item, .journey-entry, .contact-inner"
  );
  if ("IntersectionObserver" in window && !window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
    revealTargets.forEach((element) => element.classList.add("reveal"));
    const revealObserver = new IntersectionObserver((entries, observer) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) return;
        entry.target.classList.add("is-visible");
        observer.unobserve(entry.target);
      });
    }, { threshold: 0.12 });
    revealTargets.forEach((element) => revealObserver.observe(element));
  }

  const projectGrid = document.querySelector("#project-grid");
  const apiStatus = document.querySelector("#api-status");

  const createElement = (tag, className, content) => {
    const element = document.createElement(tag);
    if (className) element.className = className;
    if (content !== undefined) element.textContent = content;
    return element;
  };

  const createProjectCard = (project, index) => {
    const article = createElement("article", "project-card");
    const image = createElement("div", `project-image ${project.slug === "jatinusa" ? "project-image-b" : "project-image-a"}`);
    image.setAttribute("role", "img");
    image.setAttribute("aria-label", `${project.title} project artwork`);
    image.append(createElement("span", "", String(index + 1).padStart(2, "0")));

    const info = createElement("div", "project-info");
    info.append(createElement("p", "project-category", project.category));
    info.append(createElement("h3", "", project.title));
    info.append(createElement("p", "", project.description));

    const tags = createElement("div", "tag-list");
    project.stack.slice(0, 5).forEach((technology) => tags.append(createElement("span", "", technology)));
    info.append(tags);

    const actions = createElement("div", "project-actions");
    const detail = createElement("a", "project-link", "View case study ↗");
    detail.href = `/project.html?slug=${encodeURIComponent(project.slug)}`;
    const repository = createElement("a", "project-link", "GitHub ↗");
    repository.href = project.repository;
    repository.target = "_blank";
    repository.rel = "noopener noreferrer";
    actions.append(detail, repository);
    info.append(actions);
    article.append(image, info);
    return article;
  };

  if (projectGrid) {
    if (apiStatus) apiStatus.textContent = "Loading projects…";
    const controller = new AbortController();
    const requestTimeout = window.setTimeout(() => controller.abort(), 8000);

    fetch("/api/projects", {
      headers: { Accept: "application/json" },
      signal: controller.signal
    })
      .then((response) => {
        if (!response.ok) throw new Error(`Projects API returned ${response.status}`);
        return response.json();
      })
      .then((projects) => {
        if (!Array.isArray(projects) || projects.length === 0) {
          throw new Error("Projects API returned no project records");
        }
        const validProjects = projects.filter((project) =>
          project && typeof project.slug === "string" && /^[a-z0-9-]+$/.test(project.slug) &&
          typeof project.title === "string" && project.title.trim().length > 0 &&
          typeof project.description === "string" && project.description.trim().length > 0 &&
          typeof project.category === "string" && project.category.trim().length > 0 &&
          Array.isArray(project.stack) && project.stack.length > 0 && project.stack.every((item) => typeof item === "string" && item.trim().length > 0) &&
          typeof project.repository === "string" &&
          (() => {
            try {
              const repositoryUrl = new URL(project.repository);
              return repositoryUrl.protocol === "https:" && repositoryUrl.hostname === "github.com" && repositoryUrl.pathname.split("/").filter(Boolean).length >= 2;
            } catch {
              return false;
            }
          })()
        );
        if (validProjects.length === 0) throw new Error("Projects API returned invalid project data");

        const fragment = document.createDocumentFragment();
        validProjects.forEach((project, index) => fragment.append(createProjectCard(project, index)));
        projectGrid.replaceChildren(fragment);
        if (apiStatus) {
          apiStatus.textContent = `Live portfolio data · ${validProjects.length} projects loaded from the Go API.`;
        }
      })
      .catch((error) => {
        const reason = error.name === "AbortError"
          ? "The project service took too long to respond."
          : error.message;
        if (apiStatus) {
          apiStatus.textContent = "Showing the built-in project preview. Live project data is temporarily unavailable.";
        }
        console.warn("Project API fallback:", reason);
      })
      .finally(() => window.clearTimeout(requestTimeout));
  }

  const setHomeAsActive = () => {
    if (window.scrollY > 80) return;
    navLinks.forEach((link) => {
      const active = link.getAttribute("href") === "#top";
      link.classList.toggle("active", active);
      if (active) link.setAttribute("aria-current", "location");
      else link.removeAttribute("aria-current");
    });
  };
  window.addEventListener("scroll", setHomeAsActive, { passive: true });
  setHomeAsActive();

  const observedSections = document.querySelectorAll("main section[id]");
  if ("IntersectionObserver" in window && observedSections.length) {
    const sectionObserver = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) return;
        navLinks.forEach((link) => {
          const active = link.getAttribute("href") === `#${entry.target.id}`;
          link.classList.toggle("active", active);
          if (active) link.setAttribute("aria-current", "location");
          else link.removeAttribute("aria-current");
        });
      });
    }, { rootMargin: "-30% 0px -60% 0px", threshold: 0 });
    observedSections.forEach((section) => sectionObserver.observe(section));
  }

  const contactForm = document.querySelector("#contact-form");
  if (contactForm) {
    const nameField = contactForm.elements.namedItem("name");
    const messageField = contactForm.elements.namedItem("message");
    const feedback = document.querySelector("#contact-feedback");

    [nameField, messageField].forEach((field) => {
      field.addEventListener("input", () => {
        field.setCustomValidity("");
        field.removeAttribute("aria-invalid");
      });
    });

    contactForm.addEventListener("submit", (event) => {
      event.preventDefault();
      const name = String(nameField.value || "").trim();
      const message = String(messageField.value || "").trim();

      nameField.setCustomValidity(name.length < 2 ? "Please enter a name with at least 2 characters." : "");
      messageField.setCustomValidity(message.length < 10 ? "Please enter a message with at least 10 characters." : "");
      [nameField, messageField].forEach((field) => {
        if (!field.checkValidity()) field.setAttribute("aria-invalid", "true");
      });

      if (!contactForm.reportValidity()) {
        if (feedback) feedback.textContent = "Please review the highlighted fields before continuing.";
        return;
      }

      const channel = String(new FormData(contactForm).get("channel") || "whatsapp");
      const body = `Hi Daffa, I'm ${name}.\n\n${message}`;
      if (channel === "whatsapp") {
        const url = `https://wa.me/628977400705?text=${encodeURIComponent(body)}`;
        // Do not pass `noopener` in the window.open feature string: browsers may
        // return null even when the tab opened, which makes blocked-tab detection unreliable.
        const opened = window.open(url, "_blank");
        if (opened) opened.opener = null;
        if (feedback) {
          feedback.textContent = opened
            ? "WhatsApp opened in a new tab with your message prepared. Please review and send it there."
            : "The browser blocked the WhatsApp tab. Allow pop-ups or select Email instead.";
        }
      } else {
        const url = `mailto:daffaharimulyono@gmail.com?subject=${encodeURIComponent(`Portfolio inquiry from ${name}`)}&body=${encodeURIComponent(body)}`;
        window.location.href = url;
        if (feedback) feedback.textContent = "Your email application should open with the message prepared. Please review and send it there.";
      }
    });
  }
})();
