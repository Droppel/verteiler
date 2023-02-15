package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Group struct {
	dummy   bool
	choices []int
}

type Slot struct {
	id     int
	amount int
}

const (
	episodes      = 100000
	lineSeparator = "\n"
)

var (
	penalties = []int{0, -1, -5, -10}
)

func main() {
	options := []Slot{
		{0, 4},
		{1, 4},
		{2, 4},
		{3, 4},
		{4, 4},
		{5, 4},
		{6, 4},
		{7, 4},
	}

	groups := parseChoices()

	// rand.Seed(43)
	rand.Seed(time.Now().UnixNano())
	//Init random
	optionsForRand := make([]int, 0)
	for _, opt := range options {
		for i := 0; i < opt.amount; i++ {
			optionsForRand = append(optionsForRand, opt.id)
		}
	}

	lenRealGroups := len(groups)
	for i := 0; i < len(optionsForRand)-lenRealGroups; i++ {
		groups = append(groups, Group{true, []int{0, 1, 2}})
	}

	bestSolution := make([]int, len(groups))
	for i := range bestSolution {
		selected := rand.Intn(len(optionsForRand))
		bestSolution[i] = optionsForRand[selected]
		optionsForRand = remove(optionsForRand, selected)
	}

	bestScore, _ := calcScore(bestSolution, groups)
	fmt.Println(bestSolution)
	fmt.Println(bestScore)

	for i := 0; i < episodes; i++ {
		solution := make([]int, len(bestSolution))
		copy(solution, bestSolution)
		solution = randSwap(solution)
		solution = randSwap(solution)
		solution = randSwap(solution)
		score, _ := calcScore(solution, groups)
		if score > bestScore {
			bestScore = score
			bestSolution = solution
		}
	}
	fmt.Println(bestSolution)
	fmt.Println(bestScore)
	_, resultSpread := calcScore(bestSolution, groups)
	fmt.Printf("First: %d, Second: %d, Third: %d, None: %d\n", resultSpread[0], resultSpread[1], resultSpread[2], resultSpread[3])
}

func parseChoices() []Group {
	groups := make([]Group, 0)
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	inputString := string(input)
	for _, line := range strings.Split(inputString, lineSeparator) {
		if line[0] == '#' {
			continue
		}
		columns := strings.Split(line, " ")
		group := Group{
			dummy: false,
		}

		choices := make([]int, len(columns))
		for i, column := range columns {
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
	}
	return groups
}

func randSwap(s []int) []int {
	n1 := rand.Intn(len(s))
	n2 := rand.Intn(len(s))

	t := s[n1]
	s[n1] = s[n2]
	s[n2] = t
	return s
}

func calcScore(solution []int, groups []Group) (int, []int) {
	score := 0
	resultSpread := make([]int, len(penalties))
	for i, group := range groups {
		if group.dummy {
			continue
		}
		selectedPenalty := len(penalties) - 1
		for k, choice := range group.choices {
			if solution[i] == choice {
				selectedPenalty = k
			}
		}
		score += penalties[selectedPenalty]
		resultSpread[selectedPenalty] += 1
	}
	return score, resultSpread
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
