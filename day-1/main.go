package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	//read in the file that contains the input
	if len(os.Args) < 2 {
		panic("Usage: <exe> <input_file_of_masses>")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	totalFuelSum := 0
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) < 1 {
			continue
		}
		mass, err := strconv.Atoi(txt)
		if err != nil {
			panic(err)
		}
		initFuel := getRequiredFuel(mass)
		compensatedFuel := initFuel
		getRequiredFuelRecurse(initFuel, &compensatedFuel)
		totalFuelSum += compensatedFuel
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("total sum of fuel requirements: %d\n", totalFuelSum)
}

func getRequiredFuel(mass int) int {
	return (mass / 3) - 2
}

func getRequiredFuelRecurse(mass int, totalMass *int) {
	calc := (mass / 3) - 2
	if calc > 0 {
		*totalMass += calc
	}
	fmt.Printf("fuel mass: %d, fuel for fuel: %d, new total: %d\n", mass, calc, *totalMass)
	if calc > 2 {
		getRequiredFuelRecurse(calc, totalMass)
	}
}
