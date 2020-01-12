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
	"strings"
)

type Planet struct {
	Name     string
	Orbiting []*Planet
}

func InitPlanet(name string) *Planet {
	planet := Planet{
		Name:     name,
		Orbiting: []*Planet{},
	}
	return &planet
}

func (p *Planet) AddPlanet(planet *Planet) {
	p.Orbiting = append(p.Orbiting, planet)
}

func (p *Planet) Copy() *Planet {
	newplanet := InitPlanet(p.Name)
	for _, orbitingPlanet := range p.Orbiting {
		newplanet.Orbiting = append(newplanet.Orbiting, orbitingPlanet.Copy())
	}
	return newplanet
}

const CenterOfUniverse = "COM"
const StartLabel = "YOU"
const EndLabel = "SAN"

func main() {
	//load in the input values from a file
	if len(os.Args) < 2 {
		log.Fatalln("Usage: <exe> <input_file_of_celestial_orbits>")
	}
	input, err := ReadInput(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	//load in the orbital data
	_, rootOrbit, err := GetOrbits(input)
	if err != nil {
		log.Fatalln(err)
	}

	//find the total number of orbits
	totalOrbits := CalcTotalOrbits(0, rootOrbit)

	fmt.Printf("The total number of orbits is %d\n", totalOrbits)

	//find the number of orbital transfers required to get to Satan
	numOrbitalTransfers, err := GetMinOrbitalTransfers(rootOrbit, StartLabel, EndLabel)
	fmt.Printf("The minimum number of orbits requried to get to %s is %d", EndLabel, numOrbitalTransfers)
}

func ReadInput(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	input := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for {
		stuffWasRead := scanner.Scan()
		if stuffWasRead != true {
			return input, nil
		}
		txt := scanner.Text()
		input = append(input, txt)
	}
}

func GetOrbits(orbitsTxt []string) (map[string]*Planet, *Planet, error) {
	var inner, outer string
	var root *Planet
	lookup := make(map[string]*Planet, 0)
	for _, txt := range orbitsTxt {
		pieces := strings.Split(txt, ")")
		if len(pieces) < 2 {
			break
		}
		inner = pieces[0]
		outer = pieces[1]
		iPlanet, ok := lookup[inner]
		if !ok {
			iPlanet = InitPlanet(inner)
			lookup[inner] = iPlanet
		}
		oPlanet, ok := lookup[outer]
		if !ok {
			oPlanet = InitPlanet(outer)
			lookup[outer] = oPlanet
		}
		iPlanet.AddPlanet(oPlanet)
		if inner == CenterOfUniverse {
			root = iPlanet
		}
	}

	if root == nil {
		return nil, nil, errors.New(fmt.Sprintf("No center of universe found with name '%s'", CenterOfUniverse))
	}
	return lookup, root, nil
}

func CalcTotalOrbits(previousOrbits int, planet *Planet) int {
	temp := previousOrbits
	for _, orbiter := range planet.Orbiting {
		temp += CalcTotalOrbits(previousOrbits+1, orbiter)
	}
	return temp
}

//GetMinOrbitalTransfers calculates the minimum number of Orbital transfers it takes to get from the Orbit of the start planet
//to the orbit of the end planet
//it does this with a DepthFirstSearch, finding the target nodes from the root tree and creating pruned copies of the tree such that
//the only nodes left of the two copies of the root tree will contain nodes directly leading to the root.
//From there, we only need to find two shared nodes within the pruned trees and add the distance to the leaf nodes to get the minimum orbital transfers
func GetMinOrbitalTransfers(root *Planet, start string, end string) (int, error) {
	startChain := root.Copy()
	startChain, err := GetPrunedPlanetChain(startChain, start)
	if err != nil {
		return -1, err
	}
	planetMap := map[string]int{}
	var temp *Planet
	count := 0
	temp = startChain

	for {
		planetMap[temp.Name] = count
		count++
		if len(temp.Orbiting) == 0 {
			break
		}
		temp = temp.Orbiting[0]
	}

	endChain := root.Copy()
	endChain, err = GetPrunedPlanetChain(endChain, end)
	if err != nil {
		return -1, err
	}
	temp = endChain
	var distanceToEnd int
	for {
		if len(temp.Orbiting) == 0 {
			break
		}
		if len(temp.Orbiting) > 1 {
			panic("orbiting count greater than 1, not a pruned tree!")
		}
		if _, ok := planetMap[temp.Orbiting[0].Name]; ok {
			temp = temp.Orbiting[0]
			continue
		}
		//this is the planet that diverges
		//subtract one because we don't include the planet we are starting from
		distanceToEnd = planetMap[start] - planetMap[temp.Name] - 1
		for {
			if len(temp.Orbiting) == 0 {
				break
			}
			if len(temp.Orbiting) > 1 {
				panic("orbiting count greater than 1, not a pruned tree!")
			}
			temp = temp.Orbiting[0]
			distanceToEnd++
		}
		//substract one again for the ending planet which we want to share an orbit with, not be orbiting around
		distanceToEnd--
	}

	return distanceToEnd, nil
}

func GetPrunedPlanetChain(root *Planet, target string) (*Planet, error) {
	var foundPlanet *Planet
	for idx, planet := range root.Orbiting {
		ret, err := GetPrunedPlanetChain(planet, target)
		if err != nil {
			return nil, err
		}
		if ret != nil {
			foundPlanet = ret
			continue
		}
		if planet.Orbiting != nil {
			planet.Orbiting[idx] = nil
		}
	}
	if foundPlanet != nil {
		root.Orbiting = []*Planet{foundPlanet}
		return root, nil
	}
	if root.Name == target {
		return root, nil
	}
	root.Orbiting = nil
	return nil, nil
}
