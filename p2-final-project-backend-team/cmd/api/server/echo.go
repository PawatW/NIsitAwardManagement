package server

import (
	"fmt"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/middlewares"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/router"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	config         *configs.Config
	handlers       *handlers.Handler
	authMiddleware middlewares.AuthMiddleware
	rbacMiddleware middlewares.RBACMiddleware
}

func NewEchoServer(
	config *configs.Config,
	handlers *handlers.Handler,
	authMiddleware middlewares.AuthMiddleware,
	rbacMiddleware middlewares.RBACMiddleware,
) *EchoServer {
	return &EchoServer{
		config:         config,
		handlers:       handlers,
		authMiddleware: authMiddleware,
		rbacMiddleware: rbacMiddleware,
	}
}

func (s *EchoServer) Start() error {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     s.config.AllowOrigins,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	// Set validator
	e.Validator = utils.NewValidator()

	router := router.NewRouter(e, s.handlers, s.authMiddleware, s.rbacMiddleware)
	router.RegisterAPIRoutes()

	return e.Start(fmt.Sprintf(":%s", s.config.Port))
}
