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
	return doSprint(reflect.ValueOf(v), nil, 0)
}

func doSprint(value reflect.Value, parents []uintptr, indents int) string {
	kind := value.Kind()

	if kind == reflect.Interface {
		return doSprint(value.Elem(), parents, indents)
	}

	newParents, err := checkParents(value, parents, nil)
	if err != nil {
		return "<CYCLE>"
	}

	ind := indent(indents)
	baseInd := indent(indents - 1)
	var b strings.Builder

	switch kind {
	case reflect.Map:
		b.WriteString("{")
		b.WriteString("\n")
		for _, key := range value.MapKeys() {
			b.WriteString(ind)
			fmt.Fprintf(&b, "%#v", key)
			b.WriteString(": ")
			b.WriteString(doSprint(value.MapIndex(key), newParents, indents+1))
			b.WriteString("\n")
		}
		b.WriteString(baseInd)
		b.WriteString("}")
	case reflect.Ptr:
		fmt.Fprintf(&b, "*%v", doSprint(value.Elem(), newParents, indents+1))
	case reflect.Slice:
		b.WriteString("[")
		b.WriteString("\n")
		for i := 0; i < value.Len(); i++ {
			b.WriteString(ind)
			b.WriteString(doSprint(value.Index(i), newParents, indents+1))
			b.WriteString(",\n")
		}
		b.WriteString(baseInd)
		b.WriteString("]")
	case reflect.Struct:
		b.WriteString("{")
		b.WriteString("\n")
		for i := 0; i < value.NumField(); i++ {
			t := value.Type()
			fieldType := t.Field(i)
			b.WriteString(ind)
			b.WriteString(fieldType.Name)
			b.WriteString(": ")
			b.WriteString(doSprint(value.Field(i), newParents, indents+1))
			b.WriteString("\n")
		}
		b.WriteString(baseInd)
		b.WriteString("}")
	default:
		fmt.Fprintf(&b, "%#v", value)
	}

	return b.String()
}

func indent(i int) string {
	if i <= 0 {
		return ""
	}
	return strings.Repeat("  ", i)
}
