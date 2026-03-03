# Configuración Swagger/OpenAPI 3 - Go Points API

## Instalación de swag CLI

```powershell
# Instalar swag globalmente
go install github.com/swaggo/swag/cmd/swag@latest

# Verificar instalación
swag --version
```

## Generar Documentación

```powershell
# Desde la raíz del proyecto
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Flags explicados:
# -g: archivo main con anotaciones generales
# -o: directorio de salida (docs/)
# --parseDependency: parsear dependencias
# --parseInternal: parsear paquetes internos
```

## Estructura Generada

```
Go_Points/
├── docs/
│   ├── docs.go          # Código Go generado
│   ├── swagger.json     # Especificación JSON
│   └── swagger.yaml     # Especificación YAML
├── cmd/
│   └── server/
│       └── main.go      # Anotaciones generales @title, @version, etc.
└── internal/
    ├── domain/
    │   └── */dto.go     # Modelos documentados
    └── infrastructure/
        └── http/
            └── handler/ # Handlers con anotaciones @Summary, @Tags, etc.
```

## Acceder a Swagger UI

Una vez el servidor esté corriendo:

```
http://localhost:8080/docs/index.html
```

## Anotaciones Principales

### En main.go (Información General)

```go
// @title Go Points API
// @version 1.0
// @description API REST para sistema de puntos

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### En Handlers (Endpoints)

```go
// @Summary Breve descripción
// @Description Descripción detallada
// @Tags nombre-del-tag
// @Accept json
// @Produce json
// @Param request body dto.Request true "Descripción del body"
// @Param id path string true "ID del recurso"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /ruta [método]
```

### En DTOs (Modelos)

```go
type Request struct {
    // Descripción del campo
    Field string `json:"field" example:"valor" validate:"required"`
} // @name NombreDelModelo
```

## Tags Definidos

- **Auth**: Autenticación y gestión de tokens
- **Companies**: Gestión de empresas
- **Subscriptions**: Gestión de suscripciones
- **Consumers**: Gestión de consumidores
- **Products**: Gestión de productos
- **Points**: Sistema de puntos
- **Rewards**: Gestión de recompensas

## Seguridad JWT

Todos los endpoints protegidos incluyen:

```go
// @Security BearerAuth
```

Para usar en Swagger UI:
1. Hacer login en `/api/v1/auth/login`
2. Copiar el `access_token` de la respuesta
3. Hacer clic en "Authorize" (candado)
4. Ingresar: `Bearer {access_token}`
5. Hacer clic en "Authorize"

## Modelos de Error Estándar

```json
{
  "error": "ERROR_TYPE",
  "message": "Error message",
  "details": {
    "field": "error detail"
  }
}
```

## Paginación

Todos los endpoints con listado soportan:

```
?page=1&limit=10&sort=created_at&order=desc&search=term
```

Respuesta:
```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

## Filtros Dinámicos

```
?filter[status]=active
?filter[price__gte]=100
?filter[price__lte]=500
?filter[status__in]=active,pending
```

Operadores:
- `=`: Igual
- `__gte`: Mayor o igual
- `__gt`: Mayor que
- `__lte`: Menor o igual
- `__lt`: Menor que
- `__ne`: No igual
- `__like`: LIKE (case insensitive)
- `__in`: IN (valores separados por coma)
- `__notin`: NOT IN
- `__null`: IS NULL (true/false)

## Regenerar Documentación

Cada vez que se modifiquen anotaciones:

```powershell
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

## Validar Swagger

```powershell
# Validar swagger.yaml
swag fmt

# Ver especificación
cat docs/swagger.json | jq
```

## Integración con Frontend

La especificación generada puede ser usada con:

- **TypeScript**: openapi-generator, swagger-typescript-api
- **React**: swagger-codegen
- **Postman**: Importar swagger.json directamente
- **Insomnia**: Importar swagger.json

## Ejemplo de Uso Completo

1. **Generar docs**:
   ```powershell
   swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
   ```

2. **Iniciar servidor**:
   ```powershell
   go run cmd/server/main.go
   ```

3. **Abrir Swagger UI**:
   ```
   http://localhost:8080/docs/index.html
   ```

4. **Autenticarse**:
   - POST `/api/v1/auth/login`
   - Copiar `access_token`
   - Click "Authorize"
   - Ingresar: `Bearer {token}`

5. **Probar endpoints**:
   - Navegar por tags
   - Expandir endpoint
   - Click "Try it out"
   - Completar parámetros
   - Click "Execute"

## Troubleshooting

### Error: "swag: command not found"

```powershell
# Agregar GOPATH/bin al PATH
$env:PATH += ";$env:GOPATH\bin"

# O reinstalar
go install github.com/swaggo/swag/cmd/swag@latest
```

### Error: "docs package not found"

```powershell
# Asegurar que existe docs/docs.go
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Verificar import en main.go
# _ "github.com/Ju4n-Dieg0/Go_Points/docs"
```

### Error: "Failed to parse"

```powershell
# Formatear anotaciones
swag fmt

# Verificar sintaxis de anotaciones
# Cada línea debe empezar con //
# Sin espacios extras
```

## Scripts Útiles

### Regenerar y Ejecutar

```powershell
# PowerShell script
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
if ($?) { go run cmd/server/main.go }
```

### Ver Cambios en Tiempo Real

```powershell
# Usar air para hot reload
go install github.com/cosmtrek/air@latest
air
```

## Recursos

- [Swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [Go Swag Examples](https://github.com/swaggo/swag#declarative-comments-format)

---

**Nota**: La documentación se regenera automáticamente cada vez que se ejecuta `swag init`. No editar archivos en `docs/` manualmente.
