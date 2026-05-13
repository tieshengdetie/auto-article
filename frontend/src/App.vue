<script setup>
import { computed, onMounted, ref } from 'vue'

const platforms = [
  { type: 'nethot', name: '百度', accent: '#ff4d2d' },
  { type: 'weibohot', name: '微博', accent: '#f5a524' },
  { type: 'douyinhot', name: '抖音', accent: '#18b7b0' },
  { type: 'wxhottopic', name: '微信', accent: '#27a85a' },
  { type: 'toutiaohot', name: '头条', accent: '#d71920' }
]

const activePlatform = ref(null)
const hotSearch = ref([])
const hotLoading = ref(false)
const activeKeyword = ref('')
const searchKeyword = ref('')
const channels = ref([])
const searchMode = ref('general')
const selectedChannelId = ref(7)
const newsList = ref([])
const selectedSearchNewsKeys = ref([])
const savedNews = ref([])
const selectedNewsIds = ref([])
const newsLoading = ref(false)
const saveLoading = ref(false)
const generating = ref(false)
const userPrompt = ref('')
const llmProvider = ref('qwen')
const modelName = ref('')
const articles = ref([])
const articlePage = ref(1)
const articlePageSize = ref(6)
const articleTotal = ref(0)
const articleListMode = ref('native')
const skillPlatformFilter = ref('')
const skillArticles = ref([])
const skillArticlePage = ref(1)
const skillArticlePageSize = ref(6)
const skillArticleTotal = ref(0)
const activeArticle = ref(null)
const activeArticleType = ref('native')
const editorTextArea = ref(null)
const localImageInput = ref(null)
const publishPackage = ref(null)
const notice = ref('')
const error = ref('')

const skillPlatforms = [
  { value: '', label: '全部平台' },
  { value: 'toutiao', label: '今日头条' },
  { value: 'baijiahao', label: '百家号' },
  { value: 'xiaohongshu', label: '小红书' },
  { value: 'zhihu', label: '知乎' }
]

const selectedChannel = computed(() => {
  return channels.value.find((item) => item.channelId === Number(selectedChannelId.value))
})

const selectedSearchNews = computed(() => {
  const selected = new Set(selectedSearchNewsKeys.value)
  return newsList.value.filter((item, index) => selected.has(newsKey(item, index)))
})

const activeMarkdown = computed(() => activeArticle.value?.markdownContent || '')

const previewHasHeading = computed(() => {
  return /^#\s+/m.test(activeMarkdown.value)
})

const previewHasImage = computed(() => {
  return /!\[[^\]]*\]\([^)]+\)/.test(activeMarkdown.value)
})

