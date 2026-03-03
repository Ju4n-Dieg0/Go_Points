package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func main() {
	// Leer el swagger.json generado por swag
	swaggerJSON, err := ioutil.ReadFile("docs/swagger.json")
	if err != nil {
		log.Fatalf("Error reading swagger.json: %v", err)
	}

	// Parse JSON
	var swaggerData map[string]interface{}
	err = json.Unmarshal(swaggerJSON, &swaggerData)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Limpiar y optimizar la especificación
	cleanedSpec := cleanOpenAPISpec(swaggerData)

	// Convertir a YAML
	yamlData, err := yaml.Marshal(cleanedSpec)
	if err != nil {
		log.Fatalf("Error converting to YAML: %v", err)
	}

	// Guardar en la raíz del proyecto
	err = ioutil.WriteFile("openapi.yaml", yamlData, 0644)
	if err != nil {
		log.Fatalf("Error writing openapi.yaml: %v", err)
	}

	fmt.Println("✓ OpenAPI specification exported to: openapi.yaml")
	fmt.Println("✓ Compatible with: Postman, Insomnia, Stoplight, SwaggerHub")
}

func cleanOpenAPISpec(spec map[string]interface{}) map[string]interface{} {
	// Crear una copia limpia
	cleaned := make(map[string]interface{})

	// Copiar campos esenciales
	if val, ok := spec["openapi"]; ok {
		cleaned["openapi"] = val
	}
	if val, ok := spec["info"]; ok {
		cleaned["info"] = val
	}
	if val, ok := spec["servers"]; ok {
		cleaned["servers"] = val
	} else {
		// Agregar servidores por defecto
		cleaned["servers"] = []map[string]interface{}{
			{
				"url":         "http://localhost:8080/api/v1",
				"description": "Development server",
			},
			{
				"url":         "https://staging.gopoints.com/api/v1",
				"description": "Staging server",
			},
			{
				"url":         "https://api.gopoints.com/api/v1",
				"description": "Production server",
			},
		}
	}
	if val, ok := spec["paths"]; ok {
		cleaned["paths"] = val
	}
	if val, ok := spec["components"]; ok {
		cleaned["components"] = val
	}
	if val, ok := spec["security"]; ok {
		cleaned["security"] = val
	}
	if val, ok := spec["tags"]; ok {
		cleaned["tags"] = val
	}
	if val, ok := spec["externalDocs"]; ok {
		cleaned["externalDocs"] = val
	}

	return cleaned
}
