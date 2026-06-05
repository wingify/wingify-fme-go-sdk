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

package user

import (
	"time"

	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
)

// WingifyUserContext represents the user context for Wingify FME operations
type WingifyUserContext struct {
	ID                          string                 `json:"id"`
	UserAgent                   string                 `json:"userAgent,omitempty"`
	IPAddress                   string                 `json:"ipAddress,omitempty"`
	CustomVariables             map[string]interface{} `json:"customVariables,omitempty"`
	VariationTargetingVariables map[string]interface{} `json:"variationTargetingVariables,omitempty"`
	PostSegmentationVariables   []string               `json:"postSegmentationVariables,omitempty"`
	Wingify                     *ContextWingify        `json:"_wingify,omitempty"`
	SessionId                   int64                  `json:"sessionId,omitempty"`
	UUID                        string                 `json:"uuid,omitempty"`
	BucketingSeed               string                 `json:"bucketingSeed,omitempty"`
}

// NewWingifyUserContext creates a new WingifyUserContext from a map
func NewWingifyUserContext(context map[string]interface{}) *WingifyUserContext {
	wingifyUserContext := &WingifyUserContext{}

	if id, ok := context[enums.ContextID.GetValue()]; ok {
		if idStr, ok := id.(string); ok {
			wingifyUserContext.ID = idStr
		}
	}

	if userAgent, ok := context[enums.ContextUserAgent.GetValue()]; ok {
		if userAgentStr, ok := userAgent.(string); ok {
			wingifyUserContext.UserAgent = userAgentStr
		}
	}

	if ipAddress, ok := context[enums.ContextIPAddress.GetValue()]; ok {
		if ipAddressStr, ok := ipAddress.(string); ok {
			wingifyUserContext.IPAddress = ipAddressStr
		}
	}

	if customVariables, ok := context[enums.ContextCustomVariables.GetValue()]; ok {
		if customVarsMap, ok := customVariables.(map[string]interface{}); ok {
			wingifyUserContext.CustomVariables = customVarsMap
		}
	}

	if variationTargetingVariables, ok := context[enums.ContextVariationTargetingVariables.GetValue()]; ok {
		if vtVarsMap, ok := variationTargetingVariables.(map[string]interface{}); ok {
			wingifyUserContext.VariationTargetingVariables = vtVarsMap
		}
	}

	if postSegmentationVariables, ok := context[enums.ContextPostSegmentationVariables.GetValue()]; ok {
		if psvArray, ok := postSegmentationVariables.([]string); ok {
			wingifyUserContext.PostSegmentationVariables = psvArray
		}
	}

	if wingify, ok := context[enums.ContextWingify.GetValue()]; ok {
		if wingifyMap, ok := wingify.(map[string]interface{}); ok {
			wingifyUserContext.Wingify = NewContextWingify(wingifyMap)
		}
	}
	sessionId := time.Now().Unix()
	if contextSessionId, ok := context[enums.ContextSessionID.GetValue()]; ok {
		if id, ok := contextSessionId.(int); ok {
			sessionId = int64(id)
		} else if id, ok := contextSessionId.(int64); ok {
			sessionId = id
		}
	}
	wingifyUserContext.SessionId = sessionId

	if bucketingSeed, ok := context[enums.ContextBucketingSeed.GetValue()]; ok {
		if bucketingSeedStr, ok := bucketingSeed.(string); ok {
			wingifyUserContext.BucketingSeed = bucketingSeedStr
		}
	}

	return wingifyUserContext
}

// GetID returns the user ID
func (c *WingifyUserContext) GetID() string {
	return c.ID
}

// GetUserAgent returns the user agent
func (c *WingifyUserContext) GetUserAgent() string {
	return c.UserAgent
}

// GetIPAddress returns the IP address
func (c *WingifyUserContext) GetIPAddress() string {
	return c.IPAddress
}

// GetCustomVariables returns custom variables
func (c *WingifyUserContext) GetCustomVariables() map[string]interface{} {
	return c.CustomVariables
}

// SetCustomVariables sets custom variables
func (c *WingifyUserContext) SetCustomVariables(customVariables map[string]interface{}) {
	c.CustomVariables = customVariables
}

// GetVariationTargetingVariables returns variation targeting variables
func (c *WingifyUserContext) GetVariationTargetingVariables() map[string]interface{} {
	return c.VariationTargetingVariables
}

// SetVariationTargetingVariables sets variation targeting variables
func (c *WingifyUserContext) SetVariationTargetingVariables(variationTargetingVariables map[string]interface{}) {
	c.VariationTargetingVariables = variationTargetingVariables
}

// GetWingify returns the gateway-enriched context metadata
func (c *WingifyUserContext) GetWingify() *ContextWingify {
	return c.Wingify
}

// SetWingify sets the gateway-enriched context metadata
func (c *WingifyUserContext) SetWingify(wingify *ContextWingify) {
	c.Wingify = wingify
}

// GetPostSegmentationVariables returns post segmentation variables
func (c *WingifyUserContext) GetPostSegmentationVariables() []string {
	return c.PostSegmentationVariables
}

// SetPostSegmentationVariables sets post segmentation variables
func (c *WingifyUserContext) SetPostSegmentationVariables(postSegmentationVariables []string) {
	c.PostSegmentationVariables = postSegmentationVariables
}

// GetSessionId returns the session ID
func (c *WingifyUserContext) GetSessionId() int64 {
	return c.SessionId
}

// GetUUID returns the UUID
func (c *WingifyUserContext) GetUUID() string {
	return c.UUID
}

// SetSessionId sets the session ID
func (c *WingifyUserContext) SetSessionId(sessionId int64) {
	c.SessionId = sessionId
}

// SetUUID sets the UUID
func (c *WingifyUserContext) SetUUID(uuid string) {
	c.UUID = uuid
}

// GetBucketingSeed returns the custom bucketing seed
func (c *WingifyUserContext) GetBucketingSeed() string {
	return c.BucketingSeed
}

// SetBucketingSeed sets the custom bucketing seed
func (c *WingifyUserContext) SetBucketingSeed(bucketingSeed string) {
	c.BucketingSeed = bucketingSeed
}
