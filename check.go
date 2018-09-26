package acyclic

import (
	"fmt"
	"reflect"
)

// Check performs a DFS traversal over an interface to determine if it contains any cycles.
// If there are no cycles, nil is returned.
// If one or more cycles exist, an error will be returned. This error will contain a path to the first cycle found.
func Check(v interface{}) error {
	return doCheck(reflect.ValueOf(v), nil, nil)
}

func doCheck(value reflect.Value, parents []uintptr, names []string) error {
	kind := value.Kind()

	if kind == reflect.Interface {
		value = value.Elem()
		kind = value.Kind()
	}

	newParents, err := checkParents(value, parents, names)
	if err != nil {
		return err
	}

	if kind == reflect.Map {
		for _, key := range value.MapKeys() {
			err := doCheck(value.MapIndex(key), newParents, append(names, key.String()))
			if err != nil {
				return err
			}
		}
	}

	if kind == reflect.Ptr {
		return doCheck(value.Elem(), newParents, names)
	}

	if kind == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			err := doCheck(value.Index(i), newParents, append(names, fmt.Sprintf("[%d]", i)))
			if err != nil {
				return err
			}
		}
	}

	if kind == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			t := value.Type()
			fieldType := t.Field(i)
			err := doCheck(value.Field(i), newParents, append(names, fieldType.Name))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func checkParents(value reflect.Value, parents []uintptr, names []string) ([]uintptr, error) {
	kind := value.Kind()
	if kind == reflect.Map || kind == reflect.Ptr || kind == reflect.Slice {
		address := value.Pointer()
		for _, parent := range parents {
			if parent == address {
				return nil, fmt.Errorf("cycle found: %v", names)
			}
		}
		return append(parents, address), nil
	}
	return parents, nil
}
