Goethe parses your module directory for license, installation and command data, and generates a README.md file.

It is **very unstable** right now.

Go1.11 modules are required, and this version only scans for a [Cobra Command](https://github.com/spf13/cobra) for documentation.

Future versions will support godoc, other command line helpers, and distribution and deployment tools. It will hopefully save time for both the reader and writer of open source code.