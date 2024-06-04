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

	"H4sIAAAAAAAC/+xWUU/bMBD+K6cbD5uUFvYyaXlBSHtB2zQEaJPGuukaXxoPxw72mQJV/vvkUNq0BASI",
	"x73F9vm++767z8oCC1c3zrKVgPkCQ1FxTd3neZyytywcvtCUzXcykdO+4lB43Yh2FnM8gEsyWsHnOOUu",
	"GEyKTruRM5DrRhdkzDXEwApK58FzcNEXDJZqDiAVCRRkYcq/rLaKr1iBtiAVgyKhKQUeY4YNibBPkL/P",
	"9kYfD0Y/aXQzebufr1ejP+PJYi/78L7tRbzb38EM5bphzDGI13aGbYbOz8jqG0o0TgrXsDpe1nXMpL6y",
	"UAJPfMmYbyXmZwvc8Vxijm9215rtLgXb9UO322yBjXcNe9HcidqHPVT35TytGPoxoBVb0aVm3ymyEm/K",
	"xtlZAHHj+/TaDD1fRO1ZYX62DTpZxbvpXy4E20mbpTrT4jW0eKK299VZ1vCQMMvjV9FkDTUsx13SPvvN",
	"iu44Qb0MWTohzTgZAwdHh+vSPJMKQFbB3Gvh0E30BveN5EPsV7l6RxlIpQO4btHZjOKsTs3olEkOg7mW",
	"CmrnGQpnha9kPOSHFJpwH+vr4IuwLWyXqK/gkXeXOmhntZ2dCEkMD3Z3FQdBSBhcCbTincpmG+sEEu25",
	"dXN7K+LqVn/JCjNUvHXM3jvf6/iav38V9z8y24XnzhCnuuZhAUTXvDnQcwrQ3WOV2JfO1ySYoyLhUQof",
	"6qRiwy8B6u49B0g/YNNo9UXsJT/8NDhxzeBcPEXlgYnankKd2r8h+SDggPezZ9q8c/a4P0I/ksP/vxsv",
	"ezfStralw9xGYzJ0DVtqNOZ4+w9QhduT9l8AAAD//1B07gW6CAAA",
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