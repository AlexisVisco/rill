package rill

import (
	"reflect"
)

// Dispatch redirect a text into a function.
// Need to match the function prototype.
func (c *Commands) Dispatch(args []string) {
	flags, args := parseFlags(args)

	c.filterCommands(flags, args, func(o *NamespaceCommand, info CommandInfo) bool {
		used := o.command.Clone()
		method := reflect.ValueOf(used).MethodByName(info.MethodName)
		if params, ok := c.parameters(info, method, args); ok {
			o.injectFlags(used, flags)
			method.Call(params)
			return true
		}
		return false
	})
}

// filterCommands retrieve the namespace from the first arg.
// iterate through all command info and check if the number
// of parameters match with the number arguments.
//
// If it is the case execute the function, if the closure return
// true, no need to iterate with next commands info because there
// is already a match. Else continue to iterate.
func (c *Commands) filterCommands(flags map[string]string, args []string, foreach func(namespace *NamespaceCommand, info CommandInfo) bool) {
	if len(args) == 0 {
		Help(*c)
	} else if o, ok := c.commands[args[0]]; ok {
		for _, info := range o.infos {
			if _, ok := flags["help"]; ok {
				HelpCommand(*o)
				return
			}
			if info.numParams == len(args)-1 {
				if foreach(o, info) {
					return
				}
			}
		}
		HelpCommand(*o)
	} else {
		Help(*c)
	}
}

// parameters retrieve transform arguments passed to command line
// into the function parameters CustomTypes.
// Return the list of values of a error (during the parsing of the args)
func (c *Commands) parameters(info CommandInfo, method reflect.Value, args []string) ([]reflect.Value, bool) {
	parameters := make([]reflect.Value, method.Type().NumIn())
	ok := true
	for x := 0; x < info.numParams; x++ {
		fromString, err := valueFromString(method.Type().In(x), args[x+1])
		if err != nil {
			ok = false
			break
		}
		parameters[x] = reflect.ValueOf(fromString)
	}
	return parameters, ok
}
