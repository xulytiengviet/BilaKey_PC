package core

import "strings"

const (
	MethodCVNSS = "CVNSS4.0"
	MethodTelex = "Telex"
	MethodVNI   = "VNI"
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

func New(method string, opts Options) *Engine {
	return &Engine{Method: method, Options: opts}
}

func (e *Engine) Transform(raw string) string {
	var out string
	switch strings.ToUpper(e.Method) {
	case strings.ToUpper(MethodTelex):
		out = transformTelex(raw, e.Options.OldToneStyle)
	case strings.ToUpper(MethodVNI):
		out = transformVNI(raw, e.Options.OldToneStyle)
	default:
		out = DecodeCVNSS(raw)
	}
	if e.Options.SpellCheck && e.Options.AutoRestoreWrongKey && out != raw && !IsLikelyVietnameseSyllable(out) {
		return raw
	}
	return out
}
