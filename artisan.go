package artisan

import (
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

var RunCommand string

func InitCommand[T any](command T) (result T, err error) {
	initCommands := []*commandMain{}

	v := reflect.TypeOf(command)
	value := reflect.ValueOf(&command)
	RunCommand = os.Args[1]

	// Init and insert to memory
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag, _ := field.Tag.Lookup("artisan")
		fieldType := field.Type

		mainCommand := commandMain{}
		mainCommand.Setter(tag, i)

		for a := 0; a < fieldType.NumField(); a++ {
			subField := fieldType.Field(a)
			subTag, _ := subField.Tag.Lookup("artisan")

			mainCommand.OptSetter(subTag, a)
		}

		initCommands = append(initCommands, &mainCommand)
	}

	//Set args to initCommands
	for _, cmd := range initCommands {
		if cmd.Name == os.Args[1] {
			for c := 2; c < len(os.Args); c++ {
				c = setterOption(cmd, os.Args, c)
			}
		}
	}

	//Assign value to struct
	for _, cmd := range initCommands {
		for _, arg := range cmd.Args {
			setterField(v.Field(cmd.Index), value.Elem().FieldByIndex([]int{cmd.Index}), arg.IndexField, arg.Value)
		}

		for _, opt := range cmd.Options {
			setterField(v.Field(cmd.Index), value.Elem().FieldByIndex([]int{cmd.Index}), opt.IndexField, opt.Value)
		}
	}

	return command, nil
}

func setterField(field reflect.StructField, value reflect.Value, indexCommandOpt int, argValue string) {

	if value.FieldByIndex([]int{indexCommandOpt}).Kind() == reflect.Ptr {
		switch field.Type.Field(indexCommandOpt).Type.String() {
		case FieldTypeStringPtr:
			value.FieldByIndex([]int{indexCommandOpt}).Set(reflect.ValueOf(&argValue))
		case FieldTypeBoolPtr:
			v, _ := strconv.ParseBool(argValue)
			value.FieldByIndex([]int{indexCommandOpt}).Set(reflect.ValueOf(&v))
		case FieldTypeIntPtr:
			v, _ := strconv.ParseInt(argValue, 10, 64)
			value.FieldByIndex([]int{indexCommandOpt}).Set(reflect.ValueOf(&v))
		case FieldTypeUIntPtr:
			v, _ := strconv.ParseUint(argValue, 10, 64)
			value.FieldByIndex([]int{indexCommandOpt}).Set(reflect.ValueOf(&v))
		default:
		}
	}

	switch field.Type.Field(indexCommandOpt).Type.String() {
	case FieldTypeString:
		value.FieldByIndex([]int{indexCommandOpt}).SetString(argValue)
	case FieldTypeBool:
		v, _ := strconv.ParseBool(argValue)
		value.FieldByIndex([]int{indexCommandOpt}).SetBool(v)
	case FieldTypeInt:
		v, _ := strconv.ParseInt(argValue, 10, 64)
		value.FieldByIndex([]int{indexCommandOpt}).SetInt(v)
	case FieldTypeUInt:
		v, _ := strconv.ParseUint(argValue, 10, 64)
		value.FieldByIndex([]int{indexCommandOpt}).SetUint(v)
	default:
	}

}

func setterOption(base *commandMain, value []string, index int) int {
	var cmd *commandOpt
	var i int

	//regex : if string start with '-'
	regex, _ := regexp.Compile("^-.*")
	if regex.MatchString(value[index]) {
		opt := value[index]
		filtered := strings.TrimPrefix(opt, opt[0:1])
		cmd, i = filterSubCommand(base.Options, filtered)

		v, next := extractValue(value[index], cmd.HasValue)
		if next {
			cmd.SetValue(value[index+1])
			base.Options[i] = *cmd

			return index + 1
		}
		cmd.SetValue(v)
		base.Options[i] = *cmd
	} else {
		cmd, i = filterSubCommand(base.Args, value[index])
		v, _ := extractValue(value[index], cmd.HasValue)

		cmd.SetValue(v)
		base.Args[i] = *cmd
	}

	return index
}

func filterSubCommand(listSub []commandOpt, value string) (*commandOpt, int) {
	for i, sub := range listSub {
		if sub.Value == "" {
			if strings.Contains(value, "=") {
				split := strings.Split(value, "=")
				value = split[0]
			}

			if slices.Contains(sub.Alias, value) {
				v := sub
				return &v, i
			}

			if !strings.Contains(value, "-") && !strings.Contains(value, "=") {
				v := sub
				return &v, i
			}
		}
	}

	return nil, 0
}

func extractValue(arg string, hasValue bool) (value string, nextArg bool) {
	//if string contains "="
	if strings.Contains(arg, "=") {
		split := strings.Split(arg, "=")
		return split[1], false
	}

	//if string contains "-" (option) but not contains "=" the value are next args
	if strings.Contains(arg, "-") && !strings.Contains(arg, "=") && hasValue {
		return "", true
	}

	//if string contains "-" (option) but not contains "=" the value are next args
	if strings.Contains(arg, "-") && !strings.Contains(arg, "=") && !hasValue {
		return "true", false
	}

	//if string does not contains "-" && does not contains "=" means it's a arguments
	if !strings.Contains(arg, "-") && !strings.Contains(arg, "=") {
		return arg, true
	}

	//return empty if unknown
	return "", false
}
