#!/usr/bin/env python3
"""Smoke and contract checks against a disposable MyPills v2 test server."""

import argparse
import json
import time
import urllib.error
import urllib.request

FORM_ID = "44444444-4444-4444-8444-444444444444"
UNIT_ID = "55555555-5555-4555-8555-555555555555"
SUBSTANCE_ID = "33333333-3333-4333-8333-333333333333"
ILLNESS_ID = "66666666-6666-4666-8666-666666666666"


def call(base, method, path, payload=None, token=None, expected=200, required=()):
    data = None if payload is None else json.dumps(payload).encode()
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = "Bearer " + token
    request = urllib.request.Request(base + path, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(request, timeout=15) as response:
            status = response.status
            raw = response.read()
            location = response.headers.get("Location")
    except urllib.error.HTTPError as error:
        status = error.code
        raw = error.read()
        location = error.headers.get("Location")
    result = json.loads(raw) if raw else None
    assert status == expected, f"{method} {path}: expected {expected}, got {status}: {result}"
    for key in required:
        assert isinstance(result, dict) and key in result, f"{method} {path}: missing {key}"
    if expected == 201:
        assert location, f"{method} {path}: missing Location"
    if expected == 204:
        assert result is None, f"{method} {path}: 204 must have no body"
    print(f"{status} {method} {path}")
    return result


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--base", default="http://127.0.0.1:18080")
    parser.add_argument("--admin-login", required=True)
    parser.add_argument("--admin-password", required=True)
    args = parser.parse_args()
    base = args.base.rstrip("/")
    unique = str(int(time.time() * 1000))
    admin = call(base, "POST", "/api/v2/sessions", {"login": args.admin_login, "password": args.admin_password}, required=("token", "user"))
    admin_token = admin["token"]
    user = call(base, "POST", "/api/v2/users", {"login": "user" + unique, "email": f"user{unique}@example.test", "password": "testpass123", "age": 30, "weight": 70}, expected=201, required=("token", "user"))
    other = call(base, "POST", "/api/v2/users", {"login": "other" + unique, "email": f"other{unique}@example.test", "password": "testpass123", "age": 25, "weight": 60}, expected=201, required=("token", "user"))
    token, other_token = user["token"], other["token"]
    call(base, "GET", "/api/v2/users/me", token=token, required=("id", "login", "email", "allergies", "illnesses"))
    call(base, "PATCH", "/api/v2/users/me", {"weight": 72}, token=token, required=("weight",))
    call(base, "PATCH", "/api/v2/users/me", {"weight": None}, token=token, expected=400)
    call(base, "GET", "/api/v2/reference", required=("illnesses", "substances", "forms", "units"))
    for kind in ("illnesses", "substances", "forms", "units"):
        call(base, "GET", f"/api/v2/reference/{kind}?limit=20", required=("items", "total", "limit", "offset"))
    call(base, "GET", "/api/v2/medicines?limit=20", token=token, required=("items", "total", "limit", "offset"))
    call(base, "POST", "/api/v2/medicines", {}, token=token, expected=403)
    med = call(base, "POST", "/api/v2/medicines", {
        "name": "Проверка API " + unique, "expire_time": 24, "is_prescription": False,
        "method_of_application": "Внутрь", "effect_on_pregnant": False, "effect_on_driver": False,
        "form_id": FORM_ID, "unit_id": UNIT_ID, "indications": [ILLNESS_ID],
        "contraindications": [], "substances": [{"id": SUBSTANCE_ID, "concentration": 500}],
        "dosages": [{"type": "age", "value_from": 18, "value_to": 120,
                     "dosage_value": 1, "number_of_doses_per_day": 2}],
    }, token=admin_token, expected=201, required=("id", "indications", "dosages"))
    mid = med["id"]
    med = call(base, "PATCH", f"/api/v2/medicines/{mid}", {"expire_time": 36}, token=admin_token, required=("expire_time", "dosages"))
    assert med["expire_time"] == 36 and med["dosages"][0]["id"], "PATCH lost unchanged dosage"
    call(base, "GET", f"/api/v2/medicines?id={mid}", token=token, required=("items", "total"))
    cabinet = call(base, "POST", "/api/v2/cabinets", {"name": "Домашняя"}, token=token, expected=201, required=("id", "name"))
    cid = cabinet["id"]
    call(base, "POST", "/api/v2/cabinets", {"name": "Домашняя"}, token=token, expected=409)
    call(base, "GET", "/api/v2/cabinets?limit=20", token=token, required=("items", "total"))
    call(base, "PATCH", f"/api/v2/cabinets/{cid}", {"name": "В поездку"}, token=token, required=("id", "name"))
    call(base, "PATCH", f"/api/v2/cabinets/{cid}", {"name": "Чужая"}, token=other_token, expected=404)
    item = call(base, "POST", f"/api/v2/cabinets/{cid}/items", {
        "medicine_id": mid, "date_of_manufacture": "2026-01-01", "quantity": 5, "minimum_quantity": 6,
    }, token=token, expected=201, required=("id", "cabinet_id", "medicine_name", "low_stock"))
    iid = item["id"]
    call(base, "POST", f"/api/v2/cabinets/{cid}/items", {
        "medicine_id": mid, "date_of_manufacture": "2026-01-01", "quantity": 2,
        "minimum_quantity": 1,
    }, token=token, expected=409)
    call(base, "GET", f"/api/v2/cabinets/{cid}/items?low_stock=true", token=token, required=("items", "total"))
    call(base, "GET", f"/api/v2/cabinets/{cid}/items", token=other_token, expected=404)
    call(base, "PATCH", f"/api/v2/cabinets/{cid}/items/{iid}", {"quantity": 8}, token=token, required=("quantity",))
    call(base, "DELETE", f"/api/v2/medicines/{mid}", token=admin_token, expected=409)
    call(base, "GET", f"/api/v2/recommendations?illness_id={ILLNESS_ID}", token=token, required=("items",))
    call(base, "GET", f"/api/v2/recommendations?illness_id={ILLNESS_ID}&cabinet_id={cid}", token=other_token, expected=404)
    restock = call(base, "POST", "/api/v2/restock-list", {"medicine_id": mid, "desired_quantity": 10}, token=token, expected=201, required=("id", "medicine_id"))
    rid = restock["id"]
    call(base, "GET", "/api/v2/restock-list", token=token, required=("items", "total"))
    call(base, "PATCH", f"/api/v2/restock-list/{rid}", {"desired_quantity": 12}, token=token, required=("desired_quantity",))
    call(base, "PATCH", f"/api/v2/restock-list/{rid}", {"desired_quantity": 3}, token=other_token, expected=404)
    call(base, "DELETE", f"/api/v2/restock-list/{rid}", token=token, expected=204)
    call(base, "DELETE", f"/api/v2/cabinets/{cid}", token=token, expected=409)
    call(base, "DELETE", f"/api/v2/cabinets/{cid}/items/{iid}", token=token, expected=204)
    call(base, "DELETE", f"/api/v2/cabinets/{cid}", token=token, expected=204)
    call(base, "DELETE", f"/api/v2/medicines/{mid}", token=admin_token, expected=204)
    call(base, "GET", "/api/v1/dictionary/illnesses", token=token)
    print("v2 contract smoke passed")


if __name__ == "__main__":
    main()
