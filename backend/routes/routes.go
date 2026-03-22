package routes

import (
	"clockit/backend/controllers"
	"clockit/backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	supervisorRoutes := r.Group("/api/supervisors")
	{
		supervisorRoutes.POST("/register", controllers.SupervisorRegister)
		supervisorRoutes.POST("/login", controllers.SupervisorLogin)
		// protected endpoints - require authentication and supervisor role
		supervisorRoutes.GET("/workers/availability", services.JWTAuthMiddleware(), services.SupervisorAuthorizationMiddleware(), controllers.GetWorkerAvailability)
		supervisorRoutes.GET("/:employee_id/shifts", services.JWTAuthMiddleware(), services.SupervisorAuthorizationMiddleware(), controllers.GetShiftsByCompany)
	}

	workerRoutes := r.Group("/api/workers")
	{
		workerRoutes.POST("/register", controllers.WorkerRegister)
		workerRoutes.POST("/login", controllers.WorkerLogin)
		workerRoutes.GET("/:employee_id/shifts", controllers.GetShiftsForWorker)
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
