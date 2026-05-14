#!/usr/bin/env python3
import argparse
import hashlib
import json
import re
import sys
from datetime import datetime
from html import unescape
from pathlib import Path
from urllib.parse import urljoin, urlparse
from urllib.request import Request, urlopen


IMAGE_EXTENSIONS = {".jpg", ".jpeg", ".png", ".webp"}
SKIP_HINTS = ("logo", "avatar", "icon", "qrcode", "qr", "spacer", "blank", "loading")


def fetch(url: str, timeout: int = 15) -> bytes:
    req = Request(url, headers={"User-Agent": "Mozilla/5.0 auto-media-writer"})
    with urlopen(req, timeout=timeout) as resp:
        return resp.read()


def extract_image_urls(article_url: str, html: str) -> list[str]:
    urls = []
    patterns = [
        r"<img[^>]+(?:src|data-src|data-original|data-lazy-src)=['\"]([^'\"]+)['\"]",
        r"['\"](https?://[^'\"]+\.(?:jpg|jpeg|png|webp)(?:\?[^'\"]*)?)['\"]",
    ]
    for pattern in patterns:
        for match in re.findall(pattern, html, flags=re.I):
            url = unescape(match).strip()
            if not url or url.startswith("data:"):
                continue
            absolute = urljoin(article_url, url)
            lower = absolute.lower()
            if any(hint in lower for hint in SKIP_HINTS):
                continue
            ext = Path(urlparse(absolute).path).suffix.lower()
            if ext in IMAGE_EXTENSIONS and absolute not in urls:
                urls.append(absolute)
    return urls


def local_public_path(static_root: Path, public_prefix: str, relative_dir: Path, filename: str) -> str:
    try:
        rel = static_root.joinpath(relative_dir, filename).relative_to(static_root)
        return "/" + public_prefix.strip("/") + "/" + rel.as_posix()
    except ValueError:
        return "/" + public_prefix.strip("/") + "/" + relative_dir.as_posix() + "/" + filename


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--article-url", action="append", required=True)
    parser.add_argument("--task-id", required=True)
    parser.add_argument("--static-root", default="backend/static/article-images/uploads")
    parser.add_argument("--public-prefix", default="/static/article-images/uploads")
    parser.add_argument("--max-images", type=int, default=4)
    args = parser.parse_args()

    static_root = Path(args.static_root)
    now = datetime.now()
    relative_dir = Path(now.strftime("%Y")) / now.strftime("%m") / args.task_id
    target_dir = static_root / relative_dir
    target_dir.mkdir(parents=True, exist_ok=True)

    downloaded = []
    for article_url in args.article_url:
        if len(downloaded) >= args.max_images:
            break
        try:
            html = fetch(article_url).decode("utf-8", errors="ignore")
        except Exception as exc:
            print(f"warning: failed to fetch article {article_url}: {exc}", file=sys.stderr)
            continue
        for image_url in extract_image_urls(article_url, html):
            if len(downloaded) >= args.max_images:
                break
            try:
                data = fetch(image_url)
            except Exception as exc:
                print(f"warning: failed to fetch image {image_url}: {exc}", file=sys.stderr)
                continue
            if len(data) < 20_000:
                continue
            ext = Path(urlparse(image_url).path).suffix.lower()
            name = hashlib.sha1(image_url.encode("utf-8")).hexdigest()[:16] + ext
            output = target_dir / name
            output.write_bytes(data)
            public_path = local_public_path(static_root, args.public_prefix, relative_dir, name)
            downloaded.append({
                "role": "inline" if downloaded else "cover",
                "url": public_path,
                "localUrl": public_path,
                "type": "downloaded",
                "sourceUrl": article_url,
                "originalUrl": image_url,
                "copyrightRisk": "unknown",
                "watermark": False,
                "needsReview": False,
                "status": "ready",
            })

    print(json.dumps(downloaded, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
