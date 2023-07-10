package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
	"verteiler/genome"
	"verteiler/parser"
)

const (
	threads       = 10
	episodes      = 1000000
	lineSeparator = "\r\n"
	inputname     = "input.csv"
	maxgroupsize  = 6
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
		options := []genome.Slot{
			{Id: 0, TimeSlot: 0, Capacity: 6, Amount: 0},
			{Id: 1, TimeSlot: 0, Capacity: 6, Amount: 0},
			{Id: 2, TimeSlot: 0, Capacity: 6, Amount: 0},
			{Id: 3, TimeSlot: 0, Capacity: 6, Amount: 0},
			{Id: 4, TimeSlot: 1, Capacity: 6, Amount: 0},
			{Id: 5, TimeSlot: 1, Capacity: 6, Amount: 0},
			{Id: 6, TimeSlot: 1, Capacity: 6, Amount: 0},
			{Id: 7, TimeSlot: 1, Capacity: 6, Amount: 0},
			{Id: 8, TimeSlot: 2, Capacity: 6, Amount: 0},
			{Id: 9, TimeSlot: 2, Capacity: 6, Amount: 0},
			{Id: 10, TimeSlot: 2, Capacity: 2, Amount: 0},
			{Id: 11, TimeSlot: 2, Capacity: 6, Amount: 0},
			{Id: 12, TimeSlot: 3, Capacity: 6, Amount: 0},
			{Id: 13, TimeSlot: 3, Capacity: 6, Amount: 0},
			{Id: 14, TimeSlot: 3, Capacity: 6, Amount: 0},
			{Id: 15, TimeSlot: 3, Capacity: 6, Amount: 0},
			{Id: 16, TimeSlot: 4, Capacity: 6, Amount: 0},
			{Id: 17, TimeSlot: 4, Capacity: 6, Amount: 0},
			{Id: 18, TimeSlot: 4, Capacity: 6, Amount: 0},
			{Id: 19, TimeSlot: 4, Capacity: 6, Amount: 0},
			{Id: 20, TimeSlot: 5, Capacity: 6, Amount: 0},
			{Id: 21, TimeSlot: 5, Capacity: 6, Amount: 0},
			{Id: 22, TimeSlot: 5, Capacity: 6, Amount: 0},
			{Id: 23, TimeSlot: 5, Capacity: 6, Amount: 0},
			{Id: 24, TimeSlot: 6, Capacity: 6, Amount: 0},
			{Id: 25, TimeSlot: 6, Capacity: 6, Amount: 0},
			{Id: 26, TimeSlot: 6, Capacity: 0, Amount: 0},
			{Id: 27, TimeSlot: 6, Capacity: 6, Amount: 0},
			{Id: 28, TimeSlot: 7, Capacity: 2, Amount: 0},
			{Id: 29, TimeSlot: 7, Capacity: 6, Amount: 0},
			{Id: 30, TimeSlot: 7, Capacity: 6, Amount: 0},
			{Id: 31, TimeSlot: 7, Capacity: 6, Amount: 0},
		}

		groups, err := parser.ParseChoices(inputname, lineSeparator, maxgroupsize)
		if err != nil {
			fmt.Printf("Failed to parse choices: %v\n", err)
			return
		}

		//Init random
		seed := time.Now().UnixNano()
		// seed = 1676625230072189000
		fmt.Printf("Thread %d: Running with seed %d\n", threadId, seed)
		rand.Seed(seed)

		bestSolution := genome.Solution{
			Occupancy:     make([]genome.Slot, len(options)),
			Groups:        make([]genome.Group, len(groups)),
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
