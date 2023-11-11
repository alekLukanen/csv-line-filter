package csvLineFilter

import (
	"bufio"
	"bytes"
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

	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLines)

	return &CSVLineFilter{
		Reader:            reader,
		RegularExpression: regularExpression,
		scanner:           scanner,
	}, nil
}

func (obj *CSVLineFilter) Read(buffer []byte) (int, error) {

	bufferLength := len(buffer)
	if bufferLength == 0 {
		return 0, nil
	}

	if obj.currentLineBuffer != nil {
		currentLineBufferLength := len(*obj.currentLineBuffer)
		if obj.currentLineIndex < currentLineBufferLength {
			maxBufferLength := intMin(bufferLength, currentLineBufferLength-obj.currentLineIndex)
			copy(buffer[0:maxBufferLength], (*obj.currentLineBuffer)[obj.currentLineIndex:obj.currentLineIndex+maxBufferLength])
			obj.currentLineIndex = obj.currentLineIndex + maxBufferLength
			return maxBufferLength, nil
		}
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
	maxLineBufferLength := intMin(bufferLength, currentLineBufferLength)
	obj.currentLineIndex = maxLineBufferLength
	copy(buffer[0:maxLineBufferLength], (*obj.currentLineBuffer)[0:maxLineBufferLength])
	return maxLineBufferLength, nil
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Need to use this from the scanner package but without
// removing the newline
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0 : i+1]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
