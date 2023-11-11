package csvLineFilter

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
)

func TestFilterFile(t *testing.T) {

	data := strings.NewReader(`
column_a,column_b,column_c
"a1","b1","c1"
"a2","b2","c2"
"a3","b3","c3"
	`)
	expression := "1"

	csvLineFilterReader, err := NewCSVLineFilter(data, expression)
	if err != nil {
		t.Fatal(err)
	}

	resultData, err := ioutil.ReadAll(csvLineFilterReader)
	if err != nil {
		t.Fatal(err)
	}

	resultDataString := string(resultData)
	expectedResultDataString := "\"a1\",\"b1\",\"c1\"\n"

	if expectedResultDataString != resultDataString {
		t.Fatalf("filtered data incorrect: %s", resultDataString)
	}

}

func TestFilterWorksWithCsvReader(t *testing.T) {

	data := strings.NewReader(`
column_a,column_b,column_c
"a1","b1","c1"
"a2","b2","c2"
"a3","b3","c3"
	`)
	expression := "1"

	csvLineFilterReader, err := NewCSVLineFilter(data, expression)
	if err != nil {
		t.Fatal(err)
	}

	csvReader := csv.NewReader(csvLineFilterReader)

	items, err := csvReader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	expectedItems := [][]string{
		{"a1", "b1", "c1"},
	}

	if reflect.DeepEqual(items, expectedItems) != true {
		t.Fatalf("csv item incorrect: %s", items)
	}

}

func setupLargeFile(dir string, lineCount int) (string, error) {

	data := make([][]string, 0, lineCount+1)
	data = append(data, []string{"column_a", "column_b", "column_c"})

	for i := 0; i < lineCount; i++ {
		a := fmt.Sprintf("a%d", i)
		b := fmt.Sprintf("b%d", i)
		c := fmt.Sprintf("c%d", i)
		data = append(data, []string{a, b, c})
	}

	filePath := path.Join(dir, "example.csv")
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(data) // calls Flush internally

	if err := csvWriter.Error(); err != nil {
		return "", err
	}

	return filePath, nil
}

func BenchmarkLineFilteredCSVFile(b *testing.B) {

	lineCount := 10_000

	dir, err := ioutil.TempDir("", "benchmark")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	path, err := setupLargeFile(dir, lineCount)
	if err != nil {
		os.RemoveAll(dir)
		b.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		os.RemoveAll(dir)
		b.Fatal(err)
	}
	defer file.Close()

	b.ResetTimer()

	csvLineFilterReader, err := NewCSVLineFilter(file, `a1\d0,`)
	if err != nil {
		file.Close()
		os.RemoveAll(dir)
		b.Fatal(err)
	}

	csvReader := csv.NewReader(csvLineFilterReader)

	data := make([]string, 0, lineCount*3)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			file.Close()
			os.RemoveAll(dir)
			b.Fatal(err)
		}

		for _, el := range line {
			data = append(data, el)
		}
	}

	if len(data) != 30 {
		b.Fatalf("expected 100 items but found: %d", len(data))
	}

}

func BenchmarkUnfilteredCSVFile(b *testing.B) {

	lineCount := 10_000

	dir, err := ioutil.TempDir("", "benchmark")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	path, err := setupLargeFile(dir, lineCount)
	if err != nil {
		os.RemoveAll(dir)
		b.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		os.RemoveAll(dir)
		b.Fatal(err)
	}
	defer file.Close()

	b.ResetTimer()

	csvReader := csv.NewReader(file)

	data := make([][]string, lineCount)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			file.Close()
			os.RemoveAll(dir)
			b.Fatal(err)
		}

		data = append(data, line)
	}

}
