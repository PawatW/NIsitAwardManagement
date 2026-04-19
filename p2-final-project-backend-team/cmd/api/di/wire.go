//go:build wireinject
// +build wireinject

package di

import (
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/cmd/api/server"
	//"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	//"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers"
	//"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/database"
	//"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/router"
	"github.com/google/wire"
	//"github.com/labstack/echo/v4"
)

func InitializeAPI() (*server.EchoServer, error) {
	wire.Build(
		ConfigSet,
		InfrastructureSet,
		RepositorySet,
		ServiceSet,
		HandlerSet,
		MiddlewareSet,
		server.NewEchoServer,
	)

	//return &server.EchoServer{}
	return nil, nil
}
