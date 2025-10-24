package main

import (
	"github.com/prabeshstha11/gotodo/db"
	"github.com/prabeshstha11/gotodo/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	defer db.DB.Close()
	db.CreateTable()

	router := gin.Default()
	router.Use(cors.Default())

	routes.SetupRoutes(router)

	router.Run(":8800")
}
