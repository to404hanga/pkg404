package logger

import "testing"

func TestInfo(t *testing.T) {
	l := NewFileZapLogger("dev", "")
	l.Info("aaa")
}

func TestError(t *testing.T) {
	l := NewFileZapLogger("dev", "")
	l.Error("err")
}
