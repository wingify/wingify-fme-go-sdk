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

package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
	"github.com/wingify/wingify-fme-go-sdk/pkg/models/user"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/interfaces"
)

/*
isValidCampaignOrVariationID checks if a given value is a valid campaign ID or variation ID.
A valid ID must be an integer >= 0 or a digit-only string.

@params:
  val - interface{} : The value to be validated (can be of any type, usually int, float64, or string)

@returns:
  bool - true if the value is a non-negative integer or a digit-only string, false otherwise
*/
func isValidCampaignOrVariationID(val interface{}) bool {
	// Reject null/nil values immediately
	if val == nil {
		return false
	}

	// Use type switch to accurately check the underlying type of the interface{}
	switch v := val.(type) {
	case int:
		// Accept native integers that are non-negative
		return v >= 0
	case int8:
		return v >= 0
	case int16:
		return v >= 0
	case int32:
		return v >= 0
	case int64:
		return v >= 0
	case float64:
		// JSON numbers decode as float64 by default in Go.
		// We ensure the float is non-negative and has no fractional part (e.g. 10.0 is valid, 10.5 is invalid).
		if v >= 0 && v == float64(int64(v)) {
			return true
		}
		return false
	case float32:
		// Same logic as float64, ensuring it's a non-negative whole number
		if v >= 0 && v == float32(int32(v)) {
			return true
		}
		return false
	case string:
		// Reject empty strings
		if v == "" {
			return false
		}
		// Iterate through the string to ensure every character is a valid numeric digit (0-9)
		for _, c := range v {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	default:
		// Reject any other types (booleans, arrays, objects, etc.)
		return false
	}
}

// NormalizeWebTestingCampaignsMap normalizes Web Testing campaign map keys and variation values to strings.
func NormalizeWebTestingCampaignsMap(rawAssignments map[string]interface{}) (map[string]string, error) {
	campaignIDToVariationID := make(map[string]string)
	for campaignID, assignedVariationID := range rawAssignments {
		if !isValidCampaignOrVariationID(campaignID) || !isValidCampaignOrVariationID(assignedVariationID) {
			return nil, fmt.Errorf("invalid campaign ID or variation ID")
		}
		campaignIDToVariationID[fmt.Sprintf("%v", campaignID)] = fmt.Sprintf("%v", assignedVariationID)
	}
	return campaignIDToVariationID, nil
}

// ParseWebTestingCampaignsFromContext parses `context.platformVariables.webTestingCampaigns` (JSON string or plain object).
func ParseWebTestingCampaignsFromContext(
	context *user.WingifyUserContext,
	serviceContainer interfaces.ServiceContainerInterface,
) map[string]string {
	platformVariables := context.GetPlatformVariables()
	if platformVariables == nil {
		return nil
	}

	webTestingCampaignsInput, exists := platformVariables["webTestingCampaigns"]
	if !exists || webTestingCampaignsInput == nil {
		return nil
	}

	// If it was already parsed and normalized successfully by a previous call, return it directly.
	if normalizedMap, ok := webTestingCampaignsInput.(map[string]string); ok {
		return normalizedMap
	}

	// SDK already forwarded a plain campaignId -> variationId map.
	if reflect.TypeOf(webTestingCampaignsInput).Kind() == reflect.Map {
		// Because the map could be map[int]interface{}, map[string]int, etc.,
		// we use reflection to build a normalized map[string]interface{}
		// and let the existing NormalizeWebTestingCampaignsMap logic validate it.
		inputMap := make(map[string]interface{})
		mapValue := reflect.ValueOf(webTestingCampaignsInput)
		for _, key := range mapValue.MapKeys() {
			// Convert the key to a string for standardizing, 
			// the Normalize function will strictly validate if this string contains only digits.
			strKey := fmt.Sprintf("%v", key.Interface())
			inputMap[strKey] = mapValue.MapIndex(key).Interface()
		}

		normalizedMap, err := NormalizeWebTestingCampaignsMap(inputMap)
		if err != nil {
			serviceContainer.GetLoggerService().Error(
				"INVALID_WEB_TESTING_CAMPAIGNS",
				map[string]interface{}{"err": err.Error()},
				map[string]interface{}{"an": string(enums.ApiGetFlag), "uuid": context.GetUUID(), "sId": context.GetSessionId()},
			)
			platformVariables["webTestingCampaigns"] = nil // prevent duplicate logs
			return nil
		}
		platformVariables["webTestingCampaigns"] = normalizedMap
		return normalizedMap
	}

	// Some stacks pass JSON text (cookie, SSR prop, tag); parse it only if it's an object.
	if strInput, ok := webTestingCampaignsInput.(string); ok {
		trimmed := strings.TrimSpace(strInput)
		if trimmed == "" {
			return nil
		}

		// extract all "key": tokens and check for duplicates before parsing swallows them
		keyRegex := regexp.MustCompile(`"([^"\\]*)"\s*:`)
		matches := keyRegex.FindAllStringSubmatch(trimmed, -1)
		if matches != nil {
			keys := make(map[string]bool)
			hasDuplicate := false
			for _, match := range matches {
				if len(match) > 1 {
					key := match[1]
					if keys[key] {
						hasDuplicate = true
						break
					}
					keys[key] = true
				}
			}

			if hasDuplicate {
				serviceContainer.GetLoggerService().Error(
					"INVALID_WEB_TESTING_CAMPAIGNS_DUPLICATE_KEY",
					nil,
					map[string]interface{}{"an": string(enums.ApiGetFlag), "uuid": context.GetUUID(), "sId": context.GetSessionId()},
				)
			}
		}

		var parsedAssignments interface{}
		err := json.Unmarshal([]byte(trimmed), &parsedAssignments)
		if err != nil {
			serviceContainer.GetLoggerService().Error(
				"INVALID_WEB_TESTING_CAMPAIGNS_JSON",
				nil,
				map[string]interface{}{"an": string(enums.ApiGetFlag), "uuid": context.GetUUID(), "sId": context.GetSessionId()},
			)
			platformVariables["webTestingCampaigns"] = nil
			return nil
		}

		if parsedMap, ok := parsedAssignments.(map[string]interface{}); ok {
			normalizedMap, err := NormalizeWebTestingCampaignsMap(parsedMap)
			if err != nil {
				serviceContainer.GetLoggerService().Error(
					"INVALID_WEB_TESTING_CAMPAIGNS",
					map[string]interface{}{"err": err.Error()},
					map[string]interface{}{"an": string(enums.ApiGetFlag), "uuid": context.GetUUID(), "sId": context.GetSessionId()},
				)
				platformVariables["webTestingCampaigns"] = nil
				return nil
			}
			platformVariables["webTestingCampaigns"] = normalizedMap
			return normalizedMap
		}

		serviceContainer.GetLoggerService().Error(
			"INVALID_WEB_TESTING_CAMPAIGNS_JSON",
			nil,
			map[string]interface{}{"an": string(enums.ApiGetFlag), "uuid": context.GetUUID(), "sId": context.GetSessionId()},
		)
		platformVariables["webTestingCampaigns"] = nil
		return nil
	}

	// Booleans/numbers/other odd types are invalid.
	kind := reflect.TypeOf(webTestingCampaignsInput).Kind().String()
	serviceContainer.GetLoggerService().Error(
		"INVALID_WEB_TESTING_CAMPAIGNS_TYPE",
		map[string]interface{}{"kind": kind},
		map[string]interface{}{"an": string(enums.ApiGetFlag), "uuid": context.GetUUID(), "sId": context.GetSessionId()},
	)
	platformVariables["webTestingCampaigns"] = nil
	return nil
}

// EvaluateWebTestingCampaignVariation evaluates campaignVariation operand encoding.
func EvaluateWebTestingCampaignVariation(
	campaignVariationOperand string,
	assignedVariationsByCampaignID map[string]string,
) (result bool, invalidFormat bool) {
	assignments := assignedVariationsByCampaignID
	if assignments == nil {
		assignments = make(map[string]string)
	}

	// !123 — user should not be in campaign 123.
	match1 := regexp.MustCompile(`^!(\d+)$`).FindStringSubmatch(campaignVariationOperand)
	if match1 != nil {
		campaignID := match1[1]
		_, exists := assignments[campaignID]
		return !exists, false
	}

	// 123_!4 — in campaign 123 but not the variation 4.
	match2 := regexp.MustCompile(`^(\d+)_!(\d+)$`).FindStringSubmatch(campaignVariationOperand)
	if match2 != nil {
		campaignID := match2[1]
		variationID := match2[2]
		assignedVariationID, exists := assignments[campaignID]
		if !exists {
			return false, false
		}
		return assignedVariationID != variationID, false
	}

	// 123_4 — must be exactly that campaign and variation.
	match3 := regexp.MustCompile(`^(\d+)_(\d+)$`).FindStringSubmatch(campaignVariationOperand)
	if match3 != nil {
		campaignID := match3[1]
		variationID := match3[2]
		assignedVariationID, exists := assignments[campaignID]
		if !exists {
			return false, false
		}
		return assignedVariationID == variationID, false
	}

	// 123 — in the campaign, any variation counts.
	match4 := regexp.MustCompile(`^(\d+)$`).FindStringSubmatch(campaignVariationOperand)
	if match4 != nil {
		campaignID := match4[1]
		_, exists := assignments[campaignID]
		return exists, false
	}

	// Invalid format.
	return false, true
}
