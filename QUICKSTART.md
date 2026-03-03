# 🚀 Quick Start - Go Points API Enterprise

## Inicio Rápido (3 minutos)

### 1. Configuración Inicial

```bash
# Clonar o navegar al proyecto
cd Go_Points

# Copiar configuración de ejemplo
cp .env.example .env

# Editar .env (opcional para desarrollo local)
# Los valores por defecto funcionan con PostgreSQL local
```

### 2. Iniciar PostgreSQL

```bash
# Opción A: Docker (recomendado)
docker run -d \
  --name postgres-gopoints \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=go_points \
  -p 5432:5432 \
  postgres:16

# Opción B: PostgreSQL local
# Crear base de datos: CREATE DATABASE go_points;
```

### 3. Ejecutar API

```bash
# Instalar dependencias
make install

# Ejecutar (las migraciones se ejecutan automáticamente)
make run
```

### 4. Verificar

```bash
# Health check
curl http://localhost:8080/health

# Swagger UI
# Abrir en navegador: http://localhost:8080/docs/index.html
```

---

## Uso del API

### Registrar Usuario

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "password": "SecurePass123!",
    "role": "consumer"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan@example.com",
    "password": "SecurePass123!"
  }'
```

Respuesta:
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "user": { ... }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Usar Token

```bash
# Guardar token
export TOKEN="tu-access-token-aqui"

# Obtener perfil
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

---

## Ambientes

### Desarrollo
```bash
export APP_ENV=development
make dev  # Con hot-reload
```

### Staging
```bash
export APP_ENV=staging
make run
```

### Producción
```bash
export APP_ENV=production
make build
./bin/server
```

---

## Generar Documentación

### Swagger
```bash
make swagger
# Acceder: http://localhost:8080/docs
```

### OpenAPI
```bash
go run scripts/export_openapi.go
# Genera: openapi.yaml
```

### Postman Collection
```bash
go run scripts/generate_postman.go
# Genera: Go_Points_API.postman_collection.json
# Importar en Postman
```

---

## Comandos Útiles

```bash
# Ver todos los comandos
make help

# Desarrollo con hot-reload
make dev

# Tests
make test

# Tests con coverage
make coverage

# Linter
make lint

# Formatear código
make fmt

# Todo junto (lint + test + build)
make all

# Limpiar
make clean
```

---

## Importar en Postman

1. Generar colección:
```bash
go run scripts/generate_postman.go
```

2. En Postman:
   - File → Import
   - Seleccionar `Go_Points_API.postman_collection.json`
   - Configurar variables:
     - `base_url`: `http://localhost:8080/api/v1`
     - `access_token`: (obtenido del login)

3. Hacer login y copiar el `access_token` a la variable

4. ¡Listo! Todos los endpoints están disponibles con ejemplos

---

## Docker (Opcional)

### Build
```bash
make docker-build
```

### Run
```bash
docker run -p 8080:8080 \
  --env-file .env.production \
  go-points-api:latest
```

### Docker Compose
```bash
make docker-compose-up
```

---

## Health Checks

```bash
# Basic health
curl http://localhost:8080/health

# Readiness (incluye DB)
curl http://localhost:8080/ready

# Liveness (Kubernetes)
curl http://localhost:8080/live
```

---

## Estructura de Respuestas

### Success
```json
{
  "success": true,
  "data": { ... },
  "meta": { "pagination": { ... } },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Datos inválidos",
    "details": { ... }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Rate Limiting

El API incluye rate limiting con headers:
- `X-RateLimit-Limit`: Límite de requests
- `X-RateLimit-Remaining`: Requests restantes
- `X-RateLimit-Reset`: Cuando se resetea
- `Retry-After`: Segundos para reintentar

Por defecto:
- Auth endpoints: 5 requests/minuto
- Otros endpoints: 100 requests/minuto

---

## Troubleshooting

### Puerto ya en uso
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:8080 | xargs kill -9
```

### Error de conexión a DB
```bash
# Verificar PostgreSQL
docker ps  # Si usas Docker
psql -U postgres -d go_points  # Si es local

# Verificar variables en .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_points
```

### Swagger no se genera
```bash
# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Verificar que esté en PATH
which swag  # Linux/Mac
where swag  # Windows
```

---

## Enlaces Rápidos

- **Swagger UI**: http://localhost:8080/docs/index.html
- **Health**: http://localhost:8080/health
- **Ready**: http://localhost:8080/ready
- **API Base**: http://localhost:8080/api/v1

---

## Documentación Completa

- [README.md](README.md) - Overview del proyecto
- [ENTERPRISE.md](ENTERPRISE.md) - Features enterprise
- [CHANGELOG.md](CHANGELOG.md) - Cambios detallados
- [SWAGGER.md](SWAGGER.md) - Guía de Swagger
- [ENTERPRISE_SUMMARY.md](ENTERPRISE_SUMMARY.md) - Resumen ejecutivo

---

**¿Listo para producción?** ✅

El proyecto incluye:
- ✅ Multi-ambiente (dev, staging, prod)
- ✅ CI/CD con GitHub Actions
- ✅ Health checks para Kubernetes
- ✅ Rate limiting profesional
- ✅ Documentación OpenAPI 3.0
- ✅ Postman collection
- ✅ Security headers
- ✅ Response envelope estándar

**¡Comienza a desarrollar!** 🚀
