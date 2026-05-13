# Auto Article API

## Backend Base URL

Read the backend base URL from the user, current workspace config, or an environment variable such as `AUTO_ARTICLE_BASE_URL`. Default local development URL is usually:

`http://localhost:9001`

Never connect to MySQL directly from this skill.

## Create Skill Article

`POST /api/v1/skill-articles`

Required fields:

- `platform`: `toutiao`, `baijiahao`, `xiaohongshu`, or `zhihu`
- `keyword`
- `title`
- `markdownContent`

Recommended fields:

```json
{
  "taskId": "uuid-or-stable-id",
  "platform": "toutiao",
  "keyword": "热点关键词",
  "keywordSegments": "[\"热点\", \"关键词\"]",
  "category": "society",
  "title": "最终标题",
  "titleOptions": "[\"标题1\", \"标题2\"]",
  "summary": "100字以内摘要",
  "markdownContent": "Markdown正文",
  "tags": "[\"标签1\", \"标签2\"]",
  "coverImageUrl": "/static/article-images/uploads/2026/05/task/cover.jpg",
  "coverImageType": "ai_generated",
  "images": "[]",
  "sources": "[]",
  "hotTopics": "[]",
  "styleProfile": "{}",
  "wordCount": 1800,
  "modelProvider": "openai",
  "modelName": "model-name",
  "promptVersion": "auto-media-writer-article-v1",
  "skillVersion": "auto-media-writer-v1",
  "humanizeStatus": "done",
  "status": "draft",
  "publishStatus": "unpublished",
  "publishPayload": "{}"
}
```

All JSON-like fields are stored as strings for portability. Encode arrays/objects with valid JSON.

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

## Prepare Publish Data

`POST /api/v1/skill-articles/:id/publish-package`

This does not publish to any platform. It marks the article as `ready_to_publish`, sets `publishStatus` to `ready`, and writes a `publishPayload` JSON string for future one-click publishing.
