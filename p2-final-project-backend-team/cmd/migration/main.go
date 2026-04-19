package main

import (
	"fmt"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/database"
	"github.com/rs/zerolog/log"
)

func main() {
	config := configs.NewConfig()
	db := database.NewPostgrest(config)

	err := db.AutoMigrate(
		&models.Campus{},
		&models.Faculty{},
		&models.Department{},
		&models.User{},
		&models.AwardCategory{},
		&models.AcademicTerm{},
		&models.Request{},
		&models.RequestStatusHistory{},
		&models.RequestDocument{},
		&models.AwardCategoryChange{},
	)

	if err != nil {
		log.Fatal().Msgf("Migration Failed: %v", err)
	}

	fmt.Println("Migration completed ✅")
}
