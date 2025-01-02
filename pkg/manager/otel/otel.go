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

package otel

import (
	"context"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Options defines common controller options.
type Options struct {
	// OTLPEndpoint defines whether to ship spans to an OTLP consumer or
	// not, and where to send them to.
	OTLPEndpoint string
}

func (o *Options) AddFlags(f *pflag.FlagSet) {
	f.StringVar(&o.OTLPEndpoint, "otlp-endpoint", "", "An optional OTLP endpoint to ship spans to.")
}

// Setup creates enough infrastructure to enable span creation, shipping and
// trace contect propagation.
func (o *Options) Setup(ctx context.Context, opts ...trace.TracerProviderOption) error {
	otel.SetLogger(log.Log)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	if o.OTLPEndpoint != "" {
		exporter, err := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(o.OTLPEndpoint),
			otlptracehttp.WithInsecure(),
		)

		if err != nil {
			return err
		}

		opts = append(opts, trace.WithBatcher(exporter))
	}

	otel.SetTracerProvider(trace.NewTracerProvider(opts...))

	return nil
}
