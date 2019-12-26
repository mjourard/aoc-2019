// Copyright 2019 Adknown Inc. All rights reserved.
// Created:  2019-12-25
// Author:   matt
// Project:  aoc-2019

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const TargetOutput = 19690720

func main() {
	//read in the file that contains the input
	if len(os.Args) < 2 {
		panic("Usage: <exe> <input_file_of_masses>")
	}
	program, err := LoadIntCodeProgram(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	var i, j int
	found := false

	for i = 0; i <= 99; i++ {
		for j = 0; j <= 99; j++ {
			//clone the program
			clone := make([]int, len(program))
			copy(clone, program)
			//run the program with the noun and verb combination
			out, err := RunIntCodeProgram(clone, i, j)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("n: %d, v: %d = %d\n", i, j, out)
			if out == TargetOutput {
				fmt.Println("Found target output!! Halting...")
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	fmt.Printf("The noun %d and the verb %d produce %d\nThe final value of 100 * noun + verb = %d\n", i, j, TargetOutput, 100*i+j)
}

func LoadIntCodeProgram(filename string) ([]int, error) {
	csvfile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	// Iterate through the records
	// Read each record from csv
	tape, err := r.Read()
	if err == io.EOF {
		return nil, errors.New("input file did not contain an Intcode program")
	}
	if err != nil {
		return nil, err
	}
	intcodes := make([]int, 0)
	for _, opcode := range tape {
		curInt, err := strconv.Atoi(opcode)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to convert opcode '%s' to a number", opcode))
		}
		intcodes = append(intcodes, curInt)
	}
	return intcodes, nil
}

func RunIntCodeProgram(program []int, noun int, verb int) (int, error) {
	//load the noun and verb
	program[1] = noun
	program[2] = verb

	//run the program
	pos := 0
	var opcode, loc1, loc2, loc3 int
	for {
		if pos >= len(program) {
			return -1, errors.New("ended up at a position that is greater than the length of the passed in program")
		}
		opcode = program[pos]
		//immediately halt on opcode 99
		if opcode == 99 {
			break
		}
		loc1 = program[pos+1]
		loc2 = program[pos+2]
		loc3 = program[pos+3]
		switch opcode {
		case 1:
			program[loc3] = program[loc1] + program[loc2]
		case 2:
			program[loc3] = program[loc1] * program[loc2]
		}
		pos += 4
	}

	return program[0], nil
}
