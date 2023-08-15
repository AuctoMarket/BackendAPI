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
        "/buyers/login": {
            "post": {
                "description": "Checks to see if a buyer email exists and if supplied password matches the stored password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Logs a buyer into their account",
                "parameters": [
                    {
                        "description": "Buyers email",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Buyers password as plaintext",
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
        "/buyers/signup": {
            "post": {
                "description": "Checks to see if a buyer email exists and if not creates a new account with supplied email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Signs a new buyer up",
                "parameters": [
                    {
                        "description": "Buyers email [UNIQUE]",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Buyers password as plaintext",
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
        "/products": {
            "get": {
                "description": "Gets product information of products given query parameters provided in the Request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Gets Products with given query parameters",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Get products from a specific seller Id. Default is without any seller_id specified.",
                        "name": "seller_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Sort By a specific attribute of the product. Default is posted_date",
                        "name": "sort_by",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/data.GetProductResponseData"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/data.Message"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a new product post with the supplied data, if the data is not valid it throws and error",
                "produces": [
                    "application/json"
                ],
                "summary": "Creates a new product post",
                "parameters": [
                    {
                        "description": "The Seller who posted the product's seller_id",
                        "name": "seller_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Title of the product",
                        "name": "title",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Short description of the product",
                        "name": "description",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Price as an int of the product",
                        "name": "price",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "Condition of the product from a scale of 0 to 5",
                        "name": "condition",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "Type of product sale: Buy-Now or Pre-Order",
                        "name": "product_type",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Quantity of product to be put for sale",
                        "name": "product_quantity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/data.ProductCreateResponseData"
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
        "/products/{id}": {
            "get": {
                "description": "Checks to see if a product with a given id exists and returns its product information if it does.",
                "produces": [
                    "application/json"
                ],
                "summary": "Gets a Product by its Product ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "product_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/data.GetProductResponseData"
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
        },
        "/products/{id}/images": {
            "post": {
                "description": "Adds images to an existing product with supplied product id. If product with product id does not exist returns a",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Adds images to products",
                "parameters": [
                    {
                        "type": "string",
                        "description": "product_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Array of image files to add to the product post",
                        "name": "images",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
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
                    "415": {
                        "description": "Unsupported Media Type",
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
        "/sellers/login": {
            "post": {
                "description": "Checks to see if a sellers email exists and if supplied password matches the stored password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Logs a seller into their account",
                "parameters": [
                    {
                        "description": "Sellers email",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Sellers password as plaintext",
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
                            "$ref": "#/definitions/data.SellerLoginResponseData"
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
        "/sellers/signup": {
            "post": {
                "description": "Checks to see if a seller email does not already exists if so creates a new",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Signs a new seller up",
                "parameters": [
                    {
                        "description": "Sellers email [UNIQUE]",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Sellers password as plaintext",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Sellers seller alias that is displayed as their seller name [UNIQUE]",
                        "name": "seller_name",
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
                            "$ref": "#/definitions/data.SellerLoginResponseData"
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
        "/sellers/{id}": {
            "get": {
                "description": "Checks to see if a sellers id exists and if it does returns the specified sellers public information,",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Gets seller info based on seller id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "seller_id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/data.GetSellerByIdResponseData"
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
            "required": [
                "buyer_id",
                "email"
            ],
            "properties": {
                "buyer_id": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "data.GetProductResponseData": {
            "type": "object",
            "required": [
                "condition",
                "desc",
                "images",
                "posted_date",
                "price",
                "product_id",
                "product_quantity",
                "product_type",
                "seller_info",
                "sold_quantity",
                "title"
            ],
            "properties": {
                "condition": {
                    "type": "integer"
                },
                "desc": {
                    "type": "string"
                },
                "images": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/data.ProductImageData"
                    }
                },
                "posted_date": {
                    "type": "string",
                    "example": "2023-08-03 02:50:26.034552906 +0000 UTC m=+192.307467936"
                },
                "price": {
                    "type": "integer"
                },
                "product_id": {
                    "type": "string"
                },
                "product_quantity": {
                    "type": "integer"
                },
                "product_type": {
                    "type": "string"
                },
                "seller_info": {
                    "$ref": "#/definitions/data.GetSellerByIdResponseData"
                },
                "sold_quantity": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "data.GetSellerByIdResponseData": {
            "type": "object",
            "required": [
                "followers",
                "seller_id",
                "seller_name"
            ],
            "properties": {
                "followers": {
                    "type": "integer"
                },
                "seller_id": {
                    "type": "string"
                },
                "seller_name": {
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
        "data.ProductCreateResponseData": {
            "type": "object",
            "required": [
                "condition",
                "desc",
                "posted_date",
                "price",
                "product_id",
                "product_quantity",
                "product_type",
                "seller_id",
                "sold_quantity",
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
                    "type": "string",
                    "example": "2023-08-03 02:50:26.034552906 +0000 UTC m=+192.307467936"
                },
                "price": {
                    "type": "integer"
                },
                "product_id": {
                    "type": "string"
                },
                "product_quantity": {
                    "type": "integer"
                },
                "product_type": {
                    "type": "string"
                },
                "seller_id": {
                    "type": "string"
                },
                "sold_quantity": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "data.ProductImageData": {
            "type": "object",
            "required": [
                "image_no",
                "image_path"
            ],
            "properties": {
                "image_no": {
                    "type": "integer"
                },
                "image_path": {
                    "type": "string"
                }
            }
        },
        "data.SellerLoginResponseData": {
            "type": "object",
            "required": [
                "email",
                "followers",
                "seller_id",
                "seller_name"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "followers": {
                    "type": "integer"
                },
                "seller_id": {
                    "type": "string"
                },
                "seller_name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "*",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "AUCTO Backend API",
	Description:      "This is the REST API for Aucto's marketplace, it is currently in v1.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
