package logger

var DEBUG = false

func Error(err error) Field {
	return Field{Key: "error", Val: err}
}

func SafeString(key, val string) Field {
	if DEBUG {
		return Field{Key: key, Val: val}
	} else {
		return Field{Key: key, Val: "***"}
	}
}

func Any[T any](key string, val T) Field {
	return Field{Key: key, Val: val}
}

func Slice[T any](key string, slice []T) Field {
	return Field{Key: key, Val: slice}
}

func String(key, val string) Field {
	return Field{Key: key, Val: val}
}

func Bool(key string, val bool) Field {
	return Field{Key: key, Val: val}
}

func Int(key string, val int) Field {
	return Field{Key: key, Val: val}
}

func Int32(key string, val int32) Field {
	return Field{Key: key, Val: val}
}

func Int64(key string, val int64) Field {
	return Field{Key: key, Val: val}
}
