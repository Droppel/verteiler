package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	episodes      = 1000000
	lineSeparator = "\r\n"
	inputname     = "inputforce.txt"
)

type Group struct {
	id               int
	dummy            bool
	size             int
	members          []string
	choices          []int
	currentSelection int
}

type GroupList []Group

func (a GroupList) Len() int           { return len(a) }
func (a GroupList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a GroupList) Less(i, j int) bool { return a[i].size > a[j].size }

type Slot struct {
	id       int
	capacity int
	amount   int
}

type Solution struct {
	occupancy     []Slot
	groups        []Group
	invAllocation map[int][]int
}

func (s *Solution) Print(score int) {
	output := fmt.Sprintf("Seed: %d\n", seed)
	for _, slot := range s.occupancy {
		output += fmt.Sprintln("====================================")
		output += fmt.Sprintf("GROUP %d-%d\n", (slot.id*4 + 1), (slot.id+1)*4)
		for _, groupId := range s.invAllocation[slot.id] {
			group := s.groups[groupId]
			choice := -1
			for i, ch := range group.choices {
				if ch == group.currentSelection {
					choice = i + 1
					break
				}
			}
			output += fmt.Sprintf("%d. Wahl\n", choice)
			for _, member := range group.members {
				output += fmt.Sprintln(member)
			}
			output += fmt.Sprintln("-------------------------------------")
		}
	}

	fmt.Print(output)
	os.WriteFile(fmt.Sprintf("Score%d.txt", score), []byte(output), os.ModeAppend)
}

func (s *Solution) Copy() Solution {
	copiedSolution := Solution{
		occupancy:     make([]Slot, len(s.occupancy)),
		groups:        make([]Group, len(s.groups)),
		invAllocation: make(map[int][]int),
	}

	copy(copiedSolution.occupancy, s.occupancy)
	copy(copiedSolution.groups, s.groups)
	for k, v := range s.invAllocation {
		newArr := make([]int, len(v))
		copy(newArr, v)
		copiedSolution.invAllocation[k] = newArr
	}
	return copiedSolution
}

var (
	penalties = []int{0, -1, -5, -100}
	idCounter = 0
	seed      int64
)

func main() {
	options := []Slot{
		{0, 24, 0},
		{1, 24, 0},
		{2, 20, 0},
		{3, 24, 0},
		{4, 24, 0},
		{5, 24, 0},
		{6, 18, 0},
		{7, 20, 0},
	}

	groups, totalMembers := parseChoices()

	//Init random
	seed = time.Now().UnixNano()
	seed = 1676582312005698000
	rand.Seed(seed)

	// totalSlots := 0
	// for _, opt := range options {
	// 	totalSlots += opt.capacity
	// }
	_ = totalMembers
	// for i := 0; i < totalSlots-totalMembers; i++ {
	// 	groups = append(groups, Group{idCounter, true, 1, []string{"dummy"}, []int{-1, -1, -1}})
	// 	idCounter += 1
	// }

	bestSolution := Solution{
		occupancy:     make([]Slot, len(options)),
		groups:        make([]Group, len(groups)),
		invAllocation: make(map[int][]int),
	}
	copy(bestSolution.occupancy, options)
	copy(bestSolution.groups, groups)
	sort.Sort(groups)
	for _, group := range groups {
		selected := findPossibleSlot(group.size, bestSolution.occupancy)
		bestSolution.occupancy[selected].amount += group.size
		bestSolution.groups[group.id].currentSelection = selected
		if bestSolution.invAllocation[selected] == nil {
			bestSolution.invAllocation[selected] = []int{group.id}
		} else {
			bestSolution.invAllocation[selected] = append(bestSolution.invAllocation[selected], group.id)
		}
	}

	bestScore, _ := calcScore(bestSolution)
	fmt.Println(bestSolution)
	fmt.Println(bestScore)

	for i := 0; i < episodes; i++ {
		solution := bestSolution.Copy()
		solution.randSwap()
		solution.randSwap()
		solution.randSwap()
		score, _ := calcScore(solution)
		if score > bestScore {
			bestScore = score
			bestSolution = solution
		}
	}
	bestSolution.Print(bestScore)
	fmt.Println(bestScore)
	_, resultSpread := calcScore(bestSolution)
	fmt.Printf("First: %d, Second: %d, Third: %d, None: %d\n", resultSpread[0], resultSpread[1], resultSpread[2], resultSpread[3])
}

