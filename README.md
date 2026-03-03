# Go Points API

Backend profesional en Go usando Clean Architecture, Fiber v3, GORM y PostgreSQL.

## 🚀 Tecnologías

- **Go 1.25+**
- **Fiber v3** - Framework web
- **GORM v2** - ORM para PostgreSQL
- **PostgreSQL** - Base de datos
- **JWT** - Autenticación (access + refresh tokens)
- **UUID v7** - Identificadores únicos
- **Viper** - Gestión de configuración
- **Slog** - Logger estructurado nativo de Go
- **Docker** - Containerización

## 📁 Estructura del Proyecto

```
.
├── cmd/
│   └── server/
│       └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── config/
│   │   └── config.go            # Configuración con Viper
│   ├── database/
│   │   └── postgres.go          # Conexión a PostgreSQL con GORM
│   └── shared/
│       ├── errors/
│       │   ├── errors.go        # Errores personalizados
│       │   └── handler.go       # Middleware de manejo de errores
│       ├── logger/
│       │   └── logger.go        # Logger estructurado con slog
│       └── middleware/
│           └── middleware.go    # Middlewares globales
├── .env.example                 # Variables de entorno de ejemplo
├── .gitignore
├── docker-compose.yml           # Orquestación de servicios
├── Dockerfile                   # Imagen multistage
├── go.mod
├── go.sum
├── Makefile                     # Comandos de desarrollo
└── README.md
```

## 🏗️ Arquitectura

El proyecto sigue **Clean Architecture** con las siguientes capas:

- **Domain**: Entidades y lógica de negocio (a implementar)
- **Application**: Casos de uso y servicios (a implementar)
- **Infrastructure**: Database, config, logger, middlewares
- **Presentation**: Handlers HTTP (a implementar)

### Principios aplicados:

- ✅ Separación de responsabilidades
- ✅ Inyección de dependencias manual
- ✅ DTO Pattern (no exponer entidades)
- ✅ Repository Pattern
- ✅ Service Layer
- ✅ Manejo centralizado de errores
- ✅ Logger estructurado
- ✅ Configuración centralizada

## 🔧 Configuración

### Requisitos previos

- Go 1.25+
- Docker y Docker Compose
- Make (opcional pero recomendado)

### Variables de entorno

Copia el archivo de ejemplo y ajusta los valores:

```bash
cp .env.example .env
```

Las variables principales son:

- `APP_NAME`: Nombre de la aplicación
- `APP_ENV`: Entorno (development, production)
- `SERVER_PORT`: Puerto del servidor
- `DB_HOST`: Host de PostgreSQL
- `DB_NAME`: Nombre de la base de datos
- `JWT_ACCESS_SECRET`: Secreto para tokens de acceso
- `JWT_REFRESH_SECRET`: Secreto para tokens de refresh

## 🚀 Inicio rápido

### Opción 1: Con Docker (recomendado)

```bash
# Iniciar todos los servicios
make docker-up

# Ver logs
make docker-logs

# Detener servicios
make docker-down
```

### Opción 2: Local

```bash
# Instalar dependencias
make install

# Asegurarse de que PostgreSQL esté corriendo localmente

# Ejecutar la aplicación
make run

# O en modo desarrollo con hot reload
make dev
```

## 📝 Comandos disponibles

```bash
make help           # Ver todos los comandos disponibles
make install        # Instalar dependencias
make build          # Compilar la aplicación
make run            # Ejecutar la aplicación
make dev            # Modo desarrollo con hot reload
make test           # Ejecutar tests
make lint           # Ejecutar linters
make clean          # Limpiar archivos generados
make docker-build   # Construir imagen Docker
make docker-up      # Iniciar servicios con Docker Compose
make docker-down    # Detener servicios
```

## 🔍 Endpoints disponibles

### Health Check

```http
GET /health
```

Respuesta exitosa:
```json
{
  "status": "ok",
  "database": "healthy",
  "timestamp": "2026-03-02T10:00:00Z"
}
```

### API Base

```http
GET /api/v1/
```

Respuesta:
```json
{
  "message": "Welcome to Go Points API",
  "version": "1.0.0",
  "status": "running"
}
```

## 🛠️ Desarrollo

### Agregar nuevos módulos de negocio

Cuando agregues módulos (ej: usuarios, autenticación), sigue esta estructura:

```
internal/
├── domain/
│   └── user/
│       ├── entity.go       # Entidad de dominio
│       ├── repository.go   # Interface del repositorio
│       └── dto.go          # DTOs (Request/Response)
├── application/
│   └── user/
│       └── service.go      # Lógica de negocio
└── infrastructure/
    ├── persistence/
    │   └── user_repository.go  # Implementación del repositorio
    └── http/
        └── handler/
            └── user_handler.go  # Handlers HTTP
```

### Reglas de negocio

- ❌ **NO** escribir lógica en handlers
- ❌ **NO** acceder a DB fuera de repositories
- ❌ **NO** exponer entidades directamente
- ✅ **SÍ** usar DTOs para requests y responses
- ✅ **SÍ** usar context.Context en todas las operaciones
- ✅ **SÍ** manejar errores con los tipos personalizados

## 🐳 Docker

### Construir imagen

```bash
docker build -t go-points-api:latest .
```

### Ejecutar con Docker Compose

```bash
docker-compose up -d
```

Servicios incluidos:
- **app**: Aplicación Go en puerto 8080
- **postgres**: PostgreSQL 16 en puerto 5432

## 📊 Logs

Los logs están estructurados usando `slog`:

- **Development**: Formato texto legible
- **Production**: Formato JSON para agregación

Niveles de log: `debug`, `info`, `warn`, `error`

## 🔐 Seguridad

- Headers de seguridad configurados
- CORS habilitado
- JWT para autenticación (a implementar)
- Passwords con bcrypt (a implementar)
- Input validation con go-playground/validator (a implementar)

## 📄 Licencia

MIT

## 👥 Contribución

Este es un proyecto base. Al agregar features:

1. Mantén la arquitectura limpia
2. Sigue los principios SOLID
3. Escribe tests
4. Documenta los cambios
5. Usa DTOs siempre

---

**Próximos pasos sugeridos:**

1. Implementar módulo de autenticación
2. Implementar módulo de usuarios
3. Agregar validación de DTOs
4. Implementar tests unitarios
5. Agregar documentación con Swagger
6. Configurar CI/CD
