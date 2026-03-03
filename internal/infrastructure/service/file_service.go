package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	"github.com/Ju4n-Dieg0/Go_Points/internal/domain/files"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/logger"
	"github.com/google/uuid"
)

// LocalFileService implementa FileService usando el filesystem local
type LocalFileService struct {
	config config.FileConfig
}

// NewLocalFileService crea una nueva instancia de LocalFileService
func NewLocalFileService(config config.FileConfig) files.FileService {
	return &LocalFileService{
		config: config,
	}
}

// Upload sube un archivo y retorna su path relativo
func (s *LocalFileService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validar el archivo
	if err := s.ValidateFile(header); err != nil {
		return "", err
	}

	// Crear directorio de uploads si no existe
	if err := s.ensureUploadDir(); err != nil {
		logger.Error("Failed to create upload directory", "error", err)
		return "", errors.ErrInternal.WithError(fmt.Errorf("failed to create upload directory"))
	}

	// Generar nombre único para el archivo
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String(), ext)
	relativePath := filepath.Join(s.config.UploadDir, filename)
	fullPath := s.GetFullPath(relativePath)

	// Crear el archivo en el filesystem
	dst, err := os.Create(fullPath)
	if err != nil {
		logger.Error("Failed to create file", "error", err, "path", fullPath)
		return "", errors.ErrInternal.WithError(fmt.Errorf("failed to create file"))
	}
	defer dst.Close()

	// Copiar el contenido
	if _, err := io.Copy(dst, file); err != nil {
		// Si falla, intentar eliminar el archivo parcialmente creado
		os.Remove(fullPath)
		logger.Error("Failed to write file", "error", err, "path", fullPath)
		return "", errors.ErrInternal.WithError(fmt.Errorf("failed to write file"))
	}

	logger.Info("File uploaded successfully", "path", relativePath)
	return relativePath, nil
}

// Delete elimina un archivo por su path
func (s *LocalFileService) Delete(ctx context.Context, path string) error {
	if path == "" {
		return nil
	}

	fullPath := s.GetFullPath(path)

	// Verificar que el archivo existe
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		logger.Warn("File not found for deletion", "path", path)
		return nil // No error si el archivo no existe
	}

	// Eliminar el archivo
	if err := os.Remove(fullPath); err != nil {
		logger.Error("Failed to delete file", "error", err, "path", path)
		return errors.ErrInternal.WithError(fmt.Errorf("failed to delete file"))
	}

	logger.Info("File deleted successfully", "path", path)
	return nil
}

// GetFullPath retorna el path absoluto de un archivo
func (s *LocalFileService) GetFullPath(relativePath string) string {
	// Si ya es un path absoluto, retornarlo
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	return relativePath
}

// ValidateFile valida un archivo sin subirlo
func (s *LocalFileService) ValidateFile(header *multipart.FileHeader) error {
	// Validar tamaño
	if header.Size > s.config.MaxSize {
		return errors.ErrValidation.WithError(
			fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.config.MaxSize),
		)
	}

	// Validar MIME type
	contentType := header.Header.Get("Content-Type")
	if !s.isAllowedType(contentType) {
		return errors.ErrValidation.WithError(
			fmt.Errorf("file type '%s' is not allowed. Allowed types: %s",
				contentType,
				strings.Join(s.config.AllowedTypes, ", "),
			),
		)
	}

	return nil
}

// isAllowedType verifica si un MIME type está permitido
func (s *LocalFileService) isAllowedType(contentType string) bool {
	for _, allowed := range s.config.AllowedTypes {
		if allowed == contentType {
			return true
		}
	}
	return false
}

// ensureUploadDir crea el directorio de uploads si no existe
func (s *LocalFileService) ensureUploadDir() error {
	if err := os.MkdirAll(s.config.UploadDir, 0755); err != nil {
		return err
	}
	return nil
}