async function request(path, options = {}) {
  error.value = ''
  const res = await fetch(`/api/v1${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options
  })
  const body = await res.json()
  if (!res.ok || body.code !== 1) {
    throw new Error(body.msg || '请求失败')
  }
  return body.data
}

async function loadChannels() {
  channels.value = await request('/news/channels')
}

async function loadHot(platform) {
  if (!platform) return
  activePlatform.value = platform
  hotLoading.value = true
  try {
    const data = await request('/hotSearch/', {
      method: 'POST',
      body: JSON.stringify({ type: platform.type })
    })
    hotSearch.value = data.list || []
    notice.value = `已加载 ${data.name} ${hotSearch.value.length} 条热搜`
  } catch (err) {
    error.value = err.message
  } finally {
    hotLoading.value = false
  }
}

async function searchNews(keyword) {
  const normalizedKeyword = (keyword ?? searchKeyword.value ?? '').trim()
  searchKeyword.value = normalizedKeyword
  activeKeyword.value = normalizedKeyword
  newsLoading.value = true
  newsList.value = []
  selectedSearchNewsKeys.value = []
  savedNews.value = []
  selectedNewsIds.value = []
  publishPackage.value = null
  try {
    const payload = searchMode.value === 'channel'
      ? {
          path: '/news/all',
          body: {
            col: Number(selectedChannelId.value),
            num: 10,
            page: 1,
            rand: 0,
            word: normalizedKeyword
          }
        }
      : {
          path: '/news/general',
          body: { num: 10, page: 1, rand: 0, word: normalizedKeyword }
        }
    const data = await request(payload.path, {
      method: 'POST',
      body: JSON.stringify(payload.body)
    })
    newsList.value = data.list || []
    notice.value = normalizedKeyword
      ? `已按「${normalizedKeyword}」检索到 ${newsList.value.length} 条相关新闻`
      : `已检索到 ${newsList.value.length} 条综合新闻`
  } catch (err) {
    error.value = err.message
  } finally {
    newsLoading.value = false
  }
}

async function saveNews() {
  if (selectedSearchNews.value.length === 0) return
  saveLoading.value = true
  const archiveKeyword = activeKeyword.value || '未指定关键词'
  try {
    const data = await request('/content/news/save', {
      method: 'POST',
      body: JSON.stringify({
        keyword: archiveKeyword,
        searchType: searchMode.value === 'channel' ? 'channel' : 'general',
        channelId: searchMode.value === 'channel' ? Number(selectedChannelId.value) : 0,
        channelName: searchMode.value === 'channel' && selectedChannel.value ? selectedChannel.value.channelName : '',
        list: selectedSearchNews.value
      })
    })
    savedNews.value = data.list || []
    selectedNewsIds.value = savedNews.value.slice(0, 5).map((item) => item.id)
    notice.value = `保存 ${data.saved} 条，跳过重复 ${data.skipped} 条`
  } catch (err) {
    error.value = err.message
  } finally {
    saveLoading.value = false
  }
}

async function generateArticle() {
  if (selectedNewsIds.value.length === 0) return
  const articleKeyword = activeKeyword.value || searchKeyword.value.trim() || '综合热点'
  generating.value = true
  try {
    const article = await request('/content/articles/generate', {
      method: 'POST',
      body: JSON.stringify({
        keyword: articleKeyword,
        userPrompt: userPrompt.value,
        llmProvider: llmProvider.value,
        modelName: modelName.value,
        newsIds: selectedNewsIds.value
      })
    })
    activeArticle.value = article
    activeArticleType.value = 'native'
    await loadArticles()
    notice.value = '文章草稿已生成'
  } catch (err) {
    error.value = err.message
  } finally {
    generating.value = false
  }
}

function useHotKeyword(keyword) {
  searchKeyword.value = (keyword || '').trim()
  searchNews(searchKeyword.value)
}

async function loadArticles() {
  const params = new URLSearchParams({
    page: String(articlePage.value),
    pageSize: String(articlePageSize.value)
  })
  if (activeKeyword.value) {
    params.set('keyword', activeKeyword.value)
  }
  const data = await request(`/content/articles?${params.toString()}`)
  articles.value = data.list || []
  articleTotal.value = data.total || 0
}

async function openArticle(article) {
  activeArticle.value = await request(`/content/articles/${article.id}`)
  activeArticleType.value = 'native'
  publishPackage.value = null
}

async function loadSkillArticles() {
  const params = new URLSearchParams({
    page: String(skillArticlePage.value),
    pageSize: String(skillArticlePageSize.value)
  })
  if (activeKeyword.value) {
    params.set('keyword', activeKeyword.value)
  }
  if (skillPlatformFilter.value) {
    params.set('platform', skillPlatformFilter.value)
  }
  const data = await request(`/skill-articles?${params.toString()}`)
  skillArticles.value = data.list || []
  skillArticleTotal.value = data.total || 0
}

async function openSkillArticle(article) {
  activeArticle.value = await request(`/skill-articles/${article.id}`)
  activeArticleType.value = 'skill'
  publishPackage.value = null
}

function closeArticle() {
  activeArticle.value = null
  publishPackage.value = null
}

async function saveArticle() {
  if (!activeArticle.value) return
  const path = activeArticleType.value === 'skill'
    ? `/skill-articles/${activeArticle.value.id}`
    : `/content/articles/${activeArticle.value.id}`
  activeArticle.value = await request(path, {
    method: 'PUT',
    body: JSON.stringify(activeArticle.value)
  })
  if (activeArticleType.value === 'skill') {
    await loadSkillArticles()
  } else {
    await loadArticles()
  }
  notice.value = '文章已保存'
}

function changeArticlePage(delta) {
  const nextPage = articlePage.value + delta
  const maxPage = Math.max(1, Math.ceil(articleTotal.value / articlePageSize.value))
  if (nextPage < 1 || nextPage > maxPage) return
  articlePage.value = nextPage
  loadArticles()
}

async function createPublishPackage() {
  if (!activeArticle.value) return
  const path = activeArticleType.value === 'skill'
    ? `/skill-articles/${activeArticle.value.id}/publish-package`
    : `/content/articles/${activeArticle.value.id}/publish-package`
  publishPackage.value = await request(path, {
    method: 'POST'
  })
  activeArticle.value = publishPackage.value.article
  if (activeArticleType.value === 'skill') {
    await loadSkillArticles()
  } else {
    await loadArticles()
  }
  notice.value = '待发布包已生成'
}

function changeSkillArticlePage(delta) {
  const nextPage = skillArticlePage.value + delta
  const maxPage = Math.max(1, Math.ceil(skillArticleTotal.value / skillArticlePageSize.value))
  if (nextPage < 1 || nextPage > maxPage) return
  skillArticlePage.value = nextPage
  loadSkillArticles()
}

function switchArticleMode(mode) {
  articleListMode.value = mode
  if (mode === 'skill') {
    loadSkillArticles()
  } else {
    loadArticles()
  }
}

function changeSkillPlatform() {
  skillArticlePage.value = 1
  loadSkillArticles()
}

function platformLabel(value) {
  const platform = skillPlatforms.find((item) => item.value === value)
  return platform ? platform.label : value || '未知平台'
}

function toggleNews(id) {
  if (selectedNewsIds.value.includes(id)) {
    selectedNewsIds.value = selectedNewsIds.value.filter((item) => item !== id)
  } else {
    selectedNewsIds.value = [...selectedNewsIds.value, id]
  }
}

function insertMarkdown(before, after = '') {
  if (!activeArticle.value) return
  const target = editorTextArea.value
  const content = activeArticle.value.markdownContent || ''
  if (!target) {
    activeArticle.value.markdownContent = `${content}${before}${after}`
    return
  }
  const start = target.selectionStart
  const end = target.selectionEnd
  const selected = content.slice(start, end) || '文本'
  activeArticle.value.markdownContent = `${content.slice(0, start)}${before}${selected}${after}${content.slice(end)}`
  requestAnimationFrame(() => {
    target.focus()
    target.setSelectionRange(start + before.length, start + before.length + selected.length)
  })
}

function insertImage() {
  insertRemoteImage()
}

function insertRemoteImage() {
  const url = window.prompt('远程图片地址')
  if (!url) return
  insertMarkdown(`\n![图片说明](${normalizeAssetUrl(url)})\n`, '')
}

function chooseLocalImage() {
  localImageInput.value?.click()
}

function insertLocalImage(event) {
  const file = event.target.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    insertMarkdown(`\n![${file.name}](${reader.result})\n`, '')
    event.target.value = ''
  }
  reader.readAsDataURL(file)
}

function applyHeading(level) {
  insertMarkdown(`${'#'.repeat(level)} `, '')
}

function applyFontSize(size) {
  insertMarkdown(`<span style="font-size:${size}px">`, '</span>')
}

function formatDate(value) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

function newsKey(item, index) {
  return item.id || item.url || `${item.title || 'news'}-${index}`
}

function toggleSearchNews(item, index) {
  const key = newsKey(item, index)
  if (selectedSearchNewsKeys.value.includes(key)) {
    selectedSearchNewsKeys.value = selectedSearchNewsKeys.value.filter((itemKey) => itemKey !== key)
  } else {
    selectedSearchNewsKeys.value = [...selectedSearchNewsKeys.value, key]
  }
}

function selectAllSearchNews() {
  selectedSearchNewsKeys.value = newsList.value.map((item, index) => newsKey(item, index))
}

function clearSearchNewsSelection() {
  selectedSearchNewsKeys.value = []
}

function previewMarkdown(text = '') {
  return text
    .split('\n')
    .filter((line) => line.trim())
    .map((line) => {
      let safe = line.replace(/[&<>"']/g, (char) => ({
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
      })[char])
      safe = safe.replace(/!\[([^\]]*)\]\(([^)]+)\)/g, (_match, alt, src) => `<img src="${normalizeAssetUrl(src)}" alt="${alt}" />`)
      safe = safe.replace(/&lt;span style=&quot;font-size:(\d+)px&quot;&gt;(.+?)&lt;\/span&gt;/g, '<span style="font-size:$1px">$2</span>')
      if (safe.startsWith('### ')) return `<h3>${safe.slice(4)}</h3>`
      if (safe.startsWith('## ')) return `<h2>${safe.slice(3)}</h2>`
      if (safe.startsWith('# ')) return `<h1>${safe.slice(2)}</h1>`
      return `<p>${safe}</p>`
    })
    .join('')
}

function normalizeAssetUrl(url = '') {
  const trimmed = String(url).trim()
  if (!trimmed) return ''
  if (trimmed.startsWith('data:') || /^https?:\/\//i.test(trimmed)) return trimmed
  if (trimmed.startsWith('/static/')) return trimmed
  return trimmed
}

onMounted(async () => {
  try {
    await loadChannels()
    await loadArticles()
    await loadSkillArticles()
  } catch (err) {
    error.value = err.message
  }
})
</script>

<template>
  <main class="workspace">
    <header class="topbar">
      <div class="brand-mark">
        <span>Auto Article Studio</span>
        <small>热搜 · 新闻 · 生成 · 发布包</small>
      </div>
      <div class="status">
        <span v-if="notice">{{ notice }}</span>
        <span v-if="error" class="error">{{ error }}</span>
      </div>
    </header>

    <section class="workspace-grid">
      <aside class="panel hot-panel">
        <div class="panel-head">
          <h2>热搜选题</h2>
          <span>{{ hotLoading ? '加载中' : activePlatform ? `${hotSearch.length} 条` : '请选择平台' }}</span>
        </div>
        <section class="control-strip" aria-label="平台热搜">
          <button
            v-for="platform in platforms"
            :key="platform.type"
            class="platform-button"
            :class="{ active: activePlatform && activePlatform.type === platform.type }"
            :style="{ '--accent': platform.accent }"
            @click="loadHot(platform)"
          >
            <span>{{ platform.name }}</span>
            <small>{{ platform.type }}</small>
          </button>
        </section>
        <div class="hot-scroll">
          <div class="hot-list">
            <div v-if="!activePlatform && hotSearch.length === 0" class="empty-state">
              选择一个平台后加载热搜数据
            </div>
            <article v-for="item in hotSearch" :key="item.word || item.keyword" class="hot-item">
              <div>
                <strong>{{ item.word || item.keyword || item.hotWord }}</strong>
                <div class="brief-wrap">
                  <p class="hot-brief">
                    {{ item.brief || item.hotTag || `热度 ${item.hotIndex || item.index || item.hotWordNum || '-'}` }}
                  </p>
                  <div class="brief-popover">
                    {{ item.brief || item.hotTag || `热度 ${item.hotIndex || item.index || item.hotWordNum || '-'}` }}
                  </div>
                </div>
              </div>
              <button @click="useHotKeyword(item.word || item.keyword || item.hotWord)">去检索文章</button>
            </article>
          </div>
        </div>
      </aside>

      <section class="panel news-panel">
        <div class="panel-head">
          <h2>新闻检索</h2>
          <span>{{ activeKeyword ? `关键词：${activeKeyword}` : '未指定关键词' }}</span>
        </div>
        <div class="segmented">
          <button :class="{ active: searchMode === 'general' }" @click="searchMode = 'general'">综合新闻</button>
          <button :class="{ active: searchMode === 'channel' }" @click="searchMode = 'channel'">指定频道</button>
        </div>
        <div class="search-row">
          <input
            v-model="searchKeyword"
            type="search"
            aria-label="检索关键词"
            placeholder="输入关键词；留空检索最新综合新闻"
            @keyup.enter="searchNews(searchKeyword)"
          />
          <select v-if="searchMode === 'channel'" v-model="selectedChannelId">
            <option v-for="channel in channels" :key="channel.channelId" :value="channel.channelId">
              {{ channel.channelName }}
            </option>
          </select>
          <button :disabled="newsLoading" @click="searchNews(searchKeyword)">检索</button>
        </div>
        <div class="selection-row">
          <span>已选择 {{ selectedSearchNews.length }} / {{ newsList.length }} 条新闻入库</span>
          <div>
            <button :disabled="newsList.length === 0" @click="selectAllSearchNews">全选</button>
            <button :disabled="selectedSearchNews.length === 0" @click="clearSearchNewsSelection">清空</button>
            <button :disabled="selectedSearchNews.length === 0 || saveLoading" @click="saveNews">
              保存选中 {{ selectedSearchNews.length }}
            </button>
          </div>
        </div>
        <div class="news-list">
          <article
            v-for="(item, index) in newsList"
            :key="newsKey(item, index)"
            class="news-item"
            :class="{ selected: selectedSearchNewsKeys.includes(newsKey(item, index)) }"
          >
            <label class="news-check">
              <input
                type="checkbox"
                :checked="selectedSearchNewsKeys.includes(newsKey(item, index))"
                @change="toggleSearchNews(item, index)"
              />
            </label>
            <img v-if="item.picUrl" :src="item.picUrl" alt="" />
            <div>
              <h3>{{ item.title }}</h3>
              <p>{{ item.description }}</p>
              <footer>
                <small>{{ item.source }} · {{ item.ctime }}</small>
                <a v-if="item.url" :href="item.url" target="_blank" rel="noopener noreferrer">打开原文</a>
              </footer>
            </div>
          </article>
        </div>

        <details class="native-generator">
          <summary>
            <span>文章生成</span>
            <small>{{ selectedNewsIds.length }} 条素材</small>
          </summary>
          <div class="saved-list">
            <label v-for="item in savedNews" :key="item.id" class="saved-item">
              <input type="checkbox" :checked="selectedNewsIds.includes(item.id)" @change="toggleNews(item.id)" />
              <span>{{ item.title }}</span>
            </label>
          </div>
          <div class="model-row">
            <select v-model="llmProvider" aria-label="模型提供商">
              <option value="qwen">千问</option>
              <option value="doubao">豆包</option>
            </select>
            <input v-model="modelName" aria-label="模型名称" placeholder="模型名可选，不填走默认" />
          </div>
          <label class="prompt-box">
            <span>用户提示词</span>
            <textarea
              v-model="userPrompt"
              placeholder="例如：偏财经视角，文章要有强观点，适合头条号中年读者阅读。"
              aria-label="用户提示词"
            ></textarea>
          </label>
          <button class="primary-action" :disabled="selectedNewsIds.length === 0 || generating" @click="generateArticle">
            {{ generating ? '生成中...' : '生成新文章' }}
          </button>
        </details>
      </section>

      <section class="panel article-panel">
        <div class="panel-head">
          <h2>文章中心</h2>
          <span>{{ articleListMode === 'skill' ? `${skillArticleTotal} 篇 Skill文章` : `${articleTotal} 篇大模型文章` }}</span>
        </div>

        <div class="article-history">
          <div class="article-switcher">
            <button :class="{ active: articleListMode === 'native' }" @click="switchArticleMode('native')">
              大模型文章
            </button>
            <button :class="{ active: articleListMode === 'skill' }" @click="switchArticleMode('skill')">
              Skill文章
            </button>
          </div>

          <div v-if="articleListMode === 'skill'" class="skill-filter">
            <select v-model="skillPlatformFilter" aria-label="Skill文章平台" @change="changeSkillPlatform">
              <option v-for="platform in skillPlatforms" :key="platform.value" :value="platform.value">
                {{ platform.label }}
              </option>
            </select>
          </div>

          <template v-if="articleListMode === 'native'">
            <article v-for="article in articles" :key="article.id" class="article-card" @click="openArticle(article)">
              <div>
                <strong>{{ article.title || article.keyword }}</strong>
                <p>{{ article.summary || '暂无摘要' }}</p>
              </div>
              <footer>
                <span>{{ article.keyword }}</span>
                <span>{{ article.modelName }}</span>
                <span>{{ article.status }}</span>
                <time>{{ formatDate(article.createdAt) }}</time>
              </footer>
            </article>
            <div class="article-pager">
              <button :disabled="articlePage <= 1" @click="changeArticlePage(-1)">上一页</button>
              <span>{{ articlePage }} / {{ Math.max(1, Math.ceil(articleTotal / articlePageSize)) }} · {{ articleTotal }} 篇</span>
              <button :disabled="articlePage >= Math.ceil(articleTotal / articlePageSize)" @click="changeArticlePage(1)">下一页</button>
            </div>
          </template>

          <template v-else>
            <article
              v-for="article in skillArticles"
              :key="article.id"
              class="article-card skill-card"
              @click="openSkillArticle(article)"
            >
              <div>
                <strong>{{ article.title || article.keyword }}</strong>
                <p>{{ article.summary || '暂无摘要' }}</p>
              </div>
              <footer>
                <span>{{ platformLabel(article.platform) }}</span>
                <span>{{ article.keyword }}</span>
                <span>{{ article.humanizeStatus }}</span>
                <span>{{ article.publishStatus }}</span>
                <time>{{ formatDate(article.createdAt) }}</time>
              </footer>
            </article>
            <div class="article-pager">
              <button :disabled="skillArticlePage <= 1" @click="changeSkillArticlePage(-1)">上一页</button>
              <span>{{ skillArticlePage }} / {{ Math.max(1, Math.ceil(skillArticleTotal / skillArticlePageSize)) }} · {{ skillArticleTotal }} 篇</span>
              <button
                :disabled="skillArticlePage >= Math.ceil(skillArticleTotal / skillArticlePageSize)"
                @click="changeSkillArticlePage(1)"
              >
                下一页
              </button>
            </div>
          </template>
        </div>
      </section>
    </section>

    <section v-if="activeArticle" class="editor-band">
      <div class="editor-tools">
        <input v-model="activeArticle.title" aria-label="文章标题" />
        <button @click="saveArticle">保存编辑</button>
        <button @click="createPublishPackage">生成待发布包</button>
        <button @click="closeArticle">关闭预览</button>
        <textarea v-model="activeArticle.summary" aria-label="文章摘要" placeholder="文章摘要"></textarea>
      </div>
      <div class="markdown-toolbar">
        <button @click="applyHeading(1)">H1</button>
        <button @click="applyHeading(2)">H2</button>
        <button @click="applyHeading(3)">H3</button>
        <button @click="insertRemoteImage">远程图片</button>
        <button @click="chooseLocalImage">本地图片</button>
        <button @click="applyFontSize(16)">16px</button>
        <button @click="applyFontSize(20)">20px</button>
        <button @click="applyFontSize(24)">24px</button>
        <input ref="localImageInput" class="hidden-file" type="file" accept="image/*" @change="insertLocalImage" />
      </div>
      <div class="editor-grid">
        <textarea ref="editorTextArea" v-model="activeArticle.markdownContent" aria-label="Markdown 正文"></textarea>
        <article class="preview">
          <img v-if="activeArticle.coverImageUrl && !previewHasImage" :src="normalizeAssetUrl(activeArticle.coverImageUrl)" alt="" />
          <h1 v-if="!previewHasHeading">{{ activeArticle.title }}</h1>
          <p class="summary">{{ activeArticle.summary }}</p>
          <div v-html="previewMarkdown(activeArticle.markdownContent)"></div>
        </article>
      </div>
      <div v-if="publishPackage" class="publish-box">
        <strong>{{ activeArticleType === 'skill' ? 'Skill待发布数据' : '头条号待发布包' }}</strong>
        <span>{{ activeArticleType === 'skill' ? '已写入 publishPayload，可供后续一键发布流程使用。' : publishPackage.package.remark }}</span>
      </div>
    </section>
  </main>
</template>
