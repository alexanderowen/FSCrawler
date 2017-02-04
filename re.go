package main

import (
	"fmt"
	"os"
	"regexp"
)

func convertToRegexp(pat string) string {
	var reg string
	for _, char := range pat {
		switch char {
		case '*':
			reg = reg + ".*"
		case '.':
			reg = reg + "\\."
		case '?':
			reg = reg + "."
		default:
			reg = reg + string(char)
		}
	}
	return reg
}

func main() {
	pattern := os.Args[1]
	str := os.Args[2]
	fmt.Println("Pattern: " + pattern)
	fmt.Println("String: " + convertToRegexp(pattern))

	var validID = regexp.MustCompile(convertToRegexp(pattern))
	fmt.Printf("Match?: %t\n", validID.MatchString(str))
}
