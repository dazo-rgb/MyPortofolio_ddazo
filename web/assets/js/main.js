(() => {
  const year = document.querySelector("#year");
  if (year) year.textContent = new Date().getFullYear();

  const links = [...document.querySelectorAll(".nav a")];
  const sections = links
    .map((link) => document.querySelector(link.getAttribute("href")))
    .filter(Boolean);

  const updateActiveLink = () => {
    let current = sections[0];
    for (const section of sections) {
      if (section.getBoundingClientRect().top <= 140) current = section;
    }
    links.forEach((link) => {
      link.classList.toggle("active", link.getAttribute("href") === `#${current.id}`);
    });
  };

  window.addEventListener("scroll", updateActiveLink, { passive: true });
  updateActiveLink();
})();
