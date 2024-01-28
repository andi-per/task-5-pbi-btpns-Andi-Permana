package controllers

import (
	"api/rakamin-api/initializers"
	"api/rakamin-api/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(c *gin.Context)  {
	// Get data of req body
	var body struct {
		Username string
		Email	 string
		Password string
	}

	c.Bind(&body)

	// Create new user
	user := models.User{Username: body.Username, Email: body.Email, Password: body.Password}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
		"error" : "Bad request",
		})
		return
	}

	// return user id
	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

func Login(c *gin.Context)  {
	// Get data of req body
	var body struct {
		Email	 string
		Password string
	}

	c.Bind(&body)
	
	// user := models.User{Username: body.Username, Email: body.Email, Password: body.Password}

	// result := initializers.DB.Create(&user) // pass pointer of data to Create

	// if result.Error != nil {
	// 	c.JSON(400, gin.H{
	// 	"error" : "Bad request",
	// 	})
	// 	return
	// }

	// return user id
	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
	})
}


// // dummy controller
func GetUser(c *gin.Context)  {
	var users []models.User
	initializers.DB.Find(&users)

	c.JSON(201, gin.H{
		"users": users,
	})
}

// // dummy controller
// func GetSingleUser(c *gin.Context)  {
// 	// GET query params id
// 	id := c.Param("id")

// 	var user models.User
// 	initializers.DB.Find(&user, id)

// 	c.JSON(201, gin.H{
// 		"user": user,
// 	})
// }

func UpdateUser(c *gin.Context)  {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find the user by ID in the database
	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		// Check if the user with the given ID was not found
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Handle other database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Get data of req body
	var body struct {
		Username string
	}

	c.Bind(&body)

	// Now you have the user and can update the fields
	//update the user's Username
	newUsername := body.Username
	if newUsername != "" {
		user.Username = newUsername
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
	}

	// Update the user in the database
	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

func DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find the user by ID in the database
	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		// Check if the user with the given ID was not found
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Handle other database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Delete the user from the database
	initializers.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}