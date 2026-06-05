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
	"strconv"
	"time"

	"github.com/wingify/wingify-fme-go-sdk/pkg/constants"
	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
	log "github.com/wingify/wingify-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/wingify-fme-go-sdk/pkg/models"
	settingsModel "github.com/wingify/wingify-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/logger/core"
	loggerenums "github.com/wingify/wingify-fme-go-sdk/pkg/packages/logger/enums"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/network_layer/manager"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/storage"
	"github.com/wingify/wingify-fme-go-sdk/pkg/services"
)

// WingifyBuilder handles the construction of Wingify client instances
type WingifyBuilder struct {
	options                           *models.InitOptions
	logManager                        interfaces.LoggerServiceInterface
	settingsManager                   *services.SettingsManager
	settings                          *settingsModel.Settings
	batchEventQueue                   *services.BatchEventQueue
	wingifyClient                     *WingifyClient
	originalSettings                  string
	isSettingsFetchInProgress         bool
	isValidPollIntervalPassedFromInit bool
	pollingStopChan                   chan bool
	networkManager                    *manager.NetworkManager
}

// SetLogger sets up the logger service
func (builder *WingifyBuilder) SetLogger() *WingifyBuilder {
	loggerConfig := builder.applyBrandLoggerDefaults(builder.options.Logger)
	builder.logManager = core.NewLogManager(loggerConfig)
	builder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
		"service": "Logger",
	}))
	return builder
}

func (builder *WingifyBuilder) applyBrandLoggerDefaults(logger map[string]interface{}) map[string]interface{} {
	cfg := builder.options.GetBrandConfig()
	if logger == nil {
		logger = make(map[string]interface{})
	}
	if _, ok := logger[loggerenums.LogManagerConfigPrefix.GetValue()]; !ok {
		logger[loggerenums.LogManagerConfigPrefix.GetValue()] = cfg.LoggerPrefix
	}
	if _, ok := logger[loggerenums.LogManagerConfigName.GetValue()]; !ok {
		logger[loggerenums.LogManagerConfigName.GetValue()] = cfg.LoggerName
	}
	return logger
}

// SetSettingsManager sets up the settings manager
func (builder *WingifyBuilder) SetSettingsManager() *WingifyBuilder {
	builder.settingsManager = services.NewSettingsManager(builder.options, builder.logManager)
	builder.logManager.SetSettingsManager(builder.settingsManager)
	return builder
}

// SetNetworkManager sets up the network manager
func (builder *WingifyBuilder) SetNetworkManager() *WingifyBuilder {
	// Network manager is a singleton, just attach default client
	builder.networkManager = &manager.NetworkManager{}

	// Use retry configuration if provided
	if builder.options != nil && builder.options.RetryConfig != nil {
		builder.networkManager.AttachDefaultClientWithRetry(builder.options.RetryConfig)
		builder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Network Manager with Retry",
			"retryConfig": map[string]interface{}{
				"shouldRetry":       builder.options.RetryConfig.ShouldRetry,
				"maxRetries":        builder.options.RetryConfig.MaxRetries,
				"initialDelay":      builder.options.RetryConfig.InitialDelay,
				"backoffMultiplier": builder.options.RetryConfig.BackoffMultiplier,
			},
		}))
	} else {
		builder.networkManager.AttachDefaultClient()
		builder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Network Manager",
		}))
	}

	builder.settingsManager.SetNetworkManager(builder.networkManager)
	return builder
}

