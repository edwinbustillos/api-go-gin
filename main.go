package main

import (
	"github.com/edwinbustillos/api-go-gin/database"
	"github.com/edwinbustillos/api-go-gin/handlers"
	"github.com/edwinbustillos/api-go-gin/middlewares"
	env "github.com/edwinbustillos/api-go-gin/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	env.Load()
	port := env.GetEnv("PORT_API", "")
	if port == "" {
		port = "8080"
	}
	db := database.InitDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/token", handlers.GetToken)
	r.GET("/refreshToken", handlers.RefreshToken)

	r.Use(middlewares.JWTMiddleware())

	r.GET("/items", handlers.GetItems)
	r.GET("/item/:id", handlers.GetItem)
	r.POST("/items", handlers.CreateItem)
	r.PUT("/items/:id", handlers.UpdateItem)
	r.DELETE("/items/:id", handlers.DeleteItem)

	r.Run(":" + port)
}
