// Copyright 2019 Adknown Inc. All rights reserved.
// Created:  2019-12-27
// Author:   matt
// Project:  aoc-2019

package main

import "testing"

func TestAttemptIsValid(t *testing.T) {
	type args struct {
		attempt int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "p2_c1",
			args: args{
				attempt: 112233,
			},
			want: true,
		},
		{
			name: "p2_c2",
			args: args{
				attempt: 123444,
			},
			want: false,
		},
		{
			name: "",
			args: args{
				attempt: 111122,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AttemptIsValid(tt.args.attempt); got != tt.want {
				t.Errorf("AttemptIsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
