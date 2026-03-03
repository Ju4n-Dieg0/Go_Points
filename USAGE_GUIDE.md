# Guía de Uso Completa - Go Points API

## 🚀 Inicio del Proyecto

### 1. Configurar Variables de Entorno

```bash
cp .env.example .env
```

Edita `.env` con tus valores:

```env
APP_NAME=Go Points API
APP_ENV=development
LOG_LEVEL=info

SERVER_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_points

JWT_ACCESS_SECRET=your-super-secret-key
JWT_REFRESH_SECRET=your-super-secret-refresh-key

FILE_UPLOAD_DIR=uploads
FILE_MAX_SIZE=5242880
FILE_ALLOWED_TYPES=image/jpeg,image/jpg,image/png,image/webp
```

### 2. Iniciar Base de Datos

```bash
docker-compose up -d postgres
```

### 3. Ejecutar Aplicación

**Desarrollo:**
```bash
make run
```

**Producción (Docker):**
```bash
make docker-up
```

---

## 📋 Flujo Completo de Uso

### Paso 1: Registrar Usuario (Empresa)

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "empresa@example.com",
    "password": "SecurePass123!",
    "name": "Juan Pérez",
    "role": "COMPANY"
  }'
```

**Response:**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "empresa@example.com",
    "name": "Juan Pérez",
    "role": "COMPANY",
    "is_active": true,
    "created_at": "2026-03-02T10:00:00Z"
  },
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 900
}
```

**Guarda el `access_token` para siguientes requests.**

---

### Paso 2: Login (si ya tienes cuenta)

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "empresa@example.com",
    "password": "SecurePass123!"
  }'
```

---

### Paso 3: Crear Empresa

**IMPORTANTE:** Al crear una empresa, automáticamente se crea una suscripción activa de 30 días.

```bash
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Pizza Express",
    "address": "Av. Principal 123",
    "phone": "+51999888777",
    "description": "Pizzería artesanal"
  }'
```

**Response:**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Pizza Express",
  "address": "Av. Principal 123",
  "phone": "+51999888777",
  "description": "Pizzería artesanal",
  "is_active": true,
  "subscription": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "company_id": "660e8400-e29b-41d4-a716-446655440001",
    "start_date": "2026-03-02T10:00:00Z",
    "end_date": "2026-04-01T10:00:00Z",
    "is_active": true,
    "days_remaining": 30
  },
  "created_at": "2026-03-02T10:00:00Z"
}
```

**Guarda el `company_id` - lo necesitarás.**

---

### Paso 4: Crear Producto SIN Foto

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "name=Pizza Margarita" \
  -F "description=Pizza con tomate, mozzarella y albahaca fresca" \
  -F "price=12.99" \
  -F "is_visible=true"
```

**Response:**
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440003",
  "company_id": "660e8400-e29b-41d4-a716-446655440001",
  "name": "Pizza Margarita",
  "description": "Pizza con tomate, mozzarella y albahaca fresca",
  "price": 12.99,
  "photo": "",
  "is_visible": true,
  "created_at": "2026-03-02T10:05:00Z",
  "updated_at": "2026-03-02T10:05:00Z"
}
```

---

### Paso 5: Crear Producto CON Foto

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "name=Pizza Pepperoni" \
  -F "description=Pizza con pepperoni y queso" \
  -F "price=14.99" \
  -F "is_visible=true" \
  -F "photo=@./pepperoni.jpg"
```

**Response:**
```json
{
  "id": "990e8400-e29b-41d4-a716-446655440004",
  "company_id": "660e8400-e29b-41d4-a716-446655440001",
  "name": "Pizza Pepperoni",
  "description": "Pizza con pepperoni y queso",
  "price": 14.99,
  "photo": "uploads/1709377500_abc123-def456.jpg",
  "is_visible": true,
  "created_at": "2026-03-02T10:10:00Z",
  "updated_at": "2026-03-02T10:10:00Z"
}
```

---

### Paso 6: Listar MIS Productos

```bash
curl http://localhost:8080/api/v1/products?page=1&page_size=10 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "products": [
    {
      "id": "990e8400-e29b-41d4-a716-446655440004",
      "company_id": "660e8400-e29b-41d4-a716-446655440001",
      "name": "Pizza Pepperoni",
      "price": 14.99,
      "photo": "uploads/1709377500_abc123-def456.jpg",
      "is_visible": true
    },
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "name": "Pizza Margarita",
      "price": 12.99,
      "is_visible": true
    }
  ],
  "total": 2,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

