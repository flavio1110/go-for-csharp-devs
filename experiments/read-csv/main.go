package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func pf(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
func panicIfError(err error) {
	inner := errors.Unwrap(err)
	if inner != nil {
		panic(err)
	}
}

func main() {
	file, err := os.Open("WPP2022_Demographic_Indicators_OtherVariants.csv")
	panicIfError(err)
	defer file.Close()

	avgWithReaderAll, err := getAvgPopulationWithReadAll(file)
	panicIfError(fmt.Errorf("getAvgPopulationWithReadAll: %w", err))
	fmt.Printf("avgWithReaderAll: \t\t\t\t%f\n", avgWithReaderAll)

	_, err = file.Seek(0, 0)
	panicIfError(fmt.Errorf("seek: %w", err))

	avgWithReaderPerLine, err := getAvgPopulationWithReadPerLine(file)
	panicIfError(fmt.Errorf("getAvgPopulationWithReadPerLine: %w", err))
	fmt.Printf("getAvgPopulationWithReadPerLine: \t\t%f\n", avgWithReaderPerLine)

	_, err = file.Seek(0, 0)
	panicIfError(fmt.Errorf("seek: %w", err))

	avgWithReaderPerLineNoStruct, err := getAvgPopulationWithReadPerLineWithoutStruct(file)
	panicIfError(fmt.Errorf("getAvgPopulationWithReadPerLineWithoutStruct: %w", err))
	fmt.Printf("getAvgPopulationWithReadPerLineWithoutStruct: \t%f\n", avgWithReaderPerLineNoStruct)
}

func getAvgPopulationWithReadAll(file *os.File) (float64, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // to allow variable number of fields

	lines, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("readAll: %w", err)
	}

	isHeader := true
	pop := float64(0)

	var p population
	for _, line := range lines {
		if isHeader {
			isHeader = false
			continue
		}

		p.locTypeName = line[7]
		p.location = line[9]
		p.variant = line[11]
		p.time = pf(line[12])
		p.tPopulation1Jan = pf(line[13])

		if p.locTypeName == "Country/Area" &&
			p.time == 2022 &&
			p.variant == "Median PI" {
			pop += p.tPopulation1Jan
		}
	}
	return pop / float64(len(lines)-1) * 1000, nil
}

func getAvgPopulationWithReadPerLine(file *os.File) (float64, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true

	var p population
	isHeader := true
	pop := float64(0)
	var line []string
	var err error
	count := float64(0)

	for {
		line, err = reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return 0, fmt.Errorf("read: %w", err)
		}

		if isHeader {
			isHeader = false
			continue
		}
		count++

		p.locTypeName = line[7]
		p.location = line[9]
		p.variant = line[11]
		p.time = pf(line[12])
		p.tPopulation1Jan = pf(line[13])
		if p.locTypeName == "Country/Area" &&
			p.time == 2022 &&
			p.variant == "Median PI" {
			pop += p.tPopulation1Jan

		}
	}
	return pop / count * 1000, nil
}

func getAvgPopulationWithReadPerLineWithoutStruct(file *os.File) (float64, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true

	var p population
	isHeader := true
	pop := float64(0)
	var line []string
	var err error
	count := float64(0)

	for {
		line, err = reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return 0, fmt.Errorf("read: %w", err)
		}

		if isHeader {
			isHeader = false
			continue
		}
		count++
		p.locTypeName = line[7]
		p.variant = line[11]
		p.time = pf(line[12])
		p.tPopulation1Jan = pf(line[13])
		if line[7] == "Country/Area" &&
			line[12] == "2022" &&
			line[11] == "Median PI" {
			pop += pf(line[13])

		}
	}
	return pop / count * 1000, nil
}

type population struct {
	iso3Code        string
	locTypeName     string
	location        string
	variant         string
	time            float64
	tPopulation1Jan float64
}
