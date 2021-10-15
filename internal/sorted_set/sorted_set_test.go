package sorted_set

import (
	"fmt"
	"github.com/buraksezer/sorted"
	"testing"
)

func TestSortedSet(t *testing.T) {
	EdgeSetsWithScore = sorted.NewSortedSetWithScore(0)
	defer EdgeSetsWithScore.Close()

	keys := []string{}
	for i := 0; i < 100; i++ {
		err := EdgeSetsWithScore.Set(bkey(i), uint64(100-i))
		if err != nil {
			t.Fatalf("Expected nil. Got %v", err)
		}
		keys = append(keys, string(bkey(i)))
	}

	for i := 0; i < 100; i++ {
		ok := EdgeSetsWithScore.Check(bkey(i))
		if !ok {
			t.Fatalf("Key could not be found: %s", bkey(i))
		}
	}

	idx := 99
	EdgeSetsWithScore.Range(func(key []byte) bool {
		if keys[idx] != string(key) {
			t.Fatalf("Invalid key: %s", string(key))
		}
		idx--
		return true
	})

	for i := 0; i < 100; i++ {
		err := EdgeSetsWithScore.Delete(bkey(i))
		if err != nil {
			t.Fatalf("Expected nil. Got %v", err)
		}
	}

	if EdgeSetsWithScore.Len() != 0 {
		t.Fatalf("Expected length is zero. Got: %d", EdgeSetsWithScore.Len())
	}
}

func bkey(i int) []byte {
	return []byte(fmt.Sprintf("%09d", i))
}
