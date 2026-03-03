# Go Points API - Enterprise Features

## 🎯 Descripción General

Go Points API es una solución empresarial completa para sistemas de puntos de fidelización, diseñada con arquitectura limpia y lista para integraciones externas.

## 🚀 Características Enterprise

### ✅ Versionado de API
- **Prefijo global:** `/api/v1`
- **Estructura preparada para futuras versiones** (`/api/v2`, `/api/v3`)
- **Ruta sin hardcodeo:** Configuración centralizada

### ✅ Multi-Ambiente
- **Desarrollo:** `.env.development`
- **Staging:** `.env.staging`
- **Producción:** `.env.production`

Carga automática según variable `APP_ENV`:
```bash
export APP_ENV=production
```

### ✅ Health Checks
| Endpoint | Propósito |
|----------|-----------|
| `GET /health` | Verifica que la aplicación esté ejecutándose |
| `GET /ready` | Verifica que la aplicación y DB estén listas |
| `GET /live` | Liveness check para Kubernetes |

### ✅ Response Envelope Unificado

**Respuesta Exitosa:**
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "pagination": { ... }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Respuesta de Error:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Los datos proporcionados no son válidos",
    "details": { ... }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### ✅ Rate Limiting Documentado

Headers RFC 6585:
- `X-RateLimit-Limit`: Límite de requests
- `X-RateLimit-Remaining`: Requests restantes
- `X-RateLimit-Reset`: Tiempo de reset
- `Retry-After`: Segundos para reintentar

Configuración por ENV:
```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_AUTH_REQUESTS=5
RATE_LIMIT_AUTH_WINDOW=60
RATE_LIMIT_GENERAL_REQUESTS=100
RATE_LIMIT_GENERAL_WINDOW=60
```

### ✅ Documentación OpenAPI 3.0

**Archivo Limpio:** `openapi.yaml` en raíz del proyecto

Compatible con:
- ✅ Postman
- ✅ Insomnia
- ✅ Stoplight
- ✅ SwaggerHub
- ✅ Swagger Editor
- ✅ Redoc

**Generar:**
```bash
make swagger
go run scripts/export_openapi.go
```

### ✅ Postman Collection v2.1

Incluye:
- Todas las rutas organizadas por módulo
- Variables globales (`{{base_url}}`, `{{access_token}}`)
- Ejemplos de body con datos reales
- Carpetas por módulo: Auth, Companies, Consumers, Products, Points, Rewards

**Generar:**
```bash
go run scripts/generate_postman.go
```

## 📁 Estructura del Proyecto

```
Go_Points/
├── cmd/
│   └── server/
│       └── main.go              # Entry point con anotaciones Swagger
├── internal/
│   ├── application/             # Casos de uso
│   ├── domain/                  # Entidades y lógica de negocio
│   ├── infrastructure/
│   │   ├── http/
│   │   │   ├── handler/        # HTTP handlers
│   │   │   └── routes/         # Definición de rutas
│   │   └── persistence/        # Repositorios GORM
│   ├── shared/
│   │   ├── dto/
│   │   │   ├── envelope.go     # Response envelope unificado
│   │   │   ├── common.go       # DTOs comunes
│   │   │   └── swagger_models.go
│   │   ├── middleware/
│   │   │   ├── rate_limiter.go # Rate limiting con headers
│   │   │   ├── cors.go         # CORS configurable
│   │   │   └── security.go     # Security headers
│   │   └── validation/         # Validadores personalizados
│   └── config/
│       └── config.go            # Configuración multi-ambiente
├── scripts/
│   ├── generate_postman.go     # Genera Postman collection
│   └── export_openapi.go       # Exporta OpenAPI limpio
├── docs/                        # Swagger generado (auto)
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI/CD
├── .env.development
├── .env.staging
├── .env.production
├── .env.example
├── openapi.yaml                 # OpenAPI limpio (generado)
├── Makefile                     # Comandos útiles
└── README.md
```

## 🛠️ Comandos Make

```bash
make help              # Mostrar ayuda
make install           # Instalar dependencias
make build             # Compilar aplicación
make run               # Ejecutar aplicación
make dev               # Ejecutar con hot-reload
make test              # Ejecutar tests
make coverage          # Tests con cobertura
make lint              # Ejecutar linters
make swagger           # Generar Swagger docs
make clean             # Limpiar archivos compilados
make docker-build      # Construir imagen Docker
make all               # Ejecutar todas las verificaciones
```

## 🔐 Seguridad

### Autenticación JWT
- Access Token: 15 minutos (configurable)
- Refresh Token: 7 días (configurable)
- Bearer Authentication en Swagger

### CORS Configurable
```bash
CORS_ALLOWED_ORIGINS=https://app.gopoints.com,https://admin.gopoints.com
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=86400
```

### Security Headers
- `X-Frame-Options: DENY`
- `X-Content-Type-Options: nosniff`
- `X-XSS-Protection: 1; mode=block`
- `Content-Security-Policy`
- `Strict-Transport-Security` (HTTPS)

### Bcrypt
```bash
BCRYPT_COST=14              # Producción
MAX_LOGIN_ATTEMPTS=3
LOCKOUT_DURATION=30
```

## 🚀 Despliegue

### Variables de Entorno Críticas

**Producción:**
```bash
APP_ENV=production
JWT_ACCESS_SECRET=<64+ caracteres aleatorios>
JWT_REFRESH_SECRET=<64+ caracteres aleatorios>
DB_SSLMODE=require
BCRYPT_COST=14
RATE_LIMIT_ENABLED=true
CORS_ALLOWED_ORIGINS=https://gopoints.com
```

### Docker

```bash
# Build
make docker-build

# Run
docker run -p 8080:8080 --env-file .env.production go-points-api:latest
```

### Kubernetes Health Checks

```yaml
livenessProbe:
  httpGet:
    path: /live
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

## 📊 CI/CD Pipeline

GitHub Actions incluye:
- ✅ Linting (golangci-lint)
- ✅ Testing con PostgreSQL
- ✅ Coverage report (Codecov)
- ✅ Security scan (Gosec, Trivy)
- ✅ Build para múltiples plataformas
- ✅ Docker build y push
- ✅ Swagger generation

## 📝 Documentación API

### Swagger UI
```
http://localhost:8080/docs/index.html
```

### Postman
1. Importar `Go_Points_API.postman_collection.json`
2. Configurar variable `base_url`
3. Hacer login y copiar `access_token`

### OpenAPI
Importar `openapi.yaml` en:
- Postman
- Insomnia
- Stoplight Studio
- SwaggerHub

## 🔧 Configuración Recomendada por Ambiente

| Variable | Development | Staging | Production |
|----------|------------|---------|------------|
| `LOG_LEVEL` | `debug` | `info` | `warn` |
| `BCRYPT_COST` | `10` | `12` | `14` |
| `DB_MAX_OPEN_CONNS` | `25` | `50` | `100` |
| `RATE_LIMIT_AUTH_REQUESTS` | `10` | `5` | `5` |
| `JWT_ACCESS_EXPIRATION` | `1800` | `900` | `900` |

## 📞 Soporte

Para integraciones enterprise y soporte técnico:
- Email: support@gopoints.com
- Documentación: https://docs.gopoints.com
- API Status: https://status.gopoints.com

## 📄 Licencia

MIT License - Ver LICENSE para más detalles.
