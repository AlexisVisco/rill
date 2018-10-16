package rill

import (
	"reflect"
	"strings"
)

// The CommandInfo has all information to
// describe the command, and bind text to
// go code.
//
// Example:
// func (t) hello(s string)
// =
// t hello => t.hello("hello")
type CommandInfo struct {
	MethodName       string
	Description      string
	ShortDescription string
	numParams        int
}

// The NamespaceCommand is a entity that own a amount
// of commands, it can be represented as a sort of
// namespace.
//
// Description and short description describe what the
// "namespace" manage.
//
// He has the amount of flags described in the structure.
type NamespaceCommand struct {
	parent  *Commands // useful to reconnect to the parent and add another command namespace
	command Command

	infos            []CommandInfo
	Description      string
	ShortDescription string
	Flags            map[string]string
	FlagsDescription map[string]string
}

func newNamespaceCommand(parent *Commands, command Command, shortDesc string, desc string) *NamespaceCommand {
	namespace := &NamespaceCommand{
		parent:           parent,
		command:          command,
		Description:      desc,
		ShortDescription: shortDesc,
		infos:            make([]CommandInfo, 0),
	}
	namespace.registerFlags(command)
	return namespace
}

func (c *NamespaceCommand) Cmd(methodName string, desc string, shortDesc string) *NamespaceCommand {
	met := reflect.ValueOf(c.command).MethodByName(methodName)
	info := CommandInfo{
		MethodName:       methodName,
		Description:      desc,
		ShortDescription: shortDesc,
		numParams:        met.Type().NumIn(),
	}
	c.infos = append(c.infos, info)
	return c
}

func (c *NamespaceCommand) End() *Commands {
	return c.parent
}

// injectFlags set in the struct fields the flags founded, if an unknown flag, it is
// skipped.
func (c *NamespaceCommand) injectFlags(cmd interface{}, flags map[string]string) {
	if len(flags) > 0 {
		vStruct := reflect.ValueOf(cmd)
		for key, value := range flags {
			if fieldName, ok := c.Flags[key]; ok {
				field := vStruct.Elem().FieldByName(fieldName)
				result, err := valueFromString(field.Kind(), value)
				if err == nil {
					field.Set(reflect.ValueOf(result))
				}
			}
		}
	}
}

// registerFlags take all fields in the struct,
// take aliases with `fl` tag and take description
// with `flDesc`.
func (c *NamespaceCommand) registerFlags(command Command) {
	fl := make(map[string]string)
	flDesc := make(map[string]string)
	tStruct := reflect.TypeOf(command)
	fields := tStruct.NumField()
	for x := 0; x < fields; x++ {
		f := tStruct.Field(x)
		aliases := strings.Split(f.Tag.Get("fl"), ",")
		desc := f.Tag.Get("flDesc")
		flDesc[strings.ToLower(f.Name)] = desc
		fl[strings.ToLower(f.Name)] = f.Name
		for _, al := range aliases {
			fl[al] = f.Name
		}
	}
	c.Flags = fl
	c.FlagsDescription = flDesc
}
