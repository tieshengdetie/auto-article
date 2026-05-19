# Platform Style

## Style Selection

Build every article style from two layers:

1. Select the platform baseline: `toutiao`, `baijiahao`, `xiaohongshu`, or `zhihu`.
2. Select the article-type layer from the matrix: entertainment gossip, society hotspot, tech/finance explainer, or knowledge explainer.

When the platform and article type conflict, prioritize platform distribution logic and reader expectation. Verified facts, source limits, legal risk, and platform safety always override style. Use an eye-catching but compliant tone by default: make the conflict, contrast, suspense, public reaction, timeline, or reader impact visible, but never turn rumor into fact or imply unsupported guilt.

Record the chosen profile in `styleProfile` with `platform`, `category`, `articleType`, `toneLevel`, `platformVoice`, `typeVoice`, `titleStrategy`, and `riskNotes`.

## Shared Structure

Use this baseline unless the selected profile overrides it:

1. Title
2. Opening hook that states the event and why readers care
3. Background and timeline
4. 3-5 second-level headings
5. Analysis of causes, conflict, public emotion, and likely impact
6. Closing viewpoint and interaction question

Generate 3-5 title options for Toutiao, Baijiahao, and Zhihu. Generate 5-8 title options for Xiaohongshu.

## Eye-Catching But Compliant Titles

Hard length limit:

- Every selected `title` and every item in `titleOptions` must be no more than 30 Chinese characters. Count Chinese punctuation, Arabic numerals, English letters, and spaces toward the limit. If a high-read pattern would exceed 30 characters, compress the second clause instead of keeping the full template.

Allowed title hooks:

- Conflict: public response, disagreement, turning point, or visible contradiction.
- Contrast: before/after, expectation/reality, official response/public reaction.
- Suspense: what changed, what is still unclear, why people are arguing.
- Timeline: latest update, key moment, response after controversy.
- Reader impact: what ordinary users, fans, consumers, parents, workers, or investors should notice.

Forbidden title moves:

- State rumors, screenshots, private-life speculation, or anonymous claims as verified facts.
- Invent quotes, motives, relationships, money amounts, dates, medical details, or legal conclusions.
- Use insulting labels, vulgar innuendo, defamation, guaranteed outcomes, or guilt-by-association.
- Overpromise exclusive access with wording such as `实锤`, `内幕`, `彻底凉了`, or `全网封杀` unless a reliable source explicitly supports the claim and the article preserves that attribution.

## Proven High-Read Title Patterns

Use these patterns especially for Toutiao technology, business, consumer, society, and broad-interest hotspot articles. They are distilled from high-read Chinese self-media titles and should make the title feel like a concrete public event, not a generic report.

Core style:

- Lead with a named person, company, platform, product, region, or group. Avoid vague starts such as `近日`, `有消息称`, or `行业迎来变化` when a concrete entity is available.
- Put the sharpest hook in the first 12-18 Chinese characters: firing, fine, loss, closure, refund, delisting, outage, bankruptcy risk, price plunge, public response, product release, or responsibility dispute.
- Prefer a two-part headline joined by `，`, `：`, `；`, or `！`. The first part states the event; the second part gives the conflict, consequence, reversal, or human judgment.
- Use verified numbers when available: age, days, percentages, rankings, money amounts, store counts, employee counts, market value, sales, or time spans. Do not invent or round numbers that sources do not support.
- Make technology and finance titles human-visible. Translate product, capital, market, or company changes into effects on workers, consumers, investors, merchants, users, suppliers, or competitors.
- Use contrast and reversal: `不是X，而是Y`, `原以为X，事实却是Y`, `刚刚X，又Y`, `十年后X`, `从A到B`. Keep the contrast factual.
- Use public quotes only when the wording is sourced. Quote-style titles can work well when the quote is short, emotional, and attached to a named speaker.
- Choose strong but defensible verbs: `下架`, `关闭`, `罚`, `丢失`, `暴跌`, `退场`, `拒绝`, `回应`, `打压`, `转移`, `拉上`, `拿到`, `撑起`. Avoid guilt-implying verbs unless official facts support them.
- Keep Toutiao-style news commentary titles around 20-30 Chinese characters. Never exceed 30 characters; when the fact hook is long, keep the entity and consequence first and cut filler from the second clause.

