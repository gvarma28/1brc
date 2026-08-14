package solution

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Solution interface {
	Solve()
	saveOutput(stats map[string]Stats)
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
		if v.Min > -99.9 {
			fmt.Printf("%s: %v", k, v)
		}
		if v.Max < 99.9 {
			fmt.Printf("%s: %v", k, v)
		}
	}
	result.WriteString("}\n")

	err := os.WriteFile(s.OutputFilePath, []byte(result.String()), 0644)
	if err != nil {
		panic(err)
	}
}
