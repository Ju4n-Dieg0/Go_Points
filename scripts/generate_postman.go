package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// PostmanCollection representa la estructura de una colección de Postman v2.1
type PostmanCollection struct {
	Info     CollectionInfo `json:"info"`
	Item     []Item         `json:"item"`
	Variable []Variable     `json:"variable"`
}

type CollectionInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
	Version     string `json:"version"`
}

type Item struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Item        []Request `json:"item,omitempty"`
	Request     *Request  `json:"request,omitempty"`
}

type Request struct {
	Name        string      `json:"name"`
	Method      string      `json:"method,omitempty"`
	Header      []Header    `json:"header,omitempty"`
	Body        *Body       `json:"body,omitempty"`
	URL         URL         `json:"url,omitempty"`
	Description string      `json:"description,omitempty"`
	Request     *RequestDef `json:"request,omitempty"`
}

type RequestDef struct {
	Method      string   `json:"method"`
	Header      []Header `json:"header"`
	Body        *Body    `json:"body,omitempty"`
	URL         URL      `json:"url"`
	Description string   `json:"description,omitempty"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Body struct {
	Mode string      `json:"mode"`
	Raw  string      `json:"raw,omitempty"`
	Options *BodyOptions `json:"options,omitempty"`
}

type BodyOptions struct {
	Raw BodyRaw `json:"raw"`
}

type BodyRaw struct {
	Language string `json:"language"`
}

type URL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func main() {
	collection := PostmanCollection{
		Info: CollectionInfo{
			Name:        "Go Points API",
			Description: "API REST profesional para sistema de puntos de fidelización",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
			Version:     "1.0.0",
		},
		Variable: []Variable{
			{Key: "base_url", Value: "http://localhost:8080/api/v1", Type: "string"},
			{Key: "access_token", Value: "", Type: "string"},
		},
		Item: []Item{
			createAuthFolder(),
			createCompaniesFolder(),
			createSubscriptionsFolder(),
			createConsumersFolder(),
			createProductsFolder(),
			createPointsFolder(),
			createRewardsFolder(),
			createHealthFolder(),
		},
	}

	// Convertir a JSON
	jsonData, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Guardar archivo
	err = ioutil.WriteFile("Go_Points_API.postman_collection.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Println("✓ Postman collection generated: Go_Points_API.postman_collection.json")
}

func createAuthFolder() Item {
	return Item{
		Name:        "Auth",
		Description: "Endpoints de autenticación y gestión de tokens",
		Item: []Request{
			{
				Name:        "Register",
				Request: &RequestDef{
					Method: "POST",
					Header: []Header{
						{Key: "Content-Type", Value: "application/json", Type: "text"},
					},
					Body: &Body{
						Mode: "raw",
						Raw: `{
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "password": "SecurePass123!",
  "role": "consumer"
}`,
						Options: &BodyOptions{
							Raw: BodyRaw{Language: "json"},
						},
					},
					URL: URL{
						Raw:  "{{base_url}}/auth/register",
						Host: []string{"{{base_url}}"},
						Path: []string{"auth", "register"},
					},
					Description: "Registrar un nuevo usuario en el sistema",
				},
			},
			{
				Name:        "Login",
				Request: &RequestDef{
					Method: "POST",
					Header: []Header{
						{Key: "Content-Type", Value: "application/json", Type: "text"},
					},
					Body: &Body{
						Mode: "raw",
						Raw: `{
  "email": "juan@example.com",
  "password": "SecurePass123!"
}`,
						Options: &BodyOptions{
							Raw: BodyRaw{Language: "json"},
						},
					},
					URL: URL{
						Raw:  "{{base_url}}/auth/login",
						Host: []string{"{{base_url}}"},
						Path: []string{"auth", "login"},
					},
					Description: "Iniciar sesión y obtener tokens JWT",
				},
			},
			{
				Name:        "Get Profile",
				Request: &RequestDef{
					Method: "GET",
					Header: []Header{
						{Key: "Authorization", Value: "Bearer {{access_token}}", Type: "text"},
					},
					URL: URL{
						Raw:  "{{base_url}}/auth/profile",
						Host: []string{"{{base_url}}"},
						Path: []string{"auth", "profile"},
					},
					Description: "Obtener perfil del usuario autenticado",
				},
			},
			{
				Name:        "Refresh Token",
				Request: &RequestDef{
					Method: "POST",
					Header: []Header{
						{Key: "Content-Type", Value: "application/json", Type: "text"},
					},
					Body: &Body{
						Mode: "raw",
						Raw: `{
  "refresh_token": "your-refresh-token-here"
}`,
						Options: &BodyOptions{
							Raw: BodyRaw{Language: "json"},
						},
					},
					URL: URL{
						Raw:  "{{base_url}}/auth/refresh",
						Host: []string{"{{base_url}}"},
						Path: []string{"auth", "refresh"},
					},
					Description: "Renovar access token usando refresh token",
				},
			},
			{
				Name:        "Logout",
				Request: &RequestDef{
					Method: "POST",
					Header: []Header{
						{Key: "Authorization", Value: "Bearer {{access_token}}", Type: "text"},
					},
					URL: URL{
						Raw:  "{{base_url}}/auth/logout",
						Host: []string{"{{base_url}}"},
						Path: []string{"auth", "logout"},
					},
					Description: "Cerrar sesión (invalidar token)",
				},
			},
		},
	}
}

func createCompaniesFolder() Item {
	return Item{
		Name:        "Companies",
		Description: "Gestión de empresas del sistema",
		Item: []Request{
			{
				Name:        "Create Company",
				Request: &RequestDef{
					Method: "POST",
					Header: []Header{
						{Key: "Content-Type", Value: "application/json", Type: "text"},
						{Key: "Authorization", Value: "Bearer {{access_token}}", Type: "text"},
					},
					Body: &Body{
						Mode: "raw",
						Raw: `{
  "name": "Mi Empresa",
  "description": "Descripción de la empresa",
  "industry": "Retail",
  "email": "contacto@miempresa.com",
  "phone": "+573001234567"
}`,
						Options: &BodyOptions{
							Raw: BodyRaw{Language: "json"},
						},
					},
					URL: URL{
						Raw:  "{{base_url}}/companies",
						Host: []string{"{{base_url}}"},
						Path: []string{"companies"},
					},
				},
			},
			{
				Name:        "List Companies",
				Request: &RequestDef{
					Method: "GET",
					Header: []Header{
						{Key: "Authorization", Value: "Bearer {{access_token}}", Type: "text"},
					},
					URL: URL{
						Raw:  "{{base_url}}/companies?page=1&pageSize=10",
						Host: []string{"{{base_url}}"},
						Path: []string{"companies"},
					},
				},
			},
		},
	}
}

func createSubscriptionsFolder() Item {
	return Item{
		Name:        "Subscriptions",
		Description: "Gestión de suscripciones de empresas",
		Item: []Request{},
	}
}

func createConsumersFolder() Item {
	return Item{
		Name:        "Consumers",
		Description: "Gestión de consumidores/clientes",
		Item: []Request{},
	}
}

func createProductsFolder() Item {
	return Item{
		Name:        "Products",
		Description: "Gestión de productos de empresas",
		Item: []Request{},
	}
}

func createPointsFolder() Item {
	return Item{
		Name:        "Points",
		Description: "Sistema de puntos (ganar, redimir, consultar saldo)",
		Item: []Request{},
	}
}

func createRewardsFolder() Item {
	return Item{
		Name:        "Rewards",
		Description: "Gestión de recompensas",
		Item: []Request{},
	}
}

func createHealthFolder() Item {
	return Item{
		Name:        "Health Checks",
		Description: "Endpoints de salud del sistema",
		Item: []Request{
			{
				Name:        "Health Check",
				Request: &RequestDef{
					Method: "GET",
					Header: []Header{},
					URL: URL{
						Raw:  "http://localhost:8080/health",
						Host: []string{"localhost:8080"},
						Path: []string{"health"},
					},
					Description: "Verificar que la aplicación esté ejecutándose",
				},
			},
			{
				Name:        "Readiness Check",
				Request: &RequestDef{
					Method: "GET",
					Header: []Header{},
					URL: URL{
						Raw:  "http://localhost:8080/ready",
						Host: []string{"localhost:8080"},
						Path: []string{"ready"},
					},
					Description: "Verificar que la aplicación esté lista (DB conectada)",
				},
			},
		},
	}
}
