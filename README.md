# acyclic

[![GoDoc](https://godoc.org/github.com/theothertomelliott/acyclic?status.svg)](https://godoc.org/github.com/theothertomelliott/acyclic)

package acyclic provides the ability to quickly check for cycles within a
data structure before attempting to marshal to JSON or similar formats.

A set of functions are also provided to print these structures in a safe manner
that won't result in a stack overflow, pruning branches containing cycles and
clearly marking where they occurred.

Cycles are detected using depth first search.

# Usage example

```
import "github.com/theothertomelliott/acyclic"

func main() {
    // Create a pointer to a struct
    value := &struct {
        A string
        B interface{}
    }{
        A: "a string",
    }

    // Add a cycle
    value.B = value

    err := acyclic.Check(value)
    if err != nil {
        fmt.Println("Cycle found")
    } else {
        fmt.Println("No cycle")
    }
}
```