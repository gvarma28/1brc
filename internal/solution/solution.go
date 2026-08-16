package solution

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/gvarma28/1brc/internal/store"
)

type Stats struct {
	Min   float32
	Mean  float32
	Max   float32
	Total int
}

type Stats2 struct {
	Min   float32
	Sum   float64
	Max   float32
	Total int
}

type Solution interface {
	Solve()
	progressTracker(cnt int)
	saveOutput(stats map[string]Stats)
	saveOutput2(stats map[string]Stats2)
	saveOutput3(stats map[string]store.Stats3)
}

type DefaultSolution struct {
	OutputFilePath string
}

func (s DefaultSolution) progressTracker(cnt int) {
	progressStep := 100_000_000
	if (cnt % progressStep) == 0 {
		total := 1_000_000_000 / progressStep
		finished := cnt / progressStep
		fmt.Printf("\rProgress: [%s%s] %d%%",
			strings.Repeat("#", finished), strings.Repeat(" ", total-finished), finished*100/total)
	}
}

func (s DefaultSolution) saveOutput(stats map[string]Stats) {
	keys := []string{}
	for k := range stats {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	var result strings.Builder
	result.WriteString("{")
	for i, k := range keys {
		v := stats[k]
		s := fmt.Sprintf(", %s=%v/%s/%v", k, v.Min, fmt.Sprintf("%.1f", v.Mean), v.Max)
		if i == 0 {
			s = fmt.Sprintf("%s=%v/%s/%v", k, v.Min, fmt.Sprintf("%.1f", v.Mean), v.Max)
		}
		result.WriteString(s)
	}
	result.WriteString("}\n")

	err := os.WriteFile(s.OutputFilePath, []byte(result.String()), 0644)
	if err != nil {
		panic(err)
	}
}

func (s DefaultSolution) saveOutput2(stats map[string]Stats2) {
	keys := []string{}
	for k := range stats {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	var result strings.Builder
	result.WriteString("{")
	for i, k := range keys {
		v := stats[k]
		mean := v.Sum / float64(v.Total)
		s := fmt.Sprintf(", %s=%v/%s/%v", k, v.Min, fmt.Sprintf("%.1f", mean), v.Max)
		if i == 0 {
			s = fmt.Sprintf("%s=%v/%s/%v", k, v.Min, fmt.Sprintf("%.1f", mean), v.Max)
		}
		result.WriteString(s)
	}
	result.WriteString("}\n")

	err := os.WriteFile(s.OutputFilePath, []byte(result.String()), 0644)
	if err != nil {
		panic(err)
	}
}

func (s DefaultSolution) saveOutput3(stats map[string]store.Stats3) {
	keys := []string{}
	for k := range stats {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	var result strings.Builder
	result.WriteString("{")
	for i, k := range keys {
		v := stats[k]
		mean := (float64(v.Sum) / 10) / float64(v.Total)
		minimum, maximum := float32(v.Min)/10, float32(v.Max)/10
		s := fmt.Sprintf(", %s=%v/%s/%v", k, minimum, fmt.Sprintf("%.1f", mean), maximum)
		if i == 0 {
			s = fmt.Sprintf("%s=%v/%s/%v", k, minimum, fmt.Sprintf("%.1f", mean), maximum)
		}
		result.WriteString(s)
	}
	result.WriteString("}\n")

	err := os.WriteFile(s.OutputFilePath, []byte(result.String()), 0644)
	if err != nil {
		panic(err)
	}
}
