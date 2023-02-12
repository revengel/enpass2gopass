package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// Attachment -
type Attachment struct {
	Data string `json:"data"`
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// GetDataBase64Encoded -
func (a Attachment) GetDataBase64Encoded() string {
	return a.Data
}

// GetDataBytes -
func (a Attachment) GetDataBytes() (o []byte, err error) {
	var datab64 = a.GetDataBase64Encoded()
	var dataB = []byte(datab64)
	var sizeDecoded = base64.StdEncoding.EncodedLen(len(datab64))
	o = make([]byte, sizeDecoded)
	_, err = base64.StdEncoding.Decode(o, dataB)
	if err != nil {
		return
	}
	return
}

// GetDataString -
func (a Attachment) GetDataString() (o string, err error) {
	b, err := a.GetDataBytes()
	if err != nil {
		return
	}
	return string(b), nil
}

// GetDataContentType -
func (a Attachment) GetDataContentType() (o string, err error) {
	var dataB []byte
	dataB, err = a.GetDataBytes()
	if err != nil {
		return
	}

	o = http.DetectContentType(dataB)
	o = strings.SplitN(o, ";", 2)[0]
	return
}

// IsTextData -
func (a Attachment) IsTextData() (bool, error) {
	if a.GetKind() == "text/plain" {
		return true, nil
	}

	ct, err := a.GetDataContentType()
	if err != nil {
		return false, err
	}

	if ct == "text/plain" {
		return true, nil
	}

	return false, nil
}

// GetKind -
func (a Attachment) GetKind() string {
	return a.Kind
}

// GetName -
func (a Attachment) GetName() string {
	return a.Name
}

// GetLabelName -
func (a Attachment) GetLabelName() string {
	return transliterate(a.Name)
}

// GetNameOriginal -
func (a Attachment) GetNameOriginal() string {
	return a.Name
}
