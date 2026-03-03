# Módulo de Paginación

Sistema genérico de paginación, ordenamiento, búsqueda y filtrado para APIs REST.

## Características

✅ **Paginación**: page, limit con valores por defecto  
✅ **Ordenamiento**: sort, order (asc/desc)  
✅ **Búsqueda**: Búsqueda ILIKE en múltiples campos  
✅ **Filtros dinámicos**: Operadores avanzados (=, >=, >, <=, <, !=, LIKE, IN, NOT IN, NULL)  
✅ **Genérico**: Funciona con cualquier modelo GORM  
✅ **Seguro**: Validación de nombres de campos para prevenir SQL injection  
✅ **Type-safe**: Uso de Generics de Go 1.18+  
✅ **Respuesta estandarizada**: JSON con data y meta  

## Instalación

```bash
go get gorm.io/gorm
go get github.com/gofiber/fiber/v3
```

## Uso Básico

### 1. En tu Handler

```go
package product

import (
    "github.com/gofiber/fiber/v3"
    "your-project/internal/shared/pagination"
)

func (h *Handler) List(c fiber.Ctx) error {
    ctx := c.Context()
    
    // Extraer parámetros de paginación de la query string
    params := pagination.NewParams(c)
    
    // Obtener datos paginados del servicio
    products, meta, err := h.service.List(ctx, params)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    // Responder con formato estandarizado
    return c.Status(fiber.StatusOK).JSON(pagination.NewResponse(products, meta))
}
```

### 2. En tu Service

```go
package product

import (
    "context"
    "your-project/internal/shared/pagination"
)

func (s *Service) List(ctx context.Context, params *pagination.Params) ([]ProductResponse, pagination.Meta, error) {
    // Obtener entidades paginadas del repositorio
    products, meta, err := s.repo.List(ctx, params)
    if err != nil {
        return nil, pagination.Meta{}, err
    }
    
    // Convertir a DTOs
    responses := make([]ProductResponse, len(products))
    for i, p := range products {
        responses[i] = ToProductResponse(&p)
    }
    
    return responses, meta, nil
}
```

### 3. En tu Repository

```go
package product

import (
    "context"
    "gorm.io/gorm"
    "your-project/internal/shared/pagination"
)

func (r *Repository) List(ctx context.Context, params *pagination.Params) ([]Product, pagination.Meta, error) {
    // Campos en los que se puede buscar
    searchFields := []string{"name", "description", "code"}
    
    // Relaciones a precargar
    preloads := []string{"Company", "ProductPhotos"}
    
    // Usar helper genérico de paginación
    return pagination.Paginate[Product](
        r.db.WithContext(ctx),
        params,
        searchFields,
        preloads...,
    )
}
```

## Parámetros de Query String

### Paginación

```
GET /api/products?page=2&limit=20
```

- `page`: Número de página (default: 1, min: 1)
- `limit`: Registros por página (default: 10, min: 1, max: 100)

### Ordenamiento

```
GET /api/products?sort=name&order=asc
GET /api/products?sort=created_at&order=desc
```

- `sort`: Campo por el cual ordenar
- `order`: Dirección (asc o desc, default: desc)

### Búsqueda

```
GET /api/products?search=laptop
```

- `search`: Término de búsqueda (aplica ILIKE a los campos especificados en `searchFields`)

### Filtros Simples

```
GET /api/products?filter[status]=active
GET /api/products?filter[company_id]=123e4567-e89b-12d3-a456-426614174000
```

### Filtros con Operadores

#### Mayor o igual (gte)
```
GET /api/products?filter[price__gte]=100
```

#### Mayor que (gt)
```
GET /api/products?filter[stock__gt]=0
```

#### Menor o igual (lte)
```
GET /api/products?filter[price__lte]=500
```

#### Menor que (lt)
```
GET /api/products?filter[discount__lt]=50
```

#### No igual (ne)
```
GET /api/products?filter[status__ne]=deleted
```

#### LIKE (like)
```
GET /api/products?filter[code__like]=ABC
# Busca: %ABC%
```

#### IN (in)
```
GET /api/products?filter[status__in]=active,pending,draft
```

#### NOT IN (notin)
```
GET /api/products?filter[status__notin]=deleted,archived
```

#### IS NULL / IS NOT NULL (null)
```
GET /api/products?filter[deleted_at__null]=true   # IS NULL
GET /api/products?filter[deleted_at__null]=false  # IS NOT NULL
```

### Combinaciones

```
GET /api/products?page=1&limit=20&sort=created_at&order=desc&search=laptop&filter[status]=active&filter[price__gte]=100&filter[price__lte]=500
```

## Respuesta JSON

```json
{
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Product 1",
      "price": 150
    },
    {
      "id": "223e4567-e89b-12d3-a456-426614174001",
      "name": "Product 2",
      "price": 200
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 200,
    "total_pages": 20
  }
}
```

## Ejemplos Avanzados

### Repository con Query Previa

```go
func (r *Repository) ListByCompany(ctx context.Context, companyID uuid.UUID, params *pagination.Params) ([]Product, pagination.Meta, error) {
    // Base query con filtro de empresa
    query := r.db.WithContext(ctx).Where("company_id = ?", companyID)
    
    // Aplicar paginación
    return pagination.Paginate[Product](
        query,
        params,
        []string{"name", "description"},
        "Company", "ProductPhotos",
    )
}
```

