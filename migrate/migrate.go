package main

import (
	"api/rakamin-api/initializers"
	"api/rakamin-api/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Photo{})
}