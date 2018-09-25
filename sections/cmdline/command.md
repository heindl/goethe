
## Command-line
{{ if .HasExample -}}
#### Example
```bash
{{- .Example -}}
```
{{- end}}

{{ if .HasSubCommands -}}

{{- range .SubCommands -}}
### {{ .Name }}
{{ .Short }}
```bash
{{ .UseLine }}
{{- range .LocalFlags }}
{{- . }}
{{ end -}}
```
{{ .Long }}
{{end}}

{{- else -}}

```bash
{{ .UseLine }}
```
{{- end -}}
{{- if or .LocalFlags .PersistentFlags }}
```bash
{{if not .SubCommands -}}
{{- range .LocalFlags }}
{{- . }}
{{ end -}}
{{- end -}}
{{- range .PersistentFlags }}
{{- . }}
{{ end -}}
```
{{- end }}
