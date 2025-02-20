// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
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
