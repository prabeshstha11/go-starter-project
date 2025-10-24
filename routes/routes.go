package routes

import (
	"github.com/prabeshstha11/gotodo/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/status", func(c *gin.Context) {
		c.String(200, "Status OK")
	})

	router.POST("/create", handlers.CreateTodo)
	router.GET("/todo", handlers.GetTodos)
	router.GET("/todo/:id", handlers.GetTodoByID)
	router.PATCH("/todo/:id", handlers.UpdateTodo)
	router.DELETE("/todo/:id", handlers.DeleteTodo)
}