---

### Paso 7: Actualizar Producto (cambiar precio)

```bash
curl -X PUT http://localhost:8080/api/v1/products/880e8400-e29b-41d4-a716-446655440003 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "price=13.99"
```

---

### Paso 8: Actualizar Producto (cambiar foto)

```bash
curl -X PUT http://localhost:8080/api/v1/products/880e8400-e29b-41d4-a716-446655440003 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "photo=@./nueva-pizza.jpg"
```

**Nota:** La foto anterior se elimina automáticamente.

---

### Paso 9: Ocultar Producto

```bash
curl -X PUT http://localhost:8080/api/v1/products/880e8400-e29b-41d4-a716-446655440003 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "is_visible=false"
```

---

### Paso 10: Buscar Productos

```bash
curl "http://localhost:8080/api/v1/products/search?q=pepperoni&page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### Paso 11: Ver Catálogo Público (SIN autenticación)

```bash
curl http://localhost:8080/api/v1/products/catalog?page=1&page_size=10
```

**Response:** Solo productos con `is_visible=true`

---

### Paso 12: Eliminar Producto

```bash
curl -X DELETE http://localhost:8080/api/v1/products/880e8400-e29b-41d4-a716-446655440003 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Product deleted successfully"
}
```

**Nota:** La foto se elimina del filesystem.

---

### Paso 13: Renovar Suscripción

Cuando tu suscripción esté por expirar:

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions/660e8400-e29b-41d4-a716-446655440001/renew \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "duration_days": 30
  }'
```

**Response:**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "company_id": "660e8400-e29b-41d4-a716-446655440001",
  "start_date": "2026-03-02T10:00:00Z",
  "end_date": "2026-05-01T10:00:00Z",
  "is_active": true,
  "days_remaining": 60,
  "created_at": "2026-03-02T10:00:00Z",
  "updated_at": "2026-03-02T11:00:00Z"
}
```

---

### Paso 14: Refresh Token

Cuando tu access token expire (15 min):

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

**Response:**
```json
{
  "access_token": "new_eyJhbGc...",
  "refresh_token": "new_eyJhbGc...",
  "expires_in": 900
}
```

---

## 🛡️ Manejo de Errores

### Suscripción Expirada

Si intentas crear un producto sin suscripción activa:

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "name=Pizza" \
  -F "price=10"
```

**Response 403:**
```json
{
  "success": false,
  "error": {
    "type": "FORBIDDEN",
    "message": "Company subscription is not active",
    "status_code": 403
  }
}
```

**Solución:** Renovar suscripción (Paso 13)

---

### Archivo Demasiado Grande

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -F "name=Pizza" \
  -F "price=10" \
  -F "photo=@./imagen-10mb.jpg"
```

**Response 400:**
```json
{
  "success": false,
  "error": {
    "type": "VALIDATION_ERROR",
    "message": "file size exceeds maximum allowed size of 5242880 bytes",
    "status_code": 400
  }
}
```

---

### Tipo de Archivo No Permitido

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -F "name=Pizza" \
  -F "price=10" \
  -F "photo=@./documento.pdf"
```

**Response 400:**
```json
{
  "success": false,
  "error": {
    "type": "VALIDATION_ERROR",
    "message": "file type 'application/pdf' is not allowed. Allowed types: image/jpeg, image/jpg, image/png, image/webp",
    "status_code": 400
  }
}
```

---

### Validación de Campos

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "name=Pi" \
  -F "price=-5"
```

**Response 400:**
```json
{
  "success": false,
  "error": {
    "type": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "name": "Value is too short (minimum 3)",
      "price": "Value must be greater than 0"
    },
    "status_code": 400
  }
}
```

---

## 👤 Módulo de Consumidores

### Registrar Consumidor

```bash
curl -X POST http://localhost:8080/api/v1/consumers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "document_type": "DNI",
    "document_number": "12345678",
    "name": "María García",
    "email": "maria@example.com",
    "phone": "+51999777666"
  }'
```

### Buscar Consumidor por Documento

```bash
curl http://localhost:8080/api/v1/consumers/document/12345678 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Listar Consumidores

