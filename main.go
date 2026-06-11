package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jdcaberoy/mtgstore-backend/binder"
	"github.com/jdcaberoy/mtgstore-backend/handlers"
	"github.com/jdcaberoy/mtgstore-backend/user"
)

func main() {
	if err := Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer DB.Close()

	userSvc := user.NewUserService(DB)
	binderSvc := binder.NewBinderService(DB)

	userHandler := handlers.NewUserServiceHandler(userSvc, binderSvc)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://substationforminput-1.onrender.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "username", "token"},
		AllowCredentials: true,
	}))
	r.GET("health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	userApi := r.Group("user")
	{
		userApi.POST("/createuser", userHandler.CreateUser)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
