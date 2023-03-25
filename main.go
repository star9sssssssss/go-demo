package main

import (
	"fmt"
	"gee"
)

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.Use(gee.Recover())
	r.GET("/hello", func(c *gee.Context) {
		c.String(200, "Hello World")
	})
	r.GET("/hh", func (c *gee.Context)  {
		stus := []string{"admin"}
		c.String(c.StatusCode, stus[100])
	})
	fmt.Println()
	r.Run(":9999")
}