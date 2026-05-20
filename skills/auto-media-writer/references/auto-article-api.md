# auto-article 接口

## 后端基准 URL

从用户、当前工作区配置或 `AUTO_ARTICLE_BASE_URL` 等环境变量读取后端基准 URL。不要假设后端一定部署在本地。部署服务应设置真实来源，例如：

`AUTO_ARTICLE_BASE_URL=https://api.example.com`

本地开发兜底通常是：

`AUTO_ARTICLE_BASE_URL=http://localhost:9001`

本地后端直连要绕过代理：当基准 URL 是 `localhost`、`127.0.0.1`、`::1` 或明确的本机开发地址时，优先改用 `http://127.0.0.1:<port>`。如果请求返回 `502 Bad Gateway`、代理连接错误，或环境变量存在 `http_proxy`、`https_proxy`、`all_proxy` 指向本机代理，只在本次保存命令进程里临时清空这些代理变量，并设置 `no_proxy=localhost,127.0.0.1,::1` 后重试。不要全局修改用户代理配置；远程部署 URL 仍按用户提供的地址请求。

本技能绝不要直接连接 MySQL。

## 创建技能文章

`POST /api/v1/skill-articles`

后端必需字段：

- `platform`：`toutiao`、`baijiahao`、`xiaohongshu` 或 `zhihu`
- `keyword`
- `title`
- `markdownContent`

本技能保存前额外要求：

- `coverImageUrl`：必须指向 `/static/article-images/uploads/...` 下真实下载或上传的栅格图片
- `markdownContent` 中至少包含一张 Markdown 图片

推荐字段：

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

后端只保存网页需要展示的生成文章字段。信源包、风格档案、模型元数据、去 AI 状态和发布状态保留在工作笔记中，不要发送给 API。

## 快速保存命令

优先使用仓库级入口。它会委托给捆绑保存脚本，避免写新代码或检查后端文件：

```powershell
.\scripts\save-skill-article.ps1 C:\tmp\auto-media-writer\demo.payload.json -DryRun
.\scripts\save-skill-article.ps1 C:\tmp\auto-media-writer\demo.payload.json -BaseUrl http://localhost:9001
```

```sh
make validateSkillArticle PAYLOAD=/tmp/auto-media-writer/demo.payload.json
make saveSkillArticle PAYLOAD=/tmp/auto-media-writer/demo.payload.json BASE_URL=http://localhost:9001
```

优先通过标准输入传入构造好的载荷 JSON，并用 `-` 作为载荷参数。这样完全避免创建 `payload.json` 文件：

```sh
python skills/auto-media-writer/scripts/save_skill_article.py -
```

非本地后端必须显式设置：

```sh
python skills/auto-media-writer/scripts/save_skill_article.py - --base-url https://api.example.com
```

脚本会先用 `scripts/validate_skill_article_payload.py` 校验，再发送 `POST /api/v1/skill-articles`。校验会检查必需文章字段、`titleOptions` 的 JSON 编码，以及 `coverImageUrl` 和 Markdown 图片 URL 对应的本地图片文件是否存在。使用 `--dry-run` 可只校验并预览目标 URL，不保存。

如果 shell 或标准输入处理不方便，把临时载荷 JSON 写到仓库外，例如 Windows 的 `C:\tmp\auto-media-writer\<task-id>.payload.json` 或 Unix 类系统的 `/tmp/auto-media-writer/<task-id>.payload.json`。不要在项目根目录或技能目录创建或留下 `payload.json`、`payload_draft.json` 等临时文章文件。保存或试运行校验后删除临时载荷，除非用户明确要求保留。

在仓库根目录外运行时，显式传入上传目录：

```sh
python skills/auto-media-writer/scripts/save_skill_article.py - --static-root backend/static/article-images/uploads
```

不要为了保存文章而连接 MySQL、读取 DAO/model/service 代码，或编写一次性 Python 数据库插入脚本。

## 列出技能文章

`GET /api/v1/skill-articles?page=1&pageSize=8&platform=toutiao&keyword=...`

## 获取技能文章

`GET /api/v1/skill-articles/:id`

## 更新技能文章

`PUT /api/v1/skill-articles/:id`

用于人工编辑、替换图片或重新去 AI 后的内容更新。

## 上传本地文章图片

`POST /api/v1/skill-articles/upload-image`

使用 multipart 表单数据，并用 `image` 作为文件字段。后端会把图片保存到 `backend/static/article-images/uploads/yyyy/mm/`，并返回公开静态路径：

```json
{
  "url": "/static/article-images/uploads/2026/05/example.jpg",
  "filename": "example.jpg",
  "size": 123456
}
```

把返回的 `url` 插入 Markdown，不要嵌入 base64 数据。
