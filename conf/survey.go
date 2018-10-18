package conf

import (
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/core"
)

var customSelectTemplate = `{{- if .ShowHelp }}{{- color "cyan"}}» {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green"}}» {{ color "default"}}{{ .Message }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else}}
  {{- if and .Help (not .ShowHelp)}} {{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}}{{end}}
  {{- "\n"}}
  {{- range $ix, $choice := .PageEntries}}
    {{- if eq $ix $.SelectedIndex}}{{color "cyan"}}{{ SelectFocusIcon }} {{else}}{{color "default"}}  {{end}}
    {{- $choice}}
    {{- color "reset"}}{{"\n"}}
  {{- end}}
{{- end}}`

var customInputTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}» {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green"}}» {{ color "default"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
{{- end}}`

var customErrorTemplate = `{{color "red"}}» Invalid input: {{.Error}}{{color "reset"}}
`

var customConfirmQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}» {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green"}}» {{ color "default"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
  {{- color "white"}}{{if .Default}}(Y/n) {{else}}(y/N) {{end}}{{color "reset"}}
{{- end}}`

func init() {
	survey.InputQuestionTemplate = customInputTemplate
	survey.SelectQuestionTemplate = customSelectTemplate
	survey.ConfirmQuestionTemplate = customConfirmQuestionTemplate
	core.ErrorTemplate = customErrorTemplate
}
