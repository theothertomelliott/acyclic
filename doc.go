// Package acyclic provides the ability to quickly check for cycles within a
// data structure before attempting to marshal to JSON or similar formats.
//
// A set of functions are also provided to print these structures in a safe manner
// that won't result in a stack overflow, pruning branches containing cycles and
// clearly marking where they occurred.
//
// Cycles are detected using depth first search.
package acyclic
