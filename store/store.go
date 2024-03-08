package store

import "github.com/revengel/enpass2gopass/field"

// StoreDestination -
type StoreDestination interface {
	Close() error
	Cleanup() (bool, error)
	Save(fields []field.FieldInterface, p string) (bool, error)
}

// StoreSource -
type StoreSource interface {
	LoadData() (err error)
}
