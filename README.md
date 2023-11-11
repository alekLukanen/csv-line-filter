## CSV Line Filter

A simple prefilter for reading csv files in golang. You can use this to reduce the number 
of allocations made by the golang `encoding/csv` package when reading files containing 
lines that you do not need.

Get the package
```
go get github.com/alekLukanen/csv-line-filter
```

Then you can import and use the package
```
package main

import (
    "os"
    "log"

	"github.com/alekLukanen/csv-line-filter"
)

func main() {

    file, err := os.Open("./your-file.csv")
	if err != nil {
		log.Fatal(err)
	}

    close := func(err error) {
        file.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("issue found")
	}
    defer close()

    csvLineFilterReader, err := csvLineFilter.NewCSVLineFilter(compressedCsvIndexFileReader, `<regexp-here>`)
	if err != nil {
		close(err)
	}

	csvReader := csv.NewReader(csvLineFilterReader)

    data := make([][]string, 0, 100)
    for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			close(err)
		}

		data = append(data, line)
	}

    log.Printf("total lines: %d", len(data))
}
```