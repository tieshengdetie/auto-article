#!/usr/bin/env python3
import argparse
import json
import re
import sys


REQUIRED = ["platform", "keyword", "title", "markdownContent"]
PLATFORMS = {"toutiao", "baijiahao", "xiaohongshu", "zhihu"}


def count_cjk(text: str) -> int:
    return len(re.findall(r"[\u4e00-\u9fff]", text or ""))


def validate_json_string(payload, key):
    value = payload.get(key)
    if value in (None, ""):
        return
    if not isinstance(value, str):
        raise ValueError(f"{key} must be a JSON-encoded string")
    json.loads(value)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("payload", help="Path to a JSON payload file")
    args = parser.parse_args()

    with open(args.payload, "r", encoding="utf-8") as f:
        payload = json.load(f)

    for key in REQUIRED:
        if not str(payload.get(key, "")).strip():
            raise ValueError(f"missing required field: {key}")

    if payload["platform"] not in PLATFORMS:
        raise ValueError(f"unsupported platform: {payload['platform']}")

    for key in ["keywordSegments", "titleOptions", "tags", "images", "sources", "hotTopics", "styleProfile", "publishPayload"]:
        validate_json_string(payload, key)

    actual_count = count_cjk(payload.get("markdownContent", ""))
    declared_count = int(payload.get("wordCount") or 0)
    if declared_count and abs(actual_count - declared_count) > max(120, actual_count * 0.2):
        raise ValueError(f"wordCount looks inaccurate: declared={declared_count}, actual_cjk={actual_count}")

    print(json.dumps({"ok": True, "actualCjkCount": actual_count}, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(json.dumps({"ok": False, "error": str(exc)}, ensure_ascii=False), file=sys.stderr)
        raise SystemExit(1)
