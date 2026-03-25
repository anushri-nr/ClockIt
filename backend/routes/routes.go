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
		supervisorRoutes.GET("/workers/availability", 
			services.JWTAuthMiddleware(), 
			services.SupervisorAuthorizationMiddleware(), 
			controllers.GetWorkerAvailability,
		)
		supervisorRoutes.GET("/:employee_id/shifts", 
			services.JWTAuthMiddleware(), 
			services.SupervisorAuthorizationMiddleware(), 
			controllers.GetShiftsByCompany,
		)
	}

	workerRoutes := r.Group("/api/workers")
	{
		workerRoutes.POST("/register", controllers.WorkerRegister)
		workerRoutes.POST("/login", controllers.WorkerLogin)
		workerRoutes.GET("/:employee_id/shifts",
			services.JWTAuthMiddleware(),
			services.WorkerAuthorizationMiddleware(),
			controllers.GetShiftsForWorker,
		)
		workerRoutes.POST("/:employee_id/availability",
			services.JWTAuthMiddleware(),
			services.WorkerAuthorizationMiddleware(),
			controllers.CreateAvailability,
		)
	}

	shiftRoutes := r.Group("/api/shifts")
	{
		// We add both middlewares to ensure only logged-in supervisors can touch shifts
		shiftRoutes.POST("/create",
			services.JWTAuthMiddleware(),
			services.SupervisorAuthorizationMiddleware(),
			controllers.CreateShift,
		)
		shiftRoutes.POST("/assign",
			services.JWTAuthMiddleware(),
			services.SupervisorAuthorizationMiddleware(),
			controllers.AssignWorkerToShift,
		)
	}

	companyRoutes := r.Group("/api/companies")
	{
		companyRoutes.POST("/create", controllers.CreateCompany)
		companyRoutes.GET("/", controllers.ListCompanies)
	}
}
