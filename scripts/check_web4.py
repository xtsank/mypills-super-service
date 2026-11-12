#!/usr/bin/env python3
"""Read-only HTTP checks for the running WebLab #4 gateway."""

import argparse
import json
import sys
from html.parser import HTMLParser
from urllib.error import HTTPError, URLError
from urllib.parse import urljoin
from urllib.request import HTTPRedirectHandler, ProxyHandler, build_opener


class CheckError(Exception):
    pass


def require(condition, message):
    if not condition:
        raise CheckError(message)


class NoRedirect(HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


class Page(HTMLParser):
    def __init__(self, body):
        super().__init__()
        self.links = set()
        self.images = set()
        self.styles = set()
        self.scripts = set()
        self.feed(body.decode("utf-8"))

    def handle_starttag(self, tag, attributes):
        attrs = dict(attributes)
        if tag == "a" and attrs.get("href"):
            self.links.add(attrs["href"])
        if tag == "img" and attrs.get("src"):
            self.images.add(attrs["src"])
        if tag == "link" and attrs.get("rel") == "stylesheet" and attrs.get("href"):
            self.styles.add(attrs["href"])
        if tag == "script" and attrs.get("src"):
            self.scripts.add(attrs["src"])


class Checks:
    def __init__(self, base):
        self.base = base.rstrip("/")
        self.opener = build_opener(ProxyHandler({}), NoRedirect)
        self.requests = 0

    def get(self, path, status=200, *, media=None, contains=None):
        require(path.startswith("/") and not path.startswith("//"), f"unexpected URL: {path}")
        try:
            response = self.opener.open(self.base + path, timeout=5)
        except HTTPError as error:
            response = error
        except (URLError, OSError) as error:
            raise CheckError(f"GET {path}: {error}") from error

        with response:
            require(response.code == status, f"GET {path}: expected HTTP {status}, got {response.code}")
            headers = response.headers
            if media is not None:
                allowed = (media,) if isinstance(media, str) else media
                require(headers.get_content_type() in allowed,
                        f"GET {path}: unexpected Content-Type {headers.get_content_type()}")
            body = response.read()
            if contains is not None:
                require(contains in body, f"GET {path}: expected content is missing")
        self.requests += 1
        print(f"  {status} GET {path}")
        return headers, body

    def redirect(self, path, target, status=308):
        headers, _ = self.get(path, status)
        require(headers.get("Location") == target,
                f"GET {path}: expected Location {target}, got {headers.get('Location')}")

    def json(self, path, status=200):
        _, body = self.get(path, status, media="application/json")
        try:
            return json.loads(body)
        except ValueError as error:
            raise CheckError(f"GET {path}: invalid JSON") from error

    def assets(self, route, prefix, page):
        require(page.styles and page.scripts, f"GET {route}: CSS or JavaScript links are missing")
        for link in sorted(page.styles):
            path = urljoin(route, link)
            require(path.startswith(prefix), f"GET {route}: unexpected stylesheet URL {path}")
            self.get(path, media="text/css")
        for link in sorted(page.scripts):
            path = urljoin(route, link)
            require(path.startswith(prefix), f"GET {route}: unexpected script URL {path}")
            self.get(path, media=("application/javascript", "text/javascript"))


def check_status(http):
    print("Health and status")
    for path, expected in (("/health", b"ok\n"), ("/ready", b"ready\n")):
        headers, body = http.get(path, media="text/plain")
        require(body == expected, f"GET {path}: expected {expected.strip().decode()}")
        require(headers.get("Cache-Control") == "no-store", f"GET {path}: missing Cache-Control: no-store")
    _, body = http.get("/status", media="text/plain", contains=b"Active connections:")
    require(b"requests" in body and b"Waiting:" in body, "GET /status: Nginx counters are missing")


def check_pages(http):
    print("Home, management and React")
    _, home = http.get("/", media="text/html", contains=b"MyPills")
    _, reserved = http.get("/reserved/", media="text/html")
    require(home == reserved, "GET /reserved/: response differs from /")
    sketches = {f"/static/img/img{index}.png" for index in range(1, 4)}
    require(sketches <= Page(home).images, "GET /: one or more Lab 1 sketches are missing")
    for path in sorted(sketches):
        http.get(path, media="image/png", contains=b"\x89PNG\r\n\x1a\n")
    http.get("/static/img/web4-missing.png", 404)

    _, management = http.get("/management", media="text/html")
    expected = {
        "/", "/reserved/", "/legacy/", "/documentation", "/api/v2/docs/", "/api/v1/docs/",
        "/admin/", "/status", "/health", "/ready",
    }
    missing = expected - Page(management).links
    require(not missing, f"GET /management: missing links {', '.join(sorted(missing))}")

    http.redirect("/legacy?from=web4", "/legacy/?from=web4")
    _, legacy = http.get("/legacy/", media="text/html", contains=b'id="root"')
    http.assets("/legacy/", "/legacy/assets/", Page(legacy))
    http.get("/legacy/assets/web4-missing.js", 404)


def check_documentation(http):
    print("Project documentation")
    _, documentation = http.get(
        "/documentation", media="text/html",
        contains="Краткий перечень функциональных требований".encode(),
    )
    page = Page(documentation)
    require("/docs/openapi/v2.yaml" in page.links and "/docs/bpmn/bpmn.svg" in page.images,
            "GET /documentation: API contract or BPMN diagram is missing")
    for path in sorted(page.links | page.images):
        if path.startswith("/docs/"):
            media = None
            if path in page.images:
                media = "image/svg+xml" if path.endswith(".svg") else "image/png"
            http.get(path, media=media)
    http.get("/docs/", 404)
    http.get("/docs/lr-testing/e2e-http.pcapng", 404)


def check_api(http):
    print("API v1 and v2")
    units = http.json("/api/v1/dictionary/units")
    require(isinstance(units, list), "GET /api/v1/dictionary/units: expected a JSON list")
    links = http.json("/api/v2/reference")
    require(isinstance(links, dict) and links.get("units") == "/api/v2/reference/units",
            "GET /api/v2/reference: wrong units link")
    page = http.json("/api/v2/reference/units?limit=1&offset=1")
    require(isinstance(page, dict) and page.get("limit") == 1 and page.get("offset") == 1,
            "GET /api/v2/reference/units: query parameters did not reach the API")
    for path in ("/api/v2/users/me", "/api/v1/admin/dosage-rules"):
        error = http.json(path, 401)
        require(isinstance(error, dict) and error.get("code") == "UNAUTHORIZED",
                f"GET {path}: expected the API authorization error")
    http.get("/api/v2/web4-nonexistent", 404, media="text/plain", contains=b"404 page not found")


def check_swagger(http):
    print("Swagger and OpenAPI")
    for version in ("v1", "v2"):
        route = f"/api/{version}/docs"
        http.redirect(route, route + "/")
        _, swagger = http.get(route + "/", media="text/html", contains=b"Swagger UI")
        http.assets(route + "/", route + "/", Page(swagger))

    spec = http.json("/api/v1/docs/doc.json")
    require(isinstance(spec, dict) and spec.get("basePath") == "/api/v1",
            "Swagger v1: basePath must point to /api/v1")
    http.get("/api/v2/docs/swagger-initializer.js", contains=b"/api/spec/v2.yaml")
    _, openapi = http.get("/api/spec/v2.yaml", contains=b"openapi: 3.1.")
    require(b"/api/v2/reference:" in openapi, "OpenAPI v2: reference route is missing")
    http.redirect("/swagger", "/swagger/", 302)
    http.get("/swagger/", media="text/html", contains=b"Swagger UI")
    http.redirect("/swagger-v1", "/swagger-v1/index.html", 302)
    http.redirect("/swagger-v1/", "/swagger-v1/index.html", 302)
    http.get("/swagger-v1/index.html", media="text/html", contains=b"Swagger UI")


def check_admin(http):
    print("Adminer access")
    http.redirect("/admin", "/admin/")
    for path in ("/admin/", "/admin/?file=default.css"):
        headers, _ = http.get(path, 401)
        require(headers.get("WWW-Authenticate", "").startswith("Basic "),
                f"GET {path}: expected an HTTP Basic Auth challenge")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base", default="http://127.0.0.1:8088", help="gateway URL")
    args = parser.parse_args()
    http = Checks(args.base)
    print(f"WebLab #4: {http.base}", flush=True)
    try:
        check_status(http)
        check_pages(http)
        check_documentation(http)
        check_api(http)
        check_swagger(http)
        check_admin(http)
    except (CheckError, OSError, ValueError) as error:
        print(f"FAIL: {error}", file=sys.stderr)
        return 1
    print(f"PASS: {http.requests} HTTP checks")
    return 0


if __name__ == "__main__":
    sys.exit(main())
