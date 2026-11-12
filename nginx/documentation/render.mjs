import { readFileSync } from "node:fs";
import MarkdownIt from "markdown-it";

const markdown = new MarkdownIt({ html: false });
const normalizeLink = markdown.normalizeLink.bind(markdown);
markdown.normalizeLink = (link) => {
  const normalized = normalizeLink(link);
  if (/^(?:[/#]|[a-z][a-z\d+.-]*:)/i.test(normalized)) {
    return normalized;
  }
  const url = new URL(normalized, "http://localhost/docs/web/");
  return url.pathname + url.search + url.hash;
};

const readme = readFileSync(process.argv[2], "utf8");
console.log(`<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Документация MyPills</title>
  <style>
    :root {
      font: 1rem/1.65 Inter, system-ui, sans-serif;
      color: #182b26;
      background: #f4f8f5;
    }
    body { margin: 0; }
    main {
      box-sizing: border-box;
      max-width: 64rem;
      margin: 2rem auto;
      padding: clamp(1.25rem, 5vw, 3rem);
      background: white;
      border: 1px solid #d8e4dc;
      border-radius: .75rem;
      overflow-wrap: anywhere;
    }
    nav { margin-bottom: 1.25rem; }
    a { color: #176047; }
    h1, h2, h3 { line-height: 1.25; }
    h1 { margin-top: 0; font-size: clamp(2rem, 5vw, 3rem); }
    h2 {
      margin-top: 2.25rem;
      padding-bottom: .4rem;
      border-bottom: 1px solid #d8e4dc;
    }
    img { display: block; max-width: 100%; height: auto; }
    table, pre { display: block; max-width: 100%; overflow-x: auto; }
    table { border-collapse: collapse; }
    th, td { padding: .4rem .6rem; border: 1px solid #d8e4dc; }
    pre { padding: 1rem; border-radius: .5rem; background: #f4f8f5; }
    @media (max-width: 64rem) {
      main { margin: .75rem; }
    }
  </style>
</head>
<body>
  <main>
    <nav aria-label="Навигация"><a href="/">На главную</a></nav>
    <article>
${markdown.render(readme)}
    </article>
  </main>
</body>
</html>`);
