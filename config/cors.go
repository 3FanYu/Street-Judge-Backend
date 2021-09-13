package config

import (
	"time"
	"github.com/gin-contrib/cors"
)

func CorsConfig() cors.Config {
	config := cors.Config{
		AllowOrigins:     []string{"https://www.street-judge.club"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://www.street-judge.club"
		},
		MaxAge: 12 * time.Hour,
	}
	return config
}
