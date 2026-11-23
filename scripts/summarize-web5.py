#!/usr/bin/env python3
"""Summarize measurement points from the k6 JSON stream for WebLab 5."""
import json
import statistics
import sys
from collections import Counter, defaultdict
from pathlib import Path

points = Counter()
errors = Counter()
latencies = defaultdict(list)
ports = defaultdict(Counter)

with Path(sys.argv[1]).open() as source:
    for line in source:
        point = json.loads(line)
        if point["type"] != "Point":
            continue
        data = point["data"]
        tags = data["tags"]
        if tags.get("phase") != "measure" or "ingress" not in tags:
            continue
        ingress = tags["ingress"]
        value = data["value"]
        match point["metric"]:
            case "web5_requests":
                points[ingress, tags["operation"]] += int(value)
            case "web5_errors":
                errors[ingress, tags["operation"]] += int(value)
            case "web5_read_ms":
                latencies[ingress].append(value)
            case "web5_read_backend":
                ports[ingress][tags["port"]] += int(value)

seconds = int(sys.argv[2])
print("| Вход | Чтения | Записи | Запросов/с | Ошибки | p95 чтения, мс | Чтения по портам |")
print("|---|---:|---:|---:|---:|---:|---|")
for ingress, label in [("root", "/"), ("mirror", "/mirror")]:
    reads, writes = points[ingress, "read"], points[ingress, "write"]
    total = reads + writes
    failed = errors[ingress, "read"] + errors[ingress, "write"]
    samples = latencies[ingress]
    p95 = statistics.quantiles(samples, n=100, method="inclusive")[94] if len(samples) >= 2 else None
    histogram = ", ".join(f"{port} → {count}" for port, count in sorted(ports[ingress].items()))
    error_cell = f"{failed / total:.2%} ({failed})" if total else "—"
    p95_cell = f"{p95:.2f}" if p95 is not None else "—"
    print(f"| `{label}` | {reads} | {writes} | {total / seconds:.2f} | {error_cell} | {p95_cell} | {histogram or '—'} |")
