package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"user-registration/config"
	"user-registration/controllers"
	"user-registration/models"
)

func main() {
	config.InitConfig()
	config.InitMySQL()

	config.InitRedis()

	if err := config.DB.AutoMigrate(&models.User{}, &models.RegisterLog{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api/v1")
	{
		api.POST("/register", controllers.Register)
	}

	port := config.AppConfig.Server.Port
	if port == 0 {
		port = 8080
	}

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
