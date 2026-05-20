#!/usr/bin/env python3
import argparse
import hashlib
import json
import sys
from datetime import datetime
from pathlib import Path
from urllib.parse import urlparse
from urllib.request import Request, urlopen


EXTENSIONS = {".jpg", ".jpeg", ".png", ".webp"}


class ChineseArgumentParser(argparse.ArgumentParser):
    def format_help(self):
        return (
            super().format_help()
            .replace("usage:", "用法:")
            .replace("positional arguments:", "位置参数:")
            .replace("options:", "选项:")
            .replace("show this help message and exit", "显示帮助并退出")
        )


def fetch(url: str, timeout: int = 15):
    req = Request(url, headers={"User-Agent": "Mozilla/5.0 auto-media-writer"})
    return urlopen(req, timeout=timeout)


def local_public_path(static_root: Path, public_prefix: str, relative_dir: Path, filename: str) -> str:
    rel = static_root.joinpath(relative_dir, filename).relative_to(static_root)
    return "/" + public_prefix.strip("/") + "/" + rel.as_posix()


def infer_extension(image_url: str, content_type: str) -> str:
    ext = Path(urlparse(image_url).path).suffix.lower()
    if ext in EXTENSIONS:
        return ext
    if "png" in content_type:
        return ".png"
    if "webp" in content_type:
        return ".webp"
    if "jpeg" in content_type or "jpg" in content_type:
        return ".jpg"
    return ".jpg"


def main() -> int:
    parser = ChineseArgumentParser(description="下载图片搜索候选结果并生成本地静态地址。")
    parser.add_argument("--image-url", action="append", required=True)
    parser.add_argument("--task-id", required=True)
    parser.add_argument("--static-root", default="backend/static/article-images/uploads")
    parser.add_argument("--public-prefix", default="/static/article-images/uploads")
    parser.add_argument("--source-url", default="")
    parser.add_argument("--query", default="")
    parser.add_argument("--max-images", type=int, default=4)
    parser.add_argument("--copyright-risk", default="unknown")
    parser.add_argument("--needs-review", action="store_true")
    args = parser.parse_args()

    static_root = Path(args.static_root)
    now = datetime.now()
    relative_dir = Path(now.strftime("%Y")) / now.strftime("%m") / args.task_id
    target_dir = static_root / relative_dir
    target_dir.mkdir(parents=True, exist_ok=True)

    downloaded = []
    for image_url in args.image_url:
        if len(downloaded) >= args.max_images:
            break
        try:
            with fetch(image_url) as resp:
                data = resp.read()
                content_type = resp.headers.get("Content-Type", "")
        except Exception as exc:
            print(f"警告：获取图片失败 {image_url}: {exc}", file=sys.stderr)
            continue
        if len(data) < 20_000:
            print(f"警告：跳过过小图片 {image_url}", file=sys.stderr)
            continue
        ext = infer_extension(image_url, content_type)
        name = hashlib.sha1(image_url.encode("utf-8")).hexdigest()[:16] + ext
        output = target_dir / name
        output.write_bytes(data)
        public_path = local_public_path(static_root, args.public_prefix, relative_dir, name)
        downloaded.append({
            "role": "inline" if downloaded else "cover",
            "url": public_path,
            "localUrl": public_path,
            "type": "searched",
            "sourceUrl": args.source_url or image_url,
            "originalUrl": image_url,
            "searchQuery": args.query,
            "copyrightRisk": args.copyright_risk,
            "watermark": False,
            "needsReview": args.needs_review,
            "status": "ready",
        })

    print(json.dumps(downloaded, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
