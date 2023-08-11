package keepassstore

import (
	"os"
	"path/filepath"

	"github.com/revengel/enpass2gopass/utils"
	"github.com/sirupsen/logrus"
	"github.com/tobischo/gokeepasslib/v3"
)

// Store -
type Store struct {
	db     *gokeepasslib.Database
	prefix string
	items  *utils.UniqueStrings
	dryrun bool
	logger *logrus.Logger
}

// Close -
func (st *Store) Close() error {
	return nil
}

// Cleanup -
func (st *Store) Cleanup() (bool, error) {
	return false, nil
}

// Save -
func (st *Store) Save(s Secret, p string) (bool, error) {
	return false, nil
}

// NewStore -
func NewStore(dbPath, password, prefix string, dryrun bool, logger *logrus.Logger) (store *Store, err error) {
	absDbPath, err := filepath.Abs(dbPath)
	if err != nil {
		return
	}

	file, err := os.Open(absDbPath)
	if err != nil {
		return
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)
	err = gokeepasslib.NewDecoder(file).Decode(db)
	if err != nil {
		return
	}

	if prefix == "" {
		prefix = "enpass"
	}

	return &Store{
		db:     db,
		prefix: prefix,
		items:  utils.NewUniqueStrings(logger),
		dryrun: dryrun,
		logger: logger,
	}, nil
}
