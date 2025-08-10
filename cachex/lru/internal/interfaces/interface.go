package interfaces

type LRUCache interface {
	Add(key, value any) bool
	Get(key any) (value any, ok bool)
	Purge()
	Resize(opts ...*SizeOptions) (evicted int)
	Contains(key any) (ok bool)
	Peek(key any) (value any, ok bool)
	Remove(key any) bool
	RemoveOldest() (key, value any, ok bool)
	GetOldest() (key, value any, ok bool)
	Keys() []any
	Len() int
}

type SizeOptions struct {
	Key   string
	Value int
}

func WithSize(size int) *SizeOptions {
	return &SizeOptions{
		Key:   "size",
		Value: size,
	}
}

func WithYoungListSize(size int) *SizeOptions {
	return &SizeOptions{
		Key:   "youngListSize",
		Value: size,
	}
}
