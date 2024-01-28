package controllers

import (
	"api/rakamin-api/helpers"
	"api/rakamin-api/initializers"
	"api/rakamin-api/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func Register(c *gin.Context)  {
	// Get data of req body
	var body struct {
		Username string `valid:"type(string)"`
		Email	 string `valid:"email"`
		Password string `valid:"type(string)"`
	}

	c.Bind(&body)

	if body.Username == "" || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Validate the Email
	isEmail  := govalidator.IsEmail(body.Email)
	if !isEmail {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is invalid"})
		return
	}

	// Validate the Password length
	if len(body.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password should contain minimum 6 character"})
		return
	}


	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
		return
	}

	// Create the user in the database
	user := models.User{Username: body.Username, Email: body.Email, Password: string(hashedPassword)}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "reason": result.Error.Error()})
		return
	}

	// Generate JWT token
	token, err := helpers.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// send the token back in the response
	c.JSON(http.StatusOK, gin.H{"message": "User signed up successfully", "token": token})
}

func Login(c *gin.Context)  {
	// Get data of req body
	var body struct {
		Email	 string
		Password string
	}

	c.Bind(&body)

	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}
	
	// Validate the Email
	isEmail  := govalidator.IsEmail(body.Email)
	if !isEmail {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is invalid"})
		return
	}

	// Check if the user exists by email
	var user models.User
	result := initializers.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Passwords match, generate JWT token
	token, err := helpers.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
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