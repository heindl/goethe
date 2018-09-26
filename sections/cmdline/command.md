
## Command-line
{{ if .HasExample -}}
#### Example
```bash
{{- .Example -}}
```
{{- end}}

{{ if .HasSubCommands -}}

{{- range .SubCommands -}}
{{ if .Short }}{{.Short}}{{else}}{{.Name}}{{end}}
```bash
{{ .UseLine }}
```
{{- if .LocalFlags }}
```bash
{{- range .LocalFlags }}
{{ . }}
{{- end }}
```
{{ end }}
{{ .Long }}
{{end}}

{{- else -}}

```bash
{{ .UseLine }}
```
{{ end -}}
{{if or .LocalFlags .PersistentFlags -}}
{{if .SubCommands -}}
###### Global Flags
{{end -}}
```bash
{{if not .SubCommands -}}
{{- range .LocalFlags }}
{{- . }}
{{ end -}}
{{- end}}
{{- range .PersistentFlags }}
{{- . }}
{{ end -}}
```
{{- end }}
