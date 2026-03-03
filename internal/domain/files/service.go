package files

import (
	"context"
	"mime/multipart"
)

// FileService define la interfaz para operaciones con archivos
type FileService interface {
	// Upload sube un archivo y retorna su path relativo
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)

	// Delete elimina un archivo por su path
	Delete(ctx context.Context, path string) error

	// GetFullPath retorna el path absoluto de un archivo
	GetFullPath(relativePath string) string

	// ValidateFile valida un archivo sin subirlo
	ValidateFile(header *multipart.FileHeader) error
}
