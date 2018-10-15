package rill

import (
	"reflect"
	"fmt"
	"strings"
	"strconv"
	"errors"
)

type Command interface {
	CommandDescription() string
	CommandLabel() string
	CommandAliases() []string
}

type CommandProcessor struct {
	command          Command
	Description      string
	ShortDescription string
	MethodName       string

	numParams int
}

type Commands struct {
	commands map[string][]CommandProcessor
}

func Rill() *Commands {
	return &Commands{make(map[string][]CommandProcessor)}
}

func (c *Commands) Add(command Command, method string, shortDesc string, desc string) {
	met := reflect.ValueOf(command).MethodByName(method)
	processor := CommandProcessor{
		command:          command,
		Description:      desc,
		ShortDescription: shortDesc,
		MethodName:       method,
		numParams:        met.Type().NumIn(),
	}
	c.addToMap(command.CommandLabel(), processor)
	for _, alias := range command.CommandAliases() {
		c.addToMap(alias, processor)
	}
	fmt.Println(c)
}

func (c *Commands) addToMap(name string, processor CommandProcessor) {
	if o, ok := c.commands[name]; ok {
		c.commands[name] = append(o, processor)
	} else {
		c.commands[name] = []CommandProcessor{processor}
	}
}

func (c *Commands) Dispatch(cmd string) {
	words := strings.Split(cmd, " ")
	if len(words) > 0 {
		if o, ok := c.commands[words[0]]; ok {
			for _, processor := range o {
				if processor.numParams == len(words)-1 {
					method := reflect.ValueOf(processor.command).MethodByName(processor.MethodName)
					parameters := make([]reflect.Value, method.Type().NumIn())
					errorConvert := false
					for x := 0; x < processor.numParams; x++ {
						fromString, err := valueFromString(method.Type().In(x).Kind(), words[x+1])
						if err != nil {
							errorConvert = true
							break
						}
						parameters[x] = reflect.ValueOf(fromString)
					}
					if errorConvert {
						continue
					} else {
						method.Call(parameters)
						return
					}
				}
			}
		}
	}
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
