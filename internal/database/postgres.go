package database

import (
	"fmt"
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Database encapsula la conexión a la base de datos
type Database struct {
	DB *gorm.DB
}

// NewDatabase crea una nueva instancia de Database con configuración de GORM
func NewDatabase(cfg *config.DatabaseConfig, appEnv string) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Configurar logger de GORM
	var gormLogLevel gormLogger.LogLevel
	if appEnv == "production" {
		gormLogLevel = gormLogger.Error
	} else {
		gormLogLevel = gormLogger.Info
	}

	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt:            true,  // Preparar statements para mejor performance
		SkipDefaultTransaction: false, // Mantener transacciones automáticas por seguridad
		QueryFields:            true,  // Seleccionar campos específicos
		DisableForeignKeyConstraintWhenMigrating: false, // Mantener constraints
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configurar pool de conexiones
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute) // Cerrar conexiones idle después de 5 min

	// Verificar conexión
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connected successfully",
		"host", cfg.Host,
		"database", cfg.DBName,
		"max_open_conns", cfg.MaxOpenConns,
		"max_idle_conns", cfg.MaxIdleConns,
	)

	return &Database{DB: db}, nil
}

// Close cierra la conexión a la base de datos
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	logger.Info("Database connection closed")
	return nil
}

// AutoMigrate ejecuta las migraciones automáticas de GORM
func (d *Database) AutoMigrate(models ...interface{}) error {
	logger.Info("Running database migrations...")

	if err := d.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// HealthCheck verifica el estado de la base de datos
func (d *Database) HealthCheck() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// GetDB retorna la instancia de GORM DB
func (d *Database) GetDB() *gorm.DB {
	return d.DB
}
