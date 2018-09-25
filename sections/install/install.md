
## Install
{{ if .GoDownloaderLink }}
Download the [latest release](https://{{.ModuleRemotePath}}/releases):
```bash
bash <(wget -s {{.GoDownloaderLink}})`
```
Or build it with Go:
{{ end }}
```bash
go get -u {{.ModuleRemotePath}}
```
