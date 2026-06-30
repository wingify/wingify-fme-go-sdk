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

package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
	"github.com/wingify/wingify-fme-go-sdk/pkg/models/user"
	"github.com/wingify/wingify-fme-go-sdk/pkg/models/campaign"
	loggerCore "github.com/wingify/wingify-fme-go-sdk/pkg/packages/logger/core"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/interfaces"
	segmentationCore "github.com/wingify/wingify-fme-go-sdk/pkg/packages/segmentation_evaluator/core"
	"github.com/wingify/wingify-fme-go-sdk/pkg/packages/segmentation_evaluator/utils"
	"github.com/wingify/wingify-fme-go-sdk/pkg/services"
)

type mockSettingsManager struct {
	interfaces.SettingsManagerInterface
}

func (m *mockSettingsManager) GetIsGatewayServiceProvided() bool { return false }

type mockServiceContainer struct {
	interfaces.ServiceContainerInterface
	logger interfaces.LoggerServiceInterface
}

func (m *mockServiceContainer) GetLoggerService() interfaces.LoggerServiceInterface { return m.logger }
func (m *mockServiceContainer) GetSettingsManager() interfaces.SettingsManagerInterface {
	return &mockSettingsManager{}
}
func (m *mockServiceContainer) GetDebuggerService() interfaces.DebuggerServiceInterface {
	return services.NewDebuggerService()
}

