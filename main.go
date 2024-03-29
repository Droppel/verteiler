package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
	"verteiler/genome"
	"verteiler/parser"
)

const (
	inputname = "input.csv"
	slotname  = "slots.csv"
)

var (
	randSeed      = true
	algoSeed      = 42
	episodes      = 100000
	threads       = 6
	cutoffForSave = -100

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

	seed := 0
	//Parse arguments
	for i, arg := range os.Args {
		if rune(arg[0]) == '-' {
			switch rune(arg[1]) {
			case 's':
				seedString := os.Args[i+1]
				seed, err = strconv.Atoi(seedString)
				if err != nil {
					panic(fmt.Errorf("failed to parse seed: %w", err))
				}
				randSeed = false
			case 'e':
				episodesString := os.Args[i+1]
				episodes, err = strconv.Atoi(episodesString)
				if err != nil {
					panic(fmt.Errorf("failed to parse episode depth: %w", err))
				}
			case 't':
				threadsString := os.Args[i+1]
				threads, err = strconv.Atoi(threadsString)
				if err != nil {
					panic(fmt.Errorf("failed to parse thread count: %w", err))
				}
			case 'c':
				cutoffForSaveString := os.Args[i+1]
				cutoffForSave, err = strconv.Atoi(cutoffForSaveString)
				if err != nil {
					panic(fmt.Errorf("failed to parse cutoff number: %w", err))
				}
			}
		}
	}

	if randSeed {
		for i := 0; i < threads; i++ {
			go calcSeedLoop(i, groups, options)
			time.Sleep(10 * time.Millisecond)
		}
		<-make(chan int)
	} else {
		algoSeed = seed
		calcSeed(0, groups, options)
	}
}

func calcSeedLoop(threadId int, groupsInput genome.GroupList, optionsInput []genome.Slot) {
	for {
		calcSeed(threadId, groupsInput, optionsInput)
	}
}

func calcSeed(threadId int, groupsInput genome.GroupList, optionsInput []genome.Slot) {
	groups := make(genome.GroupList, len(groupsInput))
	copy(groups, groupsInput)
	options := make([]genome.Slot, len(optionsInput))
	copy(options, optionsInput)

	//Init random
	seed := int64(algoSeed)
	if randSeed {
		seed = time.Now().UnixNano()
	}
	fmt.Printf("Thread %d: Running with seed %d\n", threadId, seed)
	randSource := rand.New(rand.NewSource(seed))

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
		selected := findPossibleSlot(group.Size, bestSolution.Occupancy, randSource)
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
		solution.RandSwap(randSource)
		solution.RandSwap(randSource)
		solution.RandSwap(randSource)
		score, _ := calcScore(solution)
		if score > bestScore {
			bestScore = score
			bestSolution = solution
		}
	}
	solutionString := bestSolution.ToString(bestScore, seed)

	// fmt.Print(solutionString)
	if bestScore >= cutoffForSave {
		err := os.MkdirAll("scores", os.ModeAppend)
		if err != nil {
			panic(fmt.Errorf("failed to create folder: %w", err))
		}
		err = os.WriteFile(fmt.Sprintf("scores/Score%d-%d.txt", bestScore, seed), []byte(solutionString), os.ModeAppend)
		if err != nil {
			panic(fmt.Errorf("failed to store files to disk: %w", err))
		}
	}
	output := fmt.Sprintf("Thread %d: Running with seed %d finished:\n", threadId, seed)
	output += fmt.Sprintf("%v\n", bestScore)
	_, resultSpread := calcScore(bestSolution)
	output += fmt.Sprintf("First: %d, Second: %d, Third: %d, None: %d\n", resultSpread[0], resultSpread[1], resultSpread[2], resultSpread[3])
	fmt.Println(output)
}

func calcScore(solution genome.Solution) (int, []int) {
	score := 0
	resultSpread := make([]int, len(penalties))
	for _, group := range solution.Groups {
		selectedPenalty := len(penalties) - 1
		currentGroupTimeSlot := solution.Occupancy[group.CurrentSelection].TimeSlot
		for k, choice := range group.Choices {
			if choice == -1 || currentGroupTimeSlot == choice {
				selectedPenalty = k
			}
		}
		score += penalties[selectedPenalty]
		resultSpread[selectedPenalty] += 1
	}
	return score, resultSpread
}

func findPossibleSlot(size int, solution []genome.Slot, randSource *rand.Rand) int {
	solCopy := make([]genome.Slot, len(solution))
	copy(solCopy, solution)
	for {
		randSource.Int63()
		randIndex := randSource.Intn(len(solCopy))
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
