package macro

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Entry struct {
	Trigger     string
	Replacement string
}

type Table struct {
	mu      sync.RWMutex
	path    string
	entries map[string]string
}

func New(path string) *Table {
	return &Table{path: path, entries: make(map[string]string)}
}

func (t *Table) Path() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.path
}

func (t *Table) SetPath(path string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.path = path
}

func (t *Table) Load() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	f, err := os.Open(t.path)
	if errors.Is(err, os.ErrNotExist) {
		t.entries = make(map[string]string)
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	entries := make(map[string]string)
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		trigger := strings.TrimSpace(parts[0])
		if trigger == "" {
			continue
		}
		entries[trigger] = strings.TrimSpace(parts[1])
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	t.entries = entries
	return nil
}

func (t *Table) Save() error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.path == "" {
		return errors.New("macro path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(t.path), 0o755); err != nil {
		return err
	}
	keys := make([]string, 0, len(t.entries))
	for k := range t.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("# BilaKey PC macro table v1\n")
	b.WriteString("# trigger<TAB>replacement\n")
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('\t')
		b.WriteString(strings.ReplaceAll(t.entries[k], "\n", " "))
		b.WriteByte('\n')
	}
	tmp := t.path + ".tmp"
	if err := os.WriteFile(tmp, []byte(b.String()), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, t.path)
}

func (t *Table) Lookup(trigger string) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	v, ok := t.entries[trigger]
	if ok {
		return v, true
	}
	v, ok = t.entries[strings.ToLower(trigger)]
	return v, ok
}

func (t *Table) Upsert(trigger, replacement string) {
	trigger = strings.TrimSpace(trigger)
	if trigger == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[trigger] = replacement
}

func (t *Table) Delete(trigger string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, trigger)
}

func (t *Table) Entries() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for k, v := range t.entries {
		out = append(out, Entry{Trigger: k, Replacement: v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Trigger < out[j].Trigger })
	return out
}
