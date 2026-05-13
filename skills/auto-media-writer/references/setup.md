# Setup

## Required Skills

Install both skills in the target agent:

- `auto-media-writer`
- `humanizer-zh`

`auto-media-writer` must call `humanizer-zh` after the first article draft and before image insertion, publish payload creation, and database saving. If `humanizer-zh` is missing, the agent should stop and ask for installation instead of saving an article with `humanizeStatus: "done"`.

## TianAPI MCP

The skill can work without TianAPI MCP by using internet search only, but topic planning and news cross-checking are better when TianAPI MCP is installed.

Use your own TianAPI MCP service ID and API key. Do not commit real keys.

### Codex `config.toml`

Add this to your Codex config:

```toml
[mcp_servers.tianapi]
url = "https://mcp.tianapi.com/<TIANAPI_SERVICE_ID>/index?key=<TIANAPI_API_KEY>"
```

Example shape from the current environment:

```toml
[mcp_servers.tianapi-nl0tizmt]
url = "https://mcp.tianapi.com/nl0tizmt/index?key=<TIANAPI_API_KEY>"
```

### JSON MCP Config

For agents or clients that use JSON-style MCP configuration, copy this template:

```json
{
  "mcpServers": {
    "tianapi": {
      "url": "https://mcp.tianapi.com/<TIANAPI_SERVICE_ID>/index?key=<TIANAPI_API_KEY>"
    }
  }
}
```

If your agent expects the server name used by this workspace, use:

```json
{
  "mcpServers": {
    "tianapi-nl0tizmt": {
      "url": "https://mcp.tianapi.com/nl0tizmt/index?key=<TIANAPI_API_KEY>"
    }
  }
}
```

## Expected TianAPI Tools

After installation, the agent should expose tools similar to:

- `generalnews`
- `allnews`
- `networkhot`
- `nethot`
- `weibohot`
- `douyinhot`
- `toutiaohot`
- `wxhottopic`

Tool names may be prefixed by the MCP server namespace in some agents.

## Auto Article Backend

Set the backend URL for saving generated articles:

```bash
AUTO_ARTICLE_BASE_URL=http://localhost:9001
```

The backend must expose:

- `POST /api/v1/skill-articles`
- `GET /api/v1/skill-articles`
- `POST /api/v1/skill-articles/:id/publish-package`
- `/static/article-images/...`

## Quick Test

Ask the target agent:

```text
使用 auto-media-writer skill，检查天行 MCP 是否可用，并列出可调用的热搜/新闻工具。
```

If TianAPI MCP is unavailable, continue with internet search and tell the user that TianAPI cross-checking was skipped.

Then test the de-AI dependency:

```text
使用 auto-media-writer skill 和 humanizer-zh skill，生成一段 200 字今日头条热点评论样稿，但不要入库。
```
