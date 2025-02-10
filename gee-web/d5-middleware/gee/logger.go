package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// 等待执行其他的中间件或用户的Handler
		c.Next()
		// Calculate resolution time
		latency := time.Since(t)
		// Log
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, latency)
	}
}
