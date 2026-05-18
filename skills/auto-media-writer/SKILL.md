---
name: auto-media-writer
description: Generate Chinese self-media articles for Toutiao, Baijiahao, Xiaohongshu, or Zhihu from user-provided keywords or selected hot topics. Use when Codex needs to segment Chinese keywords, research latest internet and TianAPI news with a compact source pack, collect and locally save suitable images, apply platform-by-article-type writing style profiles, fall back to image search or AI-generated raster images, call humanizer-zh before saving, validate and save a minimal article record to the auto-article backend through the fixed HTTP API/script path, or plan daily self-media article topics from hot lists across platforms.
---

# Auto Media Writer

## Iron Rules

- Never save an article to the auto-article backend before calling `humanizer-zh`. The de-AI pass is mandatory, not optional. If `humanizer-zh` is not installed, automatically search for and install it with the skill installation workflow, then continue only after it is available in the active session. If it cannot be installed or loaded, stop before database insertion and tell the user what blocked installation.
- Never save an article without images. If cited source articles contain body images, download those images to `backend/static/article-images/uploads/<yyyy>/<mm>/<task-id>/` and insert their local `/static/article-images/uploads/...` URLs into the Markdown. If an image later looks unsuitable, the user will review and replace it manually; do not omit images for that reason.
- Never invent local image paths. A Markdown image such as `![alt](/static/article-images/uploads/yyyy/mm/task/file.jpg)` is valid only after that exact file exists under the backend static uploads directory or after the backend upload API returned that URL. Do not create plausible filenames by hand.
- If source articles have no usable images or image downloads fail, search the open web for relevant images using the article keywords before falling back to AI generation. Download selected search results to the same local uploads directory and insert those local URLs into the article.
- Generate concrete raster fallback images (`png`, `jpg`, or `webp`) only after both source-image download and keyword-based internet image search fail. Never use SVG placeholders.
- If image collection, downloading, conversion, or generation is likely to take a long time, pause and tell the user what will take time, how many images are needed, and whether the likely bottleneck is download, image search, or AI generation. Ask whether to continue image processing. If the user declines, return a text draft or planning result only and do not save the article to the backend.
- Do not create article payload files unless a file is genuinely needed for audit or manual replay. Prefer passing the constructed payload JSON directly to the save/validation scripts through stdin (`-`). If a temporary payload file is unavoidable, write it outside the repository and delete it after save or dry-run validation. Never leave `payload.json`, `payload_draft.json`, or similar article payload files in the project or skill directory.

## Core Workflow

1. Resolve the target platform. Require one of `toutiao`, `baijiahao`, `xiaohongshu`, or `zhihu`; ask the user only if it is missing or ambiguous.
2. Resolve the topic path:
   - If the user provides keywords, segment them first, then research each meaningful segment independently before combining results.
   - If the user asks what to write, collect hot topics from internet search and TianAPI MCP hot lists, present 5-10 topic options, and wait for the user to choose. After the user chooses, restart the keyword segmentation and research workflow with that chosen keyword.
