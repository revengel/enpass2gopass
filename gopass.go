package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/gopasspw/gopass/pkg/gopass/secrets"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNotFound is returned if an entry was not found.
	ErrNotFound = fmt.Errorf("entry is not in the password store")
)

// Gopass -
type Gopass struct {
	ctx context.Context
	api *api.Gopass
}

func (g Gopass) get(p string) (o gopass.Secret, err error) {
	o, err = g.api.Get(g.ctx, p, "latest")
	if err != nil {
		return
	}
	return
}

func (g Gopass) set(s gopass.Byter, p string) (err error) {
	err = g.api.Set(g.ctx, p, s)
	if err != nil {
		return
	}
	return
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

func (g Gopass) remove(p string) (err error) {
	err = g.api.Remove(g.ctx, p)
	if err != nil {
		return
	}
	return
}

// Close -
func (g *Gopass) Close() error {
	return g.api.Close(g.ctx)
}

func (g Gopass) diff(a, b gopass.Byter) bool {
	ahash := getHashFromBytes(a.Bytes())
	bhash := getHashFromBytes(b.Bytes())
	return ahash == bhash
}

func (g Gopass) diffWithStorage(s gopass.Byter, p string) (bool, error) {
	rSec, err := g.get(p)
	if err != nil {
		// TODO: need be refactoring
		if err.Error() == ErrNotFound.Error() {
			return false, nil
		}
		return false, err
	}

	return g.diff(s, rSec), nil
}

func (g Gopass) setIfChanged(s gopass.Byter, p string) (bool, error) {
	same, err := g.diffWithStorage(s, p)
	if err != nil {
		return false, err
	}

	if same {
		return false, nil
	}

	err = g.set(s, p)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g Gopass) sync() error {
	var err = g.api.Sync(g.ctx)
	return err
}

func newGopass(ctx context.Context) (g *Gopass, err error) {
	var gp *api.Gopass
	gp, err = api.New(ctx)
	if err != nil {
		return g, fmt.Errorf("failed to initialize gopass API: %s", err.Error())
	}

	return &Gopass{
		ctx: ctx,
		api: gp,
	}, nil
}

func gopassSecretSet(s *secrets.AKV, k, v string, multiline bool) (err error) {
	if k == "" || v == "" {
		return
	}

	if multiline {
		data := []byte(fmt.Sprintf("%s\n%s\n", k, v))
		_, err = s.Write(data)
		return err
	}

	return s.Set(k, v)
}

func gopassSecretAddYamlData(data, k, v string) string {
	if k == "" || v == "" {
		return data
	}

	if data == "" {
		data += "---\n"
	}

	data += fmt.Sprintf("%s\n\n%s\n", k, v)
	return data
}

func getGopassItemSecrets(prefix string, item DataItem) (map[string]gopass.Byter, error) {
	var (
		err    error
		data   string
		o      = make(map[string]gopass.Byter)
		s      = secrets.NewAKV()
		folder = item.GetFirstFolder()
	)

	gopassPath, err := getGopassPath(prefix, folder, item)
	if err != nil {
		return o, err
	}

	gopassDataPath, err := getGopassDataPath(gopassPath)
	if err != nil {
		return o, err
	}

	err = gopassSecretSet(s, "title", item.GetTitle(), false)
	if err != nil {
		return o, err
	}

	err = gopassSecretSet(s, "subtitle", item.GetSubtitle(), false)
	if err != nil {
		return o, err
	}

	err = gopassSecretSet(s, "category", item.GetCategoryPath(), false)
	if err != nil {
		return o, err
	}

	data = gopassSecretAddYamlData(data, "note", item.GetNote())

	err = gopassSecretSet(s, "tags", item.GetFoldersStr(), false)
	if err != nil {
		return o, err
	}

	for _, field := range item.Fields {
		var ignoreTypes = []string{"section", ".Android#"}
		if field.IsDeleted() || field.CheckTypes(ignoreTypes) {
			continue
		}

		if field.GetLabel() == "" || field.GetValue() == "" {
			continue
		}

		if s.Password() == "" && field.CheckType("password") {
			passwd := field.GetValue()
			s.SetPassword(passwd)
			continue
		}

		var labelName = field.GetLabel()
		if field.CheckType("totp") {
			labelName = "totp"
		}

		if labelName == "e_mail" {
			labelName = "email"
		}

		if field.IsMultiline() {
			data = gopassSecretAddYamlData(data, labelName, field.GetValue())
			continue
		}

		err = gopassSecretSet(s, labelName, field.GetValue(), field.IsMultiline())
		if err != nil {
			return o, err
		}
	}

	for _, attach := range item.Attachments {
		var name = attach.GetLabelName()
		isText, err := attach.IsTextData()
		if err != nil {
			return o, err
		}

		if isText {
			labelName := fmt.Sprintf("attachment - %s", attach.GetName())
			val, err := attach.GetDataString()
			if err != nil {
				return o, err
			}
			data = gopassSecretAddYamlData(data, labelName, val)
			continue
		}

		gopassAttachPath, err := getGopassAttachPath(gopassPath, name)
		if err != nil {
			return o, err
		}

		attachSec, err := getGopassItemAttachSecret(attach)
		if err != nil {
			return o, err
		}

		o[gopassAttachPath] = attachSec
	}

	_, err = s.Write([]byte(data))
	if err != nil {
		return o, err
	}

	o[gopassDataPath] = s
	return o, err
}

func getGopassItemAttachSecret(attach Attachment) (o gopass.Byter, err error) {
	var dataB64 = attach.GetDataBase64Encoded()
	var s = secrets.NewAKV()

	err = gopassSecretSet(s, "Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", attach.GetNameOriginal()), false)
	if err != nil {
		return s, err
	}

	err = gopassSecretSet(s, "Content-Transfer-Encoding", "Base64", false)
	if err != nil {
		return s, err
	}

	_, err = s.Write([]byte(dataB64))
	if err != nil {
		return
	}

	return s, nil
}

func gopassSaveSecret(s gopass.Byter, gp *Gopass, p string, dryrun bool, l *log.Entry) error {
	if dryrun {
		return nil
	}

	saved, err := gp.setIfChanged(s, p)
	if err != nil {
		return err
	}

	if saved {
		l.Info("secret has been updated")
	} else {
		l.Debug("secret already in actual state")
	}
	return nil
}
