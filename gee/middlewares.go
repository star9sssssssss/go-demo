package gee

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

//日志打印
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		log.Println("c.Next() 之前调用")
		// Process request
		c.Next()
		log.Println("c.Next() 之后调用")
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}


//错误处理，避免panic错误导致程序中止
func Recover() HandlerFunc {
	return func(c *Context) {
		defer func ()  {
			if err := recover(); err != nil {  //出现panic
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(c.StatusCode, "Internal Server Error")
			}
		}()
		c.Next()
	}
}


// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
