// Copyright 2019 Adknown Inc. All rights reserved.
// Created:  2019-12-26
// Author:   matt
// Project:  aoc-2019

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	//read in the instructions
	if len(os.Args) < 2 {
		panic("Usage: <exe> <input_file_of_wire_locations>")
	}
	wires, err := LoadWireLocations(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	//iterate over the instructions, populating a map with the UR coordinates that the wire goes through
	wireMap := make(map[string][]int, 0)
	err = PopulateWireMap(wireMap, wires)
	if err != nil {
		log.Fatalln(err)
	}

	//answer part 1
	shortest, err := GetClosestIntersectionByManhattan(wireMap)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("The closest intersection to the origin point has a Manhattan distance of %d\n", shortest)

	//answer part 2
	fewest, err := GetFewestStepsToIntersection(wireMap)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("The fewest combined steps it takes to get to an intersection is %d\n", fewest)
}

func LoadWireLocations(filename string) ([][]string, error) {
	csvfile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	wires := make([][]string, 0)
	wireCount := 0

	// Iterate through the records
	// Read each record from csv
	for {
		wire, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error reading wire on line %d", wireCount+1))
		}
		wires = append(wires, wire)
	}
	return wires, nil
}

func PopulateWireMap(wireMap map[string][]int, wires [][]string) error {
	addLocations := func(wireMap map[string]int, upFromOrigin *int, rightFromOrigin *int, dirToChange *int, mag int, increase bool, distanceTraveled *int) {
		for i := 0; i < mag; i++ {
			if increase {
				*dirToChange += 1
			} else {
				*dirToChange -= 1
			}
			*distanceTraveled++
			key := fmt.Sprintf("U%dR%d", *upFromOrigin, *rightFromOrigin)
			if _, ok := wireMap[key]; !ok {
				wireMap[key] = *distanceTraveled
			}
		}
	}

	for wireIdx, wire := range wires {
		//keep track of the distance from the origin
		upFromOrigin, rightFromOrigin, distance := 0, 0, 0
		localMap := make(map[string]int, 0)
		for lengthIdx, length := range wire {
			//split the direction and magnitude
			dir := length[:1]
			mag, err := strconv.Atoi(length[1:])
			if err != nil {
				return errors.New(fmt.Sprintf("Unknown magnitude (%s) encountered at wire %d, length %d", length[1:len(length)-1], wireIdx, lengthIdx))
			}
			switch dir {
			case "U":
				addLocations(localMap, &upFromOrigin, &rightFromOrigin, &upFromOrigin, mag, true, &distance)
			case "R":
				addLocations(localMap, &upFromOrigin, &rightFromOrigin, &rightFromOrigin, mag, true, &distance)
			case "D":
				addLocations(localMap, &upFromOrigin, &rightFromOrigin, &upFromOrigin, mag, false, &distance)
			case "L":
				addLocations(localMap, &upFromOrigin, &rightFromOrigin, &rightFromOrigin, mag, false, &distance)
			default:
				return errors.New(fmt.Sprintf("Unknown direction (%s) encountered at wire %d, length %d", dir, wireIdx, lengthIdx))
			}
		}

		//add all new positions to the map
		for k, distance := range localMap {
			if _, ok := wireMap[k]; !ok {
				wireMap[k] = make([]int, 0)
			}
			wireMap[k] = append(wireMap[k], distance)
		}
	}
	return nil
}

func GetClosestIntersectionByManhattan(wireMap map[string][]int) (int, error) {
	//last time, iterate over the map of wire locations and find the intersection with the shortest combined distance
	shortest := int(^uint(0) >> 1)
	for k, stepCounts := range wireMap {
		if len(stepCounts) < 2 {
			continue
		}
		//split the coordinates
		coordinates := strings.Split(k, "R")
		distance := 0
		mag := coordinates[0][1:]
		num, err := strconv.Atoi(mag)
		if err != nil {
			return -1, errors.Wrap(err, fmt.Sprintf("Unable to convert magnitude (%s) from coordinate (%s) to number", mag, coordinates[0]))
		}
		distance += int(math.Abs(float64(num)))
		mag = coordinates[1][1:]
		num, err = strconv.Atoi(mag)
		if err != nil {
			return -1, errors.Wrap(err, fmt.Sprintf("Unable to convert magnitude (%s) from coordinate (%s) to number", mag, coordinates[1]))
		}
		distance += int(math.Abs(float64(num)))

		if distance < shortest {
			shortest = distance
		}
	}

	return shortest, nil
}

func GetFewestStepsToIntersection(wireMap map[string][]int) (int, error) {
	//last time, iterate over the map of wire locations and find the intersection with the fewest combined distance
	fewest := int(^uint(0) >> 1)
	for _, stepCounts := range wireMap {
		if len(stepCounts) < 2 {
			continue
		}

		steps := 0
		for _, val := range stepCounts {
			steps += val
		}

		if steps < fewest {
			fewest = steps
		}
	}

	return fewest, nil
}
