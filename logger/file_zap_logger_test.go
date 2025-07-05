package logger

import "testing"

func TestInfo(t *testing.T) {
	l := NewFileZapLogger("dev", "", false)
	l.Info("aaa")
}

func TestError(t *testing.T) {
	l := NewFileZapLogger("dev", "", true)
	l.Error("err")
}

func TestWith(t *testing.T) {
	l := NewFileZapLogger("dev", "", false)
	l.Info("without any fields", String("key", "value"))
	nl := l.With(Int64("int", 64), String("k2", "v2"))
	nl.Info("with fields", String("key", "value"))
}
