package example

import "testing"

func TestAdd(t *testing.T) {
	calc := &Calculator{}
	result := calc.Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", result)
	}
}

func TestSubtract(t *testing.T) {
	calc := &Calculator{}
	result := calc.Subtract(5, 3)
	if result != 2 {
		t.Errorf("Subtract(5, 3) = %d; want 2", result)
	}
}

func TestMultiply(t *testing.T) {
	calc := &Calculator{}
	result := calc.Multiply(4, 5)
	if result != 20 {
		t.Errorf("Multiply(4, 5) = %d; want 20", result)
	}
}

// Note: Divide, Power, and Factorial are not tested
// This will create uncovered lines in the coverage report
