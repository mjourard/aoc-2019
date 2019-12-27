// Copyright 2019 Adknown Inc. All rights reserved.
// Created:  2019-12-26
// Author:   matt
// Project:  aoc-2019

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	//load in the input values from a file
	if len(os.Args) < 2 {
		log.Fatalln("Usage: <exe> <input_file_of_limits>")
	}
	lower, upper, err := GetLimits(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

}

func GetLimits(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return -1, -1, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	txt := scanner.Text()
	if len(txt) < 1 {
		return -1, -1, errors.New("first line is empty in passed in file")
	}
	lower, err := strconv.Atoi(txt)
	if err != nil {
		return -1, -1, errors.New(fmt.Sprintf("Unable to convert string (%s) to number", txt))
	}
	scanner.Scan()
	txt = scanner.Text()
	if len(txt) < 1 {
		return -1, -1, errors.New("second line is empty in passed in file")
	}
	upper, err := strconv.Atoi(txt)
	if err != nil {
		return -1, -1, errors.New(fmt.Sprintf("Unable to convert string (%s) to number", txt))
	}

	return lower, upper, nil
}

func AttemptIsValid(attempt int) bool {
	lastDigit := 10
	secondLastDigit := 10
	hasDupe, neverDecrease := false, true
	//we check for if its never decreasing from left to right by checking if it is always decreases or stays the same right to left
	for {
		digit := attempt % 10
		attempt /= 10
		if digit == lastDigit && digit != secondLastDigit && digit != attempt%10 {
			hasDupe = true
		}
		if digit > lastDigit {
			neverDecrease = false
		}
		secondLastDigit = lastDigit
		lastDigit = digit

		if attempt == 0 {
			break
		}
	}

	return hasDupe && neverDecrease
}
