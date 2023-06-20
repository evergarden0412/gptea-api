// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/auth/cred/logout": {
            "delete": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "delete refresh token",
                "tags": [
                    "token"
                ],
                "summary": "Logout",
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/auth/cred/register": {
            "post": {
                "description": "Register a credential",
                "tags": [
                    "auth"
                ],
                "summary": "Register a credential",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.credBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/server.messageResponse"
                        }
                    }
                }
            }
        },
        "/auth/cred/sign-in": {
            "post": {
                "description": "Sign in with a credential",
                "tags": [
                    "auth"
                ],
                "summary": "Sign in with a credential",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.credBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.signInHandlerOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/auth/token/refresh": {
            "post": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    },
                    {
                        "RefreshTokenAuth": []
                    }
                ],
                "description": "Refresh a token",
                "tags": [
                    "token"
                ],
                "summary": "Refresh a token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.signInHandlerOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/auth/token/verify": {
            "get": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "Verify a accesstoken",
                "tags": [
                    "token"
                ],
                "summary": "Verify a token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.tokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me": {
            "delete": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "Delete my account",
                "tags": [
                    "users"
                ],
                "summary": "Delete my account",
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/chats": {
            "get": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "Get my chats in descending order of created_at",
                "tags": [
                    "chats"
                ],
                "summary": "Get my chats",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.chatsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "Post my chat",
                "tags": [
                    "chats"
                ],
                "summary": "Post my chat",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.chatBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/server.messageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/chats/{chatID}": {
            "delete": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "Delete my chat",
                "tags": [
                    "chats"
                ],
                "summary": "Delete my chat",
                "parameters": [
                    {
                        "type": "string",
                        "description": "chatID",
                        "name": "chatID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "Patch my chat name",
                "tags": [
                    "chats"
                ],
                "summary": "Patch my chat",
                "parameters": [
                    {
                        "type": "string",
                        "description": "chatID",
                        "name": "chatID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.chatBody"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/chats/{chatID}/messages": {
            "get": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "Get my messages in descending order of created_at",
                "tags": [
                    "messages"
                ],
                "summary": "Get my messages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "chatID",
                        "name": "chatID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.messagesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "Post my message and get response when chatbot finishes processing",
                "tags": [
                    "messages"
                ],
                "summary": "Post my message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "chatID",
                        "name": "chatID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.messageBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/server.messageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/scrapbooks": {
            "get": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "Get my scrapbooks",
                "tags": [
                    "scraps"
                ],
                "summary": "Get my scrapbooks",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.scrapbooksResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "post new scrapbook",
                "tags": [
                    "scraps"
                ],
                "summary": "post new scrapbook",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.scrapbookBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/scrapbooks/{scrapbookID}": {
            "delete": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "delete scrapbook",
                "tags": [
                    "scraps"
                ],
                "summary": "delete scrapbook",
                "parameters": [
                    {
                        "type": "string",
                        "description": "scrapbookID",
                        "name": "scrapbookID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "patch scrapbook",
                "tags": [
                    "scraps"
                ],
                "summary": "patch scrapbook",
                "parameters": [
                    {
                        "type": "string",
                        "description": "scrapbookID",
                        "name": "scrapbookID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.scrapbookBody"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/scrapbooks/{scrapbookID}/scraps": {
            "get": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "get scraps on scrapbook",
                "tags": [
                    "scraps"
                ],
                "summary": "get scraps on scrapbook",
                "parameters": [
                    {
                        "type": "string",
                        "description": "scrapbookID",
                        "name": "scrapbookID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.scrapsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/scraps": {
            "post": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "post new scrap, store in default scrapbook",
                "tags": [
                    "scraps"
                ],
                "summary": "post new scrap",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.postScrapBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/scraps/{scrapID}": {
            "delete": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "delete scrap",
                "tags": [
                    "scraps"
                ],
                "summary": "delete scrap",
                "parameters": [
                    {
                        "type": "string",
                        "description": "scrapID",
                        "name": "scrapID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/scraps/{scrapID}/scrapbooks": {
            "get": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "get my scrapbooks on scrap",
                "tags": [
                    "scraps"
                ],
                "summary": "get my scrapbooks on scrap",
                "parameters": [
                    {
                        "type": "string",
                        "description": "scrapID",
                        "name": "scrapID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.scrapbooksResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        },
        "/me/scraps/{scrapID}/scrapbooks/{scrapbookID}": {
            "post": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "post scrap on scrapbook",
                "tags": [
                    "scraps"
                ],
                "summary": "post scrap on scrapbook",
                "parameters": [
                    {
                        "type": "string",
                        "description": "scrapID",
                        "name": "scrapID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "scrapbookID",
                        "name": "scrapbookID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "AccessTokenAuth": []
                    }
                ],
                "description": "delete scrap on scrapbook",
                "tags": [
                    "scraps"
                ],
                "summary": "delete scrap on scrapbook",
                "parameters": [
                    {
                        "type": "string",
                        "description": "scrapID",
                        "name": "scrapID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "scrapbookID",
                        "name": "scrapbookID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "internal.Chat": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string",
                    "example": "2021-01-01T00:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "Hjejwerhj"
                },
                "name": {
                    "type": "string",
                    "example": "basic"
                }
            }
        },
        "internal.Message": {
            "type": "object",
            "properties": {
                "chatID": {
                    "type": "string",
                    "example": "Hjejwerhj"
                },
                "content": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "seq": {
                    "description": "seq starts from 1",
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "internal.Scrap": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string",
                    "example": "2021-01-01T00:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "Hjejwerhj"
                },
                "memo": {
                    "type": "string",
                    "example": "hello"
                },
                "message": {
                    "$ref": "#/definitions/internal.Message"
                }
            }
        },
        "internal.Scrapbook": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string",
                    "example": "2021-01-01T00:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "Hjejwerhj"
                },
                "isDefault": {
                    "type": "boolean",
                    "example": true
                },
                "name": {
                    "type": "string",
                    "example": "basic"
                }
            }
        },
        "server.chatBody": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "server.chatsResponse": {
            "type": "object",
            "properties": {
                "chats": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal.Chat"
                    }
                }
            }
        },
        "server.credBody": {
            "type": "object",
            "required": [
                "accessToken",
                "cred"
            ],
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "cred": {
                    "type": "string",
                    "example": "naver"
                }
            }
        },
        "server.errorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "error message"
                }
            }
        },
        "server.messageBody": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                }
            }
        },
        "server.messageResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Hello, World!"
                }
            }
        },
        "server.messagesResponse": {
            "type": "object",
            "properties": {
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal.Message"
                    }
                }
            }
        },
        "server.postScrapBody": {
            "type": "object",
            "required": [
                "chatID",
                "seq"
            ],
            "properties": {
                "chatID": {
                    "type": "string"
                },
                "memo": {
                    "type": "string"
                },
                "seq": {
                    "type": "integer"
                }
            }
        },
        "server.scrapbookBody": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "server.scrapbooksResponse": {
            "type": "object",
            "properties": {
                "scrapbooks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal.Scrapbook"
                    }
                }
            }
        },
        "server.scrapsResponse": {
            "type": "object",
            "properties": {
                "scraps": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal.Scrap"
                    }
                }
            }
        },
        "server.signInHandlerOutput": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "server.tokenResponse": {
            "type": "object",
            "properties": {
                "exp": {
                    "type": "string"
                },
                "iat": {
                    "type": "string"
                },
                "jti": {
                    "type": "string"
                },
                "sub": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "AccessTokenAuth": {
            "description": "type ` + "`" + `Bearer {access_token}` + "`" + `",
            "type": "apiKey",
            "name": "authorization",
            "in": "header"
        },
        "RefreshTokenAuth": {
            "description": "type ` + "`" + `{refresh_token}` + "`" + `",
            "type": "apiKey",
            "name": "x-refresh-token",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "api.gptea-test.keenranger.dev",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "GPTea API",
	Description:      "This is a sample server for GPTea API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
