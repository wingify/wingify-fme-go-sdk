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

package wingify

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
	log "github.com/wingify/wingify-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/wingify-fme-go-sdk/pkg/models"
	settingsModel "github.com/wingify/wingify-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/wingify-fme-go-sdk/pkg/utils"
)

// Init initializes the Wingify FME client with the provided options
func Init(options map[string]interface{}) (clientInstance *WingifyClient, err error) {
	// handle panic and return error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to initialize Wingify FME client: %v", r)
		}
	}()

	// start time for init
	startTimeForInit := time.Now().UnixNano() / 1e6

	hostProfile, _ := options[enums.OptionHostProfile.GetValue()].(string)

	// Validate required parameters
	if options[enums.OptionSDKKey.GetValue()] == nil || options[enums.OptionSDKKey.GetValue()] == "" {
		return nil, fmt.Errorf("%s", log.BuildMessage(log.ErrorLogMessagesEnum["INVALID_SDK_KEY_IN_OPTIONS"], nil, hostProfile))
	}

	if options[enums.OptionAccountID.GetValue()] == nil || options[enums.OptionAccountID.GetValue()] == 0 {
		return nil, fmt.Errorf("%s", log.BuildMessage(log.ErrorLogMessagesEnum["INVALID_ACCOUNT_ID_IN_OPTIONS"], nil, hostProfile))
	}

	// Convert map to InitOptions using the factory function
	initOptions := models.NewInitOptions(options)

	if initOptions == nil {
		return nil, fmt.Errorf("%s", log.BuildMessage(log.ErrorLogMessagesEnum["INVALID_OPTIONS"], nil, hostProfile))
	}

	log.SetDefaultHostProfile(initOptions.HostProfile)

	// Create builder and setup services
	builder := &WingifyBuilder{
		options: initOptions,
	}

	builder.SetLogger().
		SetSettingsManager().
		SetNetworkManager().
		SetStorage().
		InitBatching().
		InitPolling()

	// Check if settings were provided in options
	builder.settingsManager.StartTimeForInit = startTimeForInit
	if initOptions.Settings != "" {
		// Parse and validate the provided settings
		builder.originalSettings = initOptions.Settings
		builder.settingsManager.IsSettingsProvidedInInit = true

		// Parse the settings JSON string into Settings object
		var settingsObj settingsModel.Settings
		err := json.Unmarshal([]byte(initOptions.Settings), &settingsObj)
		if err != nil {
			return nil, fmt.Errorf("failed to parse provided settings: %v", err)
		}

		// Set the parsed settings
		builder.settings = &settingsObj
		builder.settingsManager.SetSettings(&settingsObj, initOptions.Settings)

		clientInstance = builder.Build(builder.settings)
	} else {
		// Fetch settings from server
		builder.GetSettings(false)
		clientInstance = builder.Build(builder.settings)
	}

	return clientInstance, nil
}

// GetUUID generates a UUID for a user based on their userId and accountId
func GetUUID(userID string, accountID string) (uuid string, err error) {
	// check for valid userID and accountID
	if userID == "" || accountID == "" {
		return "", fmt.Errorf("userID and accountID are required")
	}

	// generate UUID using utils.GetUUID
	uuid = utils.GetUUID(userID, accountID)
	return uuid, nil
}
