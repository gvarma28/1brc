package solution

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Solution2 struct {
	DefaultSolution
	MeasurementsFilePath string
}

// 1. parallel workers + mutex the `map`, was very very counter-productive due to the additional overhead due goroutine management. (~6mins)
// 2. let's try to create separate workers and give each worker it's own `map` so that we wouldn't have to deal with locking, after all the workers are finished running, we will consolidate the results into a single map, would that work? edit; wow this worked great (~24secs)
func (s Solution2) Solve() {
	file, err := os.Open(s.MeasurementsFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)

	jobs := make(chan string, 1000)

	stores := []map[string]Stats2{}
	var wg sync.WaitGroup

	numCpu := runtime.NumCPU()
	fmt.Printf("Running %d workers\n", numCpu)
	for i := 0; i < numCpu; i++ {
		wg.Add(1)
		stores = append(stores, map[string]Stats2{})

		go func(i int) {
			defer wg.Done()

			for input := range jobs {
				s.calculateAverage(input, stores[i])
			}
		}(i)
	}

	cnt := 0

	var b strings.Builder
	for scanner.Scan() {
		b.WriteString(scanner.Text())
		b.WriteString("#")
		s.progressTracker(cnt)
		cnt++
		if cnt%1_000_000 == 0 {
			jobs <- b.String()
			b.Reset()
		}
	}
	if b.Len() != 0 {
		jobs <- b.String()
	}
	fmt.Printf("\n\n")

	close(jobs)
	wg.Wait()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Printf("Finished processed all records\n")

	store := s.consolidateMaps(stores)
	s.saveOutput2(store)
}

func (s Solution2) calculateAverage(inputStrs string, store map[string]Stats2) {
	inputArr := strings.Split(inputStrs, "#")
	for _, inputStr := range inputArr[:len(inputArr)-1] {
		input := strings.Split(inputStr, ";")
		station := input[0]
		temp, err := strconv.ParseFloat(input[1], 32)
		if err != nil {
			fmt.Println("Error while converting to float:", err)
			os.Exit(1)
		}

		tempF32 := float32(temp)

		if val, ok := store[station]; ok {
			val.Max = max(val.Max, tempF32)
			val.Min = min(val.Min, tempF32)
			val.Sum += float64(tempF32)
			val.Total += 1

			store[station] = val
		} else {
			store[station] = Stats2{
				Min:   tempF32,
				Max:   tempF32,
				Sum:   float64(tempF32),
				Total: 1,
			}
		}
	}
}

func (s Solution2) consolidateMaps(stores []map[string]Stats2) map[string]Stats2 {
	finalStore := map[string]Stats2{}
	for _, store := range stores {
		for k, v := range store {
			if _, ok := finalStore[k]; !ok {
				finalStore[k] = v
				continue
			}
			cur := finalStore[k]
			newStat := Stats2{
				Max:   max(cur.Max, v.Max),
				Min:   min(cur.Min, v.Min),
				Sum:   cur.Sum + v.Sum,
				Total: cur.Total + v.Total,
			}
			finalStore[k] = newStat
		}
	}
	return finalStore
}