Reusable headline skeletons:

- `“关键原话”，人物/公司回应：后半句`
- `人物/公司：现在最担心X；她/他/它曾说Y`
- `年过X、连续X天、丢失X%，人物/公司一句话点破`
- `公司/平台X年后做出Y，背后的两股势力浮出水面`
- `从A到B，人物/公司走过的路：真正难的是X`
- `X发布/下架/关闭/罚款，Y却把问题推到台前`
- `不是X，而是Y：这件事让谁最难受`
- `X占营收Y%，半年倒闭Z家：行业问题藏不住了`
- `这就是传说中“X”的Y，被A和B同时看上了`

Selection rules:

- For each article, generate at least one fact-forward title and at least one conflict/reversal title when sources support both.
- For Toutiao, prioritize concrete entity + consequence + public emotion. The title may be punchy, but it must stay inside verified facts.
- For Baijiahao, preserve searchable keywords while borrowing the high-read pattern; do not sacrifice clarity for suspense.
- For Xiaohongshu, soften hard news titles into everyday relevance, such as `这事和普通人有什么关系` or `别只看热闹，真正影响在这里`.
- For Zhihu, convert the same material into a question or judgment, such as `为什么X会走到这一步？` or `X事件真正值得讨论的是什么？`.

## Platform Baselines

### Toutiao

- Audience: broad, fast-scanning readers who expect quick context and a visible viewpoint.
- Opening: put the newest fact and central conflict in the first 100-150 Chinese characters.
- Structure: short paragraphs, direct headings, clear timeline, public reaction, and a final interaction question.
- Voice: lively and readable, with stronger hooks than Baijiahao but fewer personal notes than Xiaohongshu.
- Recommendation focus: pass review cleanly, avoid duplicate-looking title/body/cover combinations, win the first recommendation batch with click reason, read-through momentum, comments, favorites, and shares.
- Avoid: vague moralizing, slow background-first openings, and extreme shock/fear wording.

### Baijiahao

- Audience: search-driven readers who need background, timeline, and explanatory value.
- Opening: name the entity, event, latest development, and why it matters.
- Structure: descriptive headings, fact-forward timeline, background context, and broader meaning.
- Voice: restrained, explanatory, and slightly formal without becoming stiff.
- Avoid: overly colloquial wording, unsupported emotion, and headings that hide the main fact.

### Xiaohongshu

- Audience: note-style readers who prefer fast emotional entry, short paragraphs, and usable takeaways.
- Opening: first screen should answer `这事为什么值得看` or `和我有什么关系`.
- Structure: short paragraphs, compact lists when helpful, native tags at the end, and a comment-friendly closing.
- Voice: conversational, observant, and clean. Use emoji only if the user asks for them.
- Avoid: stiff transitions such as `综上所述`, `值得注意的是`, `从某种程度上来说`, and long news-report paragraphs.

### Zhihu

- Audience: readers who expect a judgment, question framing, evidence, and reasoning.
- Opening: begin with a clear answer, judgment, or question, then define what is known and unknown.
- Structure: fact boundary, dispute points, cause analysis, tradeoffs, and conclusion.
- Voice: reasoned and analytical, with less emotional push than Toutiao or Xiaohongshu.
- Avoid: pretending certainty when sources are incomplete, and writing a lively gossip tone without analysis.

## Style Matrix

### Toutiao

#### Entertainment Gossip

- Category: `entertainment`
- Article type: `gossip_quick_commentary`
- Length: 500-800 Chinese characters; hard maximum 1000 Chinese characters.
- Voice: eating-melon quick commentary; lively, skeptical, and public-reaction aware.
- Title strategy: foreground the person, conflict, reversal, latest response, or netizen dispute.
- Opening: directly state the latest development and biggest hook in 100-150 characters.
- Body rhythm: `发生了什么` -> `争议点` -> `网友反应` -> `可能影响`.
- Ending: close with a clear view and a comment question.
- Risk notes: do not write private facts, relationships, anonymous screenshots, or fan claims as verified facts. Use attribution such as `据公开报道`, `从目前信息看`, or omit weak claims.

#### Society Hotspot

