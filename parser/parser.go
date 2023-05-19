package parser

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"verteiler/datastructures"
)

func ParseChoices(filename string, lineSeparator string) datastructures.GroupList {
	idCounter := 0 //Used to give each group its own ID
	groups := make([]datastructures.Group, 0)
	input, err := os.ReadFile("inputfiles" + string(filepath.Separator) + filename)
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

		group := datastructures.Group{
			Id:               idCounter,
			Dummy:            false,
			Members:          strings.Split(membersString, ","),
			CurrentSelection: -1,
		}
		idCounter += 1

		group.Size = len(group.Members)

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
		group.Choices = choices
		groups = append(groups, group)
	}
	return groups
}
