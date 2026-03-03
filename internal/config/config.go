package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
	File     FileConfig
	Points   PointsConfig
}

type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type JWTConfig struct {
	AccessSecret      string
	RefreshSecret     string
	AccessExpiration  int
	RefreshExpiration int
}

type AppConfig struct {
	Name        string
	Environment string
	LogLevel    string
}

type FileConfig struct {
	UploadDir     string
	MaxSize       int64    // bytes
	AllowedTypes  []string // MIME types
}

type PointsConfig struct {
	ExpirationMonths       int // Meses para expiración de puntos
	InactivityMonths       int // Meses de inactividad para penalización
	InactivityPenalty      int64 // Puntos a restar por inactividad
	NotificationBeforeDays int // Días antes de notificar expiración
}

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")

	// Configurar para leer variables de entorno
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Valores por defecto
	setDefaults()

	// Intentar leer el archivo de configuración (opcional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		log.Println("No config file found, using environment variables and defaults")
	}

	config := &Config{
		Server: ServerConfig{
			Host:            viper.GetString("SERVER_HOST"),
			Port:            viper.GetString("SERVER_PORT"),
			ReadTimeout:     viper.GetInt("SERVER_READ_TIMEOUT"),
			WriteTimeout:    viper.GetInt("SERVER_WRITE_TIMEOUT"),
			ShutdownTimeout: viper.GetInt("SERVER_SHUTDOWN_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetString("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			DBName:          viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSLMODE"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetInt("DB_CONN_MAX_LIFETIME"),
		},
		JWT: JWTConfig{
			AccessSecret:      viper.GetString("JWT_ACCESS_SECRET"),
			RefreshSecret:     viper.GetString("JWT_REFRESH_SECRET"),
			AccessExpiration:  viper.GetInt("JWT_ACCESS_EXPIRATION"),
			RefreshExpiration: viper.GetInt("JWT_REFRESH_EXPIRATION"),
		},
		App: AppConfig{
			Name:        viper.GetString("APP_NAME"),
			Environment: viper.GetString("APP_ENV"),
			LogLevel:    viper.GetString("LOG_LEVEL"),
		},
		File: FileConfig{
			UploadDir:    viper.GetString("FILE_UPLOAD_DIR"),
			MaxSize:      viper.GetInt64("FILE_MAX_SIZE"),
			AllowedTypes: viper.GetStringSlice("FILE_ALLOWED_TYPES"),
		},
		Points: PointsConfig{
			ExpirationMonths:       viper.GetInt("POINT_EXPIRATION_MONTHS"),
			InactivityMonths:       viper.GetInt("INACTIVITY_MONTHS"),
			InactivityPenalty:      viper.GetInt64("INACTIVITY_PENALTY"),
			NotificationBeforeDays: viper.GetInt("NOTIFICATION_BEFORE_DAYS"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", 10)
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 10)
	viper.SetDefault("SERVER_SHUTDOWN_TIMEOUT", 15)

	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "go_points")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 300)

	// JWT defaults
	viper.SetDefault("JWT_ACCESS_SECRET", "change-me-in-production")
	viper.SetDefault("JWT_REFRESH_SECRET", "change-me-in-production-refresh")
	viper.SetDefault("JWT_ACCESS_EXPIRATION", 900)    // 15 minutos
	viper.SetDefault("JWT_REFRESH_EXPIRATION", 604800) // 7 días

	// App defaults
	viper.SetDefault("APP_NAME", "Go Points API")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("LOG_LEVEL", "info")

	// File defaults
	viper.SetDefault("FILE_UPLOAD_DIR", "uploads")
	viper.SetDefault("FILE_MAX_SIZE", 5242880) // 5MB
	viper.SetDefault("FILE_ALLOWED_TYPES", []string{"image/jpeg", "image/jpg", "image/png", "image/webp"})

	// Points defaults
	viper.SetDefault("POINT_EXPIRATION_MONTHS", 12)  // 12 meses
	viper.SetDefault("INACTIVITY_MONTHS", 6)         // 6 meses
	viper.SetDefault("INACTIVITY_PENALTY", 100)      // 100 puntos
	viper.SetDefault("NOTIFICATION_BEFORE_DAYS", 7)  // 7 días antes
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET is required")
	}
	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}
	return nil
}

func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}
