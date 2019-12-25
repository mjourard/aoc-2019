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

func main() {
	//read in the file that contains the input
	if len(os.Args) < 2 {
		panic("Usage: <exe> <input_file_of_masses>")
	}
	intcode, err := LoadIntCodeProgram(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	pos := 0
	var opcode, loc1, loc2, loc3 int
	for {
		if pos >= len(intcode) {
			break
		}
		opcode = intcode[pos]
		//immediately halt on opcode 99
		if opcode == 99 {
			break
		}
		loc1 = intcode[pos+1]
		loc2 = intcode[pos+2]
		loc3 = intcode[pos+3]
		switch opcode {
		case 1:
			intcode[loc3] = intcode[loc1] + intcode[loc2]
		case 2:
			intcode[loc3] = intcode[loc1] * intcode[loc2]
		}
		pos += 4
	}
	fmt.Printf("Final value at position 0 is %d", intcode[0])
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
