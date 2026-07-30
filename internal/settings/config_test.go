package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultSafetyAndMethod(t *testing.T) {
	cfg := Default()
	if cfg.InputMethod != "CVNSS4.0" || !cfg.Enabled || !cfg.PauseInPasswordFields {
		t.Fatalf("unsafe or unexpected defaults: %+v", cfg)
	}
}

func TestOpenMigratesMissingSafetyField(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	path := filepath.Join(dir, "BilaKeyPC", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"enabled":true,"input_method":"Telex"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	store, err := Open()
	if err != nil {
		t.Fatal(err)
	}
	cfg := store.Get()
	if cfg.InputMethod != "Telex" || !cfg.PauseInPasswordFields {
		t.Fatalf("migration lost values/defaults: %+v", cfg)
	}
}
