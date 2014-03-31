package core

import (
	"errors"
	"fmt"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	"io/ioutil"
	"reflect"
)

func Import(env *vm.Env) {
	env.Define("println", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		for i, arg := range args {
			if i != 0 {
				fmt.Print(", ")
			}
			v := arg
			if v.Kind() == reflect.Interface {
				v = v.Elem()
			}
			if arg.IsValid() {
				fmt.Print(arg.Interface())
			} else {
				fmt.Println("undefined")
			}
		}
		fmt.Println()
		return vm.NilValue, nil
	}))

	env.Define("len", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		if len(args) < 1 {
			return vm.NilValue, errors.New("Missing arguments")
		}
		if len(args) > 1 {
			return vm.NilValue, errors.New("Too many arguments")
		}
		v := args[0]
		if v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
			return vm.NilValue, errors.New("Argument should be array")
		}
		return reflect.ValueOf(v.Len()), nil
	}))

	env.Define("keys", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		if len(args) < 1 {
			return vm.NilValue, errors.New("Missing arguments")
		}
		if len(args) > 1 {
			return vm.NilValue, errors.New("Too many arguments")
		}
		v := args[0]
		if v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		if v.Kind() != reflect.Map {
			return vm.NilValue, errors.New("Argument should be map")
		}
		keys := []string{}
		mk := v.MapKeys()
		for _, key := range mk {
			keys = append(keys, key.String())
		}
		return reflect.ValueOf(keys), nil
	}))

	env.Define("bytes", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		if len(args) < 1 {
			return vm.NilValue, errors.New("Missing arguments")
		}
		if len(args) > 1 {
			return vm.NilValue, errors.New("Too many arguments")
		}
		if args[0].Kind() != reflect.String {
			return vm.NilValue, errors.New("Argument should be string")
		}
		return reflect.ValueOf([]byte(args[0].String())), nil
	}))

	env.Define("runes", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		if len(args) < 1 {
			return vm.NilValue, errors.New("Missing arguments")
		}
		if len(args) > 1 {
			return vm.NilValue, errors.New("Too many arguments")
		}
		if args[0].Kind() != reflect.String {
			return vm.NilValue, errors.New("Argument should be string")
		}
		return reflect.ValueOf([]rune(args[0].String())), nil
	}))

	env.Define("string", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		if len(args) < 1 {
			return vm.NilValue, errors.New("Missing arguments")
		}
		if len(args) > 1 {
			return vm.NilValue, errors.New("Too many arguments")
		}
		if args[0].Kind() == reflect.Invalid {
			return vm.NilValue, errors.New("Argument is undefined")
		}
		b, ok := args[0].Interface().([]byte)
		if !ok {
			return vm.NilValue, errors.New("Argument should be byte array")
		}
		return reflect.ValueOf(string(b)), nil
	}))

	env.Define("to_string", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		if len(args) < 1 {
			return vm.NilValue, errors.New("Missing arguments")
		}
		if len(args) > 1 {
			return vm.NilValue, errors.New("Too many arguments")
		}
		return reflect.ValueOf(fmt.Sprint(args[0].Interface())), nil
	}))

	env.Define("char", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		if len(args) < 1 {
			return vm.NilValue, errors.New("Missing arguments")
		}
		if len(args) > 1 {
			return vm.NilValue, errors.New("Too many arguments")
		}
		if args[0].Kind() != reflect.Int && args[0].Kind() != reflect.Int64 {
			return vm.NilValue, errors.New("Argument should be int")
		}
		return reflect.ValueOf(string(rune(args[0].Int()))), nil
	}))

	env.Define("rune", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		if len(args) < 1 {
			return vm.NilValue, errors.New("Missing arguments")
		}
		if len(args) > 1 {
			return vm.NilValue, errors.New("Too many arguments")
		}
		if args[0].Kind() != reflect.String {
			return vm.NilValue, errors.New("Argument should be string")
		}
		s := args[0].String()
		if len(s) == 0 {
			return reflect.ValueOf(0), nil
		}
		return reflect.ValueOf(s[0]), nil
	}))

	env.Define("load", vm.ToFunc(func(args ...reflect.Value) (reflect.Value, error) {
		if len(args) < 1 {
			return vm.NilValue, errors.New("Missing arguments")
		}
		if len(args) > 1 {
			return vm.NilValue, errors.New("Too many arguments")
		}
		if args[0].Kind() != reflect.String {
			return vm.NilValue, errors.New("Argument should be string")
		}
		body, err := ioutil.ReadFile(args[0].String())
		if err != nil {
			return vm.NilValue, err
		}
		scanner := new(parser.Scanner)
		scanner.Init(string(body))
		stmts, err := parser.Parse(scanner)
		if err != nil {
			return vm.NilValue, err
		}
		return vm.RunStmts(stmts, env)
	}))
}
