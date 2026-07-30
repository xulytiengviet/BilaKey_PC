package core

import (
	"strings"
	"unicode"
)

// TransformText applies the selected engine incrementally to letter tokens and
// preserves punctuation, whitespace, URLs, code identifiers and foreign words
// when SpellCheck + AutoRestoreWrongKey are enabled.
func (e *Engine) TransformText(input string) string {
	if input == "" {
		return ""
	}
	var out strings.Builder
	var token strings.Builder
	flush := func() {
		if token.Len() == 0 {
			return
		}
		out.WriteString(e.Transform(token.String()))
		token.Reset()
	}
	for _, r := range input {
		if unicode.IsLetter(r) || (strings.EqualFold(e.Method, MethodVNI) && unicode.IsDigit(r)) {
			token.WriteRune(r)
			continue
		}
		flush()
		out.WriteRune(r)
	}
	flush()
	return out.String()
}
