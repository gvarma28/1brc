package solution

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/gvarma28/1brc/internal/store"
)

type Solution3 struct {
	DefaultSolution
	MeasurementsFilePath string
}

const (
	fnvOffset = 14695981039346656037
	fnvPrime  = 1099511628211
)

// 1. in solution-2, the main bottlenecks where strings functions (Split, genSplit, etc) and float parsing. let's try to improve these since we're given that the input would have only one fractional digit
func (s Solution3) Solve() {
	file, err := os.Open(s.MeasurementsFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	jobs := make(chan string, 1000)

	// stores := []map[string]Stats3{}
	stores := []store.HashTable{}
	var wg sync.WaitGroup

	numCpu := runtime.NumCPU()
	fmt.Printf("Running %d workers\n", numCpu)
	for i := 0; i < numCpu; i++ {
		wg.Add(1)
		stores = append(stores, store.HashTable{})

		go func(i int) {
			defer wg.Done()

			for input := range jobs {
				s.calculateAverage(input, &stores[i])
			}
		}(i)
	}

	buf := make([]byte, 10*1024*1024) // 10mb
	leftover := []byte{}
	cnt := 0
	for {
		n, err := file.Read(buf)
		if n > 0 {
			chunk := append(leftover, buf[:n]...)
			start := 0
			for i := range chunk {
				if chunk[i] == '\n' {
					cnt++
					s.progressTracker(cnt)
					chunk[i] = '#'
					start = i + 1
				}
			}
			leftover = []byte{}
			if chunk[len(chunk)-1] != '#' {
				leftover = append([]byte{}, chunk[start:]...)
			}

			jobs <- string(chunk)
		}
		if err == io.EOF {
			break
		}
	}
	fmt.Printf("\n\n")

	close(jobs)
	wg.Wait()

	fmt.Printf("Finished processed all records\n")

	store := s.consolidateMaps(stores)
	s.saveOutput3(store)
}

func (s Solution3) calculateAverage(inputStrs string, store *store.HashTable) {
	idx := 0
	for i := 0; i < len(inputStrs)-1; i++ {
		if inputStrs[i] != '#' {
			continue
		}
		idx2 := idx
		hash := uint64(fnvOffset)
		for j := idx; j < i; j++ {
			b := inputStrs[j]
			if b == ';' {
				idx2 = j
				break
			}
			hash ^= uint64(b)
			hash *= fnvPrime
		}

		station := inputStrs[idx:idx2]
		tempStr := inputStrs[idx2+1 : i]

		temp := parseTemp(tempStr)
		store.GetOrInsert(station, hash, temp)

		idx = i + 1
	}
}

func (s Solution3) consolidateMaps(stores []store.HashTable) map[string]store.Stats3 {
	finalStore := map[string]store.Stats3{}
	for _, s := range stores {
		for i := 0; i < len(s.Keys); i++ {
			if s.Keys[i] == "" {
				continue
			}
			k, v := s.Keys[i], s.Stats[i]
			if _, ok := finalStore[k]; !ok {
				finalStore[k] = v
				continue
			}
			cur := finalStore[k]
			newStat := store.Stats3{
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

func parseTemp(str string) int {
	temp := 0
	sign := 1

	if str[0] == '-' {
		sign = -1
	}
	n := len(str)

	for i := 0; i < n; i++ {
		v := str[i]
		if v == '-' || v == '.' {
			continue
		}
		temp = (temp * 10) + int(v-'0')
	}

	return temp * sign
}
