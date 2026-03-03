# Arquitectura del Proyecto - Go Points

## 📐 Clean Architecture

Este proyecto implementa **Clean Architecture** real con separación estricta de responsabilidades.

## 📂 Estructura de Módulos

Cada módulo sigue el mismo patrón:

```
internal/
├── domain/{module}/
│   ├── entity.go        # Entidades de negocio (GORM models)
│   ├── dto.go           # Data Transfer Objects
│   └── repository.go    # Interfaz del repositorio
│
├── application/{module}/
│   └── service.go       # Lógica de negocio
│
└── infrastructure/
    ├── persistence/
    │   └── {module}_repository.go  # Implementación GORM
    │
    └── http/
        ├── handler/
        │   └── {module}_handler.go  # HTTP handlers
        │
        └── routes/
            └── {module}_routes.go   # Configuración de rutas
```

## 🎯 Módulos Implementados

### 1. Auth (Autenticación)

**Entidad:** `User`
- Roles: SUPER_ADMIN, COMPANY, CONSUMER
- JWT con access + refresh tokens
- Password hashing con bcrypt
- Soft delete

**Funcionalidades:**
- ✅ Registro
- ✅ Login
- ✅ Refresh token
- ✅ Logout
- ✅ Password reset
- ✅ Profile

**Middleware:**
- `AuthMiddleware` - Verifica JWT
- `RequireRole` - Control de acceso por rol
- `RequireSuperAdmin` - Solo super admin

---

### 2. Companies (Empresas)

**Entidad:** `Company`
- Relación con User
- Estado activo/inactivo
- Soft delete

**Funcionalidades:**
- ✅ Crear empresa (auto-crea suscripción)
- ✅ Actualizar empresa
- ✅ Listar empresas
- ✅ Obtener por ID

**Reglas de Negocio:**
- Al crear empresa → se crea suscripción activa de 30 días
- Transacción atómica (empresa + suscripción)

---

### 3. Subscriptions (Suscripciones)

**Entidad:** `Subscription`
- Relación con Company
- Fechas de inicio y fin
- Estado activo/inactivo
- Auto-cálculo de días restantes

**Funcionalidades:**
- ✅ Renovar suscripción (+30 días)
- ✅ Cancelar suscripción
- ✅ Validar suscripción activa
- ✅ Listar suscripciones próximas a expirar
- ✅ Desactivar suscripciones expiradas (cronjob)

**Reglas de Negocio:**
- Renovación en transacción (actualiza suscripción + activa empresa)
- Validación antes de operaciones críticas
- Expiración automática de empresas

**Middleware:**
- `RequireActiveSubscription` - Valida suscripción antes de operaciones

---

### 4. Consumers (Consumidores)

**Entidad:** `Consumer`
- Tipos de documento: DNI, CE, PASSPORT, RUC, OTHER
- **Unique constraint** en DocumentNumber
- Soft delete

**Funcionalidades:**
- ✅ Crear consumidor
- ✅ Actualizar consumidor
- ✅ Eliminar consumidor
- ✅ Listar con paginación
- ✅ Buscar por nombre/email/documento
- ✅ Obtener por ID
- ✅ Obtener por número de documento

**Reglas de Negocio:**
- Validación de tipo de documento
- Manejo de unique constraint de PostgreSQL (error 23505)
- Paginación 1-100 items

---

### 5. Files (Archivos)

**Interfaz:** `FileService`

**Implementación:** `LocalFileService`

**Funcionalidades:**
- ✅ Upload con validación MIME
- ✅ Límite de tamaño (5MB default)
- ✅ Nombres únicos (timestamp + UUID)
- ✅ Delete de archivos
- ✅ Cleanup automático

**Tipos permitidos:**
- image/jpeg
- image/jpg
- image/png
- image/webp

**Configuración:**
```env
FILE_UPLOAD_DIR=uploads
FILE_MAX_SIZE=5242880
FILE_ALLOWED_TYPES=image/jpeg,image/jpg,image/png,image/webp
```

---

### 6. Products (Productos)

**Entidad:** `Product`
- Relación con Company (FK)
- Photo (path del archivo)
- Precio (decimal 10,2)
- Visibilidad (is_visible)
- Soft delete

