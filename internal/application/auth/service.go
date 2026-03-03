package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Ju4n-Dieg0/Go_Points/internal/config"
	domainAuth "github.com/Ju4n-Dieg0/Go_Points/internal/domain/auth"
	"github.com/Ju4n-Dieg0/Go_Points/internal/shared/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service define la lógica de negocio de autenticación
type Service interface {
	Register(ctx context.Context, req *domainAuth.RegisterRequest) (*domainAuth.AuthResponse, error)
	Login(ctx context.Context, req *domainAuth.LoginRequest) (*domainAuth.AuthResponse, error)
	RefreshToken(ctx context.Context, req *domainAuth.RefreshTokenRequest) (*domainAuth.RefreshTokenResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	RequestPasswordReset(ctx context.Context, req *domainAuth.RequestPasswordResetRequest) error
	ConfirmPasswordReset(ctx context.Context, req *domainAuth.ConfirmPasswordResetRequest) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domainAuth.UserDTO, error)
}

type service struct {
	repo         domainAuth.Repository
	emailService domainAuth.EmailService
	jwtConfig    *config.JWTConfig
}

// NewService crea una nueva instancia del servicio de autenticación
func NewService(repo domainAuth.Repository, emailService domainAuth.EmailService, jwtConfig *config.JWTConfig) Service {
	return &service{
		repo:         repo,
		emailService: emailService,
		jwtConfig:    jwtConfig,
	}
}

// Register registra un nuevo usuario
func (s *service) Register(ctx context.Context, req *domainAuth.RegisterRequest) (*domainAuth.AuthResponse, error) {
	// Verificar si el email ya existe
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if exists {
		return nil, errors.ErrConflict.WithError(fmt.Errorf("email already registered"))
	}

	// Validar rol
	if !domainAuth.IsValidRole(req.Role) {
		return nil, errors.ErrValidation.WithError(fmt.Errorf("invalid role"))
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	// Crear usuario
	user := &domainAuth.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	// Generar tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	// Guardar refresh token
	refreshTokenExpiry := time.Now().Add(time.Duration(s.jwtConfig.RefreshExpiration) * time.Second)
	if err := s.repo.UpdateRefreshToken(ctx, user.ID, &refreshToken, &refreshTokenExpiry); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	// Enviar email de bienvenida (no bloqueante)
	go func() {
		_ = s.emailService.SendWelcomeEmail(context.Background(), user.Email, user.FirstName)
	}()

	return &domainAuth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtConfig.AccessExpiration,
		User:         domainAuth.ToUserDTO(user),
	}, nil
}

// Login autentica un usuario
func (s *service) Login(ctx context.Context, req *domainAuth.LoginRequest) (*domainAuth.AuthResponse, error) {
	// Buscar usuario por email
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if user == nil {
		return nil, errors.ErrAuthentication.WithError(fmt.Errorf("invalid credentials"))
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.ErrAuthentication.WithError(fmt.Errorf("invalid credentials"))
	}

	// Verificar si el usuario está activo
	if !user.IsActive {
		return nil, errors.ErrForbidden.WithError(fmt.Errorf("user account is disabled"))
	}

	// Generar tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	// Guardar refresh token y actualizar last login
	refreshTokenExpiry := time.Now().Add(time.Duration(s.jwtConfig.RefreshExpiration) * time.Second)
	if err := s.repo.UpdateRefreshToken(ctx, user.ID, &refreshToken, &refreshTokenExpiry); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	return &domainAuth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtConfig.AccessExpiration,
		User:         domainAuth.ToUserDTO(user),
	}, nil
}

// RefreshToken genera un nuevo access token usando el refresh token
func (s *service) RefreshToken(ctx context.Context, req *domainAuth.RefreshTokenRequest) (*domainAuth.RefreshTokenResponse, error) {
	// Buscar usuario por refresh token
	user, err := s.repo.FindByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if user == nil {
		return nil, errors.ErrUnauthorized.WithError(fmt.Errorf("invalid refresh token"))
	}

	// Verificar si el usuario está activo
	if !user.IsActive {
		return nil, errors.ErrForbidden.WithError(fmt.Errorf("user account is disabled"))
	}

	// Generar nuevo access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	return &domainAuth.RefreshTokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   s.jwtConfig.AccessExpiration,
	}, nil
}

// Logout invalida el refresh token del usuario
func (s *service) Logout(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.UpdateRefreshToken(ctx, userID, nil, nil); err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	return nil
}

// RequestPasswordReset solicita un reset de contraseña
func (s *service) RequestPasswordReset(ctx context.Context, req *domainAuth.RequestPasswordResetRequest) error {
	// Buscar usuario por email
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	
	// No revelar si el email existe o no (seguridad)
	if user == nil {
		return nil
	}

	// Generar token de reset
	resetToken, err := s.generateRefreshToken()
	if err != nil {
		return errors.ErrInternal.WithError(err)
	}

	// Guardar token con expiración de 1 hora
	resetTokenExpiry := time.Now().Add(1 * time.Hour)
	user.PasswordResetToken = &resetToken
	user.PasswordResetExpiry = &resetTokenExpiry

	if err := s.repo.Update(ctx, user); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	// Enviar email con token (no bloqueante)
	go func() {
		_ = s.emailService.SendPasswordResetEmail(context.Background(), user.Email, user.FirstName, resetToken)
	}()

	return nil
}

// ConfirmPasswordReset confirma el reset de contraseña
func (s *service) ConfirmPasswordReset(ctx context.Context, req *domainAuth.ConfirmPasswordResetRequest) error {
	// Buscar usuario por token
	user, err := s.repo.FindByPasswordResetToken(ctx, req.Token)
	if err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	if user == nil {
		return errors.ErrBadRequest.WithError(fmt.Errorf("invalid or expired reset token"))
	}

	// Hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrInternal.WithError(err)
	}

	// Actualizar contraseña y limpiar token
	if err := s.repo.UpdatePassword(ctx, user.ID, string(hashedPassword)); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	// Invalidar refresh token por seguridad
	if err := s.repo.UpdateRefreshToken(ctx, user.ID, nil, nil); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	// Enviar email de confirmación (no bloqueante)
	go func() {
		_ = s.emailService.SendPasswordChangedEmail(context.Background(), user.Email, user.FirstName)
	}()

	return nil
}

// GetUserByID obtiene un usuario por ID
func (s *service) GetUserByID(ctx context.Context, userID uuid.UUID) (*domainAuth.UserDTO, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if user == nil {
		return nil, errors.ErrNotFound.WithError(fmt.Errorf("user not found"))
	}

	userDTO := domainAuth.ToUserDTO(user)
	return &userDTO, nil
}

// generateAccessToken genera un JWT access token
func (s *service) generateAccessToken(user *domainAuth.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Duration(s.jwtConfig.AccessExpiration) * time.Second).Unix(),
		"iat":   time.Now().Unix(),
		"type":  "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.AccessSecret))
}

// generateRefreshToken genera un token de refresh aleatorio
func (s *service) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
