package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
	"verteiler/genome"
	"verteiler/parser"
)

const (
	randSeed = true
	seed     = 42

	threads   = 1
	episodes  = 1000000
	inputname = "input.csv"
	slotname  = "slots.csv"
)

var (
	penalties = []int{0, -1, -5, -100}
)

func main() {

	options, err := parser.ParseSlots(slotname)
	if err != nil {
		fmt.Printf("Failed to parse slots: %v\n", err)
		return
	}

	groups, err := parser.ParseChoices(inputname)
	if err != nil {
		fmt.Printf("Failed to parse choices: %v\n", err)
		return
	}

	for i := 0; i < threads; i++ {
		go calcSeed(i, groups, options)
		time.Sleep(10 * time.Millisecond)
	}
	<-make(chan int)
}

func calcSeed(threadId int, groupsInput genome.GroupList, optionsInput []genome.Slot) {
	for {
		groups := make(genome.GroupList, len(groupsInput))
		copy(groups, groupsInput)
		options := make([]genome.Slot, len(optionsInput))
		copy(options, optionsInput)

		//Init random
		seed := int64(seed)
		if randSeed {
			seed = time.Now().UnixNano()
		}
		fmt.Printf("Thread %d: Running with seed %d\n", threadId, seed)
		rand.Seed(seed)

		// Create initial random solution
		bestSolution := genome.Solution{
			Occupancy:     make([]genome.Slot, len(options)),
			Groups:        make([]genome.Group, len(groups)),
			InvAllocation: make(map[int][]int),
		}
		copy(bestSolution.Occupancy, options)
		copy(bestSolution.Groups, groups)
		sort.Sort(groups) // Sort groups by size
		for _, group := range groups {
			selected := findPossibleSlot(group.Size, bestSolution.Occupancy)
			bestSolution.Occupancy[selected].Amount += group.Size
			bestSolution.Groups[group.Id].CurrentSelection = selected
			if bestSolution.InvAllocation[selected] == nil {
				bestSolution.InvAllocation[selected] = []int{group.Id}
			} else {
				bestSolution.InvAllocation[selected] = append(bestSolution.InvAllocation[selected], group.Id)
			}
		}

		bestScore, _ := calcScore(bestSolution)

		for i := 0; i < episodes; i++ {
			solution := bestSolution.Copy()
			solution.RandSwap()
			solution.RandSwap()
			solution.RandSwap()
			score, _ := calcScore(solution)
			if score > bestScore {
				bestScore = score
				bestSolution = solution
			}
		}
		solutionString := bestSolution.ToString(bestScore, seed)

		fmt.Print(solutionString)
		if bestScore >= -20 {
			os.WriteFile(fmt.Sprintf("scores/Score%d-%d.txt", bestScore, seed), []byte(solutionString), os.ModeAppend)
		}
		fmt.Println(bestScore)
		_, resultSpread := calcScore(bestSolution)
		fmt.Printf("First: %d, Second: %d, Third: %d, None: %d\n", resultSpread[0], resultSpread[1], resultSpread[2], resultSpread[3])
	}
}

func calcScore(solution genome.Solution) (int, []int) {
	score := 0
	resultSpread := make([]int, len(penalties))
	for _, group := range solution.Groups {
		selectedPenalty := len(penalties) - 1
		for k, choice := range group.Choices {
			if solution.Occupancy[group.CurrentSelection].TimeSlot == choice {
				selectedPenalty = k
			}
		}
		score += penalties[selectedPenalty]
		resultSpread[selectedPenalty] += 1
	}
	return score, resultSpread
}

func findPossibleSlot(size int, solution []genome.Slot) int {
	solCopy := make([]genome.Slot, len(solution))
	copy(solCopy, solution)
	for {
		randIndex := rand.Intn(len(solCopy))
		randSlot := solCopy[randIndex]
		if size <= randSlot.Capacity-randSlot.Amount {
			return randSlot.Id
		} else {
			solCopy = genome.Remove(solCopy, randIndex)
			if len(solCopy) <= 0 {
				panic("no slots available")
			}
		}
	}
}
