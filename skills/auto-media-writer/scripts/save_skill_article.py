#!/usr/bin/env python3
import argparse
import json
import os
import subprocess
import sys
import urllib.error
import urllib.request
from pathlib import Path


DEFAULT_BASE_URL = "http://localhost:9001"
ENDPOINT = "/api/v1/skill-articles"


def run_validator(payload_arg: str, static_roots: list[str] | None, body: bytes | None = None) -> None:
    validator = Path(__file__).with_name("validate_skill_article_payload.py")
    command = [sys.executable, str(validator), payload_arg]
    for root in static_roots or []:
        command.extend(["--static-root", root])
    result = subprocess.run(
        command,
        input=body,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    stdout = result.stdout.decode("utf-8", errors="replace")
    stderr = result.stderr.decode("utf-8", errors="replace")
    if result.returncode != 0:
        if stderr:
            print(stderr.strip(), file=sys.stderr)
        if stdout:
            print(stdout.strip(), file=sys.stderr)
        raise SystemExit(result.returncode)
    if stdout:
        print(stdout.strip())


def load_payload(payload_arg: str) -> bytes:
    if payload_arg == "-":
        raw = sys.stdin.buffer.read()
        if not raw.strip():
            raise ValueError("stdin payload is empty")
        payload = json.loads(raw.decode("utf-8"))
    else:
        with Path(payload_arg).open("r", encoding="utf-8") as f:
            payload = json.load(f)
    return json.dumps(payload, ensure_ascii=False).encode("utf-8")


def normalize_url(base_url: str) -> str:
    return base_url.rstrip("/") + ENDPOINT


def post_payload(url: str, body: bytes, timeout: int) -> int:
    request = urllib.request.Request(
        url,
        data=body,
        method="POST",
        headers={"Content-Type": "application/json; charset=utf-8"},
    )
    try:
        with urllib.request.urlopen(request, timeout=timeout) as response:
            response_body = response.read().decode("utf-8", errors="replace")
            print(response_body or json.dumps({"ok": True, "status": response.status}))
            return 0
    except urllib.error.HTTPError as exc:
        response_body = exc.read().decode("utf-8", errors="replace")
        print(
            json.dumps(
                {"ok": False, "status": exc.code, "error": response_body},
                ensure_ascii=False,
            ),
            file=sys.stderr,
        )
        return 1
    except urllib.error.URLError as exc:
        print(
            json.dumps({"ok": False, "error": str(exc.reason)}, ensure_ascii=False),
            file=sys.stderr,
        )
        return 1


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Validate and save an auto-media-writer article payload through the HTTP API."
    )
    parser.add_argument("payload", help="Path to the JSON payload file, or '-' to read JSON from stdin")
    parser.add_argument(
        "--base-url",
        default=os.environ.get("AUTO_ARTICLE_BASE_URL", DEFAULT_BASE_URL),
        help="Backend base URL. Prefer AUTO_ARTICLE_BASE_URL or an explicit deployed URL; falls back to local development http://localhost:9001",
    )
    parser.add_argument(
        "--timeout",
        type=int,
        default=30,
        help="HTTP timeout in seconds",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Validate and print the target URL without sending the request",
    )
    parser.add_argument(
        "--static-root",
        action="append",
        help="Static uploads root that maps to /static/article-images/uploads",
    )
    args = parser.parse_args()

    body = load_payload(args.payload)
    run_validator(args.payload, args.static_root, body if args.payload == "-" else None)
    url = normalize_url(args.base_url)

    if args.dry_run:
        print(json.dumps({"ok": True, "dryRun": True, "url": url}, ensure_ascii=False))
        return 0

    return post_payload(url, body, args.timeout)


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(json.dumps({"ok": False, "error": str(exc)}, ensure_ascii=False), file=sys.stderr)
        raise SystemExit(1)
