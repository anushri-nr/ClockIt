package routes

import (
	"clockit/backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	supervisorRoutes := r.Group("/api/supervisors")
	{
		supervisorRoutes.POST("/register", controllers.SupervisorRegister)
		supervisorRoutes.GET("/workers/availability", controllers.GetWorkerAvailability)
	}

	workerRoutes := r.Group("/api/workers")
	{
		workerRoutes.POST("/register", controllers.WorkerRegister)
	}

	shiftRoutes := r.Group("/api/shifts")
	{
		shiftRoutes.POST("/create", controllers.CreateShift)
		shiftRoutes.POST("/assign", controllers.AssignWorkerToShift)
	}

	companyRoutes := r.Group("/api/companies")
	{
		companyRoutes.POST("/create", controllers.CreateCompany)
	}
}
