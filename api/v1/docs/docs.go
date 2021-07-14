// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "JonathanGzzBen",
            "url": "http://www.github.com/JonathanGzzBen",
            "email": "jonathangzzben@gmail.com"
        },
        "license": {
            "name": "MIT License",
            "url": "https://mit-license.org/"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth": {
            "get": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "tags": [
                    "auth"
                ],
                "operationId": "GetCurrentUser",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            }
        },
        "/categories": {
            "get": {
                "description": "Get all registered categories.",
                "tags": [
                    "categories"
                ],
                "summary": "Get all categories",
                "operationId": "GetAllCategories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Category"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            },
            "post": {
                "description": "Register a new category.",
                "tags": [
                    "categories"
                ],
                "summary": "Create category",
                "operationId": "CreateCategory",
                "parameters": [
                    {
                        "description": "Category",
                        "name": "category",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.CreateCategoryDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Category"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            }
        },
        "/categories/{id}": {
            "get": {
                "description": "Get category with matching ID.",
                "tags": [
                    "categories"
                ],
                "summary": "Get category",
                "operationId": "GetCategory",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Category ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Category"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates a registered category.",
                "tags": [
                    "categories"
                ],
                "summary": "Update category",
                "operationId": "UpdateCategory",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Category ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Category",
                        "name": "category",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.UpdateCategoryDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Category"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Get all registered users.",
                "tags": [
                    "users"
                ],
                "summary": "Get all users",
                "operationId": "GetAllUsers",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Get user with matching ID.",
                "tags": [
                    "users"
                ],
                "summary": "Get user",
                "operationId": "GetUser",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "description": "Update matching user with provided data.",
                "tags": [
                    "users"
                ],
                "summary": "Update user",
                "operationId": "UpdateUser",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.UpdateUserDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.CreateCategoryDTO": {
            "type": "object",
            "properties": {
                "imageUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "controllers.UpdateCategoryDTO": {
            "type": "object",
            "properties": {
                "imageUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "controllers.UpdateUserDTO": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "birthdate": {
                    "type": "string",
                    "example": "2006-01-02T15:04:05Z"
                },
                "description": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "profilePictureUrl": {
                    "type": "string"
                },
                "role": {
                    "type": "string",
                    "example": "User"
                },
                "shortDescription": {
                    "type": "string"
                }
            }
        },
        "models.APIError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "models.Category": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "imageUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "birthdate": {
                    "type": "string",
                    "example": "2006-01-02T15:04:05Z"
                },
                "description": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "profilePictureUrl": {
                    "type": "string"
                },
                "role": {
                    "type": "string",
                    "example": "User"
                },
                "shortDescription": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "AccessToken": {
            "type": "apiKey",
            "name": "AccessToken",
            "in": "header"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "/v1/auth/google-login",
            "tokenUrl": "/v1/auth/google-callback",
            "scopes": {
                "email": " Grant access to email",
                "openid": " Allow identifying account",
                "profile": " Grant access to profile"
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1.0",
	Host:        "localhost:8080",
	BasePath:    "/v1",
	Schemes:     []string{},
	Title:       "Ingenialists API V1",
	Description: "This is Ingenialist's API",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