func parseChoices() (GroupList, int) {
	groups := make([]Group, 0)
	totalMembers := 0
	input, err := os.ReadFile("inputfiles" + string(filepath.Separator) + inputname)
	if err != nil {
		panic(err)
	}
	inputString := string(input)
	for _, line := range strings.Split(inputString, lineSeparator) {
		if line[0] == '#' {
			continue
		}
		split := strings.Split(line, "|")
		membersString := split[0]
		choicesString := split[1]
		choicesStringSplit := strings.Split(choicesString, " ")

		group := Group{
			id:               idCounter,
			dummy:            false,
			members:          strings.Split(membersString, ","),
			currentSelection: -1,
		}
		idCounter += 1

		group.size = len(group.members)

		choices := make([]int, len(choicesStringSplit))
		for i, column := range choicesStringSplit {
			if column == "?" {
				choices[i] = -1
			} else {
				singleNumberString := strings.Split(column, "-")[1]
				singleNumber, _ := strconv.Atoi(singleNumberString)
				choices[i] = (singleNumber / 4) - 1
			}
		}
		group.choices = choices
		groups = append(groups, group)
		totalMembers += group.size
	}
	return groups, totalMembers
}

type Swap struct {
	slot          int
	swapParteners []int
}

func (s *Solution) randSwap() {
	groups := s.groups
	randGroup := groups[rand.Intn(len(groups))]
	possibleSwaps := make([]Swap, 0)
	for i, slot := range s.occupancy {
		if i == randGroup.currentSelection {
			continue
		}
		if slot.capacity-slot.amount >= randGroup.size {
			possibleSwaps = append(possibleSwaps, Swap{i, make([]int, 0)})
			continue
		}

		solved, set := solveSubsetSum(s, randGroup.size, s.invAllocation[i])
		if !solved {
			continue
		}

		possibleSwaps = append(possibleSwaps, Swap{slot: s.groups[set[0]].currentSelection, swapParteners: set})
	}

	if len(possibleSwaps) == 0 {
		return
	}
	choosenSwap := possibleSwaps[rand.Intn(len(possibleSwaps))]

	randGroupSelection := randGroup.currentSelection
	s.groups[randGroup.id].currentSelection = choosenSwap.slot

	//Update invAllocation and Occupancy
	s.invAllocation[randGroupSelection] = removeElement(s.invAllocation[randGroupSelection], randGroup.id)
	s.invAllocation[choosenSwap.slot] = append(s.invAllocation[choosenSwap.slot], randGroup.id)
	s.occupancy[randGroupSelection].amount -= randGroup.size
	s.occupancy[choosenSwap.slot].amount += randGroup.size
	for _, setGroupId := range choosenSwap.swapParteners {
		setGroup := s.groups[setGroupId]
		s.groups[setGroupId].currentSelection = randGroupSelection
		s.occupancy[randGroupSelection].amount += setGroup.size
		s.occupancy[choosenSwap.slot].amount -= setGroup.size
		s.invAllocation[choosenSwap.slot] = removeElement(s.invAllocation[choosenSwap.slot], setGroupId)
		s.invAllocation[randGroupSelection] = append(s.invAllocation[randGroupSelection], setGroupId)
	}
}

func solveSubsetSum(solution *Solution, sum int, groupIds []int) (bool, []int) {
	for i, groupId := range groupIds {
		summand := solution.groups[groupId].size
		if summand == sum {
			return true, []int{groupId}
		}
		solved, subGroupIds := solveSubsetSum(solution, sum-summand, groupIds[i+1:])
		if solved {
			return true, append([]int{groupId}, subGroupIds...)
		}
	}
	return false, []int{}
}

func calcScore(solution Solution) (int, []int) {
	score := 0
	resultSpread := make([]int, len(penalties))
	for _, group := range solution.groups {
		if group.dummy {
			continue
		}
		selectedPenalty := len(penalties) - 1
		for k, choice := range group.choices {
			if group.currentSelection == choice {
				selectedPenalty = k
			}
		}
		score += penalties[selectedPenalty]
		resultSpread[selectedPenalty] += 1
	}
	return score, resultSpread
}

func findPossibleSlot(size int, solution []Slot) int {
	solCopy := make([]Slot, len(solution))
	copy(solCopy, solution)
	for {
		randIndex := rand.Intn(len(solCopy))
		randSlot := solCopy[randIndex]
		if size < randSlot.capacity-randSlot.amount {
			return randSlot.id
		} else {
			solCopy = remove(solCopy, randIndex)
			if len(solCopy) <= 0 {
				panic("no slots available")
			}
		}
	}
}

func removeElement[T comparable](s []T, e T) []T {
	for i, elem := range s {
		if elem == e {
			return remove(s, i)
		}
	}
	// return false, s
	panic("Tried to remove non existant")
}

func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
