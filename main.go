package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
	"verteiler/datastructures"
	"verteiler/parser"
)

const (
	threads       = 5
	episodes      = 1000000
	lineSeparator = "\r\n"
	inputname     = "inputforce.txt"
)

var (
	penalties = []int{0, -1, -5, -100}
)

func main() {
	for i := 0; i < threads; i++ {
		go calcSeed(i)
		time.Sleep(10 * time.Millisecond)
	}
	<-make(chan int)
}

func calcSeed(threadId int) {
	for {
		options := []datastructures.Slot{
			{Id: 0, Capacity: 24, Amount: 0},
			{Id: 1, Capacity: 24, Amount: 0},
			{Id: 2, Capacity: 20, Amount: 0},
			{Id: 3, Capacity: 24, Amount: 0},
			{Id: 4, Capacity: 24, Amount: 0},
			{Id: 5, Capacity: 24, Amount: 0},
			{Id: 6, Capacity: 18, Amount: 0},
			{Id: 7, Capacity: 20, Amount: 0},
		}

		groups := parser.ParseChoices(inputname, lineSeparator)

		//Init random
		seed := time.Now().UnixNano()
		// seed = 1676625230072189000
		fmt.Printf("%d: Running with seed %d\n", threadId, seed)
		rand.Seed(seed)

		bestSolution := datastructures.Solution{
			Occupancy:     make([]datastructures.Slot, len(options)),
			Groups:        make([]datastructures.Group, len(groups)),
			InvAllocation: make(map[int][]int),
		}
		copy(bestSolution.Occupancy, options)
		copy(bestSolution.Groups, groups)
		sort.Sort(groups)
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
		bestSolution.Print(bestScore, seed)
		fmt.Println(bestScore)
		_, resultSpread := calcScore(bestSolution)
		fmt.Printf("First: %d, Second: %d, Third: %d, None: %d\n", resultSpread[0], resultSpread[1], resultSpread[2], resultSpread[3])
	}
}

func calcScore(solution datastructures.Solution) (int, []int) {
	score := 0
	resultSpread := make([]int, len(penalties))
	for _, group := range solution.Groups {
		if group.Dummy {
			continue
		}
		selectedPenalty := len(penalties) - 1
		for k, choice := range group.Choices {
			if group.CurrentSelection == choice {
				selectedPenalty = k
			}
		}
		score += penalties[selectedPenalty]
		resultSpread[selectedPenalty] += 1
	}
	return score, resultSpread
}

func findPossibleSlot(size int, solution []datastructures.Slot) int {
	solCopy := make([]datastructures.Slot, len(solution))
	copy(solCopy, solution)
	for {
		randIndex := rand.Intn(len(solCopy))
		randSlot := solCopy[randIndex]
		if size <= randSlot.Capacity-randSlot.Amount {
			return randSlot.Id
		} else {
			solCopy = datastructures.Remove(solCopy, randIndex)
			if len(solCopy) <= 0 {
				panic("no slots available")
			}
		}
	}
}
