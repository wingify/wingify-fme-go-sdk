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

package models

import (
	"strconv"

	"github.com/wingify/wingify-fme-go-sdk/pkg/brand"
	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/storage"
)

// InitOptions represents SDK initialization options.
type InitOptions struct {
	AccountID            int                    `json:"accountId"`
	SDKKey               string                 `json:"sdkKey"`
	HostProfile          string                 `json:"hostProfile,omitempty"`
	Storage              storage.Connector      `json:"-"`
	GatewayService       map[string]interface{} `json:"gatewayService,omitempty"`
	PollInterval         int                    `json:"pollInterval,omitempty"`
	Logger               map[string]interface{} `json:"logger,omitempty"`
	Integrations         *IntegrationOptions    `json:"integrations,omitempty"`
	Settings             string                 `json:"settings,omitempty"`
	IsUsageStatsDisabled bool                   `json:"isUsageStatsDisabled,omitempty"`
	WingifyMeta          map[string]interface{} `json:"_vwo_meta,omitempty"`
	RetryConfig          *RetryConfig           `json:"retryConfig,omitempty"`
	IsAliasingEnabled    bool                   `json:"isAliasingEnabled,omitempty"`
	BatchEventData       map[string]interface{} `json:"batchEventData,omitempty"`
	ProxyURL             string                 `json:"proxyUrl,omitempty"`
}

// NewInitOptions builds InitOptions from the init map.
func NewInitOptions(options map[string]interface{}) *InitOptions {
	initOptions := &InitOptions{}

	if sdkKey, ok := options[enums.OptionSDKKey.GetValue()].(string); ok {
		initOptions.SDKKey = sdkKey
	}
	accountIDVal := options[enums.OptionAccountID.GetValue()]
	switch v := accountIDVal.(type) {
	case int:
		initOptions.AccountID = v
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			initOptions.AccountID = parsed
		}
	}
	if hostProfile, ok := options[enums.OptionHostProfile.GetValue()].(string); ok {
		initOptions.HostProfile = hostProfile
	}
	if storage, ok := options[enums.OptionStorage.GetValue()].(storage.Connector); ok {
		initOptions.Storage = storage
	}
	if gatewayService, ok := options[enums.OptionGatewayService.GetValue()].(map[string]interface{}); ok {
		initOptions.GatewayService = gatewayService
	}
	if pollInterval, ok := options[enums.OptionPollInterval.GetValue()].(int); ok {
		initOptions.PollInterval = pollInterval
	}
	if logger, ok := options[enums.OptionLogger.GetValue()].(map[string]interface{}); ok {
		initOptions.Logger = logger
	}
	if integrations, ok := options[enums.OptionIntegrations.GetValue()].(map[string]interface{}); ok {
		integrationOptions := &IntegrationOptions{}
		if callback, callbackOk := integrations[enums.IntegrationCallback.GetValue()].(func(map[string]interface{})); callbackOk {
			integrationOptions.Callback = callback
		}
		initOptions.Integrations = integrationOptions
	}
	if settings, ok := options[enums.OptionSettings.GetValue()].(string); ok {
		initOptions.Settings = settings
	}
	if isUsageStatsDisabled, ok := options[enums.OptionIsUsageStatsDisabled.GetValue()].(bool); ok {
		initOptions.IsUsageStatsDisabled = isUsageStatsDisabled
	}
	if wingifyMeta, ok := options[enums.OptionVWOMeta.GetValue()].(map[string]interface{}); ok {
		initOptions.WingifyMeta = wingifyMeta
	}
	if retryConfig, ok := options[enums.OptionRetryConfig.GetValue()].(map[string]interface{}); ok {
		initOptions.RetryConfig = NewRetryConfigFromMap(retryConfig)
	} else {
		initOptions.RetryConfig = NewRetryConfig()
	}
	if isAliasingEnabled, ok := options[enums.OptionIsAliasingEnabled.GetValue()].(bool); ok {
		initOptions.IsAliasingEnabled = isAliasingEnabled
	}
	if batchEventData, ok := options[enums.OptionBatchEventData.GetValue()].(map[string]interface{}); ok {
		initOptions.BatchEventData = batchEventData
	}
	if proxyURL, ok := options[enums.OptionProxyURL.GetValue()].(string); ok {
		initOptions.ProxyURL = proxyURL
	}
	return initOptions
}

// GetIntegrations returns the integrations options.
func (o *InitOptions) GetIntegrations() *IntegrationOptions {
	return o.Integrations
}

// GetBrandConfig resolves host profile branding for this init options instance.
func (o *InitOptions) GetBrandConfig() brand.Config {
	return brand.Resolve(brand.ParseProfile(o.HostProfile))
}

// IntegrationOptions represents integration callback options.
type IntegrationOptions struct {
	Callback func(properties map[string]interface{}) `json:"-"`
}

// NetworkOptions represents network configuration options.
type NetworkOptions struct {
	Client interface{} `json:"client,omitempty"`
}
