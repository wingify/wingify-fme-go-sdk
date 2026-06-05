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

package brand

import "strings"

// Profile selects default hosts, SDK identity, and logger branding.
type Profile string

const (
	ProfileWingify Profile = "wingify"
	ProfileVWO     Profile = "vwo"
)

const (
	vwoSettingsHost     = "dev.visualwebsiteoptimizer.com"
	wingifySettingsHost = "edge.wingify.net"
	wingifyEventsHost   = "collect.wingify.net"
)

// Config holds profile-specific branding and default hosts.
type Config struct {
	Profile      Profile
	SDKName      string
	SettingsHost string
	EventsHost   string
	LoggerName   string
	LoggerPrefix string
}

// DisplayName returns the human-readable brand label for log messages.
func DisplayName(profile Profile) string {
	if profile == ProfileVWO {
		return "VWO"
	}
	return "Wingify"
}

// DisplayNameFromHostProfile resolves the brand label from a hostProfile init value.
func DisplayNameFromHostProfile(hostProfile string) string {
	return DisplayName(ParseProfile(hostProfile))
}

// ParseProfile normalizes hostProfile init option values.
func ParseProfile(value string) Profile {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(ProfileVWO):
		return ProfileVWO
	default:
		return ProfileWingify
	}
}

// Resolve returns branding and default hosts for a profile.
func Resolve(profile Profile) Config {
	switch profile {
	case ProfileVWO:
		return Config{
			Profile:      ProfileVWO,
			SDKName:      "vwo-fme-go-sdk",
			SettingsHost: vwoSettingsHost,
			EventsHost:   vwoSettingsHost,
			LoggerName:   "VWO Logger",
			LoggerPrefix: "VWO-SDK",
		}
	default:
		return Config{
			Profile:      ProfileWingify,
			SDKName:      "wingify-fme-go-sdk",
			SettingsHost: wingifySettingsHost,
			EventsHost:   wingifyEventsHost,
			LoggerName:   "Wingify Logger",
			LoggerPrefix: "Wingify-SDK",
		}
	}
}

// HostFor returns the hostname for the given request kind.
// When unifiedHost is non-empty (gateway/proxy), it is used for all kinds.
func (c Config) HostFor(kind HostKind, unifiedHost string) string {
	if unifiedHost != "" {
		return unifiedHost
	}
	switch kind {
	case HostKindSettings:
		return c.SettingsHost
	default:
		return c.EventsHost
	}
}
