package gee

import (
	"net/http"
	"strings"
)

// 路由实现方案
// 方案1. map[string]HandlerFunc
// 哈希表的形式存储 静态路由 和 处理函数 的键值对。每个路由地址对应一个处理函数, key的格式例如 GET-/, GET-/hello, POST-/hello。
// type Router struct {
// 	handlers map[string]HandlerFunc
// }

// 方案1只能支持静态路由，无法支持动态路由。
// 动态路由的定义：一条路由规则可以匹配某个类型的路由而非某一条特定路由。例如/hello/:name，可以匹配/hello/geektutu、hello/jack等。

// 因此引入方案2. 字典树 以 / 为分隔符，作为前缀树的节点
// 不同的实现方式可以支持不同的路由参数，此处实现以下路由参数
// 		- 【:】  例如 /p/:lang/doc，可以匹配 /p/c/doc 和 /p/go/doc。
//		- 【*】    例如 /static/*filepath，可以匹配/static/fav.ico，也可以匹配/static/js/jQuery.js，这种模式常用于静态服务器，能够递归地匹配子路径。

// 字典树实现
type Node struct {
	//
	pattern string
	//
	part string
	// 子节点
	children []*Node
	// 当part中含有 : 或者 * 时，为true；否则为false。
	isWild bool
}

// matchChild 返回第一个匹配的子节点; 若无匹配返回nil。
func (node *Node) matchChild(part string) *Node {
	for _, child := range node.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChild 返回所有匹配的子节点; 若无匹配返回空数组。
func (node *Node) matchChildren(part string) []*Node {
	matchChildren := make([]*Node, 0)
	for _, child := range node.children {
		if child.part == part || child.isWild {
			matchChildren = append(matchChildren, child)
		}
	}
	return matchChildren
}

func (node *Node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		node.pattern = pattern
		return
	}

	part := parts[height]
	var child *Node
	if child = node.matchChild(part); child == nil {
		// 如果该节点没有与之匹配的子节点，则新建一个
		child = &Node{
			part:   part,
			isWild: isWild(part),
		}

		node.children = append(node.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (node *Node) search(parts []string, height int) *Node {
	if len(parts) == height || strings.HasPrefix(node.part, "*") {
		if node.part == "" {
			return nil
		}
		return node
	}

	part := parts[height]
	children := node.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func isWild(part string) bool {
	if part[0] != ':' && part[0] != '*' {
		return false
	}
	return true
}

// 方案2 基于字典树实现的路由
type Router struct {
	// 每一个请求方式对应一棵字典树，此处保存所有字典树的根节点
	// roots["GET"], roots["POST"]
	roots map[string]*Node
	//
	handlers map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{
		roots:    make(map[string]*Node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 将pattern string解析成 parts []string
//
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)

	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 已经有通配符*了，不需要继续
			if item[0] == '*' {
				return parts
			}
		}
	}
	return parts
}

func getKey(method, pattern string) (key string) {
	return method + "-" + pattern
}

func (r *Router) addRoute(method, pattern string, handler HandlerFunc) {
	// 如果该方法根节点不存在, 新增新的方法根节点
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &Node{}
	}
	parts := parsePattern(pattern)
	r.roots[method].insert(pattern, parts, 0)

	r.handlers[getKey(method, pattern)] = handler
}

// getRoute
// 返回method对应的根节点以及pattern对应的
func (r *Router) getRoute(method, path string) (root *Node, params map[string]string) {
	var ok bool
	// 如果不存在对应method的字典树 返回空
	if root, ok = r.roots[method]; !ok {
		return nil, nil
	}
	searchParts := parsePattern(path)
	// 如果不存在对应的路由节点
	var matchNode *Node
	if matchNode = root.search(searchParts, 0); matchNode == nil {
		return nil, nil
	}

	parts := parsePattern(matchNode.pattern)
	params = make(map[string]string)
	for idx, part := range parts {
		// TODO 这里不可能为null吗？
		switch part[0] {
		case ':':
			params[part[1:]] = searchParts[idx]
		case '*':
			params[part[1:]] = strings.Join(searchParts[idx:], "/")
		}
	}
	return matchNode, params
}

func (r *Router) handle(ctx *Context) {
	node, params := r.getRoute(ctx.Methond, ctx.Path)
	if node == nil {
		ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
		return
	}

	ctx.Params = params
	key := ctx.Methond + "-" + ctx.Path
	r.handlers[key](ctx)
}
