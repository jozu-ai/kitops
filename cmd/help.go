package cmd

import (
	"fmt"
	"strings"
)

const helpTemplate = `
{{- if .Short -}}{{.Short}}{{- end -}}
{{- if .Long -}}
	{{sectionHead "Description:"}}
	{{- indent .Long 2 | ensureTrailingNewline}}
{{- end -}}

{{ if or .Runnable .HasSubCommands -}}
{{.UsageString}}
{{- end -}}
`

const usageTemplate = `
Usage:
{{ if .Runnable -}}{{indent .UseLine 2}}{{- end -}}

{{if .HasAvailableSubCommands -}}
	{{indent .CommandPath 2}} [command]
{{- end -}}

{{if gt (len .Aliases) 0 -}}
	{{sectionHead "Aliases:" -}}
	{{indent .NameAndAliases 2 -}}
{{- end -}}

{{if .HasExample -}}
	{{sectionHead "Examples:" -}}
	{{indent .Example 2 -}}
{{- end -}}

{{if .HasAvailableSubCommands -}}
	{{$cmds := .Commands -}}
	{{if eq (len .Groups) 0 -}}
		{{"\n\nAvailable Commands:" -}}
		{{range $cmds -}}
			{{if (or .IsAvailableCommand (eq .Name "help")) -}}
				{{"\n" -}} {{indent (rpad .Name .NamePadding) 2}} {{.Short}}
			{{- end -}}
		{{- end -}}
	{{- else -}}
		{{range $group := .Groups -}}
			{{"\n\n" -}}
			{{.Title -}}
			{{range $cmds -}}
				{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help"))) -}}
					{{"\n" -}} {{indent (rpad .Name .NamePadding) 2}} {{.Short}}
				{{- end -}}
			{{- end -}}
		{{- end -}}
		{{if not .AllChildCommandsHaveGroup -}}
			{{"\n\nAdditional Commands:" -}}
			{{range $cmds -}}
				{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help"))) -}}
					{{"\n" -}} {{indent (rpad .Name .NamePadding) 2}} {{.Short}}
				{{- end -}}
			{{- end -}}
		{{- end -}}
	{{- end -}}
{{- end -}}

{{if .HasAvailableLocalFlags -}}
	{{sectionHead "Flags:" -}}
	{{.LocalFlags.FlagUsages | trimTrailingWhitespaces -}}
{{- end -}}

{{if .HasAvailableInheritedFlags -}}
	{{sectionHead "Global Flags:" -}}
	{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}
{{- end -}}

{{if .HasHelpSubCommands}}{{sectionHead "Additional help topics:"}}
	{{range .Commands}}
		{{if .IsAdditionalHelpTopicCommand}}
			{{rpad .CommandPath .CommandPathPadding}} {{.Short}}
		{{end}}
	{{end}}
{{- end -}}

{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.
{{- end}}
`

func indentBlock(s string, indent int) string {
	lines := strings.Split(s, "\n")
	var indented []string
	for _, line := range lines {
		indented = append(indented, fmt.Sprintf("%s%s", strings.Repeat(" ", indent), line))
	}
	return strings.Join(indented, "\n")
}

func sectionHead(s string) string {
	return fmt.Sprintf("\n\n%s\n", s)
}

func ensureTrailingNewline(s string) string {
	return strings.TrimRight(s, " \n") + "\n"
}
