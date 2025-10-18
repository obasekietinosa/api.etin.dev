package main

import _ "embed"

//go:embed openapi.json
var embeddedSwagger []byte

func init() {
	if len(embeddedSwagger) == 0 {
		panic("embedded OpenAPI document is empty")
	}
}
