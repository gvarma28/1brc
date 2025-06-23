package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Stats struct {
	Min   float32
	Mean  float32
	Max   float32
	Total int
}

// output format: <weather-station>=<min>/<mean>/<max>
func main() {
	start := time.Now()

	filename := "../measurements.txt"
	file, err := os.Open(filename)
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
		calculateAverage(input, store)
		if (cnt % 1_000_000) == 0 {
			fmt.Printf("Processed %v records\n", cnt)
		}
		cnt++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Printf("Total Records: %v\n", cnt)
	fmt.Printf("Result: %v\n", store)
	fmt.Printf("Execution time: %s\n", time.Since(start))
}

func calculateAverage(inputStr string, store map[string]Stats) {
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