func TestWebTestingCampaignVariation(t *testing.T) {
	assignmentsMap := map[string]string{"1": "1", "2": "2"}

	t.Run("EvaluateWebTestingCampaignVariation", func(t *testing.T) {
		// C_V matches
		res, inv := utils.EvaluateWebTestingCampaignVariation("1_1", assignmentsMap)
		assert.True(t, res)
		assert.False(t, inv)

		// C_V differs
		res, _ = utils.EvaluateWebTestingCampaignVariation("1_2", assignmentsMap)
		assert.False(t, res)

		// C_V not in campaign
		res, _ = utils.EvaluateWebTestingCampaignVariation("99_1", assignmentsMap)
		assert.False(t, res)

		// C_!V matches
		res, _ = utils.EvaluateWebTestingCampaignVariation("1_!2", assignmentsMap)
		assert.True(t, res)

		// C_!V differs
		res, _ = utils.EvaluateWebTestingCampaignVariation("1_!1", assignmentsMap)
		assert.False(t, res)

		// !C matches
		res, _ = utils.EvaluateWebTestingCampaignVariation("!99", assignmentsMap)
		assert.True(t, res)

		// !C differs
		res, _ = utils.EvaluateWebTestingCampaignVariation("!1", assignmentsMap)
		assert.False(t, res)

		// Empty map
		res, _ = utils.EvaluateWebTestingCampaignVariation("!1", nil)
		assert.True(t, res)
		res, _ = utils.EvaluateWebTestingCampaignVariation("1_1", nil)
		assert.False(t, res)

		// Invalid encoding
		res, inv = utils.EvaluateWebTestingCampaignVariation("bogus", assignmentsMap)
		assert.False(t, res)
		assert.True(t, inv)

		// C alone
		res, _ = utils.EvaluateWebTestingCampaignVariation("100", map[string]string{"100": "1"})
		assert.True(t, res)
		res, _ = utils.EvaluateWebTestingCampaignVariation("100", map[string]string{"99": "1"})
		assert.False(t, res)
	})

	t.Run("NormalizeWebTestingCampaignsMap - Valid", func(t *testing.T) {
		raw := map[string]interface{}{
			"129": 1,
			"14":  "2",
		}
		normalized, err := utils.NormalizeWebTestingCampaignsMap(raw)
		assert.NoError(t, err)
		expected := map[string]string{"129": "1", "14": "2"}
		assert.Equal(t, expected, normalized)
	})

	t.Run("NormalizeWebTestingCampaignsMap - Invalid Inputs", func(t *testing.T) {
		invalidInputs := []map[string]interface{}{
			{"": "1"},               // empty string key
			{"1": ""},               // empty string value
			{"kaustubh": "1"},       // non-digit string key
			{"1": "some"},           // non-digit string value
			{"-1": "1"},             // negative string key
			{"1": "-1"},             // negative string value
			{"1.5": "1"},            // float string key
			{"1": "1.5"},            // float string value
			{"1": nil},              // null value
			{"1": true},             // boolean value
			{"1": []int{1}},         // array value
			{"1": map[string]int{}}, // object value
		}

		for _, invalidInput := range invalidInputs {
			normalized, err := utils.NormalizeWebTestingCampaignsMap(invalidInput)
			assert.Error(t, err)
			assert.Nil(t, normalized)
		}
	})

	t.Run("SegmentEvaluator campaignVariation DSL", func(t *testing.T) {
		logManager := loggerCore.NewLogManager(nil)
		segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)

		// DSL with webTestingCampaigns
		dsl := `{"or":[{"campaignVariation":"1_1"}]}`
		platformVars := map[string]interface{}{
			"webTestingCampaigns": `{"1":"1"}`,
		}
		ctx := user.NewWingifyUserContext(map[string]interface{}{
			enums.ContextID.GetValue():                "u1",
			enums.ContextPlatformVariables.GetValue(): platformVars,
		})
		feature := &campaign.Feature{IsGatewayServiceRequired: false}
		segmentationManager.SetContextualData(&mockServiceContainer{logger: logManager}, feature, ctx)
		assert.True(t, segmentationManager.ValidateSegmentation(dsl, map[string]interface{}{}))

		// Not in campaign
		dsl2 := `{"or":[{"campaignVariation":"!1"}]}`
		platformVars2 := map[string]interface{}{
			"webTestingCampaigns": `{}`,
		}
		ctx2 := user.NewWingifyUserContext(map[string]interface{}{
			enums.ContextID.GetValue():                "u1",
			enums.ContextPlatformVariables.GetValue(): platformVars2,
		})
		segmentationManager.SetContextualData(&mockServiceContainer{logger: logManager}, feature, ctx2)
		assert.True(t, segmentationManager.ValidateSegmentation(dsl2, map[string]interface{}{}))

		// Numeric value handling
		dsl3 := `{"or":[{"campaignVariation":104}]}`
		platformVars3 := map[string]interface{}{
			"webTestingCampaigns": `{"104":"1"}`,
		}
		ctx3 := user.NewWingifyUserContext(map[string]interface{}{
			enums.ContextID.GetValue():                "u1",
			enums.ContextPlatformVariables.GetValue(): platformVars3,
		})
		segmentationManager.SetContextualData(&mockServiceContainer{logger: logManager}, feature, ctx3)
		assert.True(t, segmentationManager.ValidateSegmentation(dsl3, map[string]interface{}{}))
	})

	t.Run("ParseWebTestingCampaignsFromContext duplicate key detection", func(t *testing.T) {
		logManager := loggerCore.NewLogManager(nil)
		sc := &mockServiceContainer{logger: logManager}
		platformVars := map[string]interface{}{
			"webTestingCampaigns": `{"1":0,"1":1}`,
		}
		ctx := user.NewWingifyUserContext(map[string]interface{}{
			enums.ContextID.GetValue():                "u1",
			enums.ContextPlatformVariables.GetValue(): platformVars,
		})

		parsed := utils.ParseWebTestingCampaignsFromContext(ctx, sc)
		// Usually json unmarshaler takes the last key
		assert.Equal(t, map[string]string{"1": "1"}, parsed)
	})

	t.Run("SegmentationManager validateSegmentation - no platformVariables", func(t *testing.T) {
		logManager := loggerCore.NewLogManager(nil)
		segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)
		sc := &mockServiceContainer{logger: logManager}

		ctx := user.NewWingifyUserContext(map[string]interface{}{
			enums.ContextID.GetValue(): "u1",
		})
		feature := &campaign.Feature{IsGatewayServiceRequired: false}
		segmentationManager.SetContextualData(sc, feature, ctx)

		// fails silently if no platformVariables for campaignVariation
		dsl := `{"or":[{"campaignVariation":"122_2"}]}`
		assert.False(t, segmentationManager.ValidateSegmentation(dsl, map[string]interface{}{}))

		// non-web DSL not affected
		dsl2 := `{"or":[{"custom_variable":{"plan":"premium"}}]}`
		// fails properly because no custom_variable match
		assert.False(t, segmentationManager.ValidateSegmentation(dsl2, map[string]interface{}{}))
	})
}
