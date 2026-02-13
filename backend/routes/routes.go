package routes

import (
	"clockit/backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	supervisorRoutes := r.Group("/api/supervisor")
	{
		supervisorRoutes.POST("/register", controllers.SupervisorRegister)
	}
}
