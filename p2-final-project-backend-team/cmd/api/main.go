package main

import (
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/cmd/api/di"
	"github.com/rs/zerolog/log"
)

func main() {
	server, err := di.InitializeAPI()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to initialize API")
	}

	if err := server.Start(); err != nil {
		log.Panic().Err(err).Msg("Failed to start server")
	}
}
