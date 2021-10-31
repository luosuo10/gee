package gee

import (
	"fmt"
	"net/http"
)


type HandlerFunc func(http.ResponseWriter, *http.Request)  // 函数签名



/*
	
*/
type Engine struct {
	router map[string]HandlerFunc  // 每个路由地址对应一个处理函数, key的格式例如 GET-/, GET-/hello, POST-/hello
}

// New() 返回一个空的router
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	engine.router[getKey(method, pattern)] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler HandlerFunc
	var ok bool
	if handler, ok = engine.router[getKey(req.Method, req.URL.Path)]; !ok {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
		return
	}
	handler(w, req)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func getKey(method, pattern string) (key string) {
	return method + "-" + pattern
}