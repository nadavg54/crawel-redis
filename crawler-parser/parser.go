package parser

import (
	_ "bufio"
	_ "io"
)

//func extractUrls()

type LineParser interface {
	isValid(line string) bool
	parse() string
}

type LineFilter interface {
	isValid(line string) bool
}

type LinePusher interface {
	push(line string)
	pushLines(lines []string)
}

type LineReader interface{
	read() string
	readAll() []string
}

func parse(input []string, filter LineFilter, linePusher LinePusher){

	for _,line := range input {
		if filter.isValid(line) {
			linePusher.push(line)
		}
	}
}
