// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package swagger

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
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "icaroribeiro@hotmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/files/upload": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "API endpoint used to store a file on disk by uploading it",
                "operationId": "Upload",
                "parameters": [
                    {
                        "type": "file",
                        "description": "The file to be uploaded",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.File"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/error.Error"
                        }
                    }
                }
            }
        },
        "/files/{fileID}/chunk": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "API endpoint used to read a chunk of bytes of a file identified by its id using the offset and limit parameters",
                "operationId": "Read",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The number of bytes to skip before starting returning the bytes. 0 is to skip no entry, 1 is to skip the first byte and so on. If not informed, the default value will be set to 0",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "The maximum number of bytes to be returned. If not informed, the default value will be set to 10",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Chunk"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/error.Error"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health check"
                ],
                "summary": "API endpoint used to verify if the service has started up correctly and is ready to accept requests",
                "operationId": "GetStatus",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/message.Message"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "error.Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "message.Message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "model.Chunk": {
            "type": "object",
            "properties": {
                "bytes_read": {
                    "type": "integer"
                },
                "bytestream_to_string": {
                    "type": "string"
                },
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                }
            }
        },
        "model.File": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "modTime": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "size": {
                    "type": "integer"
                }
            }
        }
    },
    "tags": [
        {
            "description": "It refers to the operation related to health check",
            "name": "health check"
        },
        {
            "description": "It refers to the operations related to file",
            "name": "file"
        }
    ]
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
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/",
	Schemes:     []string{"http"},
	Title:       "Avenue Securities API",
	Description: "A REST API developed using Golang for uploading and reading files in chunks of bytes defined by offset and limit parameters",
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
