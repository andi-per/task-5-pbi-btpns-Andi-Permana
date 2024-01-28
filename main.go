package main

import (
	"api/rakamin-api/initializers"
	"api/rakamin-api/routes"

	"github.com/gin-gonic/gin"
)


func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}


func main() {
	r := gin.Default()

	// Auth
	routes.SetupUserRoutes(r)

	// Photos
	routes.SetupPhotoRoutes(r)

	r.Run()
}