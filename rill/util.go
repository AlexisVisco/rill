package rill

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var CustomTypes = make(map[reflect.Type]FromString)

type FromString func(string) (interface{}, error)

// parseFlags parse flags like: --flag=value
//                              --flag
// Return a map where --key=value, if --key format
// is used, so the value is true because this format
// is for boolean fields. And a list of string that
// represent the commands without flags.
func parseFlags(s []string) (map[string]string, []string) {
	args := make([]string, 0)
	flags := make(map[string]string)

	for _, arg := range s {
		if strings.HasPrefix(arg, "--") { // check if it is a flag
			if strings.Contains(arg, "=") { // check if the format is --flag=value
				split := strings.Split(arg, "=")
				if len(split) < 2 {
					continue
				} else { // the flag is a boolean because --flag is eq to true
					key := strings.ToLower(strings.TrimPrefix(split[0], "--"))
					flags[key] = split[1]
				}
			} else { // it is not a flag
				key := strings.ToLower(strings.TrimPrefix(arg, "--"))
				flags[key] = "true"
			}
		} else {
			args = append(args, arg)
		}
	}
	return flags, args
}

// valueFronString take a value and return his real value via an
// interface{}.
// Return the value cast into an interface or an error, if the strVal
// cannot be parsed.
func valueFromString(typ reflect.Type, strVal string) (interface{}, error) {
	if f, ok := CustomTypes[typ]; ok {
		val, e := f(strVal)
		if e == nil {
			return val, nil
		}
	}
	switch typ.Kind() {
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
		return nil, errors.New("Unsupported kind: " + typ.String())
	}
	return nil, errors.New("Unsupported kind: " + typ.String())
}
