package brand

import "testing"

func TestResolveWingifyHosts(t *testing.T) {
	cfg := Resolve(ProfileWingify)
	if cfg.SettingsHost != wingifySettingsHost {
		t.Fatalf("settings host: got %s", cfg.SettingsHost)
	}
	if cfg.EventsHost != wingifyEventsHost {
		t.Fatalf("events host: got %s", cfg.EventsHost)
	}
	if cfg.HostFor(HostKindSettings, "") != wingifySettingsHost {
		t.Fatal("settings HostFor mismatch")
	}
	if cfg.HostFor(HostKindEvents, "") != wingifyEventsHost {
		t.Fatal("events HostFor mismatch")
	}
}

func TestResolveVWOSingleHost(t *testing.T) {
	cfg := Resolve(ProfileVWO)
	if cfg.SettingsHost != cfg.EventsHost {
		t.Fatal("vwo profile should use one host for settings and events")
	}
	if cfg.HostFor(HostKindEvents, "proxy.example.com") != "proxy.example.com" {
		t.Fatal("unified host override expected")
	}
}

func TestParseProfile(t *testing.T) {
	if ParseProfile("vwo") != ProfileVWO {
		t.Fatal("expected vwo profile")
	}
	if ParseProfile("") != ProfileWingify {
		t.Fatal("expected default wingify profile")
	}
}