- Category: `society`
- Article type: `hotspot_commentary`
- Length: 800-1000 Chinese characters; hard maximum 1000 Chinese characters.
- Voice: social livelihood commentary with ordinary-reader relevance; treat the account as focused on民生纠纷、劳动职场、家庭养老、消费住房、教育医疗, and practical risk avoidance rather than broad social miscellany.
- Title strategy: show the event, contradiction, responsibility question, public concern, or practical consequence for workers, families, consumers, elderly people, tenants, owners, parents, or ordinary residents.
- Opening: state the event and immediately answer why ordinary readers should care.
- Body rhythm: event timeline, who may face the same problem, responsibility or rule boundary, practical risk/reminder, public emotion.
- Ending: give a clear but moderate viewpoint; avoid slogan-like preaching.
- Recommendation tactics: avoid copying mainstream media's same timeline angle; lead with a concrete stakeholder and consequence, then add a rule, responsibility, cost, or避坑 angle that makes the piece distinct enough to survive duplicate reduction.

#### Tech And Finance

- Category: `tech_finance`
- Article type: `plain_explainer`
- Length: 800-1000 Chinese characters; hard maximum 1000 Chinese characters.
- Voice: plain-language explanation.
- Title strategy: lead with user impact, price movement, product change, company competition, or industry consequence.
- Opening: explain the newest change and who is affected.
- Body rhythm: what changed, why it matters, who benefits/loses, what to watch next.
- Ending: concise judgment with uncertainty if market or policy facts are incomplete.

#### Knowledge Explainer

- Category: `knowledge`
- Article type: `question_answer_explainer`
- Length: 800-1000 Chinese characters; hard maximum 1000 Chinese characters.
- Voice: direct problem-solving explanation.
- Title strategy: ask or answer a concrete reader question.
- Opening: state the common confusion and short answer.
- Body rhythm: answer, examples, common mistakes, final takeaway.
- Ending: simple judgment or practical suggestion.

### Baijiahao

#### Entertainment Gossip

- Category: `entertainment`
- Article type: `entertainment_news_explainer`
- Length: 700-1000 Chinese characters.
- Voice: entertainment news interpretation, more explanatory than gossipy.
- Title strategy: include searchable entity and event keywords; avoid vague emotional hooks.
- Opening: state who, what, latest response, and why the topic is being searched.
- Body rhythm: background, timeline, public response, verified risk boundary.
- Ending: summarize what is confirmed and what still needs follow-up.

#### Society Hotspot

- Category: `society`
- Article type: `structured_news_explainer`
- Length: 1500-2200 Chinese characters.
- Voice: structured social livelihood interpretation with practical value for ordinary readers.
- Title strategy: use fact-forward wording with entity, event, affected group, and key responsibility or rule question.
- Opening: provide the latest development, core fact, context, and why this matters to ordinary families, workers, consumers, or residents.
- Body rhythm: timeline, background, responsibility/rule boundary, ordinary-reader impact, future attention point.
- Ending: restrained commentary tied to facts.

#### Tech And Finance

- Category: `tech_finance`
- Article type: `search_friendly_explainer`
- Length: 1500-2200 Chinese characters.
- Voice: search-friendly explanation.
- Title strategy: include core entity, product/company/market keyword, and reader question.
- Opening: define the change or event in clear terms.
- Body rhythm: definition, background, impact, trend, uncertainty.
- Ending: summarize what readers should remember.

#### Knowledge Explainer

- Category: `knowledge`
- Article type: `evergreen_explainer`
- Length: 1500-2200 Chinese characters.
- Voice: encyclopedia-like but readable.
- Title strategy: fit `是什么`, `为什么`, or `怎么办` searches.
- Opening: answer the main question first.
- Body rhythm: concept, reason, example, method, reminder.
- Ending: short recap and practical note.

### Xiaohongshu

#### Entertainment Gossip

- Category: `entertainment`
- Article type: `sober_gossip_note`
- Length: 400-700 Chinese characters.
- Voice: clear-headed gossip note; conversational, fast, and skeptical.
- Title strategy: use a sharp but clean hook around `这事为什么吵起来`, `反转点`, `路人观感`, or `别急着站队`.
- Opening: first screen explains why the topic is noisy and what changed.
- Body rhythm: short paragraphs, 3-5 compact points, observation before judgment.
- Ending: comment-friendly question plus 3-6 topic tags.
- Risk notes: do not pretend to have insider information; separate public facts, fan interpretation, and personal observation.

