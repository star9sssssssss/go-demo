package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouteGroup               //相当与继承
	router      *router       //定义路由
	groups      []*RouteGroup //所有的分组
	// htmlTemplates *template.Template  //渲染 html
	// funcMap template.FuncMap          //渲染 html
}

//根据RouterGroup进行分组
type RouteGroup struct {
	prefix      string        //前缀路由
	middlewares []HandlerFunc //中间件
	parent      *RouteGroup   //父前缀
	engine      *Engine       //所有的分组共享一个engine
}

// //初始化系统
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouteGroup = &RouteGroup{engine: engine}
	engine.groups = []*RouteGroup{engine.RouteGroup}
	return engine
}

// 进行分组
func (g *RouteGroup) Group(prefix string) *RouteGroup {
	engine := g.engine
	newGroup := &RouteGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup) //加入到engine中，一个engine操控所有与之相关的groups
	return newGroup
}

//使用中间件, 将所有的中间件函数加入其中
func (group *RouteGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//添加路由
// func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
// 	e.router.addRoute(method, pattern, handler)
// }
func (group *RouteGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// //调用Get方法
func (group *RouteGroup) GET(pattern string, hanler HandlerFunc) {
	group.addRoute("GET", pattern, hanler)
}

//调用POST方法
func (group *RouteGroup) POST(pattern string, hanler HandlerFunc) {
	group.addRoute("POST", pattern, hanler)
}

//调用PUT方法
func (group *RouteGroup) PUT(pattern string, hanler HandlerFunc) {
	group.addRoute("PUT", pattern, hanler)
}

//调用DELETE方法
func (group *RouteGroup) DELETE(pattern string, hanler HandlerFunc) {
	group.addRoute("DELETE", pattern, hanler)
}

// //开启http 服务
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

//代替系统默认的路由处理方式
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//将前缀相同即一组的中间件加入context中
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = e
	e.router.handle(c)
}

//创建静态handler处理静态文件
//relativePath 传入相对路径
// func (group *RouteGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
// 	absolutePath := path.Join(group.prefix, relativePath) //返回绝对路径
// 	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
// 	return func(c *Context) {
// 		file := c.Param("filepath")
// 		// Check if file exists and/or if we have permission to access it
// 		if _, err := fs.Open(file); err != nil {
// 			c.Status(http.StatusNotFound)
// 			return
// 		}
// 		fileServer.ServeHTTP(c.Writer, c.Req)
// 	}
// }

// serve static files
// func (group *RouteGroup) Static(relativePath string, root string) {
// 	handler := group.createStaticHandler(relativePath, http.Dir(root))
// 	urlPattern := path.Join(relativePath, "/*filepath")
// 	// Register GET handlers
// 	group.GET(urlPattern, handler)
// }

//设置funcMap
// func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
// 	e.funcMap = funcMap
// }


// //加载全局静态资源
// func (e *Engine) LoadHTMLGlob(pattern string) {
// 	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
// }
