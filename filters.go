package csvLineFilter

import (
	"bufio"
	"errors"
	"io"
	"regexp"
)

type CSVLineFilter struct {
	Reader            io.Reader
	RegularExpression *regexp.Regexp
	scanner           *bufio.Scanner
	currentLineBuffer *[]byte
	currentLineIndex  int
}

func NewCSVLineFilter(reader io.Reader, expression string) (*CSVLineFilter, error) {

	regularExpression, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}

	return &CSVLineFilter{
		Reader:            reader,
		RegularExpression: regularExpression,
		scanner:           bufio.NewScanner(reader),
	}, nil
}

func (obj *CSVLineFilter) Read(buffer []byte) (int, error) {

	if len(buffer) == 0 {
		return 0, nil
	}

	isEndOfFile := true
	for obj.scanner.Scan() {
		lineItem := obj.scanner.Bytes()
		lineItemMatches := obj.RegularExpression.Match(lineItem)
		if len(lineItem) > 0 && lineItemMatches {
			obj.currentLineBuffer = &lineItem
			isEndOfFile = false
			break
		}
	}

	if isEndOfFile {
		return 0, io.EOF
	}

	currentLineBufferLength := len(*obj.currentLineBuffer)
	if len(buffer) < currentLineBufferLength {
		return 0, errors.New("buffer too small")
	}

	copy(buffer[0:currentLineBufferLength], (*obj.currentLineBuffer)[:])
	return currentLineBufferLength, nil
}
