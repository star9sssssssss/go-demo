package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
	roots map[string]*node
}

func newRouter() *router {
	return &router{
		make(map[string]HandlerFunc),
		make(map[string]*node),
	}
}

//解析路由
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}


func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	serachPath := parsePattern(path) //获得用户输入的访问路径
	params := make(map[string]string)  //返回模糊匹配的键值对
	root, ok := r.roots[method]
	
	if !ok {
		return nil, nil
	}

	n := root.search(serachPath, 0)
	if n != nil {
		//获得trie树上的路由
		parts := parsePattern(n.pattern)

		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = serachPath[index]
			}
			if part[0] == '*' && len(parts) > 1 {
				params[part[1:]] = strings.Join(serachPath[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

//处理请求
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)	
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		//将handle加入到其中，在处理完所有中间件的第一阶段时，运行该函数
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}