```bash
curl "http://localhost:8080/api/v1/consumers?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 🔍 Health Check

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "database": "healthy",
  "timestamp": "2026-03-02T10:00:00Z"
}
```

---

## 📊 Ejemplos con Postman

### Collection Structure

```
Go Points API
├── Auth
│   ├── Register (Company)
│   ├── Login
│   ├── Refresh Token
│   └── Logout
├── Companies
│   ├── Create Company
│   ├── List Companies
│   └── Get Company
├── Subscriptions
│   ├── Renew Subscription
│   └── Get Subscription
├── Products
│   ├── Create Product (no photo)
│   ├── Create Product (with photo)
│   ├── List My Products
│   ├── Public Catalog
│   ├── Update Product
│   ├── Search Products
│   └── Delete Product
└── Consumers
    ├── Create Consumer
    ├── List Consumers
    └── Search Consumer
```

### Variables de Entorno en Postman

```
base_url = http://localhost:8080/api/v1
access_token = (auto-set from login)
company_id = (auto-set from create company)
product_id = (auto-set from create product)
```

---

## 🐛 Debugging

### Ver Logs

```bash
docker-compose logs -f app
```

### Logs Estructurados

**Development:**
```
2026-03-02T10:00:00 INFO Product created product_id=uuid price=12.99
2026-03-02T10:00:01 ERROR Failed to upload file error="invalid mime type"
```

**Production (JSON):**
```json
{"time":"2026-03-02T10:00:00Z","level":"INFO","msg":"Product created","product_id":"uuid","price":12.99}
```

---

## 🚦 Rate Limiting (Futuro)

Implementación sugerida:

```go
// 100 requests por minuto por IP
rateLimiter := middleware.RateLimit(100, time.Minute)
app.Use(rateLimiter)
```

---

## 📦 Scripts Útiles

### Makefile Commands

```bash
make build          # Compilar binario
make run            # Ejecutar en desarrollo
make test           # Ejecutar tests (cuando existan)
make docker-build   # Build imagen Docker
make docker-up      # Iniciar con Docker Compose
make docker-down    # Detener Docker Compose
make clean          # Limpiar binarios
```

---

## 🎯 Checklist de Implementación

### Para Empresa:
- [ ] Registrar usuario con role COMPANY
- [ ] Crear empresa (auto-crea suscripción)
- [ ] Crear productos (con/sin foto)
- [ ] Gestionar visibilidad de productos
- [ ] Actualizar precios
- [ ] Renovar suscripción antes de expirar

### Para Consumidor:
- [ ] Registrar usuario con role CONSUMER
- [ ] Ver catálogo público
- [ ] (Futuro) Acumular puntos
- [ ] (Futuro) Canjear puntos

### Para Admin:
- [ ] Registrar como SUPER_ADMIN
- [ ] Gestionar todas las empresas
- [ ] Desactivar empresas
- [ ] Ver reportes

---

## ✅ Testing Manual

### Test 1: Flujo Completo

1. ✅ Registrar usuario COMPANY
2. ✅ Login
3. ✅ Crear empresa
4. ✅ Crear 3 productos con fotos
5. ✅ Actualizar precio de un producto
6. ✅ Ocultar un producto
7. ✅ Ver catálogo público (debe mostrar solo 2 visibles)
8. ✅ Eliminar un producto
9. ✅ Verificar que la foto se eliminó

### Test 2: Validaciones

1. ✅ Intentar crear producto sin nombre → Error 400
2. ✅ Intentar subir PDF → Error 400
3. ✅ Intentar subir imagen de 10MB → Error 400
4. ✅ Intentar crear producto con precio negativo → Error 400

### Test 3: Autorización

1. ✅ Intentar crear producto sin token → Error 401
2. ✅ Intentar crear producto con role CONSUMER → Error 403
3. ✅ Intentar crear producto con suscripción expirada → Error 403

---

## 🎓 Mejores Prácticas

1. **Siempre renovar suscripción antes de expirar**
2. **Usar fotos optimizadas (< 1MB)**
3. **Validar datos antes de enviar**
4. **Manejar errores en el cliente**
5. **Refresh token antes de que expire**
6. **Logout al cerrar sesión**

---

**¡El sistema está listo para usar!** 🚀
