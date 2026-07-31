package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	AppName    = "BilaKey PC"
	AppVersion = "2.5.0"
)

type Config struct {
	Enabled                   bool   `json:"enabled"`
	InputMethod               string `json:"input_method"`
	Encoding                  string `json:"encoding"`
	FreeToneMarking           bool   `json:"free_tone_marking"`
	OldToneStyle              bool   `json:"old_tone_style"`
	AlwaysUseClipboardUnicode bool   `json:"always_use_clipboard_unicode"`
	SpellCheck                bool   `json:"spell_check"`
	AutoRestoreWrongKey       bool   `json:"auto_restore_wrong_key"`
	ShowFeedback              bool   `json:"show_feedback"`
	MacroEnabled              bool   `json:"macro_enabled"`
	MacroWhileOff             bool   `json:"macro_while_off"`
	MacroFile                 string `json:"macro_file"`
	ShowDialogAtStartup       bool   `json:"show_dialog_at_startup"`
	StartWithWindows          bool   `json:"start_with_windows"`
	VietnameseInterface       bool   `json:"vietnamese_interface"`
	AutoCapInitial            bool   `json:"auto_cap_initial"`
	AutoCapSentence           bool   `json:"auto_cap_sentence"`
	DoubleShiftCaps           bool   `json:"double_shift_caps"`
	RestoreAfterDelimiter     bool   `json:"restore_after_delimiter"`
	PauseInPasswordFields     bool   `json:"pause_in_password_fields"`
}

func Default() Config {
	return Config{
		Enabled:               true,
		InputMethod:           "CVNSS4.0",
		Encoding:              "Unicode",
		FreeToneMarking:       true,
		OldToneStyle:          false,
		SpellCheck:            true,
		AutoRestoreWrongKey:   true,
		ShowFeedback:          false,
		MacroEnabled:          false,
		MacroWhileOff:         false,
		ShowDialogAtStartup:   true,
		StartWithWindows:      false,
		VietnameseInterface:   true,
		AutoCapInitial:        true,
		AutoCapSentence:       true,
		DoubleShiftCaps:       true,
		RestoreAfterDelimiter: true,
		PauseInPasswordFields: true,
	}
}

type Store struct {
	mu   sync.RWMutex
	path string
	cfg  Config
}

func ConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "BilaKeyPC"), nil
}

func DefaultPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func DefaultMacroPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "macros.tsv"), nil
}

func Open() (*Store, error) {
	path, err := DefaultPath()
	if err != nil {
		return nil, err
	}
	cfg := Default()
	macro, err := DefaultMacroPath()
	if err == nil {
		cfg.MacroFile = macro
	}

	data, readErr := os.ReadFile(path)
	if readErr == nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	} else if !errors.Is(readErr, os.ErrNotExist) {
		return nil, readErr
	}
	cfg.InputMethod = normalizeInputMethod(cfg.InputMethod)
	return &Store{path: path, cfg: cfg}, nil
}

func normalizeInputMethod(method string) string {
	key := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(method), " ", ""))
	switch key {
	case "VNI/TELEX", "TELEX/VNI", "VNI", "TELEX", "VNITELEX", "TELEXVNI":
		return "VNI/Telex"
	default:
		return "CVNSS4.0"
	}
}

func (s *Store) Get() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

func (s *Store) Update(fn func(*Config)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(&s.cfg)
	return s.saveLocked()
}

func (s *Store) Replace(cfg Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg = cfg
	return s.saveLocked()
}

func (s *Store) Path() string { return s.path }

func (s *Store) saveLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
