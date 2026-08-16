package solution

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Solution1 struct {
	DefaultSolution
	MeasurementsFilePath string
}

func (s Solution1) Solve() {
	file, err := os.Open(s.MeasurementsFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)

	store := make(map[string]Stats)

	cnt := 0
	for scanner.Scan() {
		input := scanner.Text()
		s.calculateAverage(input, store)
		s.progressTracker(cnt)
		cnt++
	}
	fmt.Printf("\n\n")

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Printf("Total Records: %v\n", cnt)
	s.saveOutput(store)
}

func (s Solution1) calculateAverage(inputStr string, store map[string]Stats) {
	input := strings.Split(inputStr, ";")
	station := input[0]
	temp, err := strconv.ParseFloat(input[1], 32)
	if err != nil {
		fmt.Println("Error while converting to float:", err)
		os.Exit(1)
	}

	tempF32 := float32(temp)

	if val, ok := store[station]; ok {
		if tempF32 > val.Max {
			val.Max = tempF32
		}
		if tempF32 < val.Min {
			val.Min = tempF32
		}
		newMean := (val.Mean*float32(val.Total) + tempF32) / float32(val.Total+1)
		val.Total += 1
		val.Mean = newMean

		store[station] = val
	} else {
		store[station] = Stats{
			Min:   tempF32,
			Max:   tempF32,
			Mean:  tempF32,
			Total: 1,
		}
	}
}
