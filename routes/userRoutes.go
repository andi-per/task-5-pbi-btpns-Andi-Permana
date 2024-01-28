package routes

import (
	"api/rakamin-api/controllers"
	"api/rakamin-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
    userGroup := router.Group("/users")
    {
        userGroup.POST("/register", controllers.Register)
        userGroup.POST("/login", controllers.Login)

        userGroup.Use(middleware.AuthenticateToken)
        userGroup.PUT("/:userId", controllers.UpdateUser)
        userGroup.DELETE("/:userId", controllers.DeleteUser)
    }
}