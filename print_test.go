package acyclic_test

import "github.com/theothertomelliott/acyclic"

func ExamplePrint() {
	// Create a pointer to a struct
	value := &struct {
		A string
		B interface{}
	}{
		A: "a string",
	}

	// Add a cycle
	value.B = value

	acyclic.Print(value)

	// Output:
	// *{
	//   A: "a string"
	//   B: <CYCLE>
	// }
}
