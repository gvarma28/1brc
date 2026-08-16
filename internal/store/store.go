package store

const bucketCount = 2048

type Stats3 struct {
	Min   int
	Sum   int64
	Max   int
	Total int
}

type HashTable struct {
	Keys  [bucketCount]string
	Stats [bucketCount]Stats3
}

func (t *HashTable) GetOrInsert(key string, hash uint64, temp int) {
	idx := hash & (bucketCount - 1) // or hash % bucketCount; if the bucketCount is a power of 2
	for {
		if t.Keys[idx] == "" {
			t.Keys[idx] = key
			t.Stats[idx] = Stats3{
				Min:   temp,
				Max:   temp,
				Sum:   int64(temp),
				Total: 1,
			}
			return
		}
		if t.Keys[idx] == key {
			val := &t.Stats[idx]
			val.Max = max(val.Max, temp)
			val.Min = min(val.Min, temp)
			val.Sum += int64(temp)
			val.Total++
			return
		}
		idx = (idx + 1) & (bucketCount - 1) // (& bucketCount - 1) is necessary for out of bounds check
	}
}
