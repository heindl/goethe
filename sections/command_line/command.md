
## Command-line
{{ if .HasExample -}}
```bash
{{ .Example }}
```
{{- end -}}
{{- if .HasSubCommands -}}

{{- range .SubCommands -}}
### {{ .Name }}
{{ .Short }}
```bash
{{ .UseLine }}
{{ .LocalFlags }}
```
{{ .Long }}
{{- end}}
{{- else -}}
```bash
{{ .UseLine }}
```
{{- end -}}
{{- if or .LocalFlags .PersistentFlags }}
###### Flags
```bash
{{if not .SubCommands -}}
{{- range .LocalFlags -}}
{{ . }}
{{- end -}}
{{- end -}}
{{- range .PersistentFlags -}}
{{ . }}
{{- end }}
```
{{- end }}
