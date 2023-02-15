package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Choice struct {
	dummy   bool
	choice1 int
	choice2 int
	choice3 int
}

type Slot struct {
	id     int
	amount int
}

const (
	episodes = 100000
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

	choices := parseChoices()

	// choices := []Choice{
	// 	{false, 2, 5, 6},
	// 	{false, 3, 6, 7},
	// 	{false, 7, 2, 0},
	// 	{false, 0, 2, 6},
	// 	{false, 7, 2, 4},
	// 	{false, 3, 4, 5},
	// 	{false, 3, 4, 5},
	// 	{false, 1, 2, 0},
	// 	{false, 7, 0, 2},
	// 	{false, 2, 5, 7},
	// 	{false, 1, 0, 2},
	// 	{false, 3, 4, 5},
	// 	{false, 7, 3, -1},
	// 	{false, 5, 0, 2},
	// 	{false, 3, 4, 5},
	// 	{false, 0, 5, 7},
	// 	{false, 6, 2, 3},
	// 	{false, 0, 4, 3},
	// 	{false, 0, 2, 1},
	// 	{false, 3, 4, 5},
	// 	{false, 2, 0, 7},
	// 	{false, 6, 5, 0},
	// 	{false, 3, 6, 7},
	// 	{false, 0, 5, 7},
	// 	{false, 6, 1, 0},
	// 	{false, 6, 0, 3},
	// }

	// rand.Seed(43)
	rand.Seed(time.Now().UnixNano())
	//Init random
	optionsForRand := make([]int, 0)
	for _, opt := range options {
		for i := 0; i < opt.amount; i++ {
			optionsForRand = append(optionsForRand, opt.id)
		}
	}

	for i := 0; i < len(optionsForRand)-len(choices); i++ {
		choices = append(choices, Choice{true, 0, 1, 2})
	}

	bestSolution := make([]int, len(choices))
	for i := range bestSolution {
		selected := rand.Intn(len(optionsForRand))
		bestSolution[i] = optionsForRand[selected]
		optionsForRand = remove(optionsForRand, selected)
	}

	bestScore, _, _, _, _ := calcScore(bestSolution, choices)
	fmt.Println(bestSolution)
	fmt.Println(bestScore)

	for i := 0; i < episodes; i++ {
		solution := make([]int, len(bestSolution))
		copy(solution, bestSolution)
		solution = randSwap(solution)
		solution = randSwap(solution)
		solution = randSwap(solution)
		score, _, _, _, _ := calcScore(solution, choices)
		if score > bestScore {
			bestScore = score
			bestSolution = solution
		}
	}
	fmt.Println(bestSolution)
	fmt.Println(bestScore)
	_, first, second, third, none := calcScore(bestSolution, choices)
	fmt.Printf("First: %d, Second: %d, Third: %d, None: %d\n", first, second, third, none)
}

func parseChoices() []Choice {
	choices := make([]Choice, 0)
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	inputString := string(input)
	for _, line := range strings.Split(inputString, "\r\n") {
		if line[0] == '#' {
			continue
		}
		cs := strings.Split(line, " ")
		choice := Choice{
			dummy: false,
		}

		if cs[0] == "?" {
			choice.choice1 = -1
		} else {
			singleNumberString := strings.Split(cs[0], "-")[1]
			singleNumber, _ := strconv.Atoi(singleNumberString)
			singleNumber /= 4
			singleNumber -= 1
			choice.choice1 = singleNumber
		}

		if cs[1] == "?" {
			choice.choice2 = -1
		} else {
			singleNumberString := strings.Split(cs[1], "-")[1]
			singleNumber, _ := strconv.Atoi(singleNumberString)
			singleNumber /= 4
			singleNumber -= 1
			choice.choice2 = singleNumber
		}

		if cs[2] == "?" {
			choice.choice3 = -1
		} else {
			singleNumberString := strings.Split(cs[2], "-")[1]
			singleNumber, _ := strconv.Atoi(singleNumberString)
			singleNumber /= 4
			singleNumber -= 1
			choice.choice3 = singleNumber
		}

		choices = append(choices, choice)
	}
	return choices
}

func randSwap(s []int) []int {
	n1 := rand.Intn(len(s))
	n2 := rand.Intn(len(s))

	t := s[n1]
	s[n1] = s[n2]
	s[n2] = t
	return s
}

func calcScore(solution []int, choices []Choice) (int, int, int, int, int) {
	score := 0
	first := 0
	second := 0
	third := 0
	none := 0
	for i, choice := range choices {
		if choice.dummy {
			continue
		}
		if solution[i] == choice.choice1 {
			score -= 0
			first += 1
		} else if solution[i] == choice.choice2 {
			score -= 1
			second += 1
		} else if solution[i] == choice.choice3 {
			score -= 5
			third += 1
		} else {
			score -= 10
			none += 1
		}
	}
	return score, first, second, third, none
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
