package simplejson

import (
	"errors"
	"log"
	"strconv"
	"sync"
)

const (
	scanBeginLiteral = iota
	scanBeginObject
	scanLiteral
	scanObjectKey
	scanArrayBegin
	scanContinue
	scanSkipSpace
	scanEndObject
	scanError
)

const (
	parserKeyObject = iota
	parserValueObject
	parserArray
)

var scannpool = sync.Pool{
	New: func() interface{} {
		return &scanner{}
	},
}

type scanner struct {
	step func(*scanner, byte) int

	err error

	parseState []int
}

func Unmarshal(data []byte, v interface{}) error {
	if err := checkValid(data); err != nil {
		return err
	}
	return nil
}

//Use the state machine like operation to valid the JSON
func checkValid(data []byte) error {
	var scan scanner
	scan.resetScanner()
	for _, c := range data {
		if scan.step(&scan, c) == scanError {
			return scan.err
		}
	}
	return nil
}

func (s *scanner) resetScanner() {
	s.step = stepBeginValue
	s.parseState = s.parseState[0:0]
}

func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\t' || c == '\r' || c == '\n')
}

func stepBeginValue(scan *scanner, b byte) int {
	log.Println("stepBeginValue " + string(b))
	if isSpace(b) {
		return scanSkipSpace
	}

	switch b {
	case '{':
		scan.step = stateBeginStringOrEmpty
		return scan.pushToPaserState(parserKeyObject, scanBeginObject)
	case '[':
		scan.step = stateBeginValueEmpty
		return scan.pushToPaserState(parserArray, scanArrayBegin)
	case '"':
		scan.step = stateScanString
		return scanBeginLiteral
	case 't':
		scan.step = stepT
		return scanBeginLiteral
	case 'f':
		scan.step = stepF
		return scanBeginLiteral
	}

	if b >= '1' && b <= '9' {
		scan.step = stepDecmial
		return scanBeginLiteral
	}

	return scanError
}

func (s *scanner) pushToPaserState(newState int, successCode int) int {
	s.parseState = append(s.parseState, newState)
	return successCode
}

func (s *scanner) popParserState() {
	n := len(s.parseState) - 1
	s.parseState = s.parseState[0:n]
	if len(s.parseState) == 0 {

	} else {
		s.step = stateScanEndValue
	}
}

func stateBeginStringOrEmpty(s *scanner, c byte) int {
	log.Println("stateBeginStringOrEmpty " + string(c))
	if isSpace(c) {
		return scanSkipSpace
	}

	return stateBeginString(s, c)
}

func stateBeginValueEmpty(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}

	if c == ']' {
		return stateScanEndValue(s, c)
	}
	return stepBeginValue(s, c)
}

func stateBeginString(s *scanner, c byte) int {
	log.Println("stateBeginString " + string(c))
	if isSpace(c) {
		return scanSkipSpace
	}

	if c == '"' {
		s.step = stateScanString
		return scanBeginLiteral
	}

	return s.error(c, "No the Begin String should be in the Json")
}

func stateScanString(s *scanner, c byte) int {
	log.Println("stateScanString " + string(c))
	if c == '"' {
		s.step = stateScanEndValue
		return scanContinue
	}

	if c < 0x20 {
		return s.error(c, "Invalid ASCII Code")
	}

	return scanContinue
}

func stateScanEndValue(s *scanner, c byte) int {
	log.Println("stateScanEndValue " + string(c))
	if isSpace(c) {
		return scanSkipSpace
	}

	n := len(s.parseState)

	ps := s.parseState[n-1]
	switch ps {
	case parserKeyObject:
		if c == ':' {
			s.parseState[n-1] = parserValueObject
			s.step = stepBeginValue
			return scanObjectKey
		}
		return s.error(c, "Key object not corrent")
	case parserValueObject:
		if c == ',' {
			s.parseState[n-1] = parserKeyObject
			s.step = stateBeginString
			return scanObjectKey
		}

		if c == '}' {
			s.popParserState()
			return scanEndObject
		}
	}

	return s.error(c, "No match char in the EndValue state")
}

func (s *scanner) error(c byte, context string) int {
	s.err = errors.New(quoteChar(c) + "is " + context)
	return scanError
}

// quoteChar formats c as a quoted character literal
func quoteChar(c byte) string {
	// special cases - different from quoted strings
	if c == '\'' {
		return `'\''`
	}
	if c == '"' {
		return `'"'`
	}

	// use quoted string with different quotation marks
	s := strconv.Quote(string(c))
	return "'" + s[1:len(s)-1] + "'"
}

func stepT(scan *scanner, b byte) int {
	if b == 'r' {
		scan.step = stepTr
		return scanContinue
	}

	return scan.error(b, "in true expect r")
}

func stepTr(scan *scanner, b byte) int {
	if b == 'u' {
		scan.step = stepTru
		return scanContinue
	}

	return scan.error(b, "in true expect u")
}

func stepTru(scan *scanner, b byte) int {
	if b == 'e' {
		scan.step = stateScanEndValue
		return scanContinue
	}

	return scan.error(b, "in true expect e")
}

func stepF(scan *scanner, b byte) int {
	if b == 'a' {
		scan.step = stepFa
		return scanContinue
	}

	return scan.error(b, "in false expect a")
}

func stepFa(scan *scanner, b byte) int {
	if b == 'l' {
		scan.step = stepFal
		return scanContinue
	}

	return scan.error(b, "in false expect l")
}

func stepFal(scan *scanner, b byte) int {
	if b == 's' {
		scan.step = stepFals
		return scanContinue
	}

	return scan.error(b, "in false expect s")
}

func stepFals(scan *scanner, b byte) int {
	if b == 'e' {
		scan.step = stateScanEndValue
		return scanContinue
	}

	return scan.error(b, "in false expect e")
}

func stepDecmial(scan *scanner, b byte) int {
	if b >= '0' && b <= '9' {
		scan.step = stepDecmial
		return scanContinue
	}

	return stateScanEndValue(scan, b)
}