func (builder *WingifyBuilder) InitBatching() *WingifyBuilder {
	// Check if batch event data is provided in options
	if builder.options.BatchEventData != nil {
		// Check if gatewayService is provided and skip SDK batching if so
		if builder.settingsManager != nil && builder.settingsManager.GetIsGatewayServiceProvided() {
			builder.logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["GATEWAY_AND_BATCH_EVENTS_CONFIG_MISMATCH"], nil))
			return builder
		}
		batchEventData := models.NewBatchEventData(builder.options.BatchEventData)
		eventsPerRequest := batchEventData.GetEventsPerRequest()
		requestTimeInterval := batchEventData.GetRequestTimeInterval()

		isEventsPerRequestValid := eventsPerRequest > 0 && eventsPerRequest <= constants.MaxEventsPerRequest
		isRequestTimeIntervalValid := requestTimeInterval > 0

		// Handle invalid data types for individual parameters
		if !isEventsPerRequestValid {
			builder.logManager.Error("INVALID_EVENTS_PER_REQUEST_VALUE", nil, map[string]interface{}{"an": enums.ApiInit})
			eventsPerRequest = constants.DefaultEventsPerRequest // Use default if invalid
		}

		if !isRequestTimeIntervalValid {
			builder.logManager.Error("INVALID_REQUEST_TIME_INTERVAL_VALUE", nil, map[string]interface{}{"an": enums.ApiInit})
			requestTimeInterval = constants.DefaultRequestTimeInterval // Use default if invalid
		}

		// Initialize BatchEventQueue for batching
		// Convert BatchFlushCallback to FlushCallback
		var flushCallback models.FlushCallback
		if batchEventData.GetFlushCallback() != nil {
			flushCallback = func(err string, events string) {
				// Convert string error to error type for BatchFlushCallback
				var errorObj error
				if err != "" {
					errorObj = fmt.Errorf("%s", err)
				}
				batchEventData.GetFlushCallback()(errorObj, events)
			}
		}

		builder.batchEventQueue = services.NewBatchEventQueue(
			eventsPerRequest,
			requestTimeInterval,
			flushCallback,
			builder.options.AccountID,
			builder.options.SDKKey,
			builder.logManager,
			builder.settingsManager,
		)
		builder.batchEventQueue.SetNetworkManager(builder.networkManager)

		builder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Batching",
		}))
	}

	return builder
}

// SetStorage sets up the storage service
func (builder *WingifyBuilder) SetStorage() *WingifyBuilder {
	if builder.options != nil && builder.options.Storage != nil {
		// Attach the storage connector to the singleton
		storageInstance := storage.GetInstance()
		storageInstance.AttachConnector(builder.options.Storage)

		builder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Storage",
		}))
	}
	return builder
}

// InitPolling initializes the polling mechanism
func (builder *WingifyBuilder) InitPolling() *WingifyBuilder {
	if builder.options.PollInterval >= 1000 && builder.options.PollInterval != 0 {
		// This is to check if the poll_interval passed in options is valid
		builder.isValidPollIntervalPassedFromInit = true
		builder.pollingStopChan = make(chan bool)
		go builder.checkAndPoll()
		return builder
	} else if builder.options.PollInterval > 0 {
		// Only log error if poll_interval is present in options but invalid
		builder.logManager.Error("INVALID_POLLING_CONFIGURATION", map[string]interface{}{
			"key":         "pollInterval",
			"correctType": "number",
		}, map[string]interface{}{"an": enums.ApiInit})
	}
	return builder
}

// InitUsageStats initializes usage statistics (placeholder for future implementation)
func (builder *WingifyBuilder) InitUsageStats() *WingifyBuilder {
	builder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
		"service": "Usage Stats",
	}))
	return builder
}

// GetSettings fetches settings from Wingify servers
func (builder *WingifyBuilder) GetSettings(forceFetch bool) string {
	settingsString := builder.settingsManager.GetSettings(forceFetch)
	builder.originalSettings = settingsString
	builder.settings = builder.settingsManager.GetSettingsObject()
	return settingsString
}

// Build creates and returns a WingifyClient instance
func (builder *WingifyBuilder) Build(settingsData *settingsModel.Settings) *WingifyClient {
	// Create Wingify client using the newWingifyClient function
	// This will process settings before creating the client
	wingifyClient := newWingifyClient(settingsData, builder)

	// Set Wingify client reference in builder
	builder.wingifyClient = wingifyClient
	// If poll_interval is not present in options, set it to the pollInterval from settings
	builder.updatePollIntervalAndCheckAndPoll(builder.originalSettings, true)
	return wingifyClient
}

