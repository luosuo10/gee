package gee

import (
	"net/http"
)

type Router struct { 
	handlers map[string]HandlerFunc  // 每个路由地址对应一个处理函数, key的格式例如 GET-/, GET-/hello, POST-/hello
}

func newRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
	}
}

func getKey(method, pattern string) (key string) {
	return method + "-" + pattern
}

func (r *Router) addRoute(method, pattern string, handler HandlerFunc) {
	r.handlers[getKey(method, pattern)] = handler
}

func (r *Router) handle(ctx *Context) {
	key := ctx.Methond + "-" + ctx.Path
	var handler HandlerFunc
	var ok bool
	if handler, ok = r.handlers[key]; !ok {
		ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
		return
	}
	handler(ctx)
}
