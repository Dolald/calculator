package main

import (
	"testing"
)

func Test_infixToRPN(t *testing.T) {

	testTable := []struct {
		input    string
		expected float64
		wantErr  bool
	}{
		{
			input:    "( ( 2 - 3 * 3 ) * 2 - ( ( 2 * 1 ) * ( 12 / 4 ) ) ) * 1",
			expected: -20,
			wantErr:  false,
		},
		{
			input:    "2 + 2",
			expected: 4,
			wantErr:  false,
		},
		{
			input:    "1",
			expected: 1,
			wantErr:  false,
		},
		{
			input:    "-1",
			expected: -1,
			wantErr:  false,
		},
		{
			input:    "0",
			expected: 0,
			wantErr:  false,
		},
		{
			input:    "-1 - 0 / 1",
			expected: -1,
			wantErr:  false,
		},
		{
			input:    "-1-2",
			expected: 0,
			wantErr:  true,
		},
		{
			input:    "-10 + 12 )",
			expected: 0,
			wantErr:  true,
		},
		{
			input:    "10 + 12 ( )",
			expected: 0,
			wantErr:  true,
		},
		{
			input:    "-10 + 12 ",
			expected: 0,
			wantErr:  true,
		},
		{
			input:    "ssd + 2",
			expected: 0,
			wantErr:  true,
		},
		{
			input:    " 12 - 1",
			expected: 0,
			wantErr:  true,
		},
		{
			input:    "12 - ( ) 2",
			expected: 0,
			wantErr:  true,
		},
		{
			input:    "12 - ( ) / 2",
			expected: 0,
			wantErr:  true,
		},
		{
			input:    "- ( 9 - 1 )",
			expected: -8,
			wantErr:  false,
		},
		{
			input:    "( 12 ) - 2",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range testTable {
		result, err := infixToRPN(tt.input)
		if err != nil && !tt.wantErr {
			t.Error("говно", err)
		}

		if tt.expected != result && !tt.wantErr {
			t.Error("input: ", tt.input, "expected: ", tt.expected, "result: ", result)
		}
	}
}
