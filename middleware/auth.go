package middleware

import (
	"ai-demo/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "uid"

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ✅ 云托管调试绕过（关键）
		if c.GetHeader("x-wx-openid") != "" {
			c.Set(ContextUserIDKey, 1) // 给个默认用户
			c.Next()
			return
		}

		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少登录凭证"})
			return
		}

		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录凭证格式错误"})
			return
		}

		token := strings.TrimSpace(authHeader[7:])
		uid, err := auth.ParseToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录状态已失效"})
			return
		}

		c.Set(ContextUserIDKey, uid)
		c.Next()
	}
}
