package datastructures

import (
	"fmt"
	"math/rand"
	"os"
)

type Group struct {
	Id               int
	Dummy            bool
	Size             int
	Members          string
	Choices          []int
	CurrentSelection int
}

type GroupList []Group

func (a GroupList) Len() int           { return len(a) }
func (a GroupList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a GroupList) Less(i, j int) bool { return a[i].Size > a[j].Size }

type Slot struct {
	Id       int
	TimeSlot int
	Capacity int
	Amount   int
}

type Solution struct {
	Occupancy     []Slot
	Groups        []Group
	InvAllocation map[int][]int
}

func (s *Solution) Print(score int, seed int64) {
	output := fmt.Sprintf("Seed: %d\n", seed)
	for _, slot := range s.Occupancy {
		output += fmt.Sprintln("====================================")
		output += fmt.Sprintf("GROUP %d Available space: %d\n", slot.Id, slot.Capacity)
		for _, groupId := range s.InvAllocation[slot.Id] {
			group := s.Groups[groupId]
			choice := -1
			for i, ch := range group.Choices {
				if ch == s.Occupancy[group.CurrentSelection].TimeSlot {
					choice = i + 1
					break
				}
			}
			output += fmt.Sprintf("%d. Wahl\n", choice)
			for _, member := range group.Members {
				output += fmt.Sprintln(member)
			}
			output += fmt.Sprintln("-------------------------------------")
		}
	}

	fmt.Print(output)
	if score >= -20 {
		os.WriteFile(fmt.Sprintf("scores/Score%d-%d.txt", score, seed), []byte(output), os.ModeAppend)
	}
}

func (s *Solution) Copy() Solution {
	copiedSolution := Solution{
		Occupancy:     make([]Slot, len(s.Occupancy)),
		Groups:        make([]Group, len(s.Groups)),
		InvAllocation: make(map[int][]int),
	}

	copy(copiedSolution.Occupancy, s.Occupancy)
	copy(copiedSolution.Groups, s.Groups)
	for k, v := range s.InvAllocation {
		newArr := make([]int, len(v))
		copy(newArr, v)
		copiedSolution.InvAllocation[k] = newArr
	}
	return copiedSolution
}

type Swap struct {
	slot         int
	swapPartners []int
}

func (s *Solution) RandSwap() {
	groups := s.Groups
	randGroup := groups[rand.Intn(len(groups))]
	possibleSwaps := make([]Swap, 0)
	for i, slot := range s.Occupancy {
		if i == randGroup.CurrentSelection {
			continue
		}
		if slot.Capacity-slot.Amount >= randGroup.Size {
			possibleSwaps = append(possibleSwaps, Swap{i, make([]int, 0)})
			continue
		}

		solved, set := solveSubsetSum(s, randGroup.Size, s.InvAllocation[i])
		if !solved {
			continue
		}

		possibleSwaps = append(possibleSwaps, Swap{slot: s.Groups[set[0]].CurrentSelection, swapPartners: set})
	}

	if len(possibleSwaps) == 0 {
		return
	}
	choosenSwap := possibleSwaps[rand.Intn(len(possibleSwaps))]

	randGroupSelection := randGroup.CurrentSelection
	s.Groups[randGroup.Id].CurrentSelection = choosenSwap.slot

	//Update invAllocation and Occupancy
	s.InvAllocation[randGroupSelection] = removeElement(s.InvAllocation[randGroupSelection], randGroup.Id)
	s.InvAllocation[choosenSwap.slot] = append(s.InvAllocation[choosenSwap.slot], randGroup.Id)
	s.Occupancy[randGroupSelection].Amount -= randGroup.Size
	s.Occupancy[choosenSwap.slot].Amount += randGroup.Size
	for _, setGroupId := range choosenSwap.swapPartners {
		setGroup := s.Groups[setGroupId]
		s.Groups[setGroupId].CurrentSelection = randGroupSelection
		s.Occupancy[randGroupSelection].Amount += setGroup.Size
		s.Occupancy[choosenSwap.slot].Amount -= setGroup.Size
		s.InvAllocation[choosenSwap.slot] = removeElement(s.InvAllocation[choosenSwap.slot], setGroupId)
		s.InvAllocation[randGroupSelection] = append(s.InvAllocation[randGroupSelection], setGroupId)
	}
}

func solveSubsetSum(solution *Solution, sum int, groupIds []int) (bool, []int) {
	for i, groupId := range groupIds {
		summand := solution.Groups[groupId].Size
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

func removeElement[T comparable](s []T, e T) []T {
	for i, elem := range s {
		if elem == e {
			return Remove(s, i)
		}
	}
	// return false, s
	panic("Tried to remove non existant")
}

func Remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
