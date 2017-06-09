package simhashlsh

import (
	"log"
	"math/rand"
	"strconv"
	"testing"
)

func Test_CosineLshQuery(t *testing.T) {
	ls := []int{20, 5, 10, 25, 4}
	ks := []int{5, 20, 10, 4, 25}
	embs := randomEmbeddings(1000, 300, 1.0)
	for j, _ := range ls {
		var avg float64
		clsh := NewCosineLsh(300, ls[j], ks[j])
		insertedEmbs := make([]string, 1000)
		for i, e := range embs {
			clsh.Insert(e, strconv.Itoa(i))
			insertedEmbs[i] = strconv.Itoa(i)
		}
		// Use the inserted embeddings as queries, and
		// verify that we can get back each query itself
		log.Printf("avg number of returned results with k=%d and l=%d for each query (out of %d indexed items): ", ks[j], ls[j], len(embs))
		for i, key := range insertedEmbs {
			results := clsh.Query(embs[i])
			avg += float64(len(results))
			found := false
			for _, foundKey := range results {
				if foundKey == key {
					found = true
				}
			}
			if !found {
				t.Error("Query fail")
			}
		}
		log.Printf(" %f", avg/float64(len(embs)))
	}
}

func randomEmbeddings(n, dim int, max float64) [][]float64 {
	random := rand.New(rand.NewSource(1))
	embs := make([][]float64, n)
	for i := 0; i < n; i++ {
		embs[i] = make([]float64, dim)
		for d := 0; d < dim; d++ {
			embs[i][d] = random.Float64() * max
		}
	}
	return embs
}
