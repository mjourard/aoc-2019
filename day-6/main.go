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
	//TODO: do the same here for the end chain
	startChain, err := GetPrunedPlanetChain(startChain, start)
	if err != nil {

	}
	return 0, nil
}

func GetPrunedPlanetChain(root *Planet, target string) (*Planet, error) {
	for idx, planet := range root.Orbiting {
		ret, err := GetPrunedPlanetChain(planet, target)
		if err != nil {
			return nil, err
		}
		if ret != nil {
			return root, nil
		}
		planet.Orbiting[idx] = nil
	}
	if root.Name == target {
		return root, nil
	}
	root.Orbiting = nil
	return nil, nil
}
