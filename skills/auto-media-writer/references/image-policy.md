# 图片策略

## 决策顺序

1. 存在相关来源文章时，检查文章页面中位于正文或正文附近的图片。
2. 如果来源文章包含正文图片，下载到后端上传目录，并在 Markdown 中引用本地静态 URL。优先使用本地副本而不是热链，保证后续发布稳定。
3. 不要虚构本地静态 URL。只有路径由捆绑下载脚本打印、由 `POST /api/v1/skill-articles/upload-image` 返回，或确实来自后端上传目录下保存的栅格文件时，`/static/article-images/uploads/...` 才有效。
4. 不要仅因为图片可能需要人工审核就省略图片。如果图片不完美但可用，仍下载并插入，然后在 `images` 元数据中标注审核关注点。
5. 只跳过明显不是文章素材的图片：追踪像素、标志、头像、二维码、空白占位图、损坏文件、明显无关推荐图，或尺寸太小不适合文章使用的文件。
6. 如果没有来源图片或下载失败，先用文章关键词搜索互联网图片，再使用 AI 生成。复用文章实体和事件词，并按主题类型加入场景限定词，例如 `editorial photo`、`scene`、`press event`、`product photo`、`still`、`poster` 或 `concept art`。
7. 把选中的网络搜索结果下载到同一套本地上传目录，并在 Markdown 中使用本地 URL。保存文章时不要直接热链搜索结果 URL。
8. 只有来源图片下载和关键词网络图片搜索都失败时，才生成 AI 图片。
9. 禁止使用空 SVG 占位图作为生成文章图片。AI 兜底图必须是有具体视觉内容的栅格资源（`.png`、`.jpg` 或 `.webp`）。

## 存储

下载或生成的图片存放在：

`backend/static/article-images/uploads/<yyyy>/<mm>/<task-id>/`

通过以下路径暴露：

`/static/article-images/uploads/<yyyy>/<mm>/<task-id>/<filename>`

保存后的文章必须使用前端可渲染、后续发布任务可复用的公开路径或绝对后端 URL。

尽量使用 `scripts/download_article_images.py`：

```bash
python scripts/download_article_images.py \
  --article-url "https://example.com/news-a" \
  --article-url "https://example.com/news-b" \
  --task-id "<task-id>" \
  --static-root "backend/static/article-images/uploads" \
  --public-prefix "/static/article-images/uploads"
```

脚本会打印 `images` JSON 数组。文章 Markdown 使用其中的本地 URL。保存前检查结果，只移除明显无关的推荐图。

当文章图片缺失或下载尝试失败时，先使用内置图片搜索工具搜索网络，再用 `scripts/download_image_candidates.py` 下载选中的结果：

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

搜索词应来自文章关键词，而不是机械复制标题。优先使用具体实体加场景词。示例搜索模式：

- 娱乐：`<person name> latest update editorial photo`
- 科技：`<brand or product> press event product photo`
- 社会：`<place or event> scene photo`
- 影视：`<title or actor> still poster`

除非图片来源权威且与文章高度匹配，否则下载的搜索结果图片应标记 `needsReview: true`。

保存前，在仓库根目录运行载荷校验，或显式传入静态根目录：

```bash
python scripts/validate_skill_article_payload.py -
python scripts/validate_skill_article_payload.py - --static-root backend/static/article-images/uploads
```

如果任何 Markdown 图片或 `coverImageUrl` 指向缺失的本地文件，校验必须失败。通过先下载/上传真实图片来修复，不要编辑成看起来合理的文件名。

## 必需图片

- 封面：首选 1 张。
- 正文图：默认 2-3 张。
- 正文图插在自然小节断点之后，不要插在开头段落内部。
- 如果来源图片可用，至少封面应来自下载的来源图片。
- 如果来源图片不可用，在考虑 AI 生成前，封面应来自下载的网络搜索结果。
- 如果下载图片适配度不确定，在元数据中设置 `needsReview: true`，并保留在草稿中供用户审计。
- 如果没有来源图或搜索图可用，生成 AI 栅格图。最终保存的文章绝不能无图。
- 保存载荷必须包含 `coverImageUrl`，且 `markdownContent` 必须至少包含一张文件真实存在的 Markdown 图片。

## 图片元数据

在 `images` JSON 字段中记录每张图片：

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

`type` 可用值：

- `source_url`
- `downloaded`
- `searched`
- `ai_generated`
- `missing`

## AI 图片指导

生成真实感编辑类栅格图片，不要生成虚假截图或虚假新闻照片。避免让真实人物做出来源不支持的事情。娱乐文章优先选择具体但不构成诽谤的视觉场景，例如舞台灯光、显示泛化评论流的手机、媒体话筒、虚化公共活动背景，或象征关系边界的构图，而不是编造偷拍照片。

不要生成黑暗、空洞、文字过多或抽象 SVG 风格的文章图片。
