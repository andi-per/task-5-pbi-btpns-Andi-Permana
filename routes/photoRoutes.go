package routes

import (
	"api/rakamin-api/controllers"
	"api/rakamin-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPhotoRoutes(router *gin.Engine) {
    photoGroup := router.Group("/photos")
    {
		photoGroup.Use(middleware.AuthenticateToken)
        photoGroup.POST("/", controllers.PostPhoto)
		photoGroup.GET("/", controllers.GetPhotos)
        photoGroup.PUT("/:photoId", controllers.UpdatePhoto)
        photoGroup.DELETE("/:photoId", controllers.DeletePhoto)
    }
}