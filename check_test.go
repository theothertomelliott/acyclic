package acyclic_test

import (
	"testing"

	"github.com/theothertomelliott/acyclic"
)

func TestCheck(t *testing.T) {
	var tests = []struct {
		name      string
		in        interface{}
		expectErr bool
	}{
		{
			name: "nil",
		},
		{
			name: "primitive type returns no error",
			in:   35,
		},
		{
			name: "acyclic map returns no error",
			in: map[string]interface{}{
				"a": "b",
				"c": []int{1, 2, 3},
				"d": map[string]interface{}{
					"e": "f",
				},
			},
		},
		{
			name: "cyclic map returns error",
			in: func() interface{} {
				m := map[string]interface{}{
					"a": "b",
				}
				// Create the cycle
				m["b"] = m
				return m
			}(),
			expectErr: true,
		},
		{
			name: "cycle in slice returns error",
			in: func() interface{} {
				m := map[string]interface{}{
					"a": "b",
				}
				// Create the cycle
				s := []interface{}{m}
				m["b"] = s
				return m
			}(),
			expectErr: true,
		},
		{
			name: "acyclic struct returns no error",
			in: testStruct{
				A: "value A",
				B: "value B",
			},
		},
		{
			name: "cyclic struct returns error",
			in: func() interface{} {
				s := &testStruct{
					A: "value A",
				}
				// Create the cycle
				s.B = s
				return s
			}(),
			expectErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := acyclic.Check(test.in)
			if test.expectErr != (err != nil) {
				if test.expectErr {
					t.Error("Expected an error, got nil")
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

type testStruct struct {
	A string
	B interface{}
}
