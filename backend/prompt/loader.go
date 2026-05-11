package prompt

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

// Loader 负责加载和渲染提示词模板
type Loader struct {
	templateDir string
	cache       map[string]*template.Template
	mu          sync.RWMutex
}

// NewLoader 创建一个新的提示词加载器
func NewLoader(templateDir string) *Loader {
	absPath, _ := filepath.Abs(templateDir)
	// 不再检查路径是否存在，延迟到实际使用时处理
	return &Loader{
		templateDir: absPath,
		cache:       make(map[string]*template.Template),
	}
}

// Load 加载并缓存一个提示词模板（如 "rag_qa" 对应 rag_qa.tmpl）
func (l *Loader) Load(name string) (*template.Template, error) {
	l.mu.RLock()
	if tmpl, ok := l.cache[name]; ok {
		l.mu.RUnlock()
		return tmpl, nil
	}
	l.mu.RUnlock()

	// 双检锁避免并发重复加载
	l.mu.Lock()
	defer l.mu.Unlock()
	if tmpl, ok := l.cache[name]; ok {
		return tmpl, nil
	}

	tmplPath := filepath.Join(l.templateDir, name+".tmpl")
	data, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(string(data))
	if err != nil {
		return nil, err
	}

	l.cache[name] = tmpl
	return tmpl, nil
}

// Render 渲染指定名称的提示词模板，并传入上下文数据
func (l *Loader) Render(name string, data interface{}) (string, error) {
	tmpl, err := l.Load(name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Preload 指定名称的模板（用于核心提示词）
func (l *Loader) Preload(names []string) error {
	for _, name := range names {
		if _, err := l.Load(name); err != nil {
			return fmt.Errorf("failed to preload prompt '%s': %w", name, err)
		}
	}
	return nil
}
