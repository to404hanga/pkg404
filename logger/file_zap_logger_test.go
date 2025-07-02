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
