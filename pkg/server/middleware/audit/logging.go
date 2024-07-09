/*
Copyright 2024 the Unikorn Authors.

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

package audit

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/routers"

	"github.com/unikorn-cloud/core/pkg/authorization/userinfo"
	"github.com/unikorn-cloud/core/pkg/openapi"
	"github.com/unikorn-cloud/core/pkg/server/errors"
	"github.com/unikorn-cloud/core/pkg/server/middleware"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Logger struct {
	// next defines the next HTTP handler in the chain.
	next http.Handler

	// openapi caches the Schema schema.
	openapi *openapi.Schema

	// application is the application name.
	application string

	// version is the application version.
	version string
}

// Ensure this implements the required interfaces.
var _ http.Handler = &Logger{}

// New returns an initialized middleware.
func New(next http.Handler, openapi *openapi.Schema, application, version string) *Logger {
	return &Logger{
		next:        next,
		openapi:     openapi,
		application: application,
		version:     version,
	}
}

// getResource will resolve to a resource type.
func getResource(route *routers.Route, params map[string]string) *Resource {
	// We are looking for "/.../resource/{idParameter}"
	// or failing that "/.../resource"
	matches := regexp.MustCompile(`/([^/]+)/{([^/}]+)}$`).FindStringSubmatch(route.Path)
	if matches == nil {
		segments := strings.Split(route.Path, "/")

		return &Resource{
			Type: segments[len(segments)-1],
		}
	}

	resource := &Resource{
		Type: matches[1],
		ID:   params[matches[2]],
	}

	return resource
}

// ServeHTTP implements the http.Handler interface.
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, params, err := l.openapi.FindRoute(r)
	if err != nil {
		errors.HandleError(w, r, errors.OAuth2ServerError("route lookup failure").WithError(err))

		return
	}

	writer := middleware.NewLoggingResponseWriter(w)

	l.next.ServeHTTP(writer, r)

	// Users and auditors care about things coming, going and changing, who did
	// those things and when?  Certainly not periodic polling that is par for the
	// course. Failures of reads may be indicative of someone trying to do
	// something they shouldn't via the API (or indeed a bug in a UI leeting them
	// attempt something they are forbidden to do).
	if r.Method == http.MethodGet {
		return
	}

	// If there is not accountibility e.g. a global call, it's not worth logging.
	userinfo := userinfo.FromContext(r.Context())
	if userinfo == nil {
		return
	}

	// If there's no scope, then discard also.
	if len(params) == 0 {
		return
	}

	logParams := []any{
		"component", &Component{
			Name:    l.application,
			Version: l.version,
		},
		"actor", &Actor{
			Subject: userinfo.Subject,
		},
		"operation", &Operation{
			Verb: r.Method,
		},
		"scope", params,
		"resource", getResource(route, params),
		"result", &Result{
			Status: writer.StatusCode(),
		},
	}

	log.FromContext(r.Context()).Info("audit", logParams...)
}

func Middleware(openapi *openapi.Schema, application, version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return New(next, openapi, application, version)
	}
}
