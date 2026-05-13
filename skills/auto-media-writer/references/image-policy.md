# Image Policy

## Decision Order

1. When a related source article exists, inspect the article page for images that appear in or near the article body.
2. Download relevant article images to the backend static directory and reference the local static URL in Markdown. Prefer local copies over hotlinking so future publishing is stable.
3. Skip images that are unrelated, watermarked, obvious thumbnails from unrelated recommendations, tracking pixels, logos, avatars, QR codes, or too small for article use.
4. Generate AI images only when no safe relevant article image can be found, the source image is broken, or copyright/watermark risk is obvious.
5. Never use empty SVG placeholder graphics as generated article images. AI fallback images should be raster assets (`.png`, `.jpg`, or `.webp`) with concrete visual content.

## Storage

Store downloaded or generated images under:

`backend/static/article-images/uploads/<yyyy>/<mm>/<task-id>/`

Expose them through:

`/static/article-images/uploads/<yyyy>/<mm>/<task-id>/<filename>`

The saved article should use public paths or absolute backend URLs that can be rendered by the frontend and reused by future publish jobs.

Use `scripts/download_article_images.py` when possible:

```bash
python scripts/download_article_images.py \
  --article-url "https://example.com/news-a" \
  --article-url "https://example.com/news-b" \
  --task-id "<task-id>" \
  --static-root "backend/static/article-images/uploads" \
  --public-prefix "/static/article-images/uploads"
```

The script prints an `images` JSON array. Review results and remove unrelated recommendation images before saving.

## Required Images

- Cover: exactly 1 preferred image.
- Inline: 2-3 images by default.
- Insert inline images after natural section breaks, not inside the opening paragraph.
- If source images are available, at least the cover should come from a downloaded source image unless there is a watermark/copyright problem.

## Image Metadata

Record each image in the `images` JSON field:

```json
{
  "role": "cover",
  "url": "/static/article-images/uploads/2026/05/task/image.jpg",
  "localUrl": "/static/article-images/uploads/2026/05/task/image.jpg",
  "type": "downloaded",
  "sourceUrl": "https://example.com/news",
  "copyrightRisk": "low",
  "watermark": false,
  "status": "ready"
}
```

Use `type` values:

- `source_url`
- `downloaded`
- `ai_generated`
- `missing`

## AI Image Guidance

Generate realistic editorial raster images, not fake screenshots or fake news photos. Avoid making real people appear to do things that sources do not support. For entertainment articles, prefer concrete but non-defamatory visual scenes such as stage lights, phones showing generic comment feeds, media microphones, blurred public-event backdrops, or symbolic relationship-boundary compositions rather than invented paparazzi photos.

Do not generate dark, empty, text-heavy, or abstract SVG artwork as article imagery.
