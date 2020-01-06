package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestIntcode_Run(t *testing.T) {

	var b bytes.Buffer
	type fields struct {
		program []int
		in      io.Reader
		out     io.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		outWant string
	}{
		{
			name: "d5_p2_example1_position-mode_true",
			fields: fields{
				program: []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
				in:      strings.NewReader("8"),
				out:     &b,
			},
			wantErr: false,
			outWant: "1\n",
		},
		{
			name: "d5_p2_example1_position-mode_false",
			fields: fields{
				program: []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
				in:      strings.NewReader("7"),
				out:     &b,
			},
			wantErr: false,
			outWant: "0\n",
		},
		{
			name: "d5_p2_example2_position-mode_true",
			fields: fields{
				program: []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8},
				in:      strings.NewReader("7"),
				out:     &b,
			},
			wantErr: false,
			outWant: "1\n",
		},
		{
			name: "d5_p2_example2_position-mode_false",
			fields: fields{
				program: []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8},
				in:      strings.NewReader("9"),
				out:     &b,
			},
			wantErr: false,
			outWant: "0\n",
		},
		{
			name: "d5_p2_example3_immediate-mode_true",
			fields: fields{
				program: []int{3, 3, 1108, -1, 8, 3, 4, 3, 99},
				in:      strings.NewReader("8"),
				out:     &b,
			},
			wantErr: false,
			outWant: "1\n",
		},
		{
			name: "d5_p2_example3_immediate-mode_false",
			fields: fields{
				program: []int{3, 3, 1108, -1, 8, 3, 4, 3, 99},
				in:      strings.NewReader("7"),
				out:     &b,
			},
			wantErr: false,
			outWant: "0\n",
		},
		{
			name: "d5_p2_example4_immediate-mode_true",
			fields: fields{
				program: []int{3, 3, 1107, -1, 8, 3, 4, 3, 99},
				in:      strings.NewReader("7"),
				out:     &b,
			},
			wantErr: false,
			outWant: "1\n",
		},
		{
			name: "d5_p2_example4_immediate-mode_false",
			fields: fields{
				program: []int{3, 3, 1107, -1, 8, 3, 4, 3, 99},
				in:      strings.NewReader("9"),
				out:     &b,
			},
			wantErr: false,
			outWant: "0\n",
		},
		{
			name: "d5_p2_example5_jumptest_zero-out",
			fields: fields{
				program: []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
				in:      strings.NewReader("0"),
				out:     &b,
			},
			wantErr: false,
			outWant: "0\n",
		},
		{
			name: "d5_p2_example5_jumptest_non-zero-out",
			fields: fields{
				program: []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
				in:      strings.NewReader("-1"),
				out:     &b,
			},
			wantErr: false,
			outWant: "1\n",
		},
		{
			name: "d5_p2_example6_jumptest_immediate_zero-out",
			fields: fields{
				program: []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1},
				in:      strings.NewReader("0"),
				out:     &b,
			},
			wantErr: false,
			outWant: "0\n",
		},
		{
			name: "d5_p2_example6_jumptest_immediate_non-zero-out",
			fields: fields{
				program: []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1},
				in:      strings.NewReader("-1"),
				out:     &b,
			},
			wantErr: false,
			outWant: "1\n",
		},
		{
			name: "d5_p2_example7_large_example_below_8",
			fields: fields{
				program: []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
					1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
					999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99},
				in:  strings.NewReader("7"),
				out: &b,
			},
			wantErr: false,
			outWant: "999\n",
		},
		{
			name: "d5_p2_example7_large_example_equal_8",
			fields: fields{
				program: []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
					1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
					999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99},
				in:  strings.NewReader("8"),
				out: &b,
			},
			wantErr: false,
			outWant: "1000\n",
		},
		{
			name: "d5_p2_example7_large_example_above_8",
			fields: fields{
				program: []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31,
					1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104,
					999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99},
				in:  strings.NewReader("9"),
				out: &b,
			},
			wantErr: false,
			outWant: "1001\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Intcode{
				program: tt.fields.program,
				in:      tt.fields.in,
				out:     tt.fields.out,
			}
			if err := i.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if output := b.String(); output != tt.outWant {
				t.Errorf("Run() output = %s, wantOut %s", output, tt.outWant)
			}
			b.Reset()
		})
	}
}
