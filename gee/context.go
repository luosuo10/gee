package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ContentType       = "Content-Type"
	ContentTypeString = "text/plain"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/HTML"
)

// 键值对的一个扩展
type H map[string]interface{}

// 将Req和ResponseWriter中常用的一部分属性抽取出来
// 提供快速构造多种类型的Response的能力, 包括String, JSON, Data, HTML
type Context struct {
	// 原始内容
	Req *http.Request
	// net/http包建议 ResponseWriter在写操作之前应该完成所有读操作，一旦执行写操作刷新header，request.body的可读性不受保证。
	Writer http.ResponseWriter

	// request related
	Path    string
	Methond string
	Params map[string]string

	// responce related
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Req:     req,
		Writer:  w,
		Path:    req.URL.Path,
		Methond: req.Method,
	}
}

func (ctx *Context) SetHeader(key string, value string) {
	ctx.Writer.Header().Set(key, value)
}

// SetStatus设置状态码
// type ResponseWriter interface { WriteHeader(statusCode int)... }
func (ctx *Context) SetStatus(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

// PostForm获取某个key对应的value[0]，如果不存在返回空
// func (r *Request) FormValue(key string) string 该key对应的value[0]
func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

// String: "text/plain"
func (ctx *Context) String(code int, format string, a ...interface{}) {
	ctx.SetHeader(ContentType, ContentTypeString)
	ctx.SetStatus(code)
	ctx.Writer.Write([]byte(fmt.Sprintf(format, a...)))
}

// JSON: "application/json"
func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader(ContentType, ContentTypeJSON)
	ctx.SetStatus(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err == nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}

// Data
func (ctx *Context) Data(code int, data []byte) {
	ctx.SetStatus(code)
	ctx.Writer.Write([]byte(data))
}

// HTML
func (ctx *Context) HTML(code int, html string) {
	ctx.SetHeader(ContentType, ContentTypeHTML)
	ctx.SetStatus(code)
	ctx.Writer.Write([]byte(html))
}
