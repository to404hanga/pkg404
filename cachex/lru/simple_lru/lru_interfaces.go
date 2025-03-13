package simple_lru

type LRUCache interface {
	Add(key, value any) bool
	Get(key any) (value any, ok bool)
	Contains(key any) (ok bool)
	Peek(key any) (value any, ok bool)
	Remove(key any) bool
	RemoveOldest() (key, value any, ok bool)
	GetOldest() (key, value any, ok bool)
	Keys() []any
	Len() int
}
