package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "REST gateway for service health and repository info.",
        "title": "Repo Stat API",
        "contact": {},
        "version": "1.0"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/ping": {
            "get": {
                "summary": "Service health check",
                "description": "Checks processor and subscriber availability.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "503": {
                        "description": "Service Unavailable"
                    }
                }
            }
        },
        "/api/repositories/info": {
            "get": {
                "summary": "Get GitHub repository info",
                "description": "Returns basic information for a GitHub repository by URL.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "repositories"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "GitHub repository URL",
                        "name": "url",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger metadata and can be overridden at runtime.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Repo Stat API",
	Description:      "REST gateway for service health and repository info.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
