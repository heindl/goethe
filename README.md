
[//]: # (Code generated by goethe DO NOT EDIT)
# goethe

*v0.0.1* | *September 25, 2018*
#### Statically generate a github flavored README.md from a Go module.

Goethe parses your module directory for license, installation and command data, and generates a Github README.md file.

It is **very unstable** right now.

Go1.11 modules are required, and this version only scans for a [Cobra Command](https://github.com/spf13/cobra) for documentation.

Future versions will support godoc, other command line helpers, and distribution and deployment tools. It will hopefully save time for both the reader and writer of open source code.

## Install

```bash
go get -u github.com/heindl/goethe
```

## Command-line
```bash
goethe [command_directory]
```
###### Flags
```bash
-p, --print   Print the template data to standard out.
```

## License
Copyright 2018 Parker Heindl. All rights reserved.
Use of this source code is governed by the [MIT License](LICENSE.md).
