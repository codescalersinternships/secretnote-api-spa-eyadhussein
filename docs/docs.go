// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Login a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User credentials to login",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseUserLoggedIn"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseUnauthorized"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseNotFound"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseInternalServerError"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "Logout a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout a user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseUserLoggedOut"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseUnauthorized"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseInternalServerError"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Register a new user with username, email, and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User credentials to register",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RegisterUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseUserRegistered"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseBadRequest"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseInternalServerError"
                        }
                    }
                }
            }
        },
        "/auth/verify-token": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "Verify token",
                "tags": [
                    "auth"
                ],
                "summary": "Verify token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseTokenVerified"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseUnauthorized"
                        }
                    }
                }
            }
        },
        "/notes": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "Create a note with title, content, max views, and expiration date",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Create a note",
                "parameters": [
                    {
                        "description": "Note details to create",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateNoteRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/util.APINote"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseUnauthorized"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseNotFound"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseInternalServerError"
                        }
                    }
                }
            }
        },
        "/notes/{id}": {
            "get": {
                "description": "Get a note by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Get a note by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Note ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseNoteRetrievedRestricted"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseNotFound"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseInternalServerError"
                        }
                    }
                }
            }
        },
        "/users/notes": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "Get notes by user ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Get notes by user ID",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/util.APINote"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseUnauthorized"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseNotFound"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/swagger.ResponseInternalServerError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.CreateNoteRequest": {
            "type": "object",
            "required": [
                "content",
                "expires_at",
                "max_views",
                "title"
            ],
            "properties": {
                "content": {
                    "type": "string"
                },
                "expires_at": {
                    "type": "string"
                },
                "max_views": {
                    "type": "integer"
                },
                "title": {
                    "type": "string",
                    "maxLength": 255
                }
            }
        },
        "models.LoginUserRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.RegisterUserRequest": {
            "description": "RegisterUserRequest represents a request to register a user",
            "type": "object",
            "required": [
                "email",
                "password",
                "password_confirmation",
                "username"
            ],
            "properties": {
                "email": {
                    "description": "@Description The email address of the user\n@Example johndoe@example.com",
                    "type": "string"
                },
                "password": {
                    "description": "@Description The password for the user\n@Example P@ssw0rd!",
                    "type": "string",
                    "minLength": 8
                },
                "password_confirmation": {
                    "description": "@Description Confirmation of the password\n@Example P@ssw0rd!",
                    "type": "string"
                },
                "username": {
                    "description": "@Description The username of the user\n@Example johndoe",
                    "type": "string"
                }
            }
        },
        "swagger.ResponseBadRequest": {
            "description": "Bad request",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "swagger.ResponseInternalServerError": {
            "description": "Internal server error",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "swagger.ResponseNotFound": {
            "description": "Not found",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "swagger.ResponseNoteRetrievedRestricted": {
            "description": "Note retrieved successfully",
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "swagger.ResponseTokenVerified": {
            "description": "Token valid",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "swagger.ResponseUnauthorized": {
            "description": "Unauthorized",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "swagger.ResponseUserLoggedIn": {
            "description": "User logged in successfully",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "swagger.ResponseUserLoggedOut": {
            "description": "User logged out successfully",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "swagger.ResponseUserRegistered": {
            "description": "User registered successfully",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "util.APINote": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "current_views": {
                    "type": "integer"
                },
                "expires_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "max_views": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "Token": {
            "type": "apiKey",
            "name": "token",
            "in": "cookie"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Secret Note API",
	Description:      "This is a sample server for managing notes.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
