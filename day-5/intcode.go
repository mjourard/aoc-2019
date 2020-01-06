// Copyright 2019 Adknown Inc. All rights reserved.
// Created:  2019-12-27
// Author:   matt
// Project:  aoc-2019

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type Intcode struct {
	program []int
	in      io.Reader
	out     io.Writer
}

func Init(program []int, in io.Reader, out io.Writer) *Intcode {
	return &Intcode{
		program: program,
		in:      in,
		out:     out,
	}
}

func (i *Intcode) Run() error {
	pos := 0
	var opcode int
	var err error
	for {
		opcode, pos, err = i.HandleInstruction(pos)
		if err != nil {
			return errors.Wrap(err, "error encountered at position "+string(pos))
		}
		if opcode == 99 {
			break
		}
	}
	return nil
}

//HandleInstruction handles a single opcode instruction. Takes in the i.program, the current position and the current opcode.
//It will return the opcode that was processed and the new position of the i.program
func (i *Intcode) HandleInstruction(pos int) (int, int, error) {
	temp := i.program[pos]
	if temp == 99 {
		return 99, pos, nil
	}
	opcode := temp % 100
	temp /= 100
	par1 := temp % 10
	temp /= 10
	par2 := temp % 10
	temp /= 10
	//par3 := temp % 10

	var loc1, loc2, loc3 int
	loc1 = i.program[pos+1]
	if par1 == 1 {
		loc1 = pos + 1
	}
	loc2 = i.program[pos+2]
	if par2 == 1 {
		loc2 = pos + 2
	}
	//parameters that an instruction writes to will never be in immediate mode
	//for now, this can stay as is
	if len(i.program) > pos+3 {
		loc3 = i.program[pos+3]
	}
	switch opcode {
	case 1:
		i.program[loc3] = i.program[loc1] + i.program[loc2]
		pos += 4
	case 2:
		i.program[loc3] = i.program[loc1] * i.program[loc2]
		pos += 4
	case 3:
		//takes a single integer as input and saves it to the position given by its only parameter
		input := make([]byte, 15)
		bytesRead, err := i.in.Read(input)
		if err != nil {
			return -1, -1, errors.Wrap(err, fmt.Sprintf("error reading input at position %d", pos))
		}
		if bytesRead == 0 {
			return -1, -1, errors.New(fmt.Sprintf("did not read any bytes from input when requested at position %d", pos))
		}
		var val int
		conversions, err := fmt.Sscanf(string(input), "%d", &val)
		if conversions == 0 {
			return -1, -1, errors.New(fmt.Sprintf("unable to parse an integer from saved input at position %d. Recorded input was %s", pos, string(input)))
		}
		if err != nil {
			return -1, -1, errors.Wrap(err, fmt.Sprintf("error parsing an integer from saved input at position %d. Recorded input was %s", pos, string(input)))
		}
		i.program[loc1] = val
		pos += 2
	case 4:
		//outputs the value of its only parameter
		_, err := i.out.Write([]byte(fmt.Sprintf("%d\n", i.program[loc1])))
		if err != nil {
			return -1, -1, errors.Wrap(err, fmt.Sprintf("io error: unable to write value %d to output", i.program[loc1]))
		}
		pos += 2
	case 5:
		//jump-if-true: if first param is non-zero, sets instruction pointer to value at second parameter. Otherwise does nothing
		pos += 3
		if i.program[loc1] != 0 {
			pos = i.program[loc2]
		}
	case 6:
		//jump-if-false: if first param is zero, sets instruction pointer to value at second parameter. Otherwise, does nothing
		pos += 3
		if i.program[loc1] == 0 {
			pos = i.program[loc2]
		}
	case 7:
		//less than: if the first param is less than the second param, it stores 1 in the position given by the third parameter. Otherwise, stores 0
		valToStore := 0
		if i.program[loc1] < i.program[loc2] {
			valToStore = 1
		}
		i.program[loc3] = valToStore
		pos += 4
	case 8:
		//equals: if first param is equal to second param, store 1 at position given by third parameter. Otherwise, store 0
		valToStore := 0
		if i.program[loc1] == i.program[loc2] {
			valToStore = 1
		}
		i.program[loc3] = valToStore
		pos += 4
	}
	return opcode, pos, nil
}
