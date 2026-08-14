package main

import (
	"fmt"
	"time"

	solution "github.com/gvarma28/1brc/internal/solution"
)

// output format: <weather-station>=<min>/<mean>/<max>
func main() {
	start := time.Now()
	const outputFilePath = "./output/output.txt"
	const measurementsFilePath = "./measurements.txt"

	var s solution.Solution = solution.Solution1{
		MeasurementsFilePath: measurementsFilePath,
		DefaultSolution: solution.DefaultSolution{
			OutputFilePath: outputFilePath,
		},
	}

	s.Solve()

	fmt.Printf("Execution time: %s\n", time.Since(start))
}