### Búsqueda en Campos Relacionados

```go
func (r *Repository) ListWithCompanySearch(ctx context.Context, params *pagination.Params) ([]Product, pagination.Meta, error) {
    // Join con tabla relacionada
    query := r.db.WithContext(ctx).
        Joins("JOIN companies ON companies.id = products.company_id")
    
    // Búsqueda en campos de ambas tablas
    searchFields := []string{
        "products.name",
        "products.description",
        "companies.name",
    }
    
    return pagination.Paginate[Product](
        query,
        params,
        searchFields,
    )
}
```

### Filtros Personalizados

```go
func (r *Repository) ListActive(ctx context.Context, params *pagination.Params) ([]Product, pagination.Meta, error) {
    // Aplicar filtros de negocio
    query := r.db.WithContext(ctx).
        Where("deleted_at IS NULL").
        Where("status = ?", "active").
        Where("stock > ?", 0)
    
    return pagination.Paginate[Product](
        query,
        params,
        []string{"name", "code"},
    )
}
```

### Uso Manual (sin helper)

```go
func (r *Repository) ListManual(ctx context.Context, params *pagination.Params) ([]Product, pagination.Meta, error) {
    var products []Product
    var total int64
    
    // Base query
    query := r.db.WithContext(ctx).Where("deleted_at IS NULL")
    
    // Aplicar búsqueda, filtros, ordenamiento
    searchFields := []string{"name", "description"}
    query = params.ApplyToQuery(query, searchFields)
    
    // Contar total (antes de aplicar paginación)
    countQuery := query.Session(&gorm.Session{})
    if err := countQuery.Model(&Product{}).Count(&total).Error; err != nil {
        return nil, pagination.Meta{}, err
    }
    
    // Obtener resultados
    if err := query.Find(&products).Error; err != nil {
        return nil, pagination.Meta{}, err
    }
    
    // Calcular metadata
    meta := params.CalculateMeta(total)
    
    return products, meta, nil
}
```

## Validación y Seguridad

### Nombres de Campo

El sistema valida automáticamente los nombres de campo para prevenir SQL injection:

✅ **Permitidos**: Letras (a-z, A-Z), números (0-9), guión bajo (_)  
❌ **Bloqueados**: Espacios, puntos, punto y coma, comillas, etc.

```go
// ✅ VÁLIDO
filter[status]=active
filter[created_at__gte]=2024-01-01
filter[company_id]=123

// ❌ INVÁLIDO (ignorados automáticamente)
filter[status; DROP TABLE]=active
filter[user.password]=test
filter[' OR '1'='1]=value
```

### Límites

- **Página mínima**: 1 (si se envía 0 o negativo, se usa 1)
- **Límite mínimo**: 1 (se fuerza a 10 si es menor)
- **Límite máximo**: 100 (se fuerza a 100 si es mayor)

Puedes modificar estas constantes en `pagination.go`:

```go
const (
    DefaultPage  = 1
    DefaultLimit = 10
    MaxLimit     = 100
    MinLimit     = 1
)
```

## Testing

El módulo incluye suite completa de tests:

```bash
go test ./internal/shared/pagination/... -v
```

Tests incluidos:
- ✅ Parámetros por defecto
- ✅ Validación de valores
- ✅ Cálculo de offset
- ✅ Ordenamiento
- ✅ Metadata
- ✅ Búsqueda
- ✅ Filtros con operadores
- ✅ Validación de nombres de campo
- ✅ Paginación genérica con GORM

## Arquitectura

```
internal/shared/pagination/
├── pagination.go          # Implementación principal
├── pagination_test.go     # Suite de tests
└── README.md             # Esta documentación
```

### Componentes

**Params**: Encapsula parámetros de paginación, ordenamiento, búsqueda y filtros  
**Meta**: Metadata de paginación (página actual, total, páginas totales)  
**Response**: Estructura de respuesta estandarizada  
**Paginate[T]**: Helper genérico para paginar cualquier modelo  

## Ventajas

✅ **Reutilizable**: Un solo código para todos los módulos  
✅ **Type-safe**: Uso de generics de Go 1.18+  
✅ **Flexible**: Soporta búsqueda, filtros, ordenamiento  
✅ **Seguro**: Validación de campos  
✅ **Testeable**: Suite completa de tests  
✅ **Documentado**: Ejemplos y casos de uso  
✅ **Productivo**: Reduce código repetitivo  
✅ **Mantenible**: Lógica centralizada  

## Futuras Mejoras

- [ ] Soporte para cursores (cursor-based pagination)
- [ ] Caché de conteos totales
- [ ] Agregaciones (sum, avg, min, max)
- [ ] Filtros por rangos de fechas
- [ ] Exportación a CSV/Excel
- [ ] GraphQL support
- [ ] Filtros guardados (saved filters)
- [ ] Profiles de ordenamiento predefinidos

---

**Creado con Clean Architecture** | Compatible con Fiber v3 + GORM v2
