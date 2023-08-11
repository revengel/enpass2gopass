package store

// Store -
type Store interface {
	Close() error
	Cleanup() (bool, error)
	Save(s Secret, p string) (bool, error)
}
