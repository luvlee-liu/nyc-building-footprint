package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func validateYear(yearStr string) bool {
	if yearStr == "" {
		return false
	}
	currentYear := time.Now().Year()
	if year, err := strconv.Atoi(yearStr); err != nil {
		return false
	} else if year < 1 || year > currentYear {
		return false
	}
	return true
}

func main() {
	inputCsvPtr := flag.String("input", "data/building.csv", "raw data csv file")
	outputCsvPtr := flag.String("output", "data/etl.csv", "output csv file")
	flag.Parse()
	startTime := time.Now()

	// Loading csv file
	rFile, err := os.Open(*inputCsvPtr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer rFile.Close()

	// create output for etl.csv
	wFile, err := os.Create(*outputCsvPtr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer wFile.Close()
	writer := csv.NewWriter(wFile)

	// create output file for invalid data entry
	const BadDataFileName = "./invalidData.csv"
	badDataFile, err := os.OpenFile(BadDataFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer badDataFile.Close()
	badDataWriter := csv.NewWriter(badDataFile)

	fmt.Println("Processing begin")

	// Read csv
	reader := csv.NewReader(rFile)

	// https://github.com/CityOfNewYork/nyc-geo-metadata/blob/master/Metadata/Metadata_BuildingFootprints.md
	const (
		BinIdx           = 1
		ConstructYearIdx = 2
		DoittIDIdx       = 6
		HeightRoof       = 7
		ShapeArea        = 10
	)
	var inputCol = []int{
		DoittIDIdx,
		BinIdx,
		ConstructYearIdx,
		HeightRoof,
		ShapeArea,
	}
	colNum := len(inputCol)

	lineCount := 0
	badLineCount := 0
	for line, err := reader.Read(); err == nil; line, err = reader.Read() {
		var outputCol []string
		for _, idx := range inputCol {
			// empty filed is invalid
			if line[idx] == "" {
				break
			}
			outputCol = append(outputCol, line[idx])
		}

		// drop and save bad records(not header) with:
		// - invalid data line with empty fields
		// - invalid year
		if lineCount > 0 && (len(outputCol) != colNum || !validateYear(line[ConstructYearIdx])) {
			badDataWriter.Write(line)
			badDataWriter.Flush()
			badLineCount++
			continue
		}

		if err = writer.Write(outputCol); err != nil {
			fmt.Println("Error:", err)
			break
		}
		writer.Flush()
		lineCount++
	}
	wFile.Sync()

	//print report
	fmt.Println("Dropped ", badLineCount, " records saved in ", BadDataFileName)
	fmt.Println("Imported ", (lineCount-1)-badLineCount, " records")
	fmt.Println("Process ", lineCount-1, " records in ", time.Since(startTime))
}
