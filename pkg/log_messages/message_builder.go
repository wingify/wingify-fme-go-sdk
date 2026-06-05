/**
 * Copyright 2025-2026 Wingify Software Pvt. Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log_messages

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wingify/wingify-fme-go-sdk/pkg/brand"
)

var placeholderRegex = regexp.MustCompile(`\{([0-9a-zA-Z_]+)\}`)

var defaultHostProfile string

// SetDefaultHostProfile sets the hostProfile used for {brand} substitution when
// BuildMessage is called without an explicit hostProfile argument.
func SetDefaultHostProfile(hostProfile string) {
	defaultHostProfile = hostProfile
}

// BuildMessage constructs a message by replacing placeholders in a template with corresponding values from a data map.
// Placeholders are in the format {key}. The {brand} placeholder is filled from hostProfile when not present in data.
//
// Parameters:
//   - template: The message template containing placeholders in the format `{key}`
//   - data: A map containing keys and values used to replace the placeholders in the template
//   - hostProfile: Optional hostProfile init value; when empty, uses SetDefaultHostProfile
//
// Returns:
//   - The constructed message with all placeholders replaced by their corresponding values from the data map
func BuildMessage(template string, data map[string]interface{}, hostProfile ...string) string {
	profile := defaultHostProfile
	if len(hostProfile) > 0 && hostProfile[0] != "" {
		profile = hostProfile[0]
	}

	merged := make(map[string]interface{})
	if data != nil {
		for key, value := range data {
			merged[key] = value
		}
	}
	if _, ok := merged["brand"]; !ok {
		merged["brand"] = brand.DisplayNameFromHostProfile(profile)
	}

	// Replace all placeholders with their corresponding values
	result := placeholderRegex.ReplaceAllStringFunc(template, func(match string) string {
		// Extract the key from the placeholder (remove { and })
		key := strings.Trim(match, "{}")

		// Retrieve the value from the data map
		value, exists := merged[key]

		// If the key does not exist or the value is nil, return an empty string
		if !exists || value == nil {
			return ""
		}

		// Convert value to string based on its type
		switch v := value.(type) {
		case string:
			return v
		case int, int8, int16, int32, int64:
			return strings.TrimSpace(strings.Replace(match, key, fmt.Sprint(v), 1))
		case uint, uint8, uint16, uint32, uint64:
			return strings.TrimSpace(strings.Replace(match, key, fmt.Sprint(v), 1))
		case float32, float64:
			return strings.TrimSpace(strings.Replace(match, key, fmt.Sprint(v), 1))
		case bool:
			if v {
				return "true"
			}
			return "false"
		case func() string:
			// If the value is a function, evaluate it
			return v()
		default:
			// For other types (including custom string types like ApiEnum), use fmt.Sprint
			return fmt.Sprint(v)
		}
	})

	return result
}
