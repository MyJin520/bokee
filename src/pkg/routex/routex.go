package routex

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Group 内嵌 gin.RouterGroup 的路由组包装：注册路由的同时记录路由描述，
// 调用方式与 gin 原生一致，仅将描述作为第二个参数：g.POST("/x", "描述", handler)。
type Group struct {
	*gin.RouterGroup
}

var (
	mu      sync.RWMutex
	descMap = make(map[string]string)
)

func key(method, path string) string {
	return method + "|" + path
}

// joinPath 拼接路由组前缀与相对路径，避免双斜杠
func joinPath(base, relative string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(relative, "/")
}

// Set 记录路由描述（启动注册阶段调用，后写覆盖）
func Set(method, path, desc string) {
	mu.Lock()
	defer mu.Unlock()
	descMap[key(method, path)] = desc
}

// Get 获取路由描述，未注册返回空串
func Get(method, path string) string {
	mu.RLock()
	defer mu.RUnlock()
	return descMap[key(method, path)]
}

func (g *Group) POST(path, desc string, handler gin.HandlerFunc) {
	g.RouterGroup.POST(path, handler)
	Set(http.MethodPost, joinPath(g.BasePath(), path), desc)
}

func (g *Group) GET(path, desc string, handler gin.HandlerFunc) {
	g.RouterGroup.GET(path, handler)
	Set(http.MethodGet, joinPath(g.BasePath(), path), desc)
}

func (g *Group) PUT(path, desc string, handler gin.HandlerFunc) {
	g.RouterGroup.PUT(path, handler)
	Set(http.MethodPut, joinPath(g.BasePath(), path), desc)
}

func (g *Group) DELETE(path, desc string, handler gin.HandlerFunc) {
	g.RouterGroup.DELETE(path, handler)
	Set(http.MethodDelete, joinPath(g.BasePath(), path), desc)
}
