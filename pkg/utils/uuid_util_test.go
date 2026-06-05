package utils

import (
	"testing"

	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
)

func TestGetUUIDFromContextNilSettings(t *testing.T) {
	context := map[string]interface{}{
		enums.ContextID.GetValue(): "user-123",
	}

	uuid, isWebUuid, err := GetUUIDFromContext(context, "getFlag", 1210020, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uuid == "" {
		t.Fatal("expected non-empty uuid")
	}
	if isWebUuid {
		t.Fatal("expected server-side uuid when settings are nil")
	}
}
