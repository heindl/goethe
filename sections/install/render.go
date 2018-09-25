// Copyright 2018 Parker Heindl. All rights reserved.
// Licensed under the MIT License. See LICENSE.md in the
// project root for information.
//
package install

import (
	"fmt"

	"github.com/heindl/goethe/utilities"
)

//go:generate parcello -m install.md
func Render(modInfo *utilities.ModuleInfo) ([]byte, error) {

	goDownloader, err := utilities.DirectoryContainsFile(modInfo.FilePath(), ".godownloader.yml")
	if err != nil {
		return nil, err
	}

	var data = struct {
		GoDownloaderLink string
		ModuleRemotePath string
	}{}

	if goDownloader {
		data.GoDownloaderLink = fmt.Sprintf(
			"https://raw.githubusercontent.com/%s/%s/master/godownloader.sh",
			modInfo.SubDomain(),
			modInfo.Name())
	}
	data.ModuleRemotePath = modInfo.RemotePath()

	return utilities.ExecuteTemplate("install.md", data)

}
