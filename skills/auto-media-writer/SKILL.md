---
name: auto-media-writer
description: Generate Chinese self-media articles for Toutiao, Baijiahao, Xiaohongshu, or Zhihu from user keywords, hot topics, internet/news research, TianAPI MCP data, source images, and AI-generated image fallbacks. Use when planning daily article topics, researching trending keywords, writing platform-specific drafts, humanizing AI-sounding prose, preparing publish-ready article payloads, or saving generated articles into the auto-article project.
---

# Auto Media Writer

## Iron Rules

- Never save an article to the auto-article backend before calling `humanizer-zh`. The de-AI pass is mandatory, not optional. If `humanizer-zh` is unavailable, stop before database insertion and tell the user to install or enable it.
- Never save an article without images. If cited source articles contain body images, download those images to `backend/static/article-images/uploads/<yyyy>/<mm>/<task-id>/`, insert their local `/static/article-images/uploads/...` URLs into the Markdown, and record them in the `images` field. If an image later looks unsuitable, the user will review and replace it manually; do not omit images for that reason.
- If source articles have no usable images or image downloads fail, search the open web for relevant images using the article keywords before falling back to AI generation. Download selected search results to the same local uploads directory and insert those local URLs into the article.
- Generate concrete raster fallback images (`png`, `jpg`, or `webp`) only after both source-image download and keyword-based internet image search fail. Never use SVG placeholders.

## Core Workflow

1. Resolve the target platform. Require one of `toutiao`, `baijiahao`, `xiaohongshu`, or `zhihu`; ask the user only if it is missing or ambiguous.
2. Resolve the topic. If the user gives keywords, segment them into useful search terms. If the user asks what to write, collect hot topics from internet search and TianAPI MCP hot lists, then present 5-10 options.
3. Research with internet search first. If internet results look stale, thin, or contradictory, cross-check with TianAPI MCP news/hot-topic tools. Prefer recent, source-attributed facts.
4. Classify the article automatically. Use `entertainment` for celebrity/film/TV/influencer gossip and public reactions; use `society` for broad social news; otherwise choose a concise category such as `finance`, `tech`, `sports`, or `general`. Ask the user when the category changes the writing direction and remains uncertain.
5. Build a source pack containing title, URL, source, publish time, summary, image URL, and reliability notes. Do not invent facts, quotes, data, or dates absent from sources.
6. Draft the article for the selected platform using `references/platform-style.md`.
7. Humanize the article before saving by calling the `humanizer-zh` skill. Pass the drafted Markdown article plus the selected platform, category, target length, title options, and source-fact constraints. Require `humanizer-zh` to preserve verified facts, URLs, dates, names, Markdown headings, image markers, and legal risk boundaries while removing AI-sounding transitions, template paragraphs, generic moralizing, and over-balanced phrasing.
8. Prepare images using `references/image-policy.md`. First download usable body images from cited source articles. If none are usable or downloads fail, search the web for relevant images with article keywords, download the selected results to the same local uploads image directory, and reference those local URLs. Use AI raster generation only if both earlier paths fail. Insert Markdown image tags at suitable positions. Do not save an image-free article.
9. Create a publish-ready payload using `references/auto-article-api.md`. Validate payload shape with `scripts/validate_skill_article_payload.py` when possible.
10. Save directly to the auto-article backend `POST /api/v1/skill-articles`. Do not write the database directly.
11. Record any improvement learned during use in this skill before future runs when the user asks the skill to evolve.

## Topic And Research Rules

- Segment Chinese keywords into entity, event, location, platform, and intent terms. Example: `某明星离婚风波` -> `某明星`, `离婚`, `风波`, `回应`, `网友热议`.
- For daily planning, compare hot lists across Baidu, Weibo, Douyin, WeChat, Toutiao, and broad web news. Rank options by recency, discussion heat, source availability, platform fit, and image availability.
- Use internet search as the primary source. Use TianAPI MCP when web results are delayed, inaccessible, or need corroboration.
- Reuse the same entity and event keywords for image-search fallback. Combine the main subject with scene qualifiers such as `scene`, `editorial photo`, `event photo`, `product photo`, `press event`, `still`, or `concept art` so the searched images stay editorially relevant.
- Keep source URLs in the final payload. If source URLs include usable images, prefer them unless watermarked, copyrighted, low quality, or unavailable.
- Include concrete dates when discussing recent events.

## Platform Defaults

- Entertainment articles: 500-800 Chinese characters by default, hotter gossip/commentary voice, but avoid fabrication and defamation. Expand only when the user asks for deeper analysis.
- Society/general news: 1500-2500 Chinese characters, hotspot commentary plus emotional resonance.
- Toutiao: clear headline, fast context, strong but moderate viewpoint, broad-reader readability.
- Baijiahao: more explanatory, slightly more structured, search-friendly headings.
- Xiaohongshu: shorter, more conversational, stronger emotional hook, practical takeaways, topic tags, and a more native note-style rhythm.
- Zhihu: more reasoned and analytical, explicit question framing, causes, tradeoffs, and conclusion.

Read `references/platform-style.md` for detailed platform style and title rules.

## Required Output

Always generate:

- one selected `title`
- multiple `titleOptions`
- `summary`
- Markdown `markdownContent`
- `tags`
- `sources` JSON
- `images` JSON
- `styleProfile` JSON
- `publishPayload` JSON for future one-click publishing

Use `humanizeStatus: "done"` only after the `humanizer-zh` rewrite pass is complete. If `humanizer-zh` is unavailable, stop before saving and tell the user to install it instead of silently skipping the de-AI pass. Never send `POST /api/v1/skill-articles` with `humanizeStatus` unset, `pending`, or `skipped`. Use `promptVersion` and `skillVersion` so future changes are traceable.

## Image Handling

Follow `references/image-policy.md`.

Default output requires one cover image and 2-3 inline images. Prefer downloaded local copies of source-article images. If source images are imperfect but usable, still download and insert them, mark review concerns in `images`, and let the user replace them later. If source images are missing or fail to download, search the internet for images that match the article keywords and select editorially relevant results before using AI generation. Record searched images distinctly in `images` metadata so the user can audit them later. Generate concrete raster AI images (`png`, `jpg`, or `webp`) only when both earlier paths fail, never SVG placeholders.

## Saving

Use the auto-article HTTP API contract in `references/auto-article-api.md`. The skill must stay portable: never require direct MySQL credentials, never hard-code secrets, and read the backend base URL from user input or environment.

For installing TianAPI MCP in another agent, use the copy-paste templates in `references/setup.md`.

## Evolution

When the user says the result is too stiff, too generic, too long, wrong platform tone, weak title, poor image choice, or otherwise not right, update the relevant style, image, API, or workflow guidance in this skill. Keep changes concise and reusable.
