// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package utilities

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/phogolabs/parcello"
	"github.com/pkg/errors"
)

// Parcello saves everything to the top level directory, which should be fixed,
// but is convenient for this utility function.

func ReadTemplate(fileName string) (*template.Template, error) {
	file, err := parcello.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open template %s", fileName)
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read template %s", fileName)
	}
	t, err := template.New(fileName).Parse(string(fileBytes))
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse template %s", fileName)
	}
	return t, nil
}

func ExecuteTemplate(fileName string, data interface{}) ([]byte, error) {
	t, err := ReadTemplate(fileName)
	if err != nil {
		return nil, err
	}
	b := bytes.Buffer{}
	if err := t.Execute(&b, data); err != nil {
		return nil, errors.Wrapf(err, "could not execute template %s", fileName)
	}
	return b.Bytes(), nil
}
