package rill

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
)

func Help(cmd Commands) {
	fmt.Println("Usage:  taskmaster COMMAND [OPTIONS]")
	fmt.Println("\nA simple management process write for 101\n")

	sort.Strings(cmd.commandsList)

	cmdDesc := make(map[string]string)
	for _, label := range cmd.commandsList {
		cmdDesc[label] = cmd.commands[label].command.CommandDescription()
	}
	alignKeyValue("Commands", cmdDesc)
	fmt.Println("\nRun 'taskmaster COMMAND --help' for more information on a command namespace.")
}

func HelpCommand(namespace NamespaceCommand) {
	fmt.Printf("Usage:  taskmaster %s [ARGS] [OPTIONS]\n", namespace.command.CommandLabel())
	fmt.Printf("\nAliases: %s\n", strings.Join(namespace.command.CommandAliases(), ", "))
	fmt.Printf("\n%s\n\n", namespace.command.CommandDescription())
	argsDesc := make(map[string]string)
	for _, v := range namespace.infos {
		argsDesc[getTypesFormatted(v, namespace.command)] = v.Description
	}
	alignKeyValue("Args:", argsDesc)
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
	fmt.Println("Options:")
	vStruct := reflect.ValueOf(command)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
	for k, v := range flags {
		field := vStruct.FieldByName(strings.Title(k))
		fmt.Fprintln(w, fmt.Sprintf("  --%s=%s\t%s", k, strings.ToLower(field.Kind().String()), v))
	}
	w.Flush()
}

func alignKeyValue(title string, commandDesc map[string]string) {
	fmt.Println(title)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
	for k, v := range commandDesc {
		fmt.Fprintln(w, fmt.Sprintf("  %s\t%s", k, v))
	}
	w.Flush()
}
