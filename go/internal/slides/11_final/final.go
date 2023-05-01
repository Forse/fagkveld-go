package main

// This package is 0.0.0 - experimental
import (
	"errors"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type operation[T Number] struct {
	a T
	b T
}

func NewOperation[T Number](a T, b T) *operation[T] {
	return &operation[T]{
		a: a,
		b: b,
	}
}

type Addition[T Number] interface {
	Add() T
}
type Subtraction[T Number] interface {
	Subtract() T
}

func (o *operation[T]) Add() T {
	return o.a + o.b
}
func (o *operation[T]) Subtract() T {
	return o.a - o.b
}

var _ Addition[int] = (*operation[int])(nil)
var _ Subtraction[int] = (*operation[int])(nil)

// --------------------------------

const iterations int = 5
const operationCount int = 50_000_000

// Global variables are allowed, but discouraged due to concurrency issues
var rng *rand.Rand = rand.New(rand.NewSource(0))

func main() {
	for i := 1; i < iterations; i++ {
		rng.Seed(0)
		calculateByReferences()
	}
	for i := 1; i < iterations; i++ {
		rng.Seed(0)
		calculatePacked()
	}
	for i := 1; i < iterations; i++ {
		rng.Seed(0)
		calculatePackedParallel()
	}
}

func calculateByReferences() (float64, error) {
	// Declaration, default initialization
	var ops [operationCount]*operation[float64]

	for i := range ops {
		ops[i] = NewOperation(rng.Float64(), rng.Float64())
	}

	// Declaration + initialization
	var sum float64 = 0

	defer timeTrack(time.Now(), "calculateByReferences")
	for i := range ops {
		sum += ops[i].Add()
	}

	avg := sum / float64(len(ops)*2)
	if math.Abs(avg-0.5) > 0.01 {
		return avg, errors.New("avg is not correct")
	}
	return avg, nil
}

func calculatePacked() (float64, error) {
	// Declaration, default initialization
	var ops [operationCount]operation[float64]

	for i := range ops {
		ops[i] = *NewOperation(rng.Float64(), rng.Float64())
	}

	// Declaration + initialization
	var sum float64 = 0

	defer timeTrack(time.Now(), "calculatePacked")
	for i := range ops {
		sum += ops[i].Add()
	}

	avg := sum / float64(len(ops)*2)
	if math.Abs(avg-0.5) > 0.01 {
		return avg, errors.New("avg is not correct")
	}
	return avg, nil
}

func calculatePackedParallel() (float64, error) {
	// Declaration, default initialization
	var ops [operationCount]operation[float64]

	for i := range ops {
		ops[i] = *NewOperation(rng.Float64(), rng.Float64())
	}

	// Declaration + initialization
	var sum float64 = 0

	defer timeTrack(time.Now(), "calculatePackedParallel")

	threads := runtime.NumCPU()
	// Allocate a channel, go's primitive for message passing
	// The second size/capacity argument is optional - bounded vs unbounded
	sums := make(chan float64, threads)
	batchSize := len(ops) / threads

	for t := 0; t < threads; t++ {
		startIndex := batchSize * t
		threadOps := ops[startIndex : startIndex+batchSize]
		go func(ops []operation[float64], sums *chan float64) {
			sum := 0.0
			for i := range ops {
				sum += ops[i].Add()
			}
			*sums <- sum
		}(threadOps, &sums)
	}

	// We expect 'threads' messages to be received
	for i := 0; i < threads; i++ {
		sum += <-sums
	}

	avg := sum / float64(len(ops)*2)

	if math.Abs(avg-0.5) > 0.01 {
		return avg, errors.New("avg is not correct")
	}
	return avg, nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
