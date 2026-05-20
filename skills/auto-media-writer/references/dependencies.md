# 依赖

## 必需的代理能力

- 能联网搜索近期网页和新闻信息。
- 可用时使用 TianAPI MCP 工具获取热榜并交叉印证新闻。
- 保存前必须使用 `humanizer-zh` 技能执行去 AI 改写。
- 具备图片生成或图片编辑能力，用于封面和正文兜底图片。
- 具备 HTTP 客户端能力，可调用 auto-article 后端 API。

TianAPI MCP 安装复制模板见 `setup.md`。

## 可选能力

- 浏览器自动化，用于检查来源图片是否可访问。
- 本地文件访问能力，用于下载或生成的图片。

## auto-article 项目要求

- 后端服务已运行。
- `skill_generated_articles` 迁移已应用。
- `/static` 路由已启用，可访问静态图片。
- `backend/static/article-images/uploads/` 可由后端或负责下载图片的代理写入。

## 环境变量

尽量使用以下名称：

- `AUTO_ARTICLE_BASE_URL`：后端基准 URL。
- `AUTO_ARTICLE_IMAGE_BASE_DIR`：可选本地图片根目录，默认 `backend/static/article-images/uploads`。
- `AUTO_ARTICLE_PUBLIC_STATIC_PREFIX`：可选公开前缀，默认 `/static/article-images/uploads`。

不要在本技能中保存 API 密钥、数据库凭据或平台账号凭据。
