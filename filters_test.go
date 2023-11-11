package csvLineFilter

import (
	"io/ioutil"
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
	expectedResultDataString := `"a1","b1","c1"`

	if expectedResultDataString != resultDataString {
		t.Fatalf("filtered data incorrect: %s", resultDataString)
	}

}
