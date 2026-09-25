package main

import (
	"errors"
	"fmt"
)

// Calculator demonstrates a struct that carries its own state (Result, History).
type Calculator struct {
	Result  float64
	History []string // in-memory log of every calculation performed
}

// NewCalculator is a constructor function. It returns a POINTER to a Calculator
// so that every part of the program shares and mutates the same underlying struct.
func NewCalculator() *Calculator {
	return &Calculator{History: make([]string, 0)}
}

// Calculate performs +, -, *, / using a switch statement.
//
// It has a POINTER RECEIVER (c *Calculator). If this were a value receiver
// (c Calculator), Go would copy the struct and c.Result = result below would
// only change the copy — the caller would never see the update. With a
// pointer receiver, we mutate the exact struct the caller holds.
func (c *Calculator) Calculate(a, b float64, op string) (float64, error) {
	var result float64

	switch op {
	case "+":
		result = c.add(a, b)
	case "-":
		result = c.subtract(a, b)
	case "*":
		result = c.multiply(a, b)
	case "/":
		if b == 0 {
			return 0, errors.New("division by zero is not allowed")
		}
		result = c.divide(a, b)
	default:
		return 0, fmt.Errorf("unknown operator: %q", op)
	}

	c.Result = result
	entry := fmt.Sprintf("%.2f %s %.2f = %.2f", a, op, b, result)
	c.History = append(c.History, entry)
	return result, nil
}

// Small private helper methods on *Calculator.
func (c *Calculator) add(a, b float64) float64      { return a + b }
func (c *Calculator) subtract(a, b float64) float64 { return a - b }
func (c *Calculator) multiply(a, b float64) float64 { return a * b }
func (c *Calculator) divide(a, b float64) float64   { return a / b }

// LastResult uses a VALUE receiver on purpose: it only reads data, never
// mutates it, so there's no need to pass a pointer.
func (c Calculator) LastResult() float64 {
	return c.Result
}

// ClearHistory demonstrates resetting a slice field through a pointer receiver.
func (c *Calculator) ClearHistory() {
	c.History = make([]string, 0)
}
