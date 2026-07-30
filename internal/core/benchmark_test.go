package core

import "testing"

func BenchmarkEngineCVNSS(b *testing.B) {
	e := New(MethodCVNSS, Options{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = e.Transform("vidf")
	}
}

func BenchmarkEngineTelex(b *testing.B) {
	e := New(MethodTelex, Options{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = e.Transform("tieengs")
	}
}

func BenchmarkEngineVNI(b *testing.B) {
	e := New(MethodVNI, Options{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = e.Transform("tieng61")
	}
}
