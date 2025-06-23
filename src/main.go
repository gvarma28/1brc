package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

const StatsFilePath = "../output/stats.json"
const MeasurementsFilePath = "../measurements.txt"

const SaveOutputFlag = false

// output format: <weather-station>=<min>/<mean>/<max>
func main() {
	start := time.Now()

	file, err := os.Open(MeasurementsFilePath)
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
		progressTracker(cnt)
		cnt++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Printf("Total Records: %v\n", cnt)
	fmt.Printf("Result: %v\n", store)
	fmt.Printf("Execution time: %s\n", time.Since(start))

	if SaveOutputFlag {
		saveOutput(store)
	}
	if cnt != 1_000_000_000 {
		validateResult(store)
	}
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

func progressTracker(cnt int) {
	progressStep := 10_000_000
	if (cnt % progressStep) == 0 {
		total := 1_000_000_000 / progressStep
		finished := cnt / progressStep
		fmt.Printf("\rProgress: [%s%s] %d%%",
			strings.Repeat("#", finished), strings.Repeat(" ", total-finished), finished*100/total)
	}
}

func saveOutput(stats map[string]Stats) {
	jsonBytes, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(StatsFilePath, jsonBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func validateResult(actual map[string]Stats) {
	jsonBytes, err := os.ReadFile(StatsFilePath)
	if err != nil {
		panic(err)
	}

	var expected map[string]Stats
	if err := json.Unmarshal(jsonBytes, &expected); err != nil {
		panic(err)
	}

	// Compare using reflect.DeepEqual
	if reflect.DeepEqual(expected, actual) {
		fmt.Println("Map matches stats.json")
	} else {
		fmt.Println("Failed: map does NOT match stats.json")
	}
}
