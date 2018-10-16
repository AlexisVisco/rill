package rill

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Command interface {
	CommandDescription() string
	CommandLabel() string
	CommandAliases() []string
	Clone() interface{}
}

type CommandProcessor struct {
	MethodName       string
	Description      string
	ShortDescription string
	numParams        int
}

type ChildCommand struct {
	parent  *Commands
	command Command

	processors       []CommandProcessor
	Description      string
	ShortDescription string
	Flags            map[string]string
	FlagsDescription map[string]string
}

func (c *ChildCommand) Cmd(methodName string, desc string, shortDesc string) *ChildCommand {
	met := reflect.ValueOf(c.command).MethodByName(methodName)
	processor := CommandProcessor{
		MethodName:       methodName,
		Description:      desc,
		ShortDescription: shortDesc,
		numParams:        met.Type().NumIn(),
	}
	c.processors = append(c.processors, processor)
	return c
}

func (c *ChildCommand) End() *Commands {
	return c.parent
}
func (c *ChildCommand) injectFlags(cmd interface{}, flags map[string]string) {
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

type Commands struct {
	commands map[string]*ChildCommand
}

func Rill() *Commands {
	return &Commands{make(map[string]*ChildCommand)}
}

func (c *Commands) Add(command Command, shortDesc string, desc string) *ChildCommand {
	child := ChildCommand{
		parent:           c,
		command:          command,
		Description:      desc,
		ShortDescription: shortDesc,
		processors:       make([]CommandProcessor, 0),
	}
	c.getFlags(command, &child)
	c.addToMap(command, &child)
	return &child
}

func (Commands) getFlags(command Command, child *ChildCommand) {
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
	child.Flags = fl
	child.FlagsDescription = flDesc
}

func (c *Commands) addToMap(command Command, child *ChildCommand) {
	c.commands[command.CommandLabel()] = child
	for _, alias := range command.CommandAliases() {
		c.commands[alias] = child
	}
}

func (c *Commands) Dispatch(args []string) {
	flags, args := c.parseFlags(args)

	if len(args) > 0 {
		if o, ok := c.commands[args[0]]; ok {
			for _, processor := range o.processors {
				if processor.numParams == len(args)-1 {
					used := o.command.Clone()
					method := reflect.ValueOf(used).MethodByName(processor.MethodName)
					parameters := make([]reflect.Value, method.Type().NumIn())
					errorConvert := false
					for x := 0; x < processor.numParams; x++ {
						fromString, err := valueFromString(method.Type().In(x).Kind(), args[x+1])
						if err != nil {
							errorConvert = true
							break
						}
						parameters[x] = reflect.ValueOf(fromString)
					}
					if errorConvert {
						continue
					} else {
						o.injectFlags(used, flags)
						method.Call(parameters)
						return
					}
				}
			}
		}
	}
}

func (Commands) parseFlags(s []string) (map[string]string, []string) {
	args := make([]string, 0)
	flags := make(map[string]string)

	for _, arg := range s {
		if strings.HasPrefix(arg, "--") { // check if it is a flag
			if strings.Contains(arg, "=") { // check if the format is --flag=value
				split := strings.Split(arg, "=")
				if len(split) < 2 {
					continue
				} else { // the flag is a boolean because --flag is eq to true
					key := strings.TrimPrefix(split[0], "--")
					flags[key] = split[1]
				}
			} else { // it is not a flag
				key := strings.TrimPrefix(arg, "--")
				flags[key] = "true"
			}
		} else {
			args = append(args, arg)
		}
	}
	return flags, args
}

func valueFromString(kind reflect.Kind, strVal string) (interface{}, error) {
	switch kind {
	case reflect.Int64:
		val, err := strconv.ParseInt(strVal, 0, 64)
		if err != nil {
			return nil, err
		}
		return val, nil
	case reflect.Uint64:
		val, err := strconv.ParseUint(strVal, 0, 64)
		if err != nil {
			return nil, err
		}
		return val, nil
	case reflect.Float32:
		val, err := strconv.ParseFloat(strVal, 32)
		if err != nil {
			return nil, err
		}
		return val, nil
	case reflect.Float64:
		val, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return nil, err
		}
		return val, nil
	case reflect.String:
		return strVal, nil
	case reflect.Bool:
		val, err := strconv.ParseBool(strVal)
		if err != nil {
			return nil, err
		}
		return val, nil
	default:
		return nil, errors.New("Unsupported kind: " + kind.String())
	}
	return nil, errors.New("Unsupported kind: " + kind.String())
}
