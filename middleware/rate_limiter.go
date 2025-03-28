package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

var limiterCache = cache.New(1*time.Minute, 2*time.Minute)

// LimitRequest middleware membatasi jumlah request per IP
func LimitRequest(limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rl:" + ip

		val, found := limiterCache.Get(key)
		if !found {
			// Belum pernah request → set ke 1
			limiterCache.Set(key, 1, duration)
		} else {
			// Sudah pernah request → cek count
			count := val.(int)
			if count >= limit {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Too many requests. Please try again later.",
				})
				c.Abort()
				return
			}
			limiterCache.Set(key, count+1, duration)
		}

		c.Next()
	}
}
