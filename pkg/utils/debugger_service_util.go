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
	"fmt"

	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/interfaces"
)

// SendDebugEventToWingify sends a debug event to Wingify.
// @param settingsManager The settings manager containing configuration
// @param eventProps The properties for the event
func SendDebugEventToWingify(settingsManager interfaces.SettingsManagerInterface, eventProps map[string]interface{}) {
	if settingsManager == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			logger := settingsManager.GetLoggerService()
			if logger == nil {
				return
			}
			logger.Error("ERROR_SENDING_DEBUG_EVENT_TO_WINGIFY", map[string]interface{}{
				"err": fmt.Sprintf("%v", r),
			}, nil, false)
		}
	}()

	// Create query parameters
	properties := GetEventsBaseProperties(
		settingsManager,
		enums.DebuggerEvent.GetValue(),
		"",
		"",
	)

	// Create payload
	payload := GetDebuggerEventPayload(
		settingsManager,
		eventProps,
	)

	// Send event
	SendEventDirectlyToDACDN(
		settingsManager,
		properties,
		payload,
		enums.DebuggerEvent.GetValue(),
	)
}
