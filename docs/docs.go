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
        "/api/v1/user-action": {
            "post": {
                "description": "Сохраняет действие пользователя в базе данных",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "UserAction"
                ],
                "summary": "Сохранение действия пользователя",
                "parameters": [
                    {
                        "description": "Данные действия пользователя",
                        "name": "action",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UserAction"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Действие сохранено",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверные данные",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка при сохранении данных",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user-actions": {
            "get": {
                "description": "Предоставление информации о действиях всех пользователей из базы данных",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "UserAction"
                ],
                "summary": "Предоставление информации о действиях всех пользователей",
                "parameters": [
                    {
                        "description": "Данные действия пользователя",
                        "name": "action",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UserAction"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка при получении действий",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/users-actions/{user_id}": {
            "get": {
                "description": "Возвращает данные о конкретном пользователе по его уникальному идентификатору (user_id)",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "UserAction"
                ],
                "summary": "Получить данные о пользователе",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID пользователя",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Данные о пользователе",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос, не указан user_id",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Пользователь с таким ID не найден",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.UserAction": {
            "type": "object",
            "properties": {
                "device_type": {
                    "type": "string"
                },
                "event_time": {
                    "type": "string"
                },
                "event_type": {
                    "type": "string"
                },
                "user_agent": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}