package pagination

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Params representa los parámetros de paginación
type Params struct {
	Page   int               `json:"page"`
	Limit  int               `json:"limit"`
	Sort   string            `json:"sort,omitempty"`
	Order  string            `json:"order,omitempty"`
	Search string            `json:"search,omitempty"`
	Filter map[string]string `json:"filter,omitempty"`
}

// Meta contiene metadata de paginación
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Response es la respuesta genérica paginada
type Response struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

// DefaultParams valores por defecto
const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
	MinLimit     = 1
)

// NewParams crea parámetros de paginación desde un contexto Fiber
func NewParams(c fiber.Ctx) *Params {
	page, _ := strconv.Atoi(c.Query("page", strconv.Itoa(DefaultPage)))
	limit, _ := strconv.Atoi(c.Query("limit", strconv.Itoa(DefaultLimit)))
	sort := c.Query("sort", "")
	order := c.Query("order", "desc")
	search := c.Query("search", "")

	// Normalizar valores
	if page < 1 {
		page = DefaultPage
	}
	if limit < MinLimit {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	// Normalizar orden
	order = strings.ToLower(order)
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Procesar filtros (filter[campo]=valor)
	filters := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		keyStr := string(key)
		if strings.HasPrefix(keyStr, "filter[") && strings.HasSuffix(keyStr, "]") {
			// Extraer nombre del campo: filter[status] -> status
			field := keyStr[7 : len(keyStr)-1]
			filters[field] = string(value)
		}
	})

	return &Params{
		Page:   page,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
		Search: search,
		Filter: filters,
	}
}

// GetOffset calcula el offset para la query
func (p *Params) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// GetOrderBy retorna la cláusula ORDER BY
func (p *Params) GetOrderBy() string {
	if p.Sort == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", p.Sort, strings.ToUpper(p.Order))
}

// CalculateMeta calcula la metadata de paginación
func (p *Params) CalculateMeta(total int64) Meta {
	totalPages := int(math.Ceil(float64(total) / float64(p.Limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	return Meta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// NewResponse crea una respuesta paginada
func NewResponse(data interface{}, meta Meta) *Response {
	return &Response{
		Data: data,
		Meta: meta,
	}
}

// ApplyToQuery aplica paginación, ordenamiento y búsqueda a una query GORM
func (p *Params) ApplyToQuery(db *gorm.DB, searchFields []string) *gorm.DB {
	// Aplicar búsqueda si está presente
	if p.Search != "" && len(searchFields) > 0 {
		db = p.applySearch(db, searchFields)
	}

	// Aplicar filtros
	if len(p.Filter) > 0 {
		db = p.applyFilters(db)
	}

	// Aplicar ordenamiento
	if p.Sort != "" {
		db = db.Order(p.GetOrderBy())
	}

	// Aplicar paginación
	db = db.Limit(p.Limit).Offset(p.GetOffset())

	return db
}

// applySearch aplica búsqueda a múltiples campos
func (p *Params) applySearch(db *gorm.DB, fields []string) *gorm.DB {
	if len(fields) == 0 {
		return db
	}

	// Construir condición OR para búsqueda
	var conditions []string
	var values []interface{}

	searchTerm := "%" + p.Search + "%"
	for _, field := range fields {
		conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
		values = append(values, searchTerm)
	}

	whereClause := strings.Join(conditions, " OR ")
	return db.Where(whereClause, values...)
}

// applyFilters aplica filtros dinámicos
func (p *Params) applyFilters(db *gorm.DB) *gorm.DB {
	for field, value := range p.Filter {
		// Sanitizar nombre de campo (solo letras, números y _)
		if !isValidFieldName(field) {
			continue
		}

		// Soportar operadores especiales
		// filter[price__gte]=100 -> price >= 100
		// filter[status__in]=active,pending -> status IN ('active', 'pending')
		parts := strings.Split(field, "__")
		fieldName := parts[0]
		operator := "="

		if len(parts) > 1 {
			switch parts[1] {
			case "gte":
				operator = ">="
			case "gt":
				operator = ">"
			case "lte":
				operator = "<="
			case "lt":
				operator = "<"
			case "ne":
				operator = "!="
			case "like":
				operator = "ILIKE"
				value = "%" + value + "%"
			case "in":
				// Soportar valores múltiples separados por coma
				values := strings.Split(value, ",")
				db = db.Where(fmt.Sprintf("%s IN (?)", fieldName), values)
				continue
			case "notin":
				values := strings.Split(value, ",")
				db = db.Where(fmt.Sprintf("%s NOT IN (?)", fieldName), values)
				continue
			case "null":
				if value == "true" || value == "1" {
					db = db.Where(fmt.Sprintf("%s IS NULL", fieldName))
				} else {
					db = db.Where(fmt.Sprintf("%s IS NOT NULL", fieldName))
				}
				continue
			}
		}

		db = db.Where(fmt.Sprintf("%s %s ?", fieldName, operator), value)
	}

	return db
}

// isValidFieldName valida que el nombre del campo sea seguro
func isValidFieldName(field string) bool {
	// Remover operador si existe
	field = strings.Split(field, "__")[0]

	// Solo permitir letras, números y guión bajo
	for _, char := range field {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}
	return len(field) > 0
}

// Paginate es un helper para paginar resultados de forma genérica
func Paginate[T any](db *gorm.DB, params *Params, searchFields []string, preload ...string) ([]T, Meta, error) {
	var results []T
	var total int64

	// Clonar query para contar
	countDB := db.Session(&gorm.Session{})

	// Aplicar búsqueda y filtros para el conteo
	if params.Search != "" && len(searchFields) > 0 {
		countDB = params.applySearch(countDB, searchFields)
	}
	if len(params.Filter) > 0 {
		countDB = params.applyFilters(countDB)
	}

	// Contar total
	if err := countDB.Model(new(T)).Count(&total).Error; err != nil {
		return nil, Meta{}, err
	}

	// Aplicar paginación, ordenamiento y búsqueda
	query := params.ApplyToQuery(db, searchFields)

	// Aplicar preloads si se especificaron
	for _, relation := range preload {
		query = query.Preload(relation)
	}

	// Obtener resultados
	if err := query.Find(&results).Error; err != nil {
		return nil, Meta{}, err
	}

	meta := params.CalculateMeta(total)

	return results, meta, nil
}
