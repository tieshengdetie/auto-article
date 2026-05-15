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

## Platform Baselines

### Toutiao

- Audience: broad, fast-scanning readers who expect quick context and a visible viewpoint.
- Opening: put the newest fact and central conflict in the first 100-150 Chinese characters.
- Structure: short paragraphs, direct headings, clear timeline, public reaction, and a final interaction question.
- Voice: lively and readable, with stronger hooks than Baijiahao but fewer personal notes than Xiaohongshu.
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
- Length: 500-800 Chinese characters.
- Voice: eating-melon quick commentary; lively, skeptical, and public-reaction aware.
- Title strategy: foreground the person, conflict, reversal, latest response, or netizen dispute.
- Opening: directly state the latest development and biggest hook in 100-150 characters.
- Body rhythm: `发生了什么` -> `争议点` -> `网友反应` -> `可能影响`.
- Ending: close with a clear view and a comment question.
- Risk notes: do not write private facts, relationships, anonymous screenshots, or fan claims as verified facts. Use attribution such as `据公开报道`, `从目前信息看`, or omit weak claims.

#### Society Hotspot

- Category: `society`
- Article type: `hotspot_commentary`
- Length: 1200-1800 Chinese characters.
- Voice: hotspot quick commentary with ordinary-reader relevance.
- Title strategy: show the event, contradiction, responsibility question, or public concern.
- Opening: state the event and why ordinary readers care.
- Body rhythm: timeline, responsibility boundary, public emotion, practical implication.
- Ending: give a clear but moderate viewpoint; avoid slogan-like preaching.

#### Tech And Finance

- Category: `tech_finance`
- Article type: `plain_explainer`
- Length: 1000-1600 Chinese characters.
- Voice: plain-language explanation.
- Title strategy: lead with user impact, price movement, product change, company competition, or industry consequence.
- Opening: explain the newest change and who is affected.
- Body rhythm: what changed, why it matters, who benefits/loses, what to watch next.
- Ending: concise judgment with uncertainty if market or policy facts are incomplete.

#### Knowledge Explainer

- Category: `knowledge`
- Article type: `question_answer_explainer`
- Length: 1000-1600 Chinese characters.
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
- Voice: structured news interpretation.
- Title strategy: use fact-forward wording with entity, event, and key question.
- Opening: provide the latest development, core fact, and context.
- Body rhythm: timeline, background, cause, impact, future attention point.
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
- Voice: personal observation plus information整理.
- Title strategy: show emotional relevance or everyday impact without fearmongering.
- Opening: explain what happened and why it hit ordinary people.
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
- Voice: rational social analysis.
- Title strategy: ask the responsibility, mechanism, or public-interest question.
- Opening: state the answer or core judgment before expanding.
- Body rhythm: fact boundary, responsibility, system context, public psychology, real-world impact.
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
