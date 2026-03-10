package console

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"text/template"
	"unicode"

	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v3"
	"github.com/xrash/smetrics"

	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/env"
)

func init() {
	cli.HelpPrinterCustom = printHelpCustom
	cli.RootCommandHelpTemplate = appHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate
	cli.SubcommandHelpTemplate = commandHelpTemplate
	cli.VersionPrinter = printVersion
	huh.ErrUserAborted = cli.Exit(color.Red().Sprint("Cancelled"), 0)
}

const maxLineLength = 10000

// Template for the help message.
var (
	appHelpTemplate = `{{$v := offset .Usage 6}}{{wrap (colorize .Usage) 3}}{{if .Version}} {{green (wrap .Version $v)}}{{end}}

{{ yellow "Usage:" }}
   {{if .UsageText}}{{wrap (colorize .UsageText) 3}}{{end}}{{if .VisibleFlags}}

{{ yellow "Global options:" }}{{template "flagTemplate" (sortVisibleFlags .)}}{{end}}{{if .VisibleCommands}}

{{ yellow "Available commands:" }}{{template "commandTemplate" .}}{{end}}
`

	commandHelpTemplate = `{{ yellow "Description:" }}
   {{ (colorize .Usage) }}

{{ yellow "Usage:" }}
   {{template "usageTemplate" .}}{{with $root := .Root}}{{end}}{{if .VisibleFlags}}

{{ yellow "Options:" }}{{template "flagTemplate" (sortVisibleFlags .)}}{{end}}
`
	commandTemplate = `{{ $cv := offsetCommands .VisibleCommands 5}}{{range .VisibleCategories}}{{if .Name}}
 {{yellow .Name}}:{{end}}{{range (sortCommands .VisibleCommands)}}
  {{$s := join .Names ", "}}{{green $s}}{{ $sp := subtract $cv (offset $s 3) }}{{ indent $sp ""}}{{wrap (colorize .Usage) $cv}}{{end}}{{end}}`
	flagTemplate = `{{ $cv := offsetFlags . 5}}{{range  .}}
   {{$s := getFlagName .}}{{green $s}}{{ $sp := subtract $cv (offset $s 1) }}{{ indent $sp ""}}{{$us := (capitalize .Usage)}}{{wrap (colorize $us) $cv}}{{$df := getFlagDefaultText . }}{{if $df}} {{yellow $df}}{{end}}{{end}}`
	usageTemplate = `{{if .UsageText}}{{wrap (colorize .UsageText) 3}}{{else}}{{(helpName .FullName)}}{{if .VisibleFlags}} [options]{{end}}{{if .ArgsUsage}}{{.ArgsUsage}}{{else}}{{if .Arguments}}{{template "argsTemplate" .}}{{end}}{{end}}{{end}}`
	argsTemplate  = `{{if .Arguments}}{{range .Arguments}}{{template "argTemplate" .}}{{end}}{{end}}`
	argTemplate   = ` {{if .Min}}<{{else}}[{{end}}{{.Name}}{{if (or (gt .Max 1) (eq .Max -1))}}...{{end}}{{if .Min}}>{{else}}]{{end}}`
)

// colorsFuncMap is a map of functions for coloring text.
var colorsFuncMap = template.FuncMap{
	"black":   color.Black().Sprint,
	"blue":    color.Blue().Sprint,
	"cyan":    color.Cyan().Sprint,
	"default": color.Default().Sprint,
	"gray":    color.Gray().Sprint,
	"green":   color.Green().Sprint,
	"magenta": color.Magenta().Sprint,
	"red":     color.Red().Sprint,
	"white":   color.White().Sprint,
	"yellow":  color.Yellow().Sprint,
}