#### Society Hotspot

- Category: `society`
- Article type: `personal_observation_digest`
- Length: 600-900 Chinese characters.
- Voice: personal observation plus social livelihood information整理.
- Title strategy: show emotional relevance, everyday impact, or a concrete避坑 reminder without fearmongering.
- Opening: explain what happened and why it may hit ordinary people in work, family, housing, consumption, healthcare, education, or retirement.
- Body rhythm: facts, personal observation, practical reminder, comment prompt.
- Ending: fewer grand conclusions; more grounded empathy or caution.

#### Tech And Finance

- Category: `tech_finance`
- Article type: `life_relevance_explainer`
- Length: 600-900 Chinese characters.
- Voice: everyday-life explanation.
- Title strategy: lead with `和我有什么关系`, price/usefulness, product experience, or consumer decision.
- Opening: translate the event into everyday impact.
- Body rhythm: what changed, what it affects, what to compare, what to watch.
- Ending: practical takeaway and tags.

#### Knowledge Explainer

- Category: `knowledge`
- Article type: `saveable_note`
- Length: 600-1000 Chinese characters.
- Voice: save-worthy note.
- Title strategy: use a question, checklist, common misunderstanding, or quick guide.
- Opening: tell readers what they will understand or solve.
- Body rhythm: concise points, examples, mistakes, summary.
- Ending: collection-friendly recap and topic tags.

### Zhihu

#### Entertainment Gossip

- Category: `entertainment`
- Article type: `opinion_analysis`
- Length: 900-1400 Chinese characters.
- Voice: opinion analysis rather than gossip performance.
- Title strategy: frame the real question behind the public argument.
- Opening: give a judgment first, then define what is confirmed and what is uncertain.
- Body rhythm: facts, dispute, public opinion, industry logic, conclusion.
- Ending: explain why the topic triggered discussion and what should not be overread.
- Risk notes: cite public information, avoid private speculation, and mark incomplete evidence clearly.

#### Society Hotspot

- Category: `society`
- Article type: `rational_hotspot_analysis`
- Length: 1500-2500 Chinese characters.
- Voice: rational social livelihood analysis with clear fact boundaries and ordinary-reader relevance.
- Title strategy: ask the responsibility, mechanism, public-interest, or ordinary-person rights question.
- Opening: state the answer or core judgment before expanding, then define how the issue connects to work, consumption, housing, family, healthcare, education, retirement, or public services.
- Body rhythm: fact boundary, responsibility/rule boundary, system context, ordinary-reader impact, public psychology, practical takeaway.
- Ending: balanced conclusion with concrete limits.

#### Tech And Finance

- Category: `tech_finance`
- Article type: `causal_chain_analysis`
- Length: 1800-2800 Chinese characters.
- Voice: causal-chain analysis.
- Title strategy: focus on why the event happened and what it changes.
- Opening: answer whether the event is important and for whom.
- Body rhythm: mechanism, business logic, market or policy context, long-term impact, uncertainty.
- Ending: clear conclusion plus what evidence would change the view.

#### Knowledge Explainer

- Category: `knowledge`
- Article type: `question_answer_argument`
- Length: 1500-2500 Chinese characters.
- Voice: question-answer reasoning.
- Title strategy: use a direct question that matches user curiosity.
- Opening: answer the question first.
- Body rhythm: thesis, evidence, counterexample, explanation, conclusion.
- Ending: concise view and usable takeaway.

## Humanizing Pass

Rewrite after the first draft:

- Keep the selected platform and article-type profile intact.
- Replace generic transitions with natural sentence-level connections.
- Cut repetitive summaries and obvious filler.
- Vary paragraph length according to the selected platform.
- Add concrete nouns and verbs from the source pack.
- Keep a human editorial stance, but do not add unsupported facts.
- Preserve legal risk boundaries and uncertainty markers.
- Remove AI tells: `在当今社会`, `不可忽视的是`, `这无疑`, repeated `引发了广泛关注`, and `我们应该理性看待` as a lazy ending.
