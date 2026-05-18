# Auto Article API

## Backend Base URL

Read the backend base URL from the user, current workspace config, or an environment variable such as `AUTO_ARTICLE_BASE_URL`. Do not assume the backend is deployed locally. For a deployed service, set the real origin, for example:

`AUTO_ARTICLE_BASE_URL=https://api.example.com`

Local development fallback is usually:

`AUTO_ARTICLE_BASE_URL=http://localhost:9001`

Never connect to MySQL directly from this skill.

## Create Skill Article

`POST /api/v1/skill-articles`

Required fields for the backend:

- `platform`: `toutiao`, `baijiahao`, `xiaohongshu`, or `zhihu`
- `keyword`
- `title`
- `markdownContent`

Required fields for this skill before saving:

- `coverImageUrl`: must point to a real downloaded/uploaded raster image under `/static/article-images/uploads/...`
- at least one Markdown image in `markdownContent`

Recommended fields:

```json
{
  "taskId": "uuid-or-stable-id",
  "platform": "toutiao",
  "keyword": "热点关键词",
  "category": "society",
  "title": "最终标题",
  "titleOptions": "[\"标题1\", \"标题2\"]",
  "summary": "100字以内摘要",
  "markdownContent": "Markdown正文",
  "coverImageUrl": "/static/article-images/uploads/2026/05/task/cover.jpg"
}
```

The backend stores only the generated article fields needed by the web page. Keep source packs, style profiles, model metadata, humanizer status, and publishing state in working notes instead of sending them to the API.

## Fast Save Command

Prefer the repository-level entrypoint, which delegates to the bundled save script instead of writing new code or inspecting backend files:

```powershell
.\scripts\save-skill-article.ps1 C:\tmp\auto-media-writer\demo.payload.json -DryRun
.\scripts\save-skill-article.ps1 C:\tmp\auto-media-writer\demo.payload.json -BaseUrl http://localhost:9001
```

```sh
make validateSkillArticle PAYLOAD=/tmp/auto-media-writer/demo.payload.json
make saveSkillArticle PAYLOAD=/tmp/auto-media-writer/demo.payload.json BASE_URL=http://localhost:9001
```

Prefer passing the constructed payload JSON through stdin by using `-` as the payload argument. This avoids creating a `payload.json` file at all:

```sh
python skills/auto-media-writer/scripts/save_skill_article.py -
```

Set the backend explicitly when it is not local:

```sh
python skills/auto-media-writer/scripts/save_skill_article.py - --base-url https://api.example.com
```

The script validates with `scripts/validate_skill_article_payload.py` before sending `POST /api/v1/skill-articles`. Validation checks required article fields, `titleOptions` JSON encoding, and local image file existence for `coverImageUrl` and Markdown image URLs. Use `--dry-run` to validate and preview the target URL without saving.

If shell/stdin handling is impractical, write the transient payload JSON file outside the repository, for example `C:\tmp\auto-media-writer\<task-id>.payload.json` on Windows or `/tmp/auto-media-writer/<task-id>.payload.json` on Unix-like systems. Do not create or leave `payload.json`, `payload_draft.json`, or similar temporary article files in the project root or skill directory. Delete the temporary payload after save or dry-run validation unless the user explicitly asks to keep it.

When running outside the repository root, pass the uploads directory explicitly:

```sh
python skills/auto-media-writer/scripts/save_skill_article.py - --static-root backend/static/article-images/uploads
```

Do not connect to MySQL, read DAO/model/service code, or write a one-off Python database insertion script for article saving.

## List Skill Articles

`GET /api/v1/skill-articles?page=1&pageSize=8&platform=toutiao&keyword=...`

## Get Skill Article

`GET /api/v1/skill-articles/:id`

## Update Skill Article

`PUT /api/v1/skill-articles/:id`

Use this for manual edits, image changes, or re-humanized content.

## Upload Local Article Image

`POST /api/v1/skill-articles/upload-image`

Use multipart form data with an `image` file field. The backend saves the image under `backend/static/article-images/uploads/yyyy/mm/` and returns a public static path:

```json
{
  "url": "/static/article-images/uploads/2026/05/example.jpg",
  "filename": "example.jpg",
  "size": 123456
}
```

Insert the returned `url` into Markdown instead of embedding base64 data.
