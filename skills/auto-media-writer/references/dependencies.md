# Dependencies

## Required Agent Capabilities

- Internet search for recent web/news information.
- TianAPI MCP tools when available for hot lists and news corroboration.
- `humanizer-zh` skill for the mandatory de-AI rewrite pass before saving.
- Image generation or image editing capability for cover and inline fallback images.
- HTTP client ability to call the auto-article backend API.

For copy-paste TianAPI MCP installation templates, see `setup.md`.

## Optional Capabilities

- Browser automation for checking source image reachability.
- Local file access for downloaded/generated images.

## Auto Article Project Requirements

- Backend service running.
- `skill_generated_articles` migration applied.
- `/static` route enabled for static image access.
- `backend/static/article-images/` writable by the backend or the agent performing image downloads.

## Environment Variables

Use these names when possible:

- `AUTO_ARTICLE_BASE_URL`: backend base URL.
- `AUTO_ARTICLE_IMAGE_BASE_DIR`: optional local image root, default `backend/static/article-images`.
- `AUTO_ARTICLE_PUBLIC_STATIC_PREFIX`: optional public prefix, default `/static/article-images`.

Do not store API keys, database credentials, or platform account credentials in this skill.
