package middleware

import (
	"net/http"
	"simple-blog-system/pkg/helper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-App-Id, X-Client-Id, X-Client-Version, X-Mock-Data")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		requestID := uuid.New().String()
		c.Set("requestID", requestID)
		c.Set("timeStart", time.Now().Format(time.RFC3339))
		c.Writer.Header().Set("X-Request-ID", requestID)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if len(authHeader) == 0 {
			helper.SaveAuditLog(c, http.StatusText(http.StatusUnauthorized))
			c.AbortWithStatusJSON(http.StatusForbidden, helper.Response{
				Message: http.StatusText(http.StatusUnauthorized),
				Success: false,
			})
			return
		}

		claims, err := ParseJWTToken(authHeader)
		if err != nil {
			requestID, _ := c.Get("requestID")
			helper.SaveAuditLog(c, "token has invalid claims")
			c.AbortWithStatusJSON(http.StatusForbidden, helper.Response{
				Message:   err.Error(),
				Success:   false,
				RequestId: requestID,
			})
			return
		}

		c.Set("id", claims.ID)
		c.Set("username", claims.Username)
	}
}
