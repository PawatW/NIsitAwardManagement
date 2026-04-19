package database

import (
	"fmt"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgrest(cfg *configs.Config) *gorm.DB {
	connection := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DB")
	}
	fmt.Println("💾 Connected to DB successfully")

	return db
}
