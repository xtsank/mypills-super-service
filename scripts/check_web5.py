#!/usr/bin/env python3
"""Short, readable demonstration of WebLab 5 on a running local Compose stack."""
import argparse
from collections import Counter
import http.client
import io
import json
import subprocess
import sys
import time
import uuid
from urllib.error import HTTPError
from urllib.request import HTTPCookieProcessor, Request, build_opener
from http.cookiejar import CookieJar


def require(condition, message):
    if not condition:
        raise RuntimeError(message)


class Demo:
    def __init__(self, base, project):
        self.base = base.rstrip("/")
        self.compose = ["docker", "compose"]
        if project:
            self.compose += ["-p", project]
        self.compose += ["-f", "docker-compose.yaml", "-f", "docker-compose.lab5.yaml"]
        self.plain = build_opener()
        self.cookies = CookieJar()
        self.browser = build_opener(HTTPCookieProcessor(self.cookies))
        self.user_id = None

    def docker(self, *args, input=None):
        result = subprocess.run(self.compose + list(args), input=input, text=True, capture_output=True)
        require(result.returncode == 0, result.stderr.strip() or "Ошибка Docker")
        return result.stdout.strip()

    def http(self, path, method="GET", payload=None, token=None, browser=False, gzip=False):
        headers = {}
        if payload is not None:
            headers["Content-Type"] = "application/json"
        if token:
            headers["Authorization"] = "Bearer " + token
        if gzip:
            headers["Accept-Encoding"] = "gzip"
        request = Request(self.base + path, method=method, headers=headers,
                          data=json.dumps(payload).encode() if payload is not None else None)
        opener = self.browser if browser else self.plain
        try:
            response = opener.open(request, timeout=8)
        except HTTPError as exc:
            response = exc
        with response:
            return response.status, response.headers, response.read()

    @staticmethod
    def port(headers):
        address = headers.get("X-Upstream-Addr", "")
        require(":" in address, "В ответе отсутствует X-Upstream-Addr")
        return int(address.rsplit(":", 1)[1])

    def balance(self, prefix, ports):
        counts = Counter()
        for _ in range(40):
            status, headers, _ = self.http(prefix + "/api/v2/reference")
            require(status == 200, f"GET {prefix}/api/v2/reference → HTTP {status}")
            counts[self.port(headers)] += 1
        require(set(counts) == set(ports), f"Неверные узлы балансировки: {dict(counts)}")
        expected = [20, 10, 10]
        require(all(abs(counts[port] - count) <= 2 for port, count in zip(ports, expected)),
                f"Ожидалось 20/10/10, получено {dict(counts)}")
        output = ", ".join(f"{port} → {counts[port]}" for port in ports)
        print(f"  {prefix or '/'}: 40 GET, {output} — OK", flush=True)

    def reject_write(self, service):
        body = '{"name":"demo"}'
        wire = (
            'POST /api/v2/cabinets HTTP/1.1\r\nHost: localhost\r\n'
            f'Content-Type: application/json\r\nContent-Length: {len(body.encode())}\r\n'
            f'Connection: close\r\n\r\n{body}'
        )
        result = subprocess.run(
            self.compose + ["exec", "-T", service, "sh", "-ec",
                            '{ cat; sleep 1; } | nc -w 5 127.0.0.1 "$SERVER_PORT"'],
            input=wire.encode(), capture_output=True,
        )
        require(result.returncode == 0 and result.stdout, f"{service}: нет HTTP-ответа")

        class Socket:
            def makefile(self, *args, **kwargs):
                return io.BytesIO(result.stdout)

        with http.client.HTTPResponse(Socket()) as response:
            response.begin()
            body = response.read()
            require(response.status == 503 and json.loads(body).get("code") == "READ_ONLY_INSTANCE",
                    f"{service}: HTTP {response.status}, {body.decode(errors='replace')}")
            print(f"  POST → {service}: HTTP {response.status}, {body.decode()}", flush=True)

    def cookie_demo(self):
        unique = uuid.uuid4().hex
        status, _, body = self.http("/api/v2/users", "POST", {
            "login": "web5demo_" + unique, "email": unique + "@example.test",
            "password": uuid.uuid4().hex, "age": 30, "weight": 70, "allergies": [], "illnesses": [],
        })
        require(status == 201, f"Создание тестового пользователя → HTTP {status}")
        auth = json.loads(body)
        self.user_id = str(uuid.UUID(auth["user"]["id"]))
        token = auth["token"]
        status, headers, body = self.http("/mirror/api/v2/cabinets", "POST",
                                         {"name": "WebLab5 demo"}, token, browser=True)
        require(status == 201 and self.port(headers) == 8084, "Создание аптечки не прошло через 8084")
        cabinet_id = str(uuid.UUID(json.loads(body)["id"]))
        cookie = next((item for item in self.cookies if item.name == "mypills_mirror_primary"), None)
        require(cookie is not None and cookie.expires is not None, "Нет cookie закрепления на primary")
        print(f"  POST /mirror/api/v2/cabinets → HTTP {status}, порт 8084")
        print(f"  Set-Cookie: {headers.get('Set-Cookie')}")
        root_path = "/api/v2/cabinets/" + cabinet_id + "/items"
        path = "/mirror" + root_path
        # Establish samples on every reader before the measured demo requests.
        # Prometheus needs two scrapes to calculate a rate.
        for prefix in ("", "/mirror"):
            for _ in range(8):
                for attempt in range(30):
                    status, _, body = self.http(prefix + root_path, token=token)
                    if status == 200:
                        break
                    time.sleep(0.1)
                require(status == 200, "Содержимое новой аптечки пока недоступно")
        time.sleep(6)
        counts = Counter()
        for _ in range(8):
            status, headers, body = self.http(path, token=token, browser=True)
            require(status == 200 and json.loads(body)["total"] == 0, "Созданная аптечка не видна")
            counts[self.port(headers)] += 1
        require(counts == {8084: 8}, f"С cookie ожидается primary: {dict(counts)}")
        print("  Сразу после записи: 8 GET → 8084 (primary), аптечка видна — OK")
        seconds = max(0, cookie.expires - time.time() + 1)
        print(f"  Ждём истечения cookie ({seconds:.0f} с)…", flush=True)
        time.sleep(seconds)
        counts = Counter()
        for _ in range(40):
            status, headers, body = self.http(path, token=token, browser=True)
            require(status == 200 and json.loads(body)["total"] == 0, "Аптечка не видна после истечения cookie")
            counts[self.port(headers)] += 1
        require(set(counts) == {8084, 8085, 8086}, f"Читатели реплики не вернулись: {dict(counts)}")
        require(abs(counts[8084] - 20) <= 2 and all(abs(counts[p] - 10) <= 2 for p in (8085, 8086)),
                f"После cookie ожидалось 20/10/10: {dict(counts)}")
        output = ", ".join(f"{port} → {counts[port]}" for port in (8084, 8085, 8086))
        print(f"  После истечения: 40 GET, {output}, оба читателя реплики видят аптечку — OK")
        for _ in range(40):
            status, _, body = self.http(root_path, token=token)
            require(status == 200 and json.loads(body)["total"] == 0, "Чтение содержимого через обычный вход не прошло")
        print("  Прочитано содержимое аптечки через оба входа — данные для графиков Grafana отправлены.")

    def run(self):
        print("1. Прямой POST на читающие backend", flush=True)
        for service in ("app-read-1", "app-read-2", "mirror-read-1", "mirror-read-2"):
            self.reject_write(service)
        print("\n2. Балансировка 2:1:1 (без cookie)", flush=True)
        writer = int(self.docker("exec", "-T", "app", "sh", "-ec", 'printf "%s" "$SERVER_PORT"'))
        self.balance("", [writer, 8082, 8083])
        self.balance("/mirror", [8084, 8085, 8086])
        print("\n3. Gzip", flush=True)
        for path in ("/legacy/", "/mirror/legacy/"):
            status, headers, body = self.http(path, gzip=True)
            require(status == 200 and headers.get("Content-Encoding") == "gzip",
                    f"{path}: нет gzip")
            import gzip
            unpacked = gzip.decompress(body)
            require(b"<html" in unpacked, f"{path}: некорректное тело gzip")
            print(f"  GET {path}: HTTP 200, Content-Encoding: gzip, "
                  f"{len(unpacked)} → {len(body)} байт — OK")
        print("\n4. Mirror: чтение после записи и возврат к реплике", flush=True)
        self.cookie_demo()

    def cleanup(self):
        if self.user_id:
            self.docker(
                "exec", "-T", "db", "sh", "-ec",
                'exec psql -X -qAt -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"',
                input=f"BEGIN; DELETE FROM cabinets WHERE user_id='{self.user_id}'; "
                      f"DELETE FROM users WHERE id='{self.user_id}'; COMMIT;",
            )
            print("\nТестовые пользователь и аптечка удалены.", flush=True)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base", default="http://127.0.0.1:8088")
    parser.add_argument("--project", default="")
    args = parser.parse_args()
    demo = Demo(args.base, args.project)
    try:
        try:
            demo.run()
        finally:
            demo.cleanup()
    except (RuntimeError, OSError, ValueError, KeyError, http.client.HTTPException) as exc:
        print("FAIL:", exc, file=sys.stderr)
        return 1
    print("PASS — демо WebLab#5 пройдено.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
