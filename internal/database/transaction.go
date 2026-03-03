package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TxFunc es una función que se ejecuta dentro de una transacción
type TxFunc func(*gorm.DB) error

// WithTransaction ejecuta una función dentro de una transacción
// Hace rollback automático en caso de error o panic
func WithTransaction(db *gorm.DB, fn TxFunc) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Defer para manejar panic y rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic in transaction: %v", r)
		}
	}()

	// Ejecutar función
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	// Commit si todo salió bien
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTransactionContext ejecuta una función dentro de una transacción con contexto
func WithTransactionContext(ctx context.Context, db *gorm.DB, fn func(*gorm.DB) error) error {
	return WithTransaction(db.WithContext(ctx), fn)
}

// BeginTx inicia una transacción con nivel de aislamiento personalizado
func BeginTx(db *gorm.DB, isolationLevel string) *gorm.DB {
	tx := db.Begin()
	if isolationLevel != "" {
		tx.Exec(fmt.Sprintf("SET TRANSACTION ISOLATION LEVEL %s", isolationLevel))
	}
	return tx
}

// IsolationLevels niveles de aislamiento PostgreSQL
const (
	IsolationLevelReadUncommitted = "READ UNCOMMITTED"
	IsolationLevelReadCommitted   = "READ COMMITTED"
	IsolationLevelRepeatableRead  = "REPEATABLE READ"
	IsolationLevelSerializable    = "SERIALIZABLE"
)
