package di

import (
	"github.com/google/wire"

	// Config
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/middlewares"

	// Infrastructure
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/auth"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/context"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/database"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/oauth"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/storage"

	// Repositories
	academicTermRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/academic_term"
	awardRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/award"
	campusRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/campus"
	departmentRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/department"
	facultyRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/faculty"
	requestRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/request"
	userRepo "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/repositories/user"

	// Services
	academicTermSvc "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/academic_term"
	adminSvc "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/admin"
	authSvc "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/auth"
	awardSvc "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/award"
	campusSvc "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/campus"
	departmentSvc "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/department"
	facultySvc "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/faculty"
	requestSvc "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/request"
	userSvc "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/services/user"

	// Handlers
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers"
	academicTermHd "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/academic_term"
	adminHd "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/admin"
	authHd "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/auth"
	awardHd "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/award"
	devHd "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/dev"
	organizationHd "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/organization"
	requestHd "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/request"
	userHd "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/user"
)

var ConfigSet = wire.NewSet(
	configs.NewConfig,
)

var InfrastructureSet = wire.NewSet(
	context.NewContext,
	database.NewPostgrest,
	oauth.NewGoogleProvider,
	auth.NewJWTManager,
	ProvideMinioClient,
)

var RepositorySet = wire.NewSet(
	userRepo.NewUserRepository,
	campusRepo.NewCampusRepository,
	departmentRepo.NewDepartmentRepository,
	facultyRepo.NewFacultyRepository,
	awardRepo.NewRepository,
	requestRepo.NewRepository,
	academicTermRepo.NewAcademicTermRepository,
)

var ServiceSet = wire.NewSet(
	adminSvc.NewAdminService,
	authSvc.NewAuthService,
	campusSvc.NewCampusService,
	facultySvc.NewFacultyService,
	departmentSvc.NewDepartmentService,
	userSvc.NewUserService,
	awardSvc.NewService,
	requestSvc.NewService,
	academicTermSvc.NewAcademicTermService,
)

var HandlerSet = wire.NewSet(
	handlers.NewHandlers,
	adminHd.NewAdminHandler,
	authHd.NewAuthHandler,
	userHd.NewUserHandler,
	awardHd.NewHandler,
	requestHd.NewHandler,
	devHd.NewDevHandler,
	organizationHd.NewOrganizationHandler,
	academicTermHd.NewHandler,
)

var MiddlewareSet = wire.NewSet(
	middlewares.NewAuthMiddleware,
	middlewares.NewRBACMiddleware,
)

// ProvideMinioClient provides MinIO client instance
func ProvideMinioClient(config *configs.Config) (storage.StorageClient, error) {
	return storage.NewMinioClient(&config.Minio)
}
