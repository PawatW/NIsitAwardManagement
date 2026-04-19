package router

import (
	"net/http"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/middlewares"
	"github.com/labstack/echo/v4"
)

type Router struct {
	echo           *echo.Echo
	handlers       *handlers.Handler
	authMiddleware middlewares.AuthMiddleware
	rbacMiddleware middlewares.RBACMiddleware
}

func NewRouter(echo *echo.Echo, handlers *handlers.Handler, authMiddleware middlewares.AuthMiddleware, rbacMiddleware middlewares.RBACMiddleware) *Router {
	return &Router{
		echo:           echo,
		handlers:       handlers,
		authMiddleware: authMiddleware,
		rbacMiddleware: rbacMiddleware,
	}
}

func (r *Router) RegisterAPIRoutes() {
	// Health Check
	r.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1 routes
	apiV1 := r.echo.Group("/api/v1")

	// Public routes
	r.RegisterAuthRoutes(apiV1)
	r.RegisterUserRoutes(apiV1)
	r.RegisterOrganizationRoutes(apiV1)
	r.RegisterHonorRollRoute(apiV1)
	r.RegisterDevRoutes(apiV1)

	// Protected routes (require authentication)
	protected := apiV1.Group("", r.authMiddleware.Middleware)
	r.RegisterAdminRoutes(protected)
	r.RegisterAwaredCategoryRoutes(protected)
	r.RegisterRequestRoutes(protected)
	r.RegisterAcademicTermRoutes(protected)
}

func (r *Router) RegisterAuthRoutes(api *echo.Group) {
	authGroup := api.Group("/auth")
	authGroup.GET("/google", r.handlers.Auth.GetGoogleLoginURL)
	authGroup.GET("/google/callback", r.handlers.Auth.GoogleCallback)
}

func (r *Router) RegisterUserRoutes(api *echo.Group) {
	userGroup := api.Group("/user")
	userGroup.POST("/register", r.handlers.User.Register)
	userGroup.GET("/me", r.handlers.User.GetCurrentUser, r.authMiddleware.Middleware)
}

func (r *Router) RegisterAdminRoutes(api *echo.Group) {
	adminGroup := api.Group("/admin", r.rbacMiddleware.RequireAdmin)
	adminGroup.POST("/users", r.handlers.Admin.CreateUser)
	adminGroup.GET("/users", r.handlers.Admin.GetAllUsers)
	adminGroup.PATCH("/users/:id/status", r.handlers.Admin.UpdateUserStatus)
	adminGroup.POST("/campuses", r.handlers.Admin.CreateCampus)
	adminGroup.POST("/faculties", r.handlers.Admin.CreateFaculty)
	adminGroup.POST("/departments", r.handlers.Admin.CreateDepartment)
}

func (r *Router) RegisterAwaredCategoryRoutes(api *echo.Group) {
	awardGroup := api.Group("/award-categories")

	// Public read endpoints (any authenticated user)
	awardGroup.GET("", r.handlers.Award.GetAllAwardCategories)
	awardGroup.GET("/:id", r.handlers.Award.GetAwardCategory)

	// Admin only endpoints (create, update, delete)
	awardGroup.POST("", r.handlers.Award.CreateAwardCategory, r.rbacMiddleware.RequireAdmin)
	awardGroup.PUT("/:id", r.handlers.Award.UpdateAwardCategory, r.rbacMiddleware.RequireAdmin)
	awardGroup.DELETE("/:id", r.handlers.Award.DeleteAwardCategory, r.rbacMiddleware.RequireAdmin)
}

func (r *Router) RegisterRequestRoutes(api *echo.Group) {
	// Request routes
	requestGroup := api.Group("/requests")
	requestGroup.POST("", r.handlers.Request.CreateRequest)
	requestGroup.GET("/my", r.handlers.Request.GetMyRequests)
	requestGroup.GET("/:id", r.handlers.Request.GetRequestDetails)
	requestGroup.POST("/:id/documents", r.handlers.Request.UploadDocument)

	// Unified approval/rejection endpoints (role-aware)
	requestGroup.POST("/:id/approve", r.handlers.Request.ApproveRequest)
	requestGroup.POST("/:id/reject", r.handlers.Request.RejectRequest)

	// Admin/committee endpoints
	requestGroup.POST("/:id/change-category", r.handlers.Request.ChangeAwardCategory)

	// Query endpoints
	// GET /api/v1/requests - Filters: student_id, award_id, academic_term_id, current_status, page, limit
	requestGroup.GET("", r.handlers.Request.GetAllRequests)
	// GET /api/v1/requests/status/:status - Filters: status (path param), page, limit (query params)
	requestGroup.GET("/status/:status", r.handlers.Request.GetRequestsByStatus)
}

func (r *Router) RegisterOrganizationRoutes(api *echo.Group) {
	organizationGroup := api.Group("/organization")
	organizationGroup.GET("/campuses", r.handlers.Organization.GetCampuses)
	organizationGroup.GET("/campuses/:campus_id/faculties", r.handlers.Organization.GetFacultiesByCampus)
	organizationGroup.GET("/faculties/:faculty_id/departments", r.handlers.Organization.GetDepartmentsByFaculty)
}

func (r *Router) RegisterDevRoutes(api *echo.Group) {
	// Development/Testing routes (only available in non-production)
	devGroup := api.Group("/dev")
	devGroup.POST("/generate-token", r.handlers.Dev.GenerateToken)
}

func (r *Router) RegisterAcademicTermRoutes(api *echo.Group) {
	termGroup := api.Group("/academic-terms")
	termGroup.GET("", r.handlers.AcademicTerm.GetAll)
	termGroup.GET("/latest", r.handlers.AcademicTerm.GetLatest)
	termGroup.GET("/:id", r.handlers.AcademicTerm.GetByID)
	termGroup.POST("", r.handlers.AcademicTerm.Create, r.rbacMiddleware.RequireAdmin)
	termGroup.PATCH("/:id", r.handlers.AcademicTerm.UpdateIsOpen, r.rbacMiddleware.RequireAdmin) // +
}

func (r *Router) RegisterHonorRollRoute(api *echo.Group) {
	// Public endpoint - no authentication required
	// GET /api/v1/honor-roll - Filters: academic_year, semester, campus_id, award_id
	api.GET("/honor-roll", r.handlers.Award.GetHonorRoll)
}
