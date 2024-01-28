package controllers

import (
	"api/rakamin-api/initializers"
	"api/rakamin-api/models"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type photo struct {
	ID		  uint      
	Title     string    
	Caption   string
	PhotoURL  string    
	UserID    uint      
}

func PostPhoto(c *gin.Context)  {
	var body struct {
		Title string `valid:"type(string)"`
		Caption	 string `valid:"-"`
		PhotoURL string `valid:"type(string)"`
	}

	c.Bind(&body)

	if body.Title == "" || body.PhotoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title & PhotoURL fields are required"})
		return
	}

	// Validate the PhotoURL
	isURL  := govalidator.IsURL(body.PhotoURL)
	if !isURL {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PhotoURL is invalid"})
		return
	}

	// Retrieve user ID from the context
	userID, exists := c.Get("userID")

	if !exists {
		// Handle case where user ID is not set in the context
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information not available"})
		return
	}

	parsedUserID := uint(userID.(float64))
	
	photo := models.Photo{Title: body.Title, Caption: body.Caption, PhotoURL: body.PhotoURL, UserID: parsedUserID }

	result := initializers.DB.Create(&photo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create photo", "reason": result.Error.Error()})
		return
	}

	// send the token back in the response
	c.JSON(http.StatusOK, gin.H{"message": "Photo is uploaded successfully"})
}

func GetPhotos(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		// Handle case where user ID is not set in the context
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not available"})
		return
	}

	parsedUserID := uint(userID.(float64))

	var photos []photo
	result := initializers.DB.Where("user_id = ?", parsedUserID).Find(&photos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve photos"})
		return
	}

	c.JSON(http.StatusOK, photos)
}

func UpdatePhoto(c *gin.Context)  {
	photoID, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get data of req body
	var body struct {
		Title string 
		Caption	 string
		PhotoURL string
	}

	c.Bind(&body)

	if body.Title == "" || body.PhotoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title & PhotoURL fields are required"})
		return
	}

	// Validate the PhotoURL
	isURL  := govalidator.IsURL(body.PhotoURL)
	if !isURL {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PhotoURL is invalid"})
		return
	}

	// Retrieve user ID from the context
	tokenUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not available"})
		return
	}

	// Query the database to check if the photo exists and belongs to the user
	var photo photo
	result := initializers.DB.Where("id = ? AND user_id = ?", photoID, tokenUserID).First(&photo)
	if result.Error != nil {
		// Handle case where the photo does not exist or does not belong to the user
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found or does not belong to the user"})
		return
	}

	photo.Title = body.Title
	photo.Caption = body.Caption
	photo.PhotoURL = body.PhotoURL
	// Update the user in the database
	initializers.DB.Save(&photo)

	c.JSON(http.StatusOK, gin.H{"message": "Photo updated successfully", "photo": photo})
}

func DeletePhoto(c *gin.Context)  {
	photoID, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve user ID from the context
	tokenUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not available"})
		return
	}

	// Query the database to check if the photo exists and belongs to the user
	var photo photo
	result := initializers.DB.Where("id = ? AND user_id = ?", photoID, tokenUserID).First(&photo)
	if result.Error != nil {
		// Handle case where the photo does not exist or does not belong to the user
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found or does not belong to the user"})
		return
	}

	initializers.DB.Delete(&photo)

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully", "photo": photo})
}