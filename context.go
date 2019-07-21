package kaa

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"reflect"
	"strconv"
	"strings"
)

// The context for a particular run of a command.
type Context interface {
	// Get the original cobra command
	GetCommand() *cobra.Command

	// Args passed into the command
	GetArgs() []string

	// Binds the struct ptr to the data passed in by the user either by args or by flags
	// Two tags are currently supported for the struct: arg and flag
	//
	//
	// Arg tag value: The index for the arg
	//   The arg tag also accepts parameters separated by a comma after the index
	//     Name string `arg:"0,optional"`
	//   Valid parameters for the arg tag are
	//     optional: If the arg cannot be found, no panic is raised and the field is kept as the default.
	// Flag tag value: The name of the flag specified in the cobra.Command
	Bind(structPtr interface{}) error

	// If an error occurred, this method will return a non-nil error.
	Error() error
}

func NewContext(cmd *cobra.Command, args []string) Context {
	return &contextImpl{
		cmd:  cmd,
		args: args,
	}
}

type cmdArg struct {
	index      int
	isOptional bool
}

func cmdArgFromTagStr(tagStr string) *cmdArg {
	if tagStr != "" {
		arg := new(cmdArg)
		argFlagVals := strings.Split(tagStr, ",")
		argInd, err := strconv.Atoi(argFlagVals[0])
		if err != nil {
			return nil
		}
		arg.index = argInd
		if len(argFlagVals) > 1 {
			arg.isOptional = strings.ToLower(argFlagVals[1]) == "optional"
		}
		return arg
	}
	return nil
}

type cmdFlag struct {
	name string
}

func cmdFlagFromTagStr(tagStr string) *cmdFlag {
	if tagStr != "" {
		flag := &cmdFlag{
			name: tagStr,
		}
		return flag
	}
	return nil
}

type contextImpl struct {
	cmd   *cobra.Command
	args  []string
	error error
}

func (c *contextImpl) Error() error {
	return c.error
}

func (c *contextImpl) Bind(structPtr interface{}) error {
	val := reflect.ValueOf(structPtr)
	if !c.isStructPtr(val) {
		return nil
	}
	val = val.Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		vField := val.Field(i)
		tField := typ.Field(i)
		if vField.CanSet() {
			arg := tField.Tag.Get("arg")
			flag := tField.Tag.Get("flag")
			cmdArg := cmdArgFromTagStr(arg)
			cmdFlag := cmdFlagFromTagStr(flag)
			if cmdArg != nil {
				if len(c.args) <= cmdArg.index {
					if !cmdArg.isOptional {
						return errors.New(fmt.Sprintf("not enough args, %s is required", strings.ToLower(tField.Name)))
					}
				} else {
					foundType := c.findType(c.args[cmdArg.index], vField)
					if foundType != nil {
						vField.Set(*foundType)
					}
				}
			} else if cmdFlag != nil {
				val := c.findFlagType(flag, vField.Type())
				if val != nil {
					vField.Set(reflect.ValueOf(val))
				}
			}
		}
	}

	return nil
}

func stripPtr(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func (c *contextImpl) findType(raw string, value reflect.Value) *reflect.Value {
	var val interface{}
	value = stripPtr(value)
	switch value.Kind() {
	case reflect.String:
		val = raw
	case reflect.Int:
		val, _ = strconv.Atoi(raw)
	case reflect.Int8:
		i, _ := strconv.Atoi(raw)
		val = int8(i)
	case reflect.Int16:
		i, _ := strconv.Atoi(raw)
		val = int16(i)
	case reflect.Int32:
		i, _ := strconv.Atoi(raw)
		val = int32(i)
	case reflect.Int64:
		i, _ := strconv.Atoi(raw)
		val = int64(i)
	case reflect.Uint:
		i, _ := strconv.Atoi(raw)
		val = uint(i)
	case reflect.Uint8:
		i, _ := strconv.Atoi(raw)
		val = uint8(i)
	case reflect.Uint16:
		i, _ := strconv.Atoi(raw)
		val = uint16(i)
	case reflect.Uint32:
		i, _ := strconv.Atoi(raw)
		val = uint32(i)
	case reflect.Uint64:
		i, _ := strconv.Atoi(raw)
		val = uint64(i)
	case reflect.Bool:
		val, _ = strconv.ParseBool(raw)
	}
	if val != nil {
		value = reflect.ValueOf(val)
		return &value
	}
	return nil
}

func (c *contextImpl) findFlagType(flag string, value reflect.Type) interface{} {
	var val interface{}
	switch value.Kind() {
	case reflect.String:
		val, _ = c.GetCommand().Flags().GetString(flag)
	case reflect.Slice:
		val = c.findSliceFlagType(flag, value.Elem())
	case reflect.Int:
		val, _ = c.GetCommand().Flags().GetInt(flag)
	case reflect.Bool:
		val, _ = c.GetCommand().Flags().GetBool(flag)
	case reflect.Float32:
		val, _ = c.GetCommand().Flags().GetFloat32(flag)
	case reflect.Float64:
		val, _ = c.GetCommand().Flags().GetFloat64(flag)
	}

	return val
}

func (c *contextImpl) findSliceFlagType(flag string, value reflect.Type) interface{} {
	var val interface{}
	switch value.Kind() {
	case reflect.String:
		val, _ = c.GetCommand().Flags().GetStringSlice(flag)
	case reflect.Int:
		val, _ = c.GetCommand().Flags().GetIntSlice(flag)
	case reflect.Bool:
		val, _ = c.GetCommand().Flags().GetBoolSlice(flag)
	}

	return val
}

func (c *contextImpl) isStructPtr(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct
}

func (c *contextImpl) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *contextImpl) GetArgs() []string {
	return c.args
}
