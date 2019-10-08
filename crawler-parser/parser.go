package parser

import (
	_ "bufio"
	_ "io"
)

//LineFilter validates lines are in a specific format like - a number, regualr text, type of HTML token
type LineFilter interface {
	isValid(line interface{}) bool
}

//LineEmitter emits lines
type LineEmitter interface {
	push(line interface{})
	pushLines(lines []interface{})
}

//FilterEmitter recieve array of lines and emit lines that were valid
func FilterEmitter(input []string, filter LineFilter, lineEmitter LineEmitter) {

	for _, line := range input {
		if filter.isValid(line) {
			lineEmitter.push(line)
		}
	}
}
