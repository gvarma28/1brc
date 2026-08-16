package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	solution "github.com/gvarma28/1brc/internal/solution"
)

const RunSolution = "3"

// output format: <weather-station>=<min>/<mean>/<max>
func main() {
	f, err := os.Create(fmt.Sprintf("cpu%s.prof", RunSolution))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	start := time.Now()
	outputFilePath := fmt.Sprintf("./output/output%s.txt", RunSolution)
	const measurementsFilePath = "./measurements.txt"

	var s solution.Solution = solution.Solution3{
		MeasurementsFilePath: measurementsFilePath,
		DefaultSolution: solution.DefaultSolution{
			OutputFilePath: outputFilePath,
		},
	}

	s.Solve()

	fmt.Printf("Execution time: %s\n", time.Since(start))
}
