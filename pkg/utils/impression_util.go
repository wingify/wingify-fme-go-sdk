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
	"net/url"
	"strings"

	"github.com/wingify/wingify-fme-go-sdk/pkg/constants"
	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
	log "github.com/wingify/wingify-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/wingify-fme-go-sdk/pkg/models/user"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/interfaces"
)

// CreateAndSendImpressionForVariationShown creates and sends an impression for a variation shown event.
// This function constructs the necessary properties and payload for the event
// and uses the NetworkUtil to send a POST API request.
func CreateAndSendImpressionForVariationShown(
	serviceContainer interfaces.ServiceContainerInterface,
	campaignID int,
	variationID int,
	context *user.WingifyUserContext,
	featureKey string,
) {
	// Get base properties for the event
	properties := GetEventsBaseProperties(
		serviceContainer.GetSettingsManager(),
		enums.VariationShown.GetValue(),
		EncodeURIComponent(context.GetUserAgent()),
		context.GetIPAddress(),
	)

	// Construct payload data for tracking the user
	payload := GetTrackUserPayloadData(
		serviceContainer,
		enums.VariationShown.GetValue(),
		campaignID,
		variationID,
		context,
	)

	// get campaign key and variation name
	campaignKeyWithVariationName := GetCampaignKeyFromCampaignID(serviceContainer.GetSettings(), campaignID)
	variationName := GetVariationNameFromCampaignIdAndVariationId(serviceContainer.GetSettings(), campaignID, variationID)

	campaignKey := campaignKeyWithVariationName // default campaign key

	if campaignKeyWithVariationName == featureKey {
		campaignKey = constants.IMPACT_ANALYSIS
	} else {
		// split campaignKeyWithVariationName by featureKey_ and get the first part
		campaignKey = strings.Split(campaignKeyWithVariationName, featureKey+"_")[1]
	}

	// get campaign type
	campaignType := GetCampaignTypeFromCampaignId(serviceContainer.GetSettings(), campaignID)

	// Check if batch event queue is available
	if serviceContainer.GetBatchEventQueue().IsInitialized() {
		// Enqueue the event to the batch queue for future processing
		serviceContainer.GetBatchEventQueue().Enqueue(payload)
	} else {
		// Send the event immediately if batch event queue is not available
		SendPostAPIRequest(serviceContainer, properties, payload, context, map[string]interface{}{
			"campaignKey":   campaignKey,
			"variationName": variationName,
			"featureKey":    featureKey,
			"campaignType":  campaignType,
		})
	}
}

// CreateAndSendImpressionForUsageTracking creates and sends an impression for a usage tracking event.
func CreateAndSendImpressionForUsageTracking(
	serviceContainer interfaces.ServiceContainerInterface,
	context *user.WingifyUserContext,
	featureKey string,
) {
	// Get base properties for the event
	properties := GetEventsBaseProperties(
		serviceContainer.GetSettingsManager(),
		enums.FeTrackUsage.GetValue(),
		EncodeURIComponent(context.GetUserAgent()),
		context.GetIPAddress(),
	)

	// Construct payload data for tracking the usage
	var vwoMeta map[string]interface{}
	if serviceContainer.GetInitOptions() != nil {
		vwoMeta = serviceContainer.GetInitOptions().WingifyMeta
	}

	payload := GetTrackUsagePayloadData(
		serviceContainer,
		context,
		vwoMeta,
	)

	// Log info about track usage dispatch
	serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["TRACK_USAGE_DISPATCHED"], map[string]interface{}{
		"accountId":  serviceContainer.GetSettingsManager().GetAccountID(),
		"userId":     context.GetID(),
		"featureKey": featureKey,
	}))

	// Check if batch event queue is available, if not available then send the event immediately
	if serviceContainer.GetBatchEventQueue() != nil && serviceContainer.GetBatchEventQueue().IsInitialized() {
		// Enqueue the event to the batch queue for future processing
		serviceContainer.GetBatchEventQueue().Enqueue(payload)
	} else {
		// Send the event immediately if batch event queue is not available
		SendPostAPIRequest(serviceContainer, properties, payload, context, map[string]interface{}{
			"featureKey": featureKey,
		})
	}
}

// EncodeURIComponent encodes a string to be URL-safe
func EncodeURIComponent(value string) string {
	return url.QueryEscape(value)
}