**Funcionalidades:**
- ✅ Crear producto (con foto opcional)
- ✅ Actualizar producto (con foto opcional)
- ✅ Eliminar producto (+ cleanup de foto)
- ✅ Listar productos de empresa
- ✅ Listar catálogo público (solo visibles)
- ✅ Buscar productos

**Reglas de Negocio:**
- **Requiere suscripción activa** para crear/actualizar/eliminar
- Solo role COMPANY puede gestionar productos
- Al actualizar foto → elimina la anterior
- Al eliminar producto → elimina la foto
- Catálogo público solo muestra `is_visible=true`

**Dependencias:**
- FileService (manejo de fotos)
- SubscriptionService (validación)

---

## 🔐 Autenticación y Autorización

### JWT Strategy

**Access Token:**
- Expira en 15 minutos (configurable)
- Claims: user_id, email, role
- Usado en header `Authorization: Bearer {token}`

**Refresh Token:**
- Expira en 7 días (configurable)
- Permite renovar access token
- Invalidado en logout

### Control de Acceso

**Roles:**
- `SUPER_ADMIN` - Acceso completo
- `COMPANY` - Gestión de productos, empresas
- `CONSUMER` - Solo lectura (catálogo público)

**Middleware de Autorización:**
```go
// Requiere autenticación
protected := api.Use(middleware.AuthMiddleware(jwtConfig))

// Requiere rol específico
protected.Use(middleware.RequireRole("COMPANY"))

// Solo super admin
protected.Use(middleware.RequireSuperAdmin())
```

---

## 💾 Base de Datos

### Migraciones Automáticas

```go
db.AutoMigrate(
    &auth.User{},
    &company.Company{},
    &subscription.Subscription{},
    &consumer.Consumer{},
    &product.Product{},
)
```

### Relaciones

```
User (1) ──┬─→ (1) Company
           │
           └─→ (N) Consumer

Company (1) ──┬─→ (1) Subscription
              │
              └─→ (N) Product
```

### Índices

- `users.email` - Unique
- `consumers.document_number` - Unique
- `companies.user_id` - Index
- `subscriptions.company_id` - Index
- `products.company_id` - Index
- Soft delete: `deleted_at` - Index en todas las entidades

---

## 🛡️ Validación y Errores

### Validación de DTOs

```go
type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=255"`
    Description string  `json:"description" validate:"max=1000"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    IsVisible   *bool   `json:"is_visible"`
}
```

### Manejo de Errores Centralizado

```go
// Tipos de errores
- ErrValidation     (400)
- ErrNotFound       (404)
- ErrUnauthorized   (401)
- ErrForbidden      (403)
- ErrConflict       (409)
- ErrInternal       (500)
- ErrBadRequest     (400)
- ErrDatabase       (500)
```

**Formato de respuesta:**
```json
{
  "success": false,
  "error": {
    "type": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "name": "This field is required"
    }
  }
}
```

---

## 🔄 Flujo de Datos

### Ejemplo: Crear Producto

```
1. HTTP Request (multipart/form-data)
   ↓
2. ProductHandler.Create()
   - Extrae user_id del context (middleware)
   - Valida DTO
   - Obtiene archivo multipart
   ↓
3. ProductService.Create()
   - Valida suscripción activa (SubscriptionService)
   - Sube foto (FileService)
   - Crea entity Product
   - Guarda en DB (Repository)
   - Retorna DTO
   ↓
4. ProductHandler retorna response JSON
```

### Transacciones

**Caso 1: Crear Empresa**
```go
db.Transaction(func(tx *gorm.DB) error {
    // 1. Crear empresa
    companyRepo.Create(tx, company)
    
    // 2. Crear suscripción
    subscriptionRepo.Create(tx, subscription)
    
    return nil
})
```

**Caso 2: Renovar Suscripción**
```go
db.Transaction(func(tx *gorm.DB) error {
    // 1. Actualizar suscripción (+30 días)
    subscriptionRepo.Update(tx, subscription)
    
    // 2. Activar empresa
    companyRepo.Update(tx, company)
    
    return nil
})
```

---

## 📦 Inyección de Dependencias Manual

