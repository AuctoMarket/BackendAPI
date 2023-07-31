// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/buyer/login/": {
            "post": {
                "description": "Checks to see if a user email exists and if supplied password matches the stored password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Logs a buyer into their account",
                "parameters": [
                    {
                        "description": "email",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "password",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/data.BuyerLoginResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    }
                }
            }
        },
        "/buyer/signup/": {
            "post": {
                "description": "Checks to see if a user email exists and if not creates a new account with supplied email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Signs a new buyer up",
                "parameters": [
                    {
                        "description": "email",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "password",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/data.BuyerLoginResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    }
                }
            }
        },
        "/product/create": {
            "post": {
                "description": "Creates a new product post with the supplied data, if the data is not valid it throws and error",
                "produces": [
                    "application/json"
                ],
                "summary": "Creates a new product post",
                "parameters": [
                    {
                        "description": "seller_id",
                        "name": "seller_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "title",
                        "name": "title",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "desc",
                        "name": "description",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "price",
                        "name": "price",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "condition",
                        "name": "condition",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "product_type",
                        "name": "product_type",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/data.ProductResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    }
                }
            }
        },
        "/product/{id}": {
            "get": {
                "description": "Checks to see if a product with a given id exists and returns its product information if it does.",
                "produces": [
                    "application/json"
                ],
                "summary": "Gets a Product by its Product ID",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/data.ProductResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "data.BuyerLoginResponseData": {
            "type": "object",
            "properties": {
                "buyer_id": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "data.Message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "data.ProductResponseData": {
            "type": "object",
            "required": [
                "condition",
                "desc",
                "posted_date",
                "price",
                "product_id",
                "product_type",
                "seller_id",
                "title"
            ],
            "properties": {
                "condition": {
                    "type": "integer"
                },
                "desc": {
                    "type": "string"
                },
                "posted_date": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "product_id": {
                    "type": "string"
                },
                "product_type": {
                    "type": "string"
                },
                "seller_id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "AUCTO Backend API",
	Description:      "This is the backend REST API for Aucto's marketplace, it is currently in v1.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
