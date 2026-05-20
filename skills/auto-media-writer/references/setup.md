# 设置

## 必需技能

目标代理中需要安装两个技能：

- `auto-media-writer`
- `humanizer-zh`

`auto-media-writer` 必须在第一版文章草稿后、图片插入、发布载荷创建和数据库保存前调用 `humanizer-zh`。如果缺少 `humanizer-zh`，代理应停止并请求安装，而不是用 `humanizeStatus: "done"` 保存文章。

## TianAPI MCP

没有 TianAPI MCP 时，本技能仍可只用互联网搜索工作；但安装后，选题规划和新闻交叉核验会更好。

使用你自己的 TianAPI MCP 服务 ID 和 API key。不要提交真实密钥。

### Codex `config.toml`

把以下配置加入 Codex 配置：

```toml
[mcp_servers.tianapi]
url = "https://mcp.tianapi.com/<TIANAPI_SERVICE_ID>/index?key=<TIANAPI_API_KEY>"
```

当前环境中的示例形态：

```toml
[mcp_servers.tianapi-nl0tizmt]
url = "https://mcp.tianapi.com/nl0tizmt/index?key=<TIANAPI_API_KEY>"
```

### JSON MCP 配置

对使用 JSON 风格 MCP 配置的代理或客户端，复制以下模板：

```json
{
  "mcpServers": {
    "tianapi": {
      "url": "https://mcp.tianapi.com/<TIANAPI_SERVICE_ID>/index?key=<TIANAPI_API_KEY>"
    }
  }
}
```

如果代理需要使用当前工作区的服务名，使用：

```json
{
  "mcpServers": {
    "tianapi-nl0tizmt": {
      "url": "https://mcp.tianapi.com/nl0tizmt/index?key=<TIANAPI_API_KEY>"
    }
  }
}
```

## 预期 TianAPI 工具

安装后，代理应暴露类似工具：

- `generalnews`
- `allnews`
- `networkhot`
- `nethot`
- `weibohot`
- `douyinhot`
- `toutiaohot`
- `wxhottopic`

在某些代理中，工具名可能带有 MCP 服务命名空间前缀。

## auto-article 后端

设置保存生成文章所需的后端 URL。使用真实部署后端来源；localhost 只用于本地开发：

```bash
AUTO_ARTICLE_BASE_URL=https://api.example.com
# 仅本地开发：
# AUTO_ARTICLE_BASE_URL=http://localhost:9001
```

后端必须暴露：

- `POST /api/v1/skill-articles`
- `GET /api/v1/skill-articles`
- `/static/article-images/uploads/...`

## 快速测试

向目标代理发送：

```text
使用 auto-media-writer 技能，检查天行 MCP 是否可用，并列出可调用的热搜/新闻工具。
```

如果 TianAPI MCP 不可用，继续使用互联网搜索，并告诉用户已跳过 TianAPI 交叉核验。

然后测试去 AI 依赖：

```text
使用 auto-media-writer 技能和 humanizer-zh 技能，生成一段 200 字今日头条热点评论样稿，但不要入库。
```
