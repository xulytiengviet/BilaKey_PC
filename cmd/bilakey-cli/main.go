package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xulytiengviet/BilaKey_PC/internal/core"
)

func main() {
	method := flag.String("method", "cvnss", "cvnss, telex hoặc vni")
	inspect := flag.Bool("inspect", false, "in chi tiết candidate graph CVNSS dạng JSON")
	audit := flag.Bool("audit", false, "in audit lõi CVNSS dạng JSON")
	text := flag.Bool("text", false, "chuyển toàn bộ văn bản, giữ nguyên nội dung hỗn hợp")
	flag.Parse()

	if *audit {
		writeJSON(core.AuditCVNSS())
		return
	}
	input := strings.Join(flag.Args(), " ")
	if input == "" {
		fmt.Fprintln(os.Stderr, "usage: bilakey-cli [-method cvnss|telex|vni] [-text|-inspect|-audit] <input>")
		os.Exit(2)
	}
	if *inspect {
		writeJSON(core.InspectCVNSS(input))
		return
	}
	selected := core.MethodCVNSS
	switch strings.ToLower(*method) {
	case "telex":
		selected = core.MethodTelex
	case "vni":
		selected = core.MethodVNI
	case "cvnss", "cvnss4.0":
	default:
		fmt.Fprintf(os.Stderr, "method không hợp lệ: %s\n", *method)
		os.Exit(2)
	}
	engine := core.New(selected, core.Options{SpellCheck: true, AutoRestoreWrongKey: true})
	if *text {
		fmt.Println(engine.TransformText(input))
		return
	}
	fmt.Println(engine.Transform(input))
}

func writeJSON(value any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
