// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"Yux2/w0AAP//h3igmNwaAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
