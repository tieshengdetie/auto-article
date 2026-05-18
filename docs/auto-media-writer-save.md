# Auto Media Writer 一键入库

这个仓库的固定入口只做一件事：复用 `skills/auto-media-writer/scripts/save_skill_article.py`，先跑 `validate_skill_article_payload.py`，再调用 `POST /api/v1/skill-articles`。以后 agent 不需要临时写 Python、直连 MySQL 或重新推断后端入库逻辑。

## 推荐命令

Windows / PowerShell：

```powershell
.\scripts\save-skill-article.ps1 C:\tmp\auto-media-writer\demo.payload.json -DryRun
.\scripts\save-skill-article.ps1 C:\tmp\auto-media-writer\demo.payload.json -BaseUrl http://localhost:9001
```

Unix-like shell：

```sh
sh scripts/save-skill-article.sh /tmp/auto-media-writer/demo.payload.json --dry-run
sh scripts/save-skill-article.sh /tmp/auto-media-writer/demo.payload.json --base-url http://localhost:9001
```

Make：

```sh
make validateSkillArticle PAYLOAD=/tmp/auto-media-writer/demo.payload.json
make saveSkillArticle PAYLOAD=/tmp/auto-media-writer/demo.payload.json BASE_URL=http://localhost:9001
```

`PAYLOAD=-` 可以从 stdin 读取 JSON。后端地址优先使用 `AUTO_ARTICLE_BASE_URL`，也可以通过命令参数显式传入。图片静态目录不在默认位置时，传 `--static-root` / `-StaticRoot`，或设置 `AUTO_ARTICLE_STATIC_ROOT`。

脚本需要 Python 3.9+。默认会优先尝试 `python3`，必要时可以用 `PYTHON` 环境变量指定解释器。

## Payload 约定

示例模板见 `examples/auto-media-writer.payload.example.json`。真正入库前必须满足：

- 已完成 `humanizer-zh` 去 AI 改写。
- `platform` 是 `toutiao`、`baijiahao`、`xiaohongshu` 或 `zhihu`。
- `coverImageUrl` 和正文 Markdown 图片都指向真实存在的本地静态图片，路径形如 `/static/article-images/uploads/.../*.jpg|png|webp`。
- `titleOptions` 如果提供，必须是 JSON 编码后的字符串。
- 不把 source pack、style profile、模型信息、发布状态等过程元数据发送给后端。

## 建议流程

1. 生成文章和图片，把图片落到 `backend/static/article-images/uploads/<yyyy>/<mm>/<task-id>/`。
2. 组装最小 payload，临时文件放在仓库外，例如 `C:\tmp\auto-media-writer\` 或 `/tmp/auto-media-writer/`。
3. 先执行 dry-run 命令，确认校验通过并看到目标 URL。
4. 去掉 `-DryRun` / `--dry-run` 正式入库。

不要新增一次性入库脚本，不要绕过校验脚本，不要直连数据库。
