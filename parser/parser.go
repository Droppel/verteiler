package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"verteiler/genome"
)

const (
	dontCareWord = "Egal"
)

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file %s: %w", filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse file as CSV for %s: %w", filePath, err)
	}
	return records, nil
}

func ParseChoices(filename string, lineSeparator string, maxgroupsize int) (genome.GroupList, error) {
	groups := make([]genome.Group, 0)
	records, err := readCsvFile(filename)
	if err != nil {
		return nil, err
	}

	for id, groupSettings := range records[1:] {
		size, _ := strconv.Atoi(groupSettings[1])
		group := genome.Group{
			Id:               id,
			Members:          groupSettings[2],
			Size:             size,
			CurrentSelection: -1,
		}

		choices := make([]int, 3)
		for i, column := range groupSettings[3:] {
			if column == dontCareWord {
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

	return groups, nil
}
