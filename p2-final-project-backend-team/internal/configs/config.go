package configs

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AppEnv       string      `env:"APP_ENV" envDefault:"development"`
	Port         string      `env:"PORT"`
	AllowOrigins []string    `env:"ALLOW_ORIGINS" envSeparator:","`
	Database     DBConfig    `envPrefix:"DATABASE_"`
	OAuth        OAuthConfig `envPrefix:"OAUTH_"`
	JWT          JWTConfig   `envPrefix:"JWT_"`
	Minio        MinioConfig `envPrefix:"MINIO_"`
}

// IsDevelopment returns true if app is in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv != "production"
}

type DBConfig struct {
	Host     string `env:"HOST"`
	Name     string `env:"NAME"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Port     string `env:"PORT"`
}

type OAuthConfig struct {
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL"`
}

type JWTConfig struct {
	Secret      string `env:"SECRET"`
	ExpiryHours int    `env:"EXPIRY_HOURS" envDefault:"24"`
}

type MinioConfig struct {
	Endpoint        string `env:"ENDPOINT"`
	PublicEndpoint  string `env:"PUBLIC_ENDPOINT"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `env:"USE_SSL" envDefault:"false"`
	BucketName      string `env:"BUCKET_NAME" envDefault:"student-uploads"`
	Region          string `env:"REGION" envDefault:"ap-southeast-1"`
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found or error loading it.")
	}

	config := &Config{}

	if err := env.Parse(config); err != nil {
		log.Panic().Err(err).Msg("Failed to parse environment variables")
	}

	return config
}
