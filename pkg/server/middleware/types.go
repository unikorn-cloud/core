/*
Copyright 2024-2025 the Unikorn Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package middleware

import (
	"bytes"
	"net/http"
)

// LoggingResponseWriter is the ubiquitous reimplementation of a response
// writer that allows access to the HTTP status code in middleware.
type LoggingResponseWriter struct {
	next http.ResponseWriter
	code int
	body *bytes.Buffer
}

func NewLoggingResponseWriter(next http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{
		next: next,
	}
}

// Check the correct interface is implmented.
var _ http.ResponseWriter = &LoggingResponseWriter{}

func (w *LoggingResponseWriter) Header() http.Header {
	return w.next.Header()
}

func (w *LoggingResponseWriter) Write(body []byte) (int, error) {
	if w.body == nil {
		w.body = &bytes.Buffer{}
	}

	w.body.Write(body)

	return w.next.Write(body)
}

func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.next.WriteHeader(statusCode)
}

func (w *LoggingResponseWriter) StatusCode() int {
	if w.code == 0 {
		return http.StatusOK
	}

	return w.code
}

func (w *LoggingResponseWriter) Body() *bytes.Buffer {
	return w.body
}
