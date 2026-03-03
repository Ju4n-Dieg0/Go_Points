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
	DocumentType   DocumentType   `gorm:"type:varchar(20);not null;index:idx_document"`
	DocumentNumber string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_document_number"`
	Name           string         `gorm:"type:varchar(255);not null;index"`
	Email          string         `gorm:"type:varchar(255);not null;index"`
	Phone          *string        `gorm:"type:varchar(20)"`
	Photo          *string        `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;index"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
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
