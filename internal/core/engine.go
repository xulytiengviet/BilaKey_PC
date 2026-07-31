package core

import "strings"

const (
	MethodCVNSS    = "CVNSS4.0"
	MethodVNITelex = "VNI/Telex"
	MethodTelex    = "Telex" // legacy configuration alias
	MethodVNI      = "VNI"   // legacy configuration alias
)

type Options struct {
	OldToneStyle        bool
	FreeToneMarking     bool
	SpellCheck          bool
	AutoRestoreWrongKey bool
}

type Engine struct {
	Method  string
	Options Options
}

// NormalizeMethod migrates legacy Telex/VNI selections to the unified mode.
func NormalizeMethod(method string) string {
	key := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(method), " ", ""))
	switch key {
	case "VNI/TELEX", "TELEX/VNI", "VNI", "TELEX", "VNITELEX", "TELEXVNI":
		return MethodVNITelex
	default:
		return MethodCVNSS
	}
}

func New(method string, opts Options) *Engine {
	return &Engine{Method: NormalizeMethod(method), Options: opts}
}

func (e *Engine) Transform(raw string) string {
	var out string
	switch NormalizeMethod(e.Method) {
	case MethodVNITelex:
		out = transformVNITelex(raw, e.Options.OldToneStyle)
	default:
		out = DecodeCVNSS(raw)
	}
	if e.Options.SpellCheck && e.Options.AutoRestoreWrongKey && out != raw && !IsLikelyVietnameseSyllable(out) {
		return raw
	}
	return out
}
