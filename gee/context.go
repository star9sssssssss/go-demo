package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct { //代表上下文
	// origin objects
	Writer http.ResponseWriter
	Req *http.Request
	//request info
	Path string
	Method string
	Params map[string]string
	//response info
	StatusCode int
	//middleware 实现中间件
	handlers []HandlerFunc  //存储中间件的函数
	index int  //记录当前访问的中间件
	engine *Engine
}

//执行下一个中间件, 这样实现可以进行回调实现(即实现该函数后，继续实现该函数后面的部分)
func (c *Context) Next() {  
	c.index ++;
	totalNum := len(c.handlers)  //当前中间件总数
	for ; c.index < totalNum; c.index ++ {
		c.handlers[c.index](c)
	}
}


func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
		index: -1,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//设置状态
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//返回json数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), code)
	}
}

//返回正常数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//返回失败
func (c *Context) Fail(code int, data string) {
	c.Status(code)
	c.Writer.Write([]byte(data))
}

//返回html格式的内容
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
	// if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
	// 	c.Fail(500, err.Error())
	// }
}