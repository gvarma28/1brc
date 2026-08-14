package solution

import (
	"fmt"
	"os"
	"slices"
	"strings"
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
	saveOutput(stats map[string]Stats)
	saveOutput2(stats map[string]Stats2)
}

type DefaultSolution struct {
	OutputFilePath string
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
