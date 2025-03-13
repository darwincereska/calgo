package calculator

import (
 "testing"
)

func TestCalculate(t *testing.T) {
 tests := []struct {
  input    string
  expected Value
 }{
  {
   "2 + 2",
   Value{Number: 4, Type: NumberType},
  },
  {
   "3 * 4",
   Value{Number: 12, Type: NumberType},
  },
  {
   "10 / 2",
   Value{Number: 5, Type: NumberType},
  },
  {
   "2 ^ 3", 
   Value{Number: 8, Type: NumberType},
  },
  {
   "sqrt(16)",
   Value{Number: 4, Type: NumberType},
  },
  {
   "sin(0)",
   Value{Number: 0, Type: NumberType},
  },
  {
   "# this is a comment",
   Value{Type: CommentType},
  },
  {
   "x = 5",
   Value{Number: 5, Type: NumberType},
  },
 }

 for _, tt := range tests {
  result, err := Calculate(tt.input)
  if err != nil {
   t.Errorf("Calculate(%q) returned unexpected error: %v", tt.input, err)
  }
  if result.Number != tt.expected.Number || result.Type != tt.expected.Type {
   t.Errorf("Calculate(%q) = %v, expected %v", tt.input, result, tt.expected)
  }
 }
}