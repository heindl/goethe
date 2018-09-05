# {{.ModuleName}}

### {{.RootCommand.Short}}

{{.RootCommand.Long}}

### Install
{{ if .GoDownloader && .Domain eq "github.com" }}
```bash
bash <(wget -s https://raw.githubusercontent.com/{{.SubDomain}}/{{.ModuleName}}/master/godownloader.sh)`
```
{{ else if .GoReleaser && .Domain eq "github.com" }}
Download the latest release: https://github.com/{{.SubDomain}}/{{.ModuleName}}/releases
{{ end }}

buf := new(bytes.Buffer)
	buf.WriteString("# " + rootCmd.Name() + "\n\n")

	buf.WriteString("### " + rootCmd.Short + "\n\n")

	buf.WriteString("#### " + rootCmd.Long + "\n\n")

	buf.WriteString("### Installation \n")
	if options.UsesGoreleaser && strings.HasPrefix(options.PackageName, "github/") {

		buf.WriteString(
			fmt.Sprintf(`
				bash <(wget -s https://raw.githubusercontent.com/%s/master/godownloader.sh)`, strings.TrimPrefix(options.PackageName, "github/")),
		)
	}
