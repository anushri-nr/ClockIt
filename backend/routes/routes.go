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
		supervisorRoutes.GET("/shifts/requested",
			services.JWTAuthMiddleware(),
			services.SupervisorAuthorizationMiddleware(),
			controllers.GetRequestedShifts,
		)
		supervisorRoutes.GET("/workers/overtime",
			services.JWTAuthMiddleware(),
			services.SupervisorAuthorizationMiddleware(),
			controllers.GetWorkersWithOvertimeHours,
		)
		supervisorRoutes.GET("/shifts/assigned",
			services.JWTAuthMiddleware(),
			services.SupervisorAuthorizationMiddleware(),
			controllers.GetAssignedShifts,
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
		workerRoutes.GET("/shifts/released",
			services.JWTAuthMiddleware(),
			services.WorkerAuthorizationMiddleware(),
			controllers.GetReleasedShiftsForWorkerCompany,
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
		shiftRoutes.POST("/release",
			services.JWTAuthMiddleware(),
			services.WorkerAuthorizationMiddleware(),
			controllers.ReleaseShiftForWorker,
		)
		supervisorRoutes.PATCH("/shifts/:shift_id/reject",
			services.JWTAuthMiddleware(),
			services.SupervisorAuthorizationMiddleware(),
			controllers.RejectShiftRequest,
		)
		workerRoutes.POST("/shifts/:shift_id/request",
			services.JWTAuthMiddleware(),
			services.WorkerAuthorizationMiddleware(),
			controllers.RequestShift,
		)
		shiftRoutes.DELETE("/:shift_id",
			services.JWTAuthMiddleware(),
			services.SupervisorAuthorizationMiddleware(),
			controllers.DeleteShift,
		)
		shiftRoutes.PATCH("/:shift_id/unassign",
			services.JWTAuthMiddleware(),
			services.SupervisorAuthorizationMiddleware(),
			controllers.UnassignWorker,
		)
	}

	companyRoutes := r.Group("/api/companies")
	{
		companyRoutes.POST("/create", controllers.CreateCompany)
		companyRoutes.GET("/", controllers.ListCompanies)
	}
}
