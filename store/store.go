package store

import "github.com/revengel/enpass2gopass/field"

// Store -
type Store interface {
	Close() error
	Cleanup() (bool, error)
	Save(fields []field.FieldInterface, p string) (bool, error)
}
