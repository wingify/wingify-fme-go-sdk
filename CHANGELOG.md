# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.60.0] - 2026-06-29

### Added

- Support for **Web Testing pre-segmentation**: campaign segmentation can use the `campaignVariation` operand. The SDK evaluates it against **`context.platformVariables.webTestingCampaigns`**, a map of Web Testing campaign ID → variation ID (plain object or JSON string). The customer must pass this data in the context to enable web testing pre-segmentation. Supported operand values in settings: `122` (user in campaign), `122_2` (exact variation), `122_!1` (in campaign but not variation 1), `!122` (not in campaign).

  Example usage:

  ```go
  context := map[string]interface{}{
      "id": "user-123",
      "platformVariables": map[string]interface{}{
          // This is an example, replace with actual map
          "webTestingCampaigns": map[string]interface{}{
              "123": "4",
              "456": "1",
          },
      },
  }

  flag, err := wingifyClient.GetFlag("feature-key", context)
  ```

## [1.55.0] - 2026-06-15

### Added
- Added user tracking support: sends a `vwo_feUserTracked` event when user tracking is enabled for the account and no variation-shown impression was dispatched for the evaluation.

## [1.50.0] - 2026-06-05

### Added

- This release introduces Wingify as the primary SDK branding and package namespace

	```go
	package main

	import (
		"fmt"
		"log"

		wingify "github.com/wingify/wingify-fme-go-sdk"
	)

	func main() {
		// Initialize Wingify FME SDK with your account details
		options := map[string]interface{}{
			"sdkKey":    "32-alpha-numeric-sdk-key", // Replace with your SDK key
			"accountId": "123456",                   // Replace with your account ID
		}

		client, err := wingify.Init(options)
		if err != nil {
			log.Fatalf("Failed to initialize Wingify client: %v", err)
		}

		context := map[string]interface{}{
			"id": "unique_user_id", // Set a unique user identifier
		}

		getFlag, err := client.GetFlag("feature_key", context)
		if err != nil {
			log.Printf("Error getting feature flag: %v", err)
		} else {
			isFeatureEnabled := getFlag.IsEnabled()
			fmt.Println("Is feature enabled?", isFeatureEnabled)

			variableValue := getFlag.GetVariable("feature_variable", "default_value")
			fmt.Println("Variable value:", variableValue)
		}

		trackResponse, err := client.TrackEvent("event_name", context, nil)
		if err != nil {
			log.Printf("Error tracking event: %v", err)
		} else {
			fmt.Println("Event tracked:", trackResponse)
		}

		attributeMap := map[string]interface{}{
			"attribute-name": "attribute-value",
		}
		err = client.SetAttribute(attributeMap, context)
		if err != nil {
			log.Printf("Error setting attributes: %v", err)
		}
	}
	```
