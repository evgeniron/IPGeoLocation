package db

import (
    "encoding/csv"
    "log"
    "os"
)

// Defining sort of enum of keys
const (
	ip = iota
	country
	city
)

type CsvDb struct {
	locationMap map[string]Location
}

// Map IP to locations: IP = {Country, City}
func (db *CsvDb) mapLocations(records [][]string) {
	locations := make(map[string]Location)

	for _, csvRecord := range records {
		record := csvRecord[ip]
		locations[record] = Location{Country: csvRecord[country], City: csvRecord[city]}
	}

	db.locationMap = locations
}

// Implement db interface
func (db *CsvDb) GetLocation(Ip string) Location {
	return db.locationMap[Ip]
}

// Initialize new instance of CSV database
func NewCsvDb() DB {
	conf := NewDBConfig("DB")
	db := CsvDb{}
	db.mapLocations(readCsv(conf.Path))
    return &db
}

// readCsv accepts a file and returns its content as a multi-dimentional type
// with lines and each column. Only parses to string type.
func readCsv(filename string) ([][]string) {

    // Open CSV file
    f, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    // Read File into a Variable
    lines, err := csv.NewReader(f).ReadAll()
    if err != nil {
        log.Fatal(err)   
    }

    return lines
}