3. Segment keywords into at most 5 meaningful search terms: entity, event, time clue, audience intent, and controversy/angle. Keep only terms that change the search result. Example: `某明星离婚风波` -> `某明星`, `离婚`, `回应`, `网友热议`.
4. Research the main combined phrase plus the 2-3 strongest segments first. Expand to more segments only when source facts are thin, stale, or contradictory. Query TianAPI MCP only for freshness, corroboration, or hot-list planning instead of mirroring every web search. Useful TianAPI tools include hot lists (`networkhot`, `nethot`, `weibohot`, `douyinhot`, `toutiaohot`, `wxhottopic`) and news search (`allnews`, `generalnews`).
5. Merge useful information into a compact source pack. Keep 2-4 best sources total for normal articles and 5-6 only for complex society/finance topics. For each source, keep only title, URL, source, publish time, 1-2 verified facts, image URL candidates, and reliability notes. Remove duplicate, low-quality, stale, or unsupported claims. Do not invent facts, quotes, data, or dates absent from sources.
6. Classify the article automatically. Use `entertainment` for celebrity/film/TV/influencer gossip and public reactions; use `society` for broad social news; use `tech_finance` for technology, business, consumer products, AI, finance, markets, companies, or industry competition; use `knowledge` for evergreen explainers, how-to, education, culture, health, workplace, or other non-breaking explanatory topics; otherwise choose `general`. Ask the user when the category changes the writing direction and remains uncertain.
7. Select the style profile before drafting. Read `references/platform-style.md`, combine the selected platform with the selected category/article type, and follow that profile for title strategy, opening hook, paragraph rhythm, length, ending, and risk boundaries. Use the platform baseline first, then layer the article-type profile; when they conflict, platform distribution logic wins, and verified facts/legal risk boundaries override all style choices.
8. Draft the article for the selected platform and article type, integrating the source pack into a coherent self-media article instead of summarizing search results mechanically. Make the title and opening more eye-catching when the facts support it, but never turn rumors into facts or add unsupported private details.
9. Prepare images using `references/image-policy.md`. First download usable body images from cited source articles. If none are usable or downloads fail, search the internet for images using the segmented keywords and article angle, select editorially relevant results, download them to the local uploads directory, and reference their local URLs. Use AI raster generation only when both source images and keyword-based image search fail. Insert Markdown image tags at suitable positions. Do not save an image-free article.
10. Before saving, call `humanizer-zh`. Pass the drafted Markdown article plus only the selected platform, category, target length, chosen title options, compact style profile, and essential source-fact constraints. Require `humanizer-zh` to preserve verified facts, URLs, dates, names, Markdown headings, image markers, and legal risk boundaries while removing AI-sounding transitions, template paragraphs, generic moralizing, and over-balanced phrasing.
11. Create a minimal article payload object using `references/auto-article-api.md`. Prefer validating and saving that JSON through stdin (`-`) without writing a payload file. If a JSON handoff file is needed, write it outside the repository, such as an OS temp directory or `C:\tmp\auto-media-writer\`, never under `skills/auto-media-writer/` and never as a lingering `payload.json` or `payload_draft.json` in the project tree. Delete temporary payload files after saving or dry-run validation.
12. Save through the low-freedom fast path: when working inside this repository, prefer the repository-level entrypoint (`scripts/save-skill-article.ps1`, `scripts/save-skill-article.sh`, or `make saveSkillArticle`) that delegates to `scripts/save_skill_article.py`; otherwise use `scripts/save_skill_article.py` directly or make the exact documented HTTP `POST /api/v1/skill-articles`. Do not inspect backend source files, connect to MySQL, or write one-off database insertion scripts.
13. Record any improvement learned during use in this skill before future runs when the user asks the skill to evolve.

## Topic And Research Rules

- Always segment user-provided Chinese keywords before research, but cap the first search pass at one combined phrase plus 2-3 strongest segments. Expand only when necessary.
- For each segment, prefer latest credible sources. Capture concrete dates for recent events and distinguish verified facts from commentary, rumors, and platform reactions.
- For daily planning, compare hot lists across Baidu, Weibo, Douyin, WeChat, Toutiao, all-network hot lists, and broad web news. Rank options by recency, discussion heat, source availability, platform fit, image availability, and legal risk. Present 5-8 choices with one short reason for each.
- Use internet search as the primary source. Use TianAPI MCP when web results are delayed, inaccessible, or need corroboration. If TianAPI provides the freshest relevant item, use it as the lead but still seek at least one corroborating source when possible.
- Reuse the same entity and event keywords for image-search fallback. Combine the main subject with scene qualifiers such as `scene`, `editorial photo`, `event photo`, `product photo`, `press event`, `still`, or `concept art` so the searched images stay editorially relevant.
- Keep source URLs in working notes and cite them in the article when useful, but do not send bulky source metadata to the backend article record. If source URLs include usable images, prefer them unless watermarked, copyrighted, low quality, or unavailable.
- Include concrete dates when discussing recent events.

## Cost Controls

- Default source budget: 2-4 sources total, 1-2 verified facts per source, and no long pasted excerpts.
- Default search budget: one combined query plus 2-3 segment queries. Stop once the timeline, main dispute, and image path are clear.
- Default hot-topic planning budget: 5-8 options, one-line rationale, no long source pack until the user chooses a topic.
- Default image budget: one cover image and 1-2 inline images. Use 2-3 inline images only when the platform/article length benefits from them.
- Keep style profile, source constraints, and humanizer instructions compact; do not pass the full style matrix or full source pack into the humanizer.

## Topic Planning Mode

When the user asks for article ideas, daily topic planning, or says they do not know what to write:

1. Fetch current hot topics from internet search and TianAPI MCP hot-list tools when available.
2. Compare topics across platforms and remove entries that lack credible sources, are too risky, or cannot support images.
3. Return 5-8 topic options with suggested keywords, platform fit, likely category, why it is hot, source availability, and image availability.
4. Wait for the user to choose a topic or keyword. After selection, run the normal keyword segmentation, research, style-profile selection, image, humanization, and save workflow.

## Platform And Style Defaults

- Always read `references/platform-style.md` before drafting platform-specific content. It contains the authoritative platform x article-type matrix.
- Default tone level: eye-catching but compliant. Use conflict, contrast, suspense, public reaction, and timeline hooks when supported by sources; do not use clickbait that invents facts or implies unverified guilt.
- Entertainment articles default to an `entertainment` / `gossip_quick_commentary` profile: lively, gossip-aware, skeptical, and readable. Expand only when the user asks for deeper analysis or the source pack has a verified timeline requiring more context.
- Society/general news should combine hotspot commentary with ordinary-reader relevance and emotional resonance, without empty preaching.
- Toutiao prioritizes fast context, broad-reader readability, clear conflict, and interactive endings.
- Baijiahao prioritizes search-friendly structure, descriptive headings, background, timeline, and explanatory value.
- Xiaohongshu prioritizes native note rhythm, short paragraphs, personal observation, practical takeaways, and tags.
- Zhihu prioritizes a clear judgment or question, structured reasoning, causes, tradeoffs, uncertainty, and conclusion.

### Proven Toutiao Natural Commentary Pattern

Use this pattern for Toutiao society, creator-economy, consumer, and broad-interest hotspot commentary when the topic supports a human judgment rather than a dry explainer.

- Open fast with a concrete sentence that names the core event and emotion in plain words; avoid slow background setup.
- Use only a few verified facts: time, person/entity, public action, cited source, and one or two concrete details. Do not pile up source summaries.
- Let the article sound like a real editor thinking aloud: short paragraphs, direct judgments, and everyday phrases such as `真实情况常常反过来` or `这不是敬业，是损耗` when supported by the facts.
- Build the body as event -> phenomenon -> ordinary-reader relevance -> author judgment -> comment question. Keep headings clear, but avoid classroom structures such as `第一、第二、第三` unless the user asks for a checklist.
- Keep emotion restrained but present. Do not sensationalize illness, private life, tragedy, or unverified motives; move the point toward ordinary people, work rhythm, consumption decisions, or public responsibility.
- End with a discussion question close to the reader's situation instead of a grand slogan.
- Before saving, run an explicit de-AI rewrite pass that removes template transitions, numbered lecture structure, over-balanced phrasing, generic moralizing, and stiff summary paragraphs.

## Style Profile Notes

Keep a compact style profile in working notes and pass it to `humanizer-zh`. Do not save it to the backend unless the user explicitly asks for auditing metadata. The profile should include:

- `platform`: one of `toutiao`, `baijiahao`, `xiaohongshu`, `zhihu`
- `category`: selected category such as `entertainment`, `society`, `tech_finance`, `knowledge`, or `general`
- `articleType`: concise article-type key such as `gossip_quick_commentary`, `hotspot_commentary`, `plain_explainer`, or `question_answer_explainer`
- `toneLevel`: default `eye-catching-but-compliant` unless the user asks for a different tone
- `platformVoice`: short description of the selected platform baseline
- `typeVoice`: short description of the selected article-type layer
- `titleStrategy`: concrete title strategy used for the article
- `riskNotes`: legal, factual, source, or platform-risk boundaries that the draft and humanizer must preserve

## Required Output

Always generate:

- one selected `title`
- multiple `titleOptions`
- `summary`
- Markdown `markdownContent`
- `coverImageUrl`
- a compact working-note source list for factual audit, not for backend storage
- a compact working-note style profile, not for backend storage

Only save after the `humanizer-zh` rewrite pass is complete. If `humanizer-zh` is unavailable, stop before saving and tell the user to install it instead of silently skipping the de-AI pass. Never send process metadata such as `humanizeStatus`, `status`, `publishStatus`, `publishPayload`, model names, prompt versions, full source packs, or style profiles to `POST /api/v1/skill-articles`.

## Image Handling

Follow `references/image-policy.md`.

Default output requires one cover image and 1-2 inline images. Prefer downloaded local copies of source-article images. If source images are imperfect but usable, still download and insert them, and let the user replace them later. If source images are missing or fail to download, search the internet for images that match the article keywords and select editorially relevant results before using AI generation. Generate concrete raster AI images (`png`, `jpg`, or `webp`) only when both earlier paths fail, never SVG placeholders.

### Proven Toutiao Hotspot Image Pattern

Use this pattern for Toutiao hotspot commentary when source material contains public screenshots, declarations, account pages, comment sections, event photos, or product visuals.

- Prefer real source-article images over polished generic illustrations. A relevant screenshot with platform UI is usually better than a clean but empty stock-style image.
- Make the cover image immediately identify the event: account screenshot, public statement, long-post title, scene photo, product poster, or official product image.
- Match inline images to article rhythm: first image shows what happened, second shows the key statement or evidence, third shows public reaction, product detail, or useful context.
- Keep 3 images as the default for short Toutiao commentary: 1 cover plus 2 inline images. Add more only when the article length and evidence genuinely need them.
- Avoid unrelated decorative images such as generic phones, keyboards, office desks, backs of people, city night scenes, or abstract AI art when a concrete source image exists.
- For personal, health, dispute, and legal-risk topics, only use images from public reports or public account material. Do not generate fake scenes, illness visuals, paparazzi-style images, or privacy-invasive pictures.
- After downloading, inspect candidate images and reject QR codes, logos, avatars-only images, tracking assets, and unrelated recommendation graphics even if the download script accepted them.

Before doing slow image work, estimate the effort. If the workflow needs many downloads, manual source-page inspection, image conversion, or AI generation, tell the user that image handling may take significant time and ask whether to continue. If the user asks to skip images, do not save the article to the backend because image-free saves are forbidden by this skill.

## Humanizer Dependency

Before drafting final payloads, verify that `humanizer-zh` is available. If it is missing, automatically use the skill discovery/installation flow to find and install it. If installation succeeds but the active session cannot load it until Codex restarts, tell the user to restart Codex and stop before saving. Never emulate `humanizer-zh` with ad hoc rewriting when the skill is required but unavailable.

## Saving

Use the auto-article HTTP API contract in `references/auto-article-api.md`. The skill must stay portable: never require direct MySQL credentials, never hard-code secrets, and read the backend base URL from user input or environment.

### Fast Save Protocol

Treat saving as a fixed terminal API call, not as a code-discovery task:

1. Build one JSON payload object with only the fields allowed by `references/auto-article-api.md`.
2. Prefer passing the constructed JSON to the bundled scripts through stdin by using `-` as the payload argument. This avoids creating `payload.json` at all.
3. Validate the payload with `scripts/validate_skill_article_payload.py`; this also checks that `coverImageUrl` and Markdown image URLs point to real downloaded local files.
4. Save with the bundled script:

```sh
python <skill-dir>/scripts/save_skill_article.py -
```

When working inside this repository, prefer the repository-level one-click entrypoint, which delegates to the bundled script and keeps validation centralized:

```powershell
.\scripts\save-skill-article.ps1 payload.json -DryRun
.\scripts\save-skill-article.ps1 payload.json -BaseUrl http://localhost:9001
```

```sh
make validateSkillArticle PAYLOAD=/tmp/auto-media-writer/payload.json
make saveSkillArticle PAYLOAD=/tmp/auto-media-writer/payload.json BASE_URL=http://localhost:9001
```

If shell/stdin handling is impractical, use a temporary payload file outside the repository, for example `C:\tmp\auto-media-writer\<task-id>.payload.json` on Windows or `/tmp/auto-media-writer/<task-id>.payload.json` on Unix-like systems. Do not create or leave `payload.json`, `payload_draft.json`, or similar temporary article files in the project root or skill directory. After a successful save or dry-run validation, remove the temporary payload file unless the user explicitly asks to keep it for audit.

Resolve the backend base URL before saving. Prefer `AUTO_ARTICLE_BASE_URL` or the user's explicit deployment URL. Use `--base-url http://localhost:9001` only for local development, not as an assumption for deployed environments. Use `--dry-run` to validate and preview the target URL without sending the request.

Use `--static-root backend/static/article-images/uploads` only when running outside the repository root and automatic image-path discovery cannot find the uploads directory.

If the bundled script is unavailable, validate first and then send the exact HTTP request documented in `references/auto-article-api.md`. Do not create a new save script unless the bundled script itself is broken and must be repaired.

Hard stops during saving:

- Do not read backend DAO, model, service, router, generated GORM, or `.pb.go` files to infer database behavior.
- Do not connect to MySQL, require database credentials, or write SQL.
- Do not write temporary Python scripts for database insertion.
- Do not send source packs, style profiles, model metadata, process state, or publishing state to the API.
- Do not insert `/static/article-images/uploads/...` paths unless they came from a download script, AI image generation output saved to disk, or the upload-image API response and the local file exists.

For installing TianAPI MCP in another agent, use the copy-paste templates in `references/setup.md`.

## Evolution

When the user says the result is too stiff, too generic, too long, wrong platform tone, weak title, poor image choice, or otherwise not right, update the relevant style, image, API, or workflow guidance in this skill. Keep changes concise and reusable.
