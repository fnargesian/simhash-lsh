package simhashlsh

import (
	"sort"
	"strconv"
	"sync"
)

// DistanceFunc is a function for calculate distance between two vectors
type DistanceFunc func(p1 []float64, p2 []float64) float64

func euclideanDistSquare(p1 []float64, p2 []float64) (sum float64) {
	for i := range p1 {
		d := p2[i] - p1[i]
		sum += d * d
	}
	return sum
}

type cosineLshParam struct {
	// Dimensionality of the input data.
	dim int
	// Number of hash tables.
	l int
	// Number of hash functions for each table.
	m int
	// Hyperplanes
	hyperplanes [][]float64
	// Number of hash functions
	h int
	// Func for calculate distance between vectors
	dFunc DistanceFunc
}

// NewLshParams initializes the LSH settings.
func newCosineLshParam(dim, l, m, h int, hyperplanes [][]float64) *cosineLshParam {
	return &cosineLshParam{
		dim:         dim,
		l:           l,
		m:           m,
		hyperplanes: hyperplanes,
		h:           h,
		dFunc:       euclideanDistSquare,
	}
}

// Hash returns all combined hash values for all hash tables.
func (clsh *cosineLshParam) hash(point []float64) []hashTableKey {
	simhash := newSimhash(clsh.hyperplanes, point)
	hvs := make([]hashTableKey, clsh.l)
	for i := range hvs {
		s := make(hashTableKey, clsh.m)
		for j := 0; j < clsh.m; j++ {
			s[j] = uint8(simhash.sig[i*clsh.m+j])
		}
		hvs[i] = s
	}
	return hvs
}

// CosineLsh is an implementation of Random projection LSH
// https://en.wikipedia.org/wiki/Locality-sensitive_hashing#Random_projection
type CosineLsh struct {
	// Param type
	*cosineLshParam
	// Tables
	tables []hashTable
}

// NewCosineLsh created an instance of Cosine LSH.
// dim is the number of dimensions of the input points (also the number of dimensions of each hyperplane)
// l is the number of hash tables, m is the number of hash values in each hash table.
func NewCosineLsh(dim, l, m int) *CosineLsh {
	h := m * l
	hyperplanes := newHyperplanes(h, dim)
	tables := make([]hashTable, l)
	for i := range tables {
		tables[i] = make(hashTable)
	}
	return &CosineLsh{
		cosineLshParam: newCosineLshParam(dim, l, m, h, hyperplanes),
		tables:         tables,
	}
}

// Insert adds a new data point to the Cosine LSH.
// point is a data point being inserted into the index and
// id is the unique identifier for the data point.
func (index *CosineLsh) Insert(point []float64, id string) {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(index.hash(point))
	// Insert key into all hash tables
	var wg sync.WaitGroup
	wg.Add(len(index.tables))
	for i := range index.tables {
		hv := hvs[i]
		table := index.tables[i]
		go func(table hashTable, hv uint64) {
			if _, exist := table[hv]; !exist {
				table[hv] = make(hashTableBucket, 0)
			}
			table[hv] = append(table[hv], Point{Vector: point, ID: id})
			wg.Done()
		}(table, hv)
	}
	wg.Wait()
}

type distanceTuple struct {
	id   string
	dist float64
}

// Query finds the ids of approximate nearest neighbour candidates,
// in un-sorted order, given the query point.
func (index *CosineLsh) Query(q []float64) []Point {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(index.hash(q))
	// Keep track of keys seen
	seen := make(map[string]Point)
	for i, table := range index.tables {
		if candidates, exist := table[hvs[i]]; exist {
			for _, id := range candidates {
				if _, exist := seen[id.ID]; exist {
					continue
				}
				seen[id.ID] = id
			}
		}
	}

	distances := make([]distanceTuple, 0, len(seen))
	// is it matrix on vector multiplication?
	for key, value := range seen {
		dist := index.dFunc(q, value.Vector)
		distances = append(distances, distanceTuple{id: key, dist: dist})
	}
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].dist < distances[j].dist
	})

	// Collect results
	values := make([]Point, 0, len(seen))
	for _, distProduct := range distances {
		values = append(values, seen[distProduct.id])
	}

	return values
}

func (index *CosineLsh) toBasicHashTableKeys(keys []hashTableKey) []uint64 {
	basicKeys := make([]uint64, index.cosineLshParam.l)
	for i, key := range keys {
		s := ""
		for _, hashVal := range key {
			switch hashVal {
			case uint8(0):
				s += "0"
			case uint8(1):
				s += "1"
			default:
				panic("Hash value is not 0 or 1")
			}
		}
		v, err := strconv.ParseUint(s, 2, 64)
		if err != nil {
			panic(err)
		}
		basicKeys[i] = v
	}
	return basicKeys
}
