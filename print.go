package acyclic

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Print pretty-prints v to stdout.
// If the structure contains any cycles, these will be pruned and marked in the output.
func Print(v interface{}) {
	fmt.Print(Sprint(v))
}

// Fprint formats v and outputs it to w.
// The format matches that of Print.
func Fprint(w io.Writer, v interface{}) {
	fmt.Fprint(w, Sprint(v))
}

// Sprint returns the pretty-printed form of v as a string.
// The form matches that of Print.
func Sprint(v interface{}) string {
	return doSprint(reflect.ValueOf(v), nil)
}

func doSprint(value reflect.Value, parents []uintptr) string {
	kind := value.Kind()

	if kind == reflect.Interface {
		value = value.Elem()
		kind = value.Kind()
	}

	newParents, err := checkParents(value, parents, nil)
	if err != nil {
		return "<CYCLE>"
	}

	if kind == reflect.Map {
		var lines []string
		for _, key := range value.MapKeys() {
			lines = append(lines, fmt.Sprintf("%v: %v", key, doSprint(value.MapIndex(key), newParents)))
		}
		return fmt.Sprintf("{\n  %v\n}", strings.Join(lines, "\n  "))
	}

	if kind == reflect.Ptr {
		return fmt.Sprintf("*%v", doSprint(value.Elem(), newParents))
	}

	if kind == reflect.Slice {
		var lines []string
		for i := 0; i < value.Len(); i++ {
			lines = append(lines, value.Index(i).String())
		}
		return fmt.Sprintf("[\n  %v\n]", strings.Join(lines, "\n  "))
	}

	if kind == reflect.Struct {
		var lines []string
		for i := 0; i < value.NumField(); i++ {
			t := value.Type()
			fieldType := t.Field(i)
			lines = append(lines, fmt.Sprintf("%v: %v", fieldType.Name, doSprint(value.Field(i), newParents)))
		}
		return fmt.Sprintf("{\n  %v\n}", strings.Join(lines, "\n  "))
	}

	return fmt.Sprintf("%#v", value)
}