// updatePollIntervalAndCheckAndPoll updates the poll interval from settings and starts polling if needed
func (builder *WingifyBuilder) updatePollIntervalAndCheckAndPoll(settingsJSON string, shouldCheckAndPoll bool) {
	// Only update the poll_interval if poll_interval is not valid or not present in options
	var processedSettings *settingsModel.Settings
	if settingsJSON != "" {
		err := json.Unmarshal([]byte(settingsJSON), &processedSettings)
		if err != nil {
			// Ignore error, processedSettings will be nil
			processedSettings = nil
		}
	}

	if !builder.isValidPollIntervalPassedFromInit && processedSettings != nil {
		builder.options.PollInterval = processedSettings.GetPollInterval()

		if processedSettings.GetPollInterval() == constants.DefaultPollInterval {
			builder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["USING_POLL_INTERVAL_FROM_SETTINGS"], map[string]interface{}{
				"source":       "default",
				"pollInterval": strconv.Itoa(constants.DefaultPollInterval),
			}))
		} else {
			builder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["USING_POLL_INTERVAL_FROM_SETTINGS"], map[string]interface{}{
				"source":       "settings",
				"pollInterval": strconv.Itoa(builder.options.PollInterval),
			}))
		}
	}

	if shouldCheckAndPoll && !builder.isValidPollIntervalPassedFromInit && processedSettings != nil && builder.options.PollInterval >= 1000 {
		builder.pollingStopChan = make(chan bool)
		go builder.checkAndPoll()
	}
}

// checkAndPoll checks for settings updates at the configured interval
func (builder *WingifyBuilder) checkAndPoll() {
	ticker := time.NewTicker(time.Duration(builder.options.PollInterval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			latestSettings := builder.fetchSettings(true)
			if builder.originalSettings != "" && latestSettings != "" {
				// Compare settings
				if !builder.areSettingsEqual(builder.originalSettings, latestSettings) {
					builder.updateSettingsOnBuilder(latestSettings)
				} else {
					builder.logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["POLLING_NO_CHANGE_IN_SETTINGS"], map[string]interface{}{}))
				}
			} else if builder.originalSettings == "" && latestSettings != "" {
				builder.updateSettingsOnBuilder(latestSettings)
			}

		case <-builder.pollingStopChan:
			return
		}
	}
}

// fetchSettings fetches settings from Wingify servers
func (builder *WingifyBuilder) fetchSettings(forceFetch bool) string {
	// Check if a fetch operation is already in progress
	if builder.isSettingsFetchInProgress || builder.settingsManager == nil {
		return ""
	}

	apiName := string(enums.ApiInit)
	if forceFetch {
		apiName = constants.POLLING
	}

	// Set the flag to indicate that a fetch operation is in progress
	builder.isSettingsFetchInProgress = true
	defer func() {
		if r := recover(); r != nil {
			builder.logManager.Error("ERROR_FETCHING_SETTINGS", map[string]interface{}{
				"err": r,
			}, map[string]interface{}{"an": apiName})
			builder.isSettingsFetchInProgress = false
		}
	}()

	// Retrieve the settings
	settingsString := builder.settingsManager.GetSettings(forceFetch)

	if !forceFetch {
		// Store the original settings
		builder.originalSettings = settingsString
	}
	builder.isSettingsFetchInProgress = false
	return settingsString
}

// areSettingsEqual compares two settings JSON strings
func (builder *WingifyBuilder) areSettingsEqual(settings1, settings2 string) bool {
	var obj1, obj2 interface{}

	if err := json.Unmarshal([]byte(settings1), &obj1); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(settings2), &obj2); err != nil {
		return false
	}

	// Deep comparison using JSON marshaling
	json1, _ := json.Marshal(obj1)
	json2, _ := json.Marshal(obj2)

	return string(json1) == string(json2)
}

// updateSettingsOnBuilder updates the settings on the WingifyBuilder instance
func (builder *WingifyBuilder) updateSettingsOnBuilder(latestSettings string) {
	if builder.wingifyClient != nil {
		err := builder.wingifyClient.updateSettingsInternal(latestSettings)
		if err == nil {
			builder.logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["POLLING_SET_SETTINGS"], map[string]interface{}{}))
			builder.originalSettings = latestSettings
			builder.updatePollIntervalAndCheckAndPoll(builder.originalSettings, false)
		} else {
			builder.logManager.Error("ERROR_UPDATING_SETTINGS", map[string]interface{}{
				"err":              err.Error(),
				"originalSettings": builder.originalSettings,
				"latestSettings":   latestSettings,
			}, map[string]interface{}{"an": constants.POLLING})
		}
	}
}

// StopPolling stops the polling goroutine
func (builder *WingifyBuilder) StopPolling() {
	if builder.pollingStopChan != nil {
		close(builder.pollingStopChan)
	}
}
