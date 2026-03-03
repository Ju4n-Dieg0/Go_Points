package consumer

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentType representa los tipos de documento válidos
type DocumentType string

const (
	DocumentTypeDNI      DocumentType = "DNI"
	DocumentTypeCE       DocumentType = "CE"       // Carnet de Extranjería
	DocumentTypePassport DocumentType = "PASSPORT"
	DocumentTypeRUC      DocumentType = "RUC"
	DocumentTypeOther    DocumentType = "OTHER"
)

// Consumer representa la entidad de consumidor en el dominio
type Consumer struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	DocumentType   DocumentType   `gorm:"type:varchar(20);not null;index:idx_consumer_document,priority:1"`
	DocumentNumber string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_consumer_document_unique;index:idx_consumer_document,priority:2"`
	Name           string         `gorm:"type:varchar(255);not null;index:idx_consumer_name"`
	Email          string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_consumer_email"`
	Phone          *string        `gorm:"type:varchar(20);index:idx_consumer_phone"`
	Photo          *string        `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;not null;index:idx_consumer_created"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime;not null"`
	DeletedAt      gorm.DeletedAt `gorm:"index:idx_consumer_deleted"`
}

// TableName especifica el nombre de la tabla
func (Consumer) TableName() string {
	return "consumers"
}

// IsValidDocumentType verifica si el tipo de documento es válido
func IsValidDocumentType(docType DocumentType) bool {
	switch docType {
	case DocumentTypeDNI, DocumentTypeCE, DocumentTypePassport, DocumentTypeRUC, DocumentTypeOther:
		return true
	default:
		return false
	}
}

// GetFullIdentification retorna la identificación completa del consumidor
func (c *Consumer) GetFullIdentification() string {
	return string(c.DocumentType) + "-" + c.DocumentNumber
}
