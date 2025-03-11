package cookie

import (
	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, key string, val string, ttl int) {
	println(c.Request.Host)
	c.SetCookie(key, val, ttl, "/", c.Request.Host, false, false)
}
