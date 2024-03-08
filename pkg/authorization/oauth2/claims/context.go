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

package claims

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrContextError is raised when a required value cannot be retrieved
	// from a context.
	ErrContextError = errors.New("value missing from context")
)

// contextKey defines a new context key type unique to this package.
type contextKey int

const (
	// claimsKey is used to store claims in a context.
	claimsKey contextKey = iota
)

// NewContext injects the given claims into a new context.
func NewContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// FromContext extracts the claims from a context.
func FromContext(ctx context.Context) (*Claims, error) {
	value := ctx.Value(claimsKey)
	if value == nil {
		return nil, fmt.Errorf("%w: unable to find claims", ErrContextError)
	}

	claims, ok := value.(*Claims)
	if !ok {
		return nil, fmt.Errorf("%w: unable to assert claims", ErrContextError)
	}

	return claims, nil
}
