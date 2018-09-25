// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package render

import (
	"fmt"
	"io"

	"github.com/heindl/goethe/sections/cmdline"
	"github.com/heindl/goethe/sections/header"
	"github.com/heindl/goethe/sections/install"
	"github.com/heindl/goethe/sections/license"
	"github.com/heindl/goethe/utilities"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Section int

const (
	Header        Section = 0
	Install       Section = 1
	Usage         Section = 2
	License       Section = 3
	totalSections         = 4
)

func (Ω Section) Render(info *utilities.ModuleInfo) ([]byte, error) {
	switch Ω {
	case Header:
		return header.Render(info)
	case Install:
		return install.Render(info)
	case Usage:
		return cmdline.Render(info)
	case License:
		return license.Render(info)
	default:
		return nil, nil
	}
}
func Render(filePath string, writer io.Writer) error {

	modInfo, err := utilities.ReadModule(filePath)
	if err != nil {
		return err
	}

	rendered := make([][]byte, totalSections)

	eg := errgroup.Group{}

	for _i := 0; _i < totalSections; _i++ {
		i := _i
		eg.Go(func() error {
			var err error

			rendered[i], err = Section(i).Render(modInfo)
			if err != nil {
				fmt.Println("err", i, err)
			}
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	for i, b := range rendered {
		if _, err := writer.Write(b); err != nil {
			return errors.Wrapf(err, "could not write %d", i)
		}
	}

	return nil
}
