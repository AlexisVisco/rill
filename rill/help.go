package rill

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
)

type ArgList []Arg

func (p ArgList) Len() int           { return len(p) }
func (p ArgList) Less(i, j int) bool { return p[i].name < p[j].name }
func (p ArgList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Arg struct {
	name string
	desc string
}

func Help(cmd Commands) {
	fmt.Println("Usage:  taskmaster COMMAND [OPTIONS]")
	fmt.Println("\nA simple management process write for 101\n")

	sort.Strings(cmd.commandsList)

	arglist := make(ArgList, 0)
	for _, label := range cmd.commandsList {
		arglist = append(arglist, Arg{
			name: label,
			desc: cmd.commands[label].command.CommandDescription(),
		})
	}
	alignKeyValue("Commands", arglist)
	fmt.Println("\nRun 'taskmaster COMMAND --help' for more information on a command namespace.")
}

func HelpCommand(namespace NamespaceCommand) {
	fmt.Printf("Usage:  taskmaster %s [ARGS] [OPTIONS]\n", namespace.command.CommandLabel())
	fmt.Printf("\nAliases: %s\n", strings.Join(namespace.command.CommandAliases(), ", "))
	fmt.Printf("\n%s\n\n", namespace.command.CommandDescription())
	arglist := make(ArgList, 0)
	for _, v := range namespace.infos {
		arglist = append(arglist, Arg{
			name: getTypesFormatted(v, namespace.command),
			desc: v.Description,
		})
	}
	alignKeyValue("Args:", arglist)
	println()
	alignFlags(namespace.FlagsDescription, namespace.command)
}

func getTypesFormatted(info CommandInfo, c Command) string {
	str := ""
	vStruct := reflect.ValueOf(c)
	met := vStruct.MethodByName(info.MethodName).Type()
	for i := 0; i < info.numParams; i++ {
		name := ""
		if len(info.namesParams)-1 >= i {
			name = info.namesParams[i]
		}
		if strings.HasPrefix(name, "*") {
			str += strings.Trim(name, "*")
		} else {
			str += "<" + name + " " + strings.ToLower(met.In(i).String()) + ">"
		}
		if i != info.numParams-1 {
			str += " "
		}
	}
	if info.numParams == 0 {
		str = "default"
	}
	return str
}

func alignFlags(flags map[string]string, command Command) {
	if len(flags) == 0 {
		return
	}
	fmt.Println("Options:")
	vStruct := reflect.ValueOf(command)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
	for k, v := range flags {
		field := vStruct.FieldByName(strings.Title(k))
		fmt.Fprintln(w, fmt.Sprintf("  --%s=%s\t%s", k, strings.ToLower(field.Kind().String()), v))
	}
	w.Flush()
}

func alignKeyValue(title string, commandDesc ArgList) {
	sort.Sort(commandDesc)
	fmt.Println(title)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
	for _, v := range commandDesc {
		fmt.Fprintln(w, fmt.Sprintf("  %s\t%s", v.name, v.desc))
	}
	w.Flush()
}
