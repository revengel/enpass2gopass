package keepassstore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/revengel/enpass2gopass/field"
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

// Creates recursive groups under a given parent group
func (st *Store) createRecursiveGroups(parentGroup *gokeepasslib.Group, groups []string) *gokeepasslib.Group {
	if len(groups) <= 0 {
		return parentGroup
	}

	for _, g := range parentGroup.Groups {
		if g.Name == groups[0] {
			return st.createRecursiveGroups(&g, groups[1:])
		}
	}

	// Create a new group
	group := gokeepasslib.NewGroup()
	group.Name = groups[0]

	// Add the group to the parent group
	parentGroup.Groups = append(parentGroup.Groups, group)
	return st.createRecursiveGroups(&group, groups[1:])
}

// Save -
func (st *Store) Save(fields []field.FieldInterface, p string) (bool, error) {
	dir := filepath.Dir(p)
	groups := strings.Split(dir, "/")
	var root *gokeepasslib.Group

	var exists bool
	for i, g := range st.db.Content.Root.Groups {
		if g.Name == st.prefix {
			root = &st.db.Content.Root.Groups[i]
			exists = true
		}
	}

	if !exists {
		group := gokeepasslib.NewGroup()
		group.Name = st.prefix
	}

	var mainSecret = NewSecret()
	var attachments []field.FieldInterface
	for _, f := range fields {
		switch f.GetType() {
		case field.SecretTitleField:
			mainSecret.setKeyOrAlt("Title", f.GetKey(), f.GetValueString(), false)
		case field.SecretUsernameField:
			mainSecret.setKeyOrAlt("Username", f.GetKey(), f.GetValueString(), false)
		case field.SecretPasswordField:
			mainSecret.setKeyOrAlt("Password", f.GetKey(), f.GetValueString(), true)
		case field.SecretURLField:
			mainSecret.setKeyOrAlt("URL", f.GetKey(), f.GetValueString(), false)
		case field.SecretTagsField:
			mainSecret.setKeyOrAlt("Tags", f.GetKey(), f.GetValueString(), false)
		case field.SecretAttachmentField:
			attachments = append(attachments, f)
		default:
			if f.IsMultiline() {
				var notes = mainSecret.GetContent("Notes")
				notes += fmt.Sprintf("%s\n\n%s\n", f.GetKey(), f.GetValueString())
				mainSecret.setKey("Notes", notes, false)
				continue
			}
			mainSecret.setKey(f.GetKey(), f.GetValueString(), f.IsSensitive())
		}
	}

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
