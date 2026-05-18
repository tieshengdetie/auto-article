# Image Policy

## Decision Order

1. When a related source article exists, inspect the article page for images that appear in or near the article body.
2. If source articles contain body images, download them to the backend uploads directory and reference the local static URL in Markdown. Prefer local copies over hotlinking so future publishing is stable.
3. Do not invent local static URLs. A path under `/static/article-images/uploads/...` is valid only when it was printed by a bundled download script, returned by `POST /api/v1/skill-articles/upload-image`, or generated from an actual raster file saved under the backend uploads directory.
4. Do not omit images merely because they may need manual review. If an image is imperfect but usable, still download and insert it, then mark review concerns in the `images` metadata.
5. Skip only images that are clearly not article material: tracking pixels, logos, avatars, QR codes, blank placeholders, broken files, obvious unrelated recommendations, or files too small for article use.
6. If no source image exists or downloads fail, search the internet for images using the article keywords before using AI generation. Reuse the article entities and event words, and add scene qualifiers such as `editorial photo`, `scene`, `press event`, `product photo`, `still`, `poster`, or `concept art` based on topic type.
7. Download selected internet-search results to the same local uploads directory and use local URLs in Markdown. Do not hotlink search result URLs directly in saved articles.
8. Generate AI images only when both source-image download and keyword-based internet image search fail.
9. Never use empty SVG placeholder graphics as generated article images. AI fallback images should be raster assets (`.png`, `.jpg`, or `.webp`) with concrete visual content.

## Storage

Store downloaded or generated images under:

`backend/static/article-images/uploads/<yyyy>/<mm>/<task-id>/`

Expose them through:

`/static/article-images/uploads/<yyyy>/<mm>/<task-id>/<filename>`

The saved article must use public paths or absolute backend URLs that can be rendered by the frontend and reused by future publish jobs.

Use `scripts/download_article_images.py` when possible:

```bash
python scripts/download_article_images.py \
  --article-url "https://example.com/news-a" \
  --article-url "https://example.com/news-b" \
  --task-id "<task-id>" \
  --static-root "backend/static/article-images/uploads" \
  --public-prefix "/static/article-images/uploads"
```

The script prints an `images` JSON array. Use those local URLs in the article Markdown. Review results and remove only clearly unrelated recommendation images before saving.

When article images are missing or download attempts fail, search the web with the built-in image search tool first, then download the chosen results with `scripts/download_image_candidates.py`:

```bash
python scripts/download_image_candidates.py \
  --image-url "https://example.com/search-result-a.jpg" \
  --image-url "https://example.com/search-result-b.webp" \
  --task-id "<task-id>" \
  --query "celebrity name press event editorial photo" \
  --source-url "https://images.example.com/result-page" \
  --needs-review \
  --static-root "backend/static/article-images/uploads" \
  --public-prefix "/static/article-images/uploads"
```

The search query should be derived from the article keywords rather than copied mechanically from the title. Prefer specific entities plus scene words. Example query patterns:

- Entertainment: `<person name> latest update editorial photo`
- Tech: `<brand or product> press event product photo`
- Society: `<place or event> scene photo`
- Film/TV: `<title or actor> still poster`

Mark downloaded search-result images with `needsReview: true` unless they are clearly authoritative and tightly matched to the article.

Before saving, run payload validation from the repository root or pass the static root explicitly:

```bash
python scripts/validate_skill_article_payload.py -
python scripts/validate_skill_article_payload.py - --static-root backend/static/article-images/uploads
```

Validation must fail if any Markdown image or `coverImageUrl` points to a missing local file. Fix failures by downloading/uploading the real image first, not by editing the filename to look plausible.

## Required Images

- Cover: exactly 1 preferred image.
- Inline: 2-3 images by default.
- Insert inline images after natural section breaks, not inside the opening paragraph.
- If source images are available, at least the cover should come from a downloaded source image.
- If source images are unavailable, the cover should come from a downloaded web-search result before AI generation is considered.
- If downloaded images have uncertain fit, set `needsReview: true` in metadata and keep them in the draft for the user's audit.
- If no source image or searched image is available, generate AI raster images instead. The final saved article must never be image-free.
- The saved payload must include `coverImageUrl`, and `markdownContent` must include at least one Markdown image whose file exists locally.

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
  "needsReview": false,
  "status": "ready"
}
```

Use `type` values:

- `source_url`
- `downloaded`
- `searched`
- `ai_generated`
- `missing`

## AI Image Guidance

Generate realistic editorial raster images, not fake screenshots or fake news photos. Avoid making real people appear to do things that sources do not support. For entertainment articles, prefer concrete but non-defamatory visual scenes such as stage lights, phones showing generic comment feeds, media microphones, blurred public-event backdrops, or symbolic relationship-boundary compositions rather than invented paparazzi photos.

Do not generate dark, empty, text-heavy, or abstract SVG artwork as article imagery.
