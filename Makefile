PYTHON ?= python3
PAYLOAD ?=
BASE_URL ?= $(AUTO_ARTICLE_BASE_URL)
TIMEOUT ?= 30
STATIC_ROOT ?=
DRY_RUN ?=

AUTO_MEDIA_WRITER_SAVE := skills/auto-media-writer/scripts/save_skill_article.py

SAVE_ARGS = $(PAYLOAD) $(if $(strip $(BASE_URL)),--base-url $(BASE_URL),) $(if $(strip $(TIMEOUT)),--timeout $(TIMEOUT),) $(if $(strip $(STATIC_ROOT)),--static-root $(STATIC_ROOT),) $(if $(strip $(DRY_RUN)),--dry-run,)

.PHONY: saveSkillArticle validateSkillArticle

saveSkillArticle:
	$(PYTHON) $(AUTO_MEDIA_WRITER_SAVE) $(SAVE_ARGS)

validateSkillArticle: DRY_RUN=1
validateSkillArticle: saveSkillArticle
