#!/usr/bin/env python3
import argparse
import json
import os
import re
import sys
from pathlib import Path
from urllib.parse import urlparse


REQUIRED = ["platform", "keyword", "title", "markdownContent"]
PLATFORMS = {"toutiao", "baijiahao", "xiaohongshu", "zhihu"}
LOCAL_IMAGE_PREFIX = "/static/article-images/uploads/"
IMAGE_EXTENSIONS = {".jpg", ".jpeg", ".png", ".webp"}
MARKDOWN_IMAGE_RE = re.compile(r"!\[[^\]]*\]\(([^)\s]+)(?:\s+['\"][^'\"]*['\"])?\)")


def count_cjk(text: str) -> int:
    return len(re.findall(r"[\u4e00-\u9fff]", text or ""))


def validate_json_string(payload, key):
    value = payload.get(key)
    if value in (None, ""):
        return
    if not isinstance(value, str):
        raise ValueError(f"{key} must be a JSON-encoded string")
    json.loads(value)


def extract_markdown_image_urls(markdown: str) -> list[str]:
    urls = []
    for match in MARKDOWN_IMAGE_RE.findall(markdown or ""):
        url = match.strip().strip("<>").strip()
        if url:
            urls.append(url)
    return urls


def image_path(url: str) -> str:
    parsed = urlparse(url)
    if parsed.scheme in ("http", "https"):
        return parsed.path
    if parsed.scheme:
        raise ValueError(f"unsupported image URL scheme: {url}")
    return parsed.path or url


def candidate_static_roots(explicit_roots: list[str] | None) -> list[Path]:
    roots = []
    env_root = os.environ.get("AUTO_ARTICLE_STATIC_ROOT")
    if env_root:
        roots.append(Path(env_root))
    for root in explicit_roots or []:
        roots.append(Path(root))

    cwd = Path.cwd()
    for base in [cwd, *cwd.parents]:
        roots.append(base / "backend" / "static" / "article-images" / "uploads")
        roots.append(base / "static" / "article-images" / "uploads")

    unique = []
    seen = set()
    for root in roots:
        try:
            resolved = root.resolve()
        except OSError:
            resolved = root.absolute()
        key = str(resolved).lower()
        if key not in seen:
            seen.add(key)
            unique.append(resolved)
    return unique


def resolve_local_image(public_url: str, static_roots: list[Path]) -> Path | None:
    path = image_path(public_url)
    if not path.startswith(LOCAL_IMAGE_PREFIX):
        raise ValueError(
            f"image must be a downloaded local static URL under {LOCAL_IMAGE_PREFIX}: {public_url}"
        )
    ext = Path(path).suffix.lower()
    if ext not in IMAGE_EXTENSIONS:
        raise ValueError(f"image must be a raster jpg/png/webp file: {public_url}")

    rel = path[len(LOCAL_IMAGE_PREFIX):].lstrip("/")
    for root in static_roots:
        candidate = root / rel
        if candidate.is_file():
            return candidate
    return None


def validate_images(payload, static_roots):
    cover = str(payload.get("coverImageUrl", "")).strip()
    if not cover:
        raise ValueError("missing coverImageUrl: skill saves must include a real downloaded cover image")

    markdown_urls = extract_markdown_image_urls(payload.get("markdownContent", ""))
    if not markdown_urls:
        raise ValueError("markdownContent must include at least one Markdown image from downloaded local files")

    for url in [cover, *markdown_urls]:
        local_file = resolve_local_image(url, static_roots)
        if local_file is None:
            raise ValueError(
                "image URL points to no local file; use download_article_images.py, "
                f"download_image_candidates.py, or upload-image before inserting it: {url}"
            )


def load_payload(payload_arg: str):
    if payload_arg == "-":
        raw = sys.stdin.buffer.read()
        if not raw.strip():
            raise ValueError("stdin payload is empty")
        return json.loads(raw.decode("utf-8"))

    with open(payload_arg, "r", encoding="utf-8") as f:
        return json.load(f)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("payload", help="Path to a JSON payload file, or '-' to read JSON from stdin")
    parser.add_argument(
        "--static-root",
        action="append",
        help="Static uploads root that maps to /static/article-images/uploads",
    )
    args = parser.parse_args()

    payload = load_payload(args.payload)

    for key in REQUIRED:
        if not str(payload.get(key, "")).strip():
            raise ValueError(f"missing required field: {key}")

    if payload["platform"] not in PLATFORMS:
        raise ValueError(f"unsupported platform: {payload['platform']}")

    for key in ["titleOptions"]:
        validate_json_string(payload, key)

    validate_images(payload, candidate_static_roots(args.static_root))

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
