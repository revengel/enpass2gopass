package store

// Store -
type Store interface {
	SetIfChanged(s Secret, p string) (bool, error)
	Close() error
	Remove(p string) error
	List(keyRe string) (keys []string, err error)
}