```go
// Repositories
authRepo := persistence.NewAuthRepository(db)
productRepo := persistence.NewProductRepository(db)

// Services
fileService := service.NewLocalFileService(cfg.File)
subscriptionSvc := subscription.NewService(subscriptionRepo, companyRepo, db)
productSvc := product.NewService(productRepo, fileService, subscriptionSvc)

// Handlers
productHandler := handler.NewProductHandler(productSvc)

// Routes
routes.SetupProductRoutes(api, productHandler, &cfg.JWT)
```

**Ventajas:**
- ✅ Sin reflexión
- ✅ Type-safe
- ✅ Testeable (fácil mock)
- ✅ Explícito y claro

---

## 🧪 Testing Strategy (Sugerido)

### Unit Tests

```go
// Mockear dependencias
mockRepo := mocks.NewMockProductRepository()
mockFileService := mocks.NewMockFileService()
mockSubService := mocks.NewMockSubscriptionService()

// Inyectar mocks
service := product.NewService(mockRepo, mockFileService, mockSubService)

// Test
result, err := service.Create(ctx, companyID, req, file, header)
```

### Integration Tests

- Database real con testcontainers
- Endpoints completos
- Validación de transacciones

---

## 📊 Monitoreo y Logging

### Structured Logging (slog)

```go
logger.Info("Product created",
    "product_id", product.ID,
    "company_id", companyID,
    "price", product.Price,
)

logger.Error("Failed to upload file",
    "error", err,
    "path", path,
)
```

**Formato Development:**
```
2026-03-02T10:00:00 INFO Product created product_id=uuid company_id=uuid price=99.99
```

**Formato Production (JSON):**
```json
{
  "time": "2026-03-02T10:00:00Z",
  "level": "INFO",
  "msg": "Product created",
  "product_id": "uuid",
  "company_id": "uuid",
  "price": 99.99
}
```

---

## 🚀 Performance

### Optimizaciones Implementadas

1. **Connection Pooling:**
   - MaxOpenConns: 25
   - MaxIdleConns: 5
   - ConnMaxLifetime: 300s

2. **Paginación:**
   - Max 100 items por request
   - Default 10 items

3. **Índices:**
   - FK: company_id, user_id
   - Unique: email, document_number
   - Soft delete: deleted_at

4. **Lazy Loading:**
   - Archivos servidos bajo demanda
   - No preload innecesario de relaciones

---

## 🔧 Configuración

### Centralizada con Viper

```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
    App      AppConfig
    File     FileConfig
}
```

**Prioridad:**
1. Variables de entorno
2. Archivo `.env`
3. Defaults

**Validación en startup:**
```go
if err := config.Validate(); err != nil {
    log.Fatal(err)
}
```

---

## 🎯 Próximos Módulos Sugeridos

### 7. Points (Puntos)

- Acumulación de puntos por compras
- Canje de puntos por productos
- Transacciones (history)
- Validación de suscripción

### 8. Transactions (Transacciones)

- Registro de compras
- Historial de canjes
- Reportes

### 9. Notifications (Notificaciones)

- Email (ya hay interfaz EmailService)
- SMS
- Push notifications

---

## ✅ Checklist de Calidad

- ✅ Clean Architecture
- ✅ SOLID principles
- ✅ DTO Pattern (no entity exposure)
- ✅ Repository Pattern
- ✅ Service Layer
- ✅ Dependency Injection (manual)
- ✅ Context propagation
- ✅ Error handling centralizado
- ✅ Structured logging
- ✅ Input validation
- ✅ JWT authentication
- ✅ Role-based access control
- ✅ Transaction management
- ✅ Soft deletes
- ✅ Pagination
- ✅ File upload validation
- ✅ Unique constraint handling
- ✅ Código compilable
- ✅ Zero business logic en handlers
- ✅ Zero DB access fuera de repositories

---

## 📝 Documentación

- `README.md` - Setup y quickstart
- `API_PRODUCTS.md` - Documentación de API Products
- `ARCHITECTURE.md` - Este archivo
- `.env.example` - Variables de entorno

---

**Este proyecto está listo para producción en términos de arquitectura.**

Próximos pasos recomendados:
1. Tests (unit + integration)
2. Swagger/OpenAPI
3. CI/CD pipeline
4. Monitoring (Prometheus, Grafana)
5. Rate limiting
6. Caching (Redis)
