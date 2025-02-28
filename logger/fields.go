package logger

import "strings"

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

// SafePhoneZH 安全地返回中国手机号，eg: 135****1234
func SafePhoneZH(phone string) Field {
	if DEBUG {
		return Field{Key: "phone_zh", Val: phone}
	} else {
		return Field{Key: "phone_zh", Val: phone[:3] + "****" + phone[len(phone)-4:]}
	}
}

// SafeEmail 安全地返回邮箱，eg: ***@example.com
func SafeEmail(email string) Field {
	if DEBUG {
		return Field{Key: "email", Val: email}
	} else {
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			panic("Invalid email")
		}
		return Field{Key: "email", Val: "@" + parts[1]}
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
