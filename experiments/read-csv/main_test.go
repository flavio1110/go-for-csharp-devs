package main

import (
	"fmt"
	"os"
	"testing"
)

func Benchmark_getAvgPopulationWithReadAll(b *testing.B) {
	file, err := os.Open("WPP2022_Demographic_Indicators_OtherVariants.csv")
	panicIfError(err)
	defer file.Close()

	for i := 0; i < b.N; i++ {
		_, err = file.Seek(0, 0)
		panicIfError(fmt.Errorf("seek: %w", err))

		_, err = getAvgPopulationWithReadAll(file)
		panicIfError(fmt.Errorf("getAvgPopulationWithReadAll: %w", err))
	}
}

func Benchmark_getAvgPopulationWithReadPerLine(b *testing.B) {
	file, err := os.Open("WPP2022_Demographic_Indicators_OtherVariants.csv")
	panicIfError(err)
	defer file.Close()

	for i := 0; i < b.N; i++ {
		_, err = file.Seek(0, 0)
		panicIfError(fmt.Errorf("seek: %w", err))

		_, err = getAvgPopulationWithReadPerLine(file)
		panicIfError(fmt.Errorf("getAvgPopulationWithReadPerLine: %w", err))
	}
}

func Benchmark_getAvgPopulationWithReadPerLineWithoutStruct(b *testing.B) {
	file, err := os.Open("WPP2022_Demographic_Indicators_OtherVariants.csv")
	panicIfError(err)
	defer file.Close()

	for i := 0; i < b.N; i++ {
		_, err = file.Seek(0, 0)
		panicIfError(fmt.Errorf("seek: %w", err))

		_, err = getAvgPopulationWithReadPerLineWithoutStruct(file)
		panicIfError(fmt.Errorf("Benchmark_getAvgPopulationWithReadPerLineWithoutStruct: %w", err))
	}
}
