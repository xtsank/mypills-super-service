#!/usr/bin/env python3
"""Merge the saved CI runs from before the Allure history split."""

import json
import os
from collections import Counter
from datetime import datetime
from pathlib import Path


TREND_FILES = (
    "history-trend.json",
    "duration-trend.json",
    "retry-trend.json",
    "categories-trend.json",
)


def read_json(path: Path):
    return json.loads(path.read_text())


def timestamp_ms(value: str) -> int:
    return int(datetime.fromisoformat(value.replace("Z", "+00:00")).timestamp() * 1000)


def main() -> None:
    current_dir = Path(os.environ.get("CI_HISTORY_DIR", "allure-history-v2"))
    legacy_root = Path(os.environ["LEGACY_ROOT"])
    runs = sorted(json.loads(os.environ["LEGACY_RUNS"]), key=lambda run: run["number"])
    if [run["number"] for run in runs] != [4, 5, 6]:
        raise ValueError("Expected saved CI runs 4, 5 and 6")

    old_points = {name: [] for name in TREND_FILES}
    old_items = {}
    for run in runs:
        source = legacy_root / str(run["number"])
        for name in TREND_FILES:
            trend = read_json(source / name)
            if not trend:
                raise ValueError(f"CI run {run['number']} has no {name}")
            old_points[name].append(trend[0])  # Allure stores the current run first.

        start = timestamp_ms(run["start"]) - 60_000
        end = timestamp_ms(run["end"]) + 60_000
        count = 0
        for history_id, test in read_json(source / "history.json").items():
            for item in test["items"]:
                if start <= item["time"]["start"] <= end:
                    old_items.setdefault(history_id, []).append(item)
                    count += 1
        expected = old_points["history-trend.json"][-1]["data"]["total"]
        if count != expected:
            raise ValueError(
                f"CI run {run['number']} has {count} historical tests, expected {expected}"
            )

    merged_trends = {}
    for name in TREND_FILES:
        current = read_json(current_dir / name) if (current_dir / name).exists() else []
        merged_trends[name] = (current + list(reversed(old_points[name])))[:20]

    history_path = current_dir / "history.json"
    merged_history = read_json(history_path) if history_path.exists() else {}
    for history_id, items in old_items.items():
        existing = merged_history.setdefault(history_id, {"statistic": {}, "items": []})
        by_uid = {item["uid"]: item for item in existing["items"]}
        by_uid.update({item["uid"]: item for item in items})
        existing["items"] = sorted(
            by_uid.values(), key=lambda item: item["time"]["start"], reverse=True
        )[:20]
        counts = Counter(item["status"] for item in existing["items"])
        existing["statistic"] = {
            **{status: counts[status] for status in ("failed", "broken", "skipped", "passed", "unknown")},
            "total": len(existing["items"]),
        }

    current_dir.mkdir(parents=True, exist_ok=True)
    for name, data in merged_trends.items():
        (current_dir / name).write_text(json.dumps(data))
    history_path.write_text(json.dumps(merged_history))
    (current_dir / "ci-history-restored.txt").write_text("CI runs 4, 5, 6\n")
    print("Restored Allure history from CI runs 4, 5 and 6")


if __name__ == "__main__":
    main()
