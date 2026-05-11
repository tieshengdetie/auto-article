package prompt

import (
	"sync"
)

var (
	globalLoader *Loader
	once         sync.Once
	initErr      error
)

// InitGlobal 初始化全局提示词加载器
func InitGlobal(templateDir string) error {
	once.Do(func() {
		globalLoader = NewLoader(templateDir)
		// 可选：预加载核心模板
		// corePrompts := []string{"chat", "rag_qa"}
		// globalLoader.Preload(corePrompts)
	})
	return initErr
}

// Global 返回全局 Loader 实例
func Global() *Loader {
	return globalLoader
}
