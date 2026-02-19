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
		supervisorRoutes.GET("/:employee_id/shifts", controllers.GetShiftsCreatedBySupervisor)
	}

	workerRoutes := r.Group("/api/workers")
	{
		workerRoutes.POST("/register", controllers.WorkerRegister)
		workerRoutes.GET(":employee_id/shifts", controllers.GetShiftsForWorker)
	}

	shiftRoutes := r.Group("/api/shifts")
	{
		shiftRoutes.POST("/create", controllers.CreateShift)
		shiftRoutes.POST("/assign", controllers.AssignWorkerToShift)
	}

	companyRoutes := r.Group("/api/companies")
	{
		companyRoutes.POST("/create", controllers.CreateCompany)
		companyRoutes.GET("/", controllers.ListCompanies)
	}
}