func capitalize(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// colorize wraps the text in the default color
// support style tags like <fg=red>text</>
// more details in https://gookit.github.io/color/#/?id=tag-attributes
func colorize(text string) string {
	return color.Default().Sprint(text)

}

func commandNotFound(_ context.Context, cmd *cli.Command, command string) {
	var (
		msgTxt     = fmt.Sprintf("Command '%s' is not defined.", command)
		suggestion string
	)
	if alternatives := findAlternatives(command, func() (collection []string) {
		for i := range cmd.Commands {
			collection = append(collection, cmd.Commands[i].Names()...)
		}
		return
	}()); len(alternatives) > 0 {
		if len(alternatives) == 1 {
			msgTxt = msgTxt + " Did you mean this?"
		} else {
			msgTxt = msgTxt + " Did you mean one of these?"
		}
		suggestion = "\n  " + strings.Join(alternatives, "\n  ")
	}
	color.Errorln(msgTxt)
	color.Gray().Println(suggestion)
}

func findAlternatives(name string, collection []string) (result []string) {
	var (
		threshold       = 1e3
		alternatives    = make(map[string]float64)
		collectionParts = make(map[string][]string)
	)
	for i := range collection {
		collectionParts[collection[i]] = strings.Split(collection[i], ":")
	}
	for i, sub := range strings.Split(name, ":") {
		for collectionName, parts := range collectionParts {
			exists := alternatives[collectionName] != 0
			if len(parts) <= i {
				if exists {
					alternatives[collectionName] += threshold
				}
				continue
			}
			lev := smetrics.WagnerFischer(sub, parts[i], 1, 1, 1)
			if float64(lev) <= float64(len(sub))/3 || strings.Contains(parts[i], sub) {
				if exists {
					alternatives[collectionName] += float64(lev)
				} else {
					alternatives[collectionName] = float64(lev)
				}
			} else if exists {
				alternatives[collectionName] += threshold
			}
		}
	}
	for _, item := range collection {
		lev := smetrics.WagnerFischer(name, item, 1, 1, 1)
		if float64(lev) <= float64(len(name))/3 || strings.Contains(item, name) {
			if alternatives[item] != 0 {
				alternatives[item] -= float64(lev)
			} else {
				alternatives[item] = float64(lev)
			}
		}
	}
	type scoredItem struct {
		name  string
		score float64
	}
	var sortedAlternatives []scoredItem
	for item, score := range alternatives {
		if score < 2*threshold {
			sortedAlternatives = append(sortedAlternatives, scoredItem{item, score})
		}
	}
	sort.Slice(sortedAlternatives, func(i, j int) bool {
		if sortedAlternatives[i].score == sortedAlternatives[j].score {
			return sortedAlternatives[i].name < sortedAlternatives[j].name
		}
		return sortedAlternatives[i].score < sortedAlternatives[j].score
	})
	for _, item := range sortedAlternatives {
		result = append(result, item.name)
	}
	return result
}

func getFlagDefaultText(flag cli.DocGenerationFlag) string {
	defaultValueString := ""
	if bf, ok := flag.(*cli.BoolFlag); !ok || bf.IsDefaultVisible() {
		if s := flag.GetDefaultText(); s != "" {
			defaultValueString = fmt.Sprintf(`[default: %s]`, s)
		}
	}
	return defaultValueString
}

func getFlagName(flag cli.Flag) string {
	names := flag.Names()
	sort.Slice(names, func(i, j int) bool {
		return len(names[i]) < len(names[j])
	})
	prefixed := cli.FlagNamePrefixer(names, "")
	// If there is no short name, add some padding to align flag name.
	if len(names) == 1 {
		prefixed = "    " + prefixed
	}

	return prefixed
}

func handleNoANSI() {
	if noANSI || env.IsNoANSI() {
		color.Disable()
	} else {
		color.Enable()
	}
}

func helpName(fullName string) string {
	return fullName
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.ReplaceAll(v, "\n", "\n"+pad)
}

// lexicographicLess compares strings alphabetically considering case.
func lexicographicLess(i, j string) bool {
	iRunes := []rune(i)
	jRunes := []rune(j)

	lenShared := min(len(iRunes), len(jRunes))

	for index := range lenShared {
		ir := iRunes[index]
		jr := jRunes[index]

		if lir, ljr := unicode.ToLower(ir), unicode.ToLower(jr); lir != ljr {
			return lir < ljr
		}

		if ir != jr {
			return ir < jr
		}
	}

	return i < j
}

func offset(input string, fixed int) int {
	return len(input) + fixed
}

func offsetCommands(cmd []*cli.Command, fixed int) int {
	var maxLen = 0
	for i := range cmd {
		if s := strings.Join(cmd[i].Names(), ", "); len(s) > maxLen {
			maxLen = len(s)
		}
	}
	return maxLen + fixed
}

func offsetFlags(flags []cli.Flag, fixed int) int {
	var maxLen = 0
	for i := range flags {
		if s := getFlagName(flags[i]); len(s) > maxLen {
			maxLen = len(s)
		}
	}
	return maxLen + fixed
}

func onUsageError(_ context.Context, _ *cli.Command, err error, _ bool) error {
	if flag, ok := strings.CutPrefix(err.Error(), "flag provided but not defined: -"); ok {
		color.Red().Printfln("The '%s' option does not exist.", flag)
		return nil
	}
	if flag, ok := strings.CutPrefix(err.Error(), "flag needs an argument: "); ok {
		color.Red().Printfln("The '%s' option requires a value.", flag)
		return nil
	}
	if errMsg := err.Error(); strings.HasPrefix(errMsg, "invalid value") && strings.Contains(errMsg, "for flag -") {
		var value, flag string
		if _, parseErr := fmt.Sscanf(errMsg, "invalid value %q for flag -%s", &value, &flag); parseErr == nil {
			color.Red().Printfln("Invalid value '%s' for option '%s'.", value, strings.TrimSuffix(flag, ":"))
			return nil
		}
	}
	if errMsg := err.Error(); strings.HasPrefix(errMsg, "invalid value") && strings.Contains(errMsg, "for argument") {
		var value, argument string
		if _, parseErr := fmt.Sscanf(errMsg, "invalid value %q for argument %s", &value, &argument); parseErr == nil {
			var subErrMsg string
			subErrMsgPos := strings.Index(errMsg, ":")
			if subErrMsgPos != -1 {
				subErrMsg = errMsg[subErrMsgPos+2:]
			}
			color.Red().Printfln("Invalid value '%s' for argument '%s'. Error: %s", value, strings.TrimSuffix(argument, ":"), subErrMsg)
			return nil
		}
	}

	if errMsg := err.Error(); strings.HasPrefix(errMsg, "sufficient count of arg") && strings.Contains(errMsg, "not provided") {
		var argument string
		var given, expected int
		if _, parseErr := fmt.Sscanf(errMsg, "sufficient count of arg %s not provided, given %d expected %d", &argument, &given, &expected); parseErr == nil {
			if expected == 1 {
				color.Red().Printfln("The '%s' argument requires a value.", argument)
			} else {
				color.Red().Printfln("The '%s' argument requires at least %d values.", argument, expected)
			}
			return nil
		}
	}

	return err
}

func printHelpCustom(out io.Writer, templ string, data any, _ map[string]any) {
	funcMap := template.FuncMap{
		"capitalize":         capitalize,
		"colorize":           colorize,
		"getFlagName":        getFlagName,
		"getFlagDefaultText": getFlagDefaultText,
		"indent":             indent,
		"helpName":           helpName,
		"join":               strings.Join,
		"offset":             offset,
		"offsetCommands":     offsetCommands,
		"offsetFlags":        offsetFlags,
		"sortVisibleFlags":   sortVisibleFlags,
		"sortCommands":       sortCommands,
		"subtract":           subtract,
		"trim":               strings.TrimSpace,
		"wrap":               wrap,
	}

	w := tabwriter.NewWriter(out, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Funcs(colorsFuncMap).Parse(templ))
	templates := map[string]string{
		"usageTemplate":   usageTemplate,
		"commandTemplate": commandTemplate,
		"flagTemplate":    flagTemplate,
		"argsTemplate":    argsTemplate,
		"argTemplate":     argTemplate,
	}
	for name, value := range templates {
		if _, err := t.New(name).Parse(value); err != nil {
			printTemplateError(err)
		}
	}

	handleNoANSI()

	err := t.Execute(w, data)
	if err != nil {
		// If the writer is closed, t.Execute will fail, and there's nothing
		// we can do to recover.
		printTemplateError(err)
		return
	}
	_ = w.Flush()
}

func printTemplateError(err error) {
	if os.Getenv("CLI_TEMPLATE_ERROR_DEBUG") != "" {
		_, _ = fmt.Fprintf(cli.ErrWriter, "CLI TEMPLATE ERROR: %+v\n", err)
	}
}

func printVersion(cmd *cli.Command) {
	handleNoANSI()

	_, _ = fmt.Fprintf(cmd.Writer, "%v %v\n", cmd.Usage, color.Green().Sprint(cmd.Version))
}

func sortCommands(commands []*cli.Command) []*cli.Command {
	sort.Slice(commands, func(i, j int) bool {
		return lexicographicLess(commands[i].Name, commands[j].Name)
	})

	return commands
}

func sortVisibleFlags(cmd *cli.Command) []cli.Flag {
	var (
		flags       = cmd.VisibleFlags()
		globalFlags = cmd.Root().VisibleFlags()
	)
	sort.Sort(cli.FlagsByName(flags))

	globalFlagNames := make(map[string]struct{})
	for i := range globalFlags {
		globalFlagNames[getFlagName(globalFlags[i])] = struct{}{}
	}
	sort.Slice(flags, func(i, j int) bool {
		_, isGlobalI := globalFlagNames[getFlagName(flags[i])]
		_, isGlobalJ := globalFlagNames[getFlagName(flags[j])]
		return !isGlobalI && isGlobalJ
	})

	return flags
}

func subtract(a, b int) int {
	return a - b
}

func wrap(input string, offset int) string {
	var ss []string

	lines := strings.Split(input, "\n")

	padding := strings.Repeat(" ", offset)

	for i, line := range lines {
		if line == "" {
			ss = append(ss, line)
		} else {
			wrapped := wrapLine(line, offset, padding)
			if i == 0 {
				ss = append(ss, wrapped)
			} else {
				ss = append(ss, padding+wrapped)

			}
		}
	}

	return strings.Join(ss, "\n")
}

func wrapLine(input string, offset int, padding string) string {
	if maxLineLength <= offset || len(input) <= maxLineLength-offset {
		return input
	}

	lineWidth := maxLineLength - offset
	words := strings.Fields(input)
	if len(words) == 0 {
		return input
	}

	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + padding + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}

	return wrapped
}
