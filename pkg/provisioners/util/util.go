/*
Copyright 2022-2024 EscherCloud.
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

package util

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrNamespaceLookup = errors.New("unable to lookup namespace")
)

// GetConfigurationHash is used to restart badly behaved apps that don't respect configuration
// changes.
func GetConfigurationHash(config any) (string, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	configHash := sha256.Sum256(configJSON)

	return fmt.Sprintf("%x", configHash), nil
}
