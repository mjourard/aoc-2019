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
	if len(os.Args) < 3 {
		panic("Usage: <exe> <input_file_of_intcode_program> <input_to_program>")
	}
	//load the program
	program, err := LoadIntCodeProgram(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	//create the io values for the program
	in, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	out := os.Stdout
	intcode := Init(program, in, out)

	err = intcode.Run()
	if err != nil {
		log.Fatalln(err)
	}
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
