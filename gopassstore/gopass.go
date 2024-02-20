package gopassstore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/gopasspw/gopass/pkg/gopass/secrets"
	"github.com/revengel/enpass2gopass/field"
	"github.com/revengel/enpass2gopass/utils"
	"github.com/sirupsen/logrus"
)

var (
	// ErrNotFound is returned if an entry was not found.
	ErrNotFound = fmt.Errorf("entry is not in the password store")
)

// Gopass -
type Gopass struct {
	ctx            context.Context
	api            *api.Gopass
	prefix         string
	uniqueKeys     *utils.UniqueStrings
	uniquePrefixes *utils.UniqueStrings
	dryrun         bool
	logger         *logrus.Logger
}

// Get -
func (g Gopass) get(p string) (o gopass.Secret, err error) {
	return g.api.Get(g.ctx, p, "latest")
}

// Set -
func (g Gopass) set(s gopass.Byter, p string) (err error) {
	return g.api.Set(g.ctx, p, s)
}

func (g Gopass) list(keyRe string) (keys []string, err error) {
	keys, err = g.api.List(g.ctx)
	if err != nil {
		return
	}

	if keyRe == "" {
		return
	}

	gopassKeyRe := regexp.MustCompile(keyRe)
	var filteredKeys []string
	for _, k := range keys {
		if !gopassKeyRe.MatchString(k) {
			continue
		}
		filteredKeys = append(filteredKeys, k)
	}

	return filteredKeys, nil
}

// Remove -
func (g Gopass) remove(p string) error {
	return g.api.Remove(g.ctx, p)
}

// Close -
func (g *Gopass) Close() error {
	return g.api.Close(g.ctx)
}

// Diff -
func (g Gopass) diff(a, b gopass.Byter) bool {
	ahash := utils.GetHashFromBytes(a.Bytes())
	bhash := utils.GetHashFromBytes(b.Bytes())
	return ahash == bhash
}

func (g Gopass) saveSecret(s gopass.Byter, p string) (bool, error) {
	p = g.uniqueKeys.Unique(p)
	var l = g.logger.WithField("gopasskey", p)

	rSec, err := g.get(p)
	if err != nil && err.Error() != ErrNotFound.Error() {
		return false, err
	}

	if rSec != nil && g.diff(s, rSec) {
		l.Debug("gopass secret already in actual state")
		return false, nil
	}

	l.Info("secret will be updated")
	if g.logger.IsLevelEnabled(logrus.DebugLevel) {
		fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
		if rSec != nil {
			_, err = io.Copy(os.Stdout, bytes.NewReader(rSec.Bytes()))
			if err != nil {
				return false, err
			}
		}
		fmt.Println("=============")
		_, err = io.Copy(os.Stdout, bytes.NewReader(s.Bytes()))
		if err != nil {
			return false, err
		}
	}

	if g.dryrun {
		return true, nil
	}

	err = g.set(s, p)
	if err != nil {
		return false, err
	}

	l.Info("secret has been updated")
	return true, nil
}

// Cleanup -
func (g Gopass) Cleanup() (bool, error) {
	var deletesCount = 0
	ll, err := g.list(`^` + g.prefix + `/`)
	if err != nil {
		return false, err
	}

	for _, k := range ll {
		if g.uniqueKeys.Has(k) {
			continue
		}

		var lc = g.logger.WithField("type", "cleaner").
			WithField("gopasskey", k)

		lc.Info("gopass key will be deleted")

		if g.dryrun {
			continue
		}

		err = g.remove(k)
		if err != nil {
			return false, err
		}

		deletesCount++
	}

	return deletesCount > 0, nil
}

func (g Gopass) getMainSecretPath(p string) string {
	return filepath.Join(g.prefix, p, "data")
}

func (g Gopass) getAttachmentSecretPath(p, attachmentName string) string {
	return filepath.Join(g.prefix, p, "attachments", attachmentName)
}

// Save -
func (g Gopass) Save(fields []field.FieldInterface, p string) (bool, error) {
	var err error
	var out bool
	p = g.uniquePrefixes.Unique(p)
	var keyPath = g.getMainSecretPath(p)

	// create gopass secrets
	var mainSecret = secrets.NewAKV()
	var attachments = make(map[string]*secrets.AKV)
	var multilineFields []field.FieldInterface
	for _, f := range fields {
		switch f.GetType() {
		case field.SecretAttachmentField:
			// create separate secrets for attachments
			var secret = secrets.NewAKV()
			err = secret.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", f.GetKey()))
			if err != nil {
				return false, err
			}

			err = secret.Set("Content-Transfer-Encoding", "Base64")
			if err != nil {
				return false, err
			}

			_, err = secret.Write(f.GetValue())
			if err != nil {
				return false, err
			}

			attachments[f.GetKey()] = secret
		case field.SecretPasswordField:
			// SetPassword -
			if mainSecret.Password() == "" {
				mainSecret.SetPassword(f.GetValueString())
				continue
			}
			fallthrough
		default:
			if f.IsMultiline() {
				multilineFields = append(multilineFields, f)
				continue
			}
			err = mainSecret.Set(f.GetKey(), f.GetValueString())
			if err != nil {
				return false, nil
			}
		}
	}

	// writing multiline fields in end of decret
	if len(multilineFields) > 0 {
		var data = "---\n"
		for _, f := range multilineFields {
			data += fmt.Sprintf("%s\n\n%s\n", f.GetKey(), f.GetValueString())
		}

		_, err = mainSecret.Write([]byte(data))
		if err != nil {
			return false, nil
		}
	}

	same, err := g.saveSecret(mainSecret, keyPath)
	if err != nil {
		return out, err
	}

	out = out || same

	for attachName, secret := range attachments {
		var keyPath = g.getAttachmentSecretPath(p, attachName)
		same, err := g.saveSecret(secret, keyPath)
		if err != nil {
			return out, err
		}

		out = out || same
	}

	return out, nil
}

// NewStore -
func NewStore(ctx context.Context, prefix string, dryrun bool, logger *logrus.Logger) (g *Gopass, err error) {
	var gp *api.Gopass
	gp, err = api.New(ctx)
	if err != nil {
		return g, fmt.Errorf("failed to initialize gopass API: %s", err.Error())
	}

	if prefix == "" {
		prefix = "enpass"
	}

	return &Gopass{
		ctx:            ctx,
		api:            gp,
		prefix:         prefix,
		uniqueKeys:     utils.NewUniqueStrings(logger),
		uniquePrefixes: utils.NewUniqueStrings(logger),
		dryrun:         dryrun,
		logger:         logger,
	}, nil
}
