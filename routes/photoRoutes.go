package routes

import (
	"api/rakamin-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupPhotoRoutes(router *gin.Engine) {
    photoGroup := router.Group("/photos")
    {
        photoGroup.GET("/", controllers.GetUser)
        photoGroup.POST("/", controllers.Login)
        photoGroup.GET("/:id", controllers.Login)
        photoGroup.PUT("/:id", controllers.Login)
        photoGroup.DELETE("/:id", controllers.Login)
    }
}