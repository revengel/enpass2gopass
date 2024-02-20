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
	file   *os.File
}

// Close -
func (st *Store) Close() error {
	return nil
}

// Cleanup -
func (st *Store) Cleanup() (bool, error) {
	return false, nil
}

func (st *Store) createGroupRecursive(group *gokeepasslib.Group, path []string) *gokeepasslib.Group {
	for _, g := range group.Groups {
		if g.Name == path[0] {
			return st.createGroupRecursive(&g, path[1:])
		}
	}

	var sg = gokeepasslib.NewGroup()
	sg.Name = path[0]
	group.Groups = append(group.Groups, sg)
	return st.createGroupRecursive(&sg, path[1:])
}

// Save -
func (st *Store) Save(fields []field.FieldInterface, p string) (bool, error) {
	var err error
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
			_ = append(attachments, f)
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

	var secretPath = filepath.Join(st.prefix, p)
	// Split the group path into individual levels
	groupLevels := strings.Split(secretPath, "/")

	// Start from the root group
	var rootGroup *gokeepasslib.Group
	for _, gr := range st.db.Content.Root.Groups {
		if gr.Name == groupLevels[0] {
			rootGroup = &gr
			break
		}
	}

	if rootGroup == nil {
		gr := gokeepasslib.NewGroup()
		gr.Name = groupLevels[0]
		rootGroup = &gr
		st.db.Content.Root.Groups = append(st.db.Content.Root.Groups, *rootGroup)
	}

	// Create the groups with the specified path recursively
	lastGroup := st.createGroupRecursive(rootGroup, groupLevels[1:])
	lastGroup.Entries = append(lastGroup.Entries, mainSecret.Entry)

	defer st.db.LockProtectedEntries()

	encoder := gokeepasslib.NewEncoder(st.file)
	err = encoder.Encode(st.db)
	if err != nil {
		return false, err
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
