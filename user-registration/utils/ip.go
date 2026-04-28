package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Real-IP")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
		if ip != "" {
			ips := strings.Split(ip, ",")
			if len(ips) > 0 {
				ip = strings.TrimSpace(ips[0])
			}
		}
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